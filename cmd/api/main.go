package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-backend-boilerplate/internal/cache"
	"go-backend-boilerplate/internal/config"
	"go-backend-boilerplate/internal/database"
	"go-backend-boilerplate/internal/observability"
	"go-backend-boilerplate/internal/server"
	"go-backend-boilerplate/internal/shared/logging"
	"go-backend-boilerplate/internal/storage"
	"go-backend-boilerplate/internal/worker"
)

func main() {
	cfg := config.Load()
	logging.Setup(cfg.AppEnv)

	if err := cfg.Validate(); err != nil {
		slog.Error("config validation", "err", err)
		os.Exit(1)
	}

	shutdownOTel, err := observability.Setup(context.Background(), cfg)
	if err != nil {
		slog.Error("otel setup", "err", err)
		os.Exit(1)
	}

	if err := database.EnsureDatabase(cfg); err != nil {
		slog.Error("ensure database", "err", err)
		os.Exit(1)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("connect database", "err", err)
		os.Exit(1)
	}

	if cfg.AutoMigrate {
		if err := database.AutoMigrate(db); err != nil {
			slog.Error("auto-migrate", "err", err)
			os.Exit(1)
		}
		slog.Info("auto-migrate completed")
	}

	rdb, err := cache.Connect(cfg)
	if err != nil {
		slog.Warn("redis not reachable", "addr", cfg.RedisAddr(), "err", err)
	}

	store, err := storage.New(context.Background(), cfg)
	if err != nil {
		slog.Error("storage init", "err", err)
		os.Exit(1)
	}

	pool := worker.New(4, 128)
	pool.Start()

	app := server.New(cfg, db, rdb, store, pool)

	go func() {
		if err := app.Listen(cfg.Addr()); err != nil {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()
	slog.Info("server listening", "addr", cfg.Addr(), "env", cfg.AppEnv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server")

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		slog.Error("forced shutdown", "err", err)
	}
	drainCtx, cancelDrain := context.WithTimeout(context.Background(), 10*time.Second)
	pool.Shutdown(drainCtx)
	cancelDrain()
	if rdb != nil {
		_ = rdb.Close()
	}
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = shutdownOTel(shutdownCtx)
	slog.Info("server stopped")
}

package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-backend-boilerplate/internal/cache"
	"go-backend-boilerplate/internal/config"
	"go-backend-boilerplate/internal/database"
	"go-backend-boilerplate/internal/server"
	"go-backend-boilerplate/internal/shared/logging"
)

func main() {
	cfg := config.Load()
	logging.Setup(cfg.AppEnv)

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

	app := server.New(cfg, db, rdb)

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
	if rdb != nil {
		_ = rdb.Close()
	}
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	slog.Info("server stopped")
}

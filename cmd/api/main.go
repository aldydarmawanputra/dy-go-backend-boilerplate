package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-backend-boilerplate/internal/cache"
	"go-backend-boilerplate/internal/config"
	"go-backend-boilerplate/internal/database"
	"go-backend-boilerplate/internal/server"
)

func main() {
	cfg := config.Load()

	if err := database.EnsureDatabase(cfg); err != nil {
		log.Fatalf("failed to ensure database exists: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if cfg.AutoMigrate {
		if err := database.AutoMigrate(db); err != nil {
			log.Fatalf("auto-migrate failed: %v", err)
		}
		log.Println("auto-migrate completed")
	}

	rdb, err := cache.Connect(cfg)
	if err != nil {
		log.Printf("warning: redis not reachable at %s: %v", cfg.RedisAddr(), err)
	}

	app := server.New(cfg, db, rdb)

	go func() {
		if err := app.Listen(cfg.Addr()); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
	log.Printf("server listening on %s (env=%s)", cfg.Addr(), cfg.AppEnv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Printf("forced shutdown: %v", err)
	}
	if rdb != nil {
		_ = rdb.Close()
	}
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	log.Println("server stopped")
}

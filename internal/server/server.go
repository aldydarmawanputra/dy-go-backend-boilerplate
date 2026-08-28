package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"go-backend-boilerplate/internal/config"
	"go-backend-boilerplate/internal/middleware"
	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/jwtutil"
	"go-backend-boilerplate/internal/shared/response"
	"go-backend-boilerplate/internal/storage"
)

func New(cfg *config.Config, db *gorm.DB, rdb *redis.Client, store storage.Storage) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "go-backend-boilerplate",
		ErrorHandler: errorHandler,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSec) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeoutSec) * time.Second,
		BodyLimit:    cfg.BodyLimitBytes,
	})

	middleware.Setup(app, cfg)

	jwtMgr := jwtutil.New(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAudience, cfg.JWTExpireHours)
	registerRoutes(app, cfg, db, rdb, store, jwtMgr)

	return app
}

func errorHandler(c *fiber.Ctx, err error) error {
	if fe, ok := err.(*fiber.Error); ok {
		return response.Error(c, apperror.New(fe.Code, "HTTP_ERROR", fe.Message))
	}
	return response.Error(c, err)
}

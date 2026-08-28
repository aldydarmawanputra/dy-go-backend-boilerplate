package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"go-backend-boilerplate/internal/config"
	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/response"
)

func Setup(app *fiber.App, cfg *config.Config) {
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(requestLogger())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowMethods:     cfg.CORSAllowMethods,
		AllowHeaders:     cfg.CORSAllowHeaders,
		AllowCredentials: cfg.CORSAllowCredentials,
	}))
	app.Use(RateLimit(cfg.RateLimitMax, cfg.RateLimitWindowSec))
}

func requestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		attrs := []any{
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"latency_ms", time.Since(start).Milliseconds(),
			"request_id", c.Locals("requestid"),
			"ip", c.IP(),
		}
		if err != nil {
			attrs = append(attrs, "error", err.Error())
		}
		slog.Info("request", attrs...)
		return err
	}
}

func RateLimit(max, windowSec int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: time.Duration(windowSec) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return response.Error(c, apperror.New(fiber.StatusTooManyRequests, "TOO_MANY_REQUESTS", "rate limit exceeded"))
		},
	})
}

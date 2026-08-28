package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} ${latency} ${method} ${path} (req=${locals:requestid})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowMethods:     cfg.CORSAllowMethods,
		AllowHeaders:     cfg.CORSAllowHeaders,
		AllowCredentials: cfg.CORSAllowCredentials,
	}))
	app.Use(RateLimit(cfg.RateLimitMax, cfg.RateLimitWindowSec))
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

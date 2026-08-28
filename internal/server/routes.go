package server

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"go-backend-boilerplate/internal/config"
	"go-backend-boilerplate/internal/mailer"
	"go-backend-boilerplate/internal/middleware"
	apikeymod "go-backend-boilerplate/internal/modules/apikey"
	authmod "go-backend-boilerplate/internal/modules/auth"
	docsmod "go-backend-boilerplate/internal/modules/docs"
	filemod "go-backend-boilerplate/internal/modules/file"
	paymentmod "go-backend-boilerplate/internal/modules/payment"
	rolemod "go-backend-boilerplate/internal/modules/role"
	usermod "go-backend-boilerplate/internal/modules/user"
	paygw "go-backend-boilerplate/internal/payment"
	"go-backend-boilerplate/internal/realtime"
	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/jwtutil"
	"go-backend-boilerplate/internal/shared/response"
	"go-backend-boilerplate/internal/storage"
	"go-backend-boilerplate/internal/worker"
)

func registerRoutes(app *fiber.App, cfg *config.Config, db *gorm.DB, rdb *redis.Client, store storage.Storage, hub *realtime.Hub, pool *worker.Pool, jwtMgr *jwtutil.Manager) {
	roleRepo := rolemod.NewRepository(db)

	userRepo := usermod.NewRepository(db)
	userSvc := usermod.NewService(userRepo, roleRepo)
	userHandler := usermod.NewHandler(userSvc)

	refreshStore := authmod.NewRefreshStore(rdb, cfg.RefreshExpireHours)
	tokenStore := authmod.NewTokenStore(rdb)
	mail := mailer.New(cfg)
	authSvc := authmod.NewService(cfg, userSvc, userRepo, roleRepo, jwtMgr, refreshStore, tokenStore, mail, pool)
	authHandler := authmod.NewHandler(authSvc)

	app.Get("/health", healthHandler(db, rdb))

	docsHandler := docsmod.NewHandler()
	app.Get("/documentation", func(c *fiber.Ctx) error {
		return c.Redirect("/documentation/en", fiber.StatusFound)
	})
	app.Get("/documentation/en", docsHandler.English)
	app.Get("/documentation/id", docsHandler.Indonesian)

	app.Use("/ws", realtime.Upgrade)
	app.Get("/ws", hub.Handler())

	v1 := app.Group("/api/v1")

	v1.Post("/broadcast", middleware.Auth(jwtMgr), middleware.RequireRole(rolemod.Admin), func(c *fiber.Ctx) error {
		var body struct {
			Message string `json:"message"`
		}
		if err := c.BodyParser(&body); err != nil {
			return response.Error(c, apperror.BadRequest("invalid request body"))
		}
		hub.Broadcast([]byte(body.Message))
		return response.OK(c, fiber.Map{"clients": hub.Count()})
	})

	auth := v1.Group("/auth", middleware.RateLimit(cfg.AuthRateLimitMax, cfg.RateLimitWindowSec, middleware.RedisStorage(rdb)))
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)
	auth.Post("/verify-email", authHandler.VerifyEmail)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)

	if cfg.StorageDriver == "local" || cfg.StorageDriver == "" {
		app.Static("/storage", cfg.StorageLocalPath)
	}
	fileHandler := filemod.NewHandler(store)
	files := v1.Group("/files", middleware.Auth(jwtMgr))
	files.Post("/", fileHandler.Upload)

	if gw, err := paygw.New(cfg); err != nil {
		slog.Warn("payment gateway disabled", "err", err)
	} else {
		paymentSvc := paymentmod.NewService(paymentmod.NewRepository(db), gw)
		paymentHandler := paymentmod.NewHandler(paymentSvc, cfg.PaymentCurrency)
		payments := v1.Group("/payments")
		payments.Post("/", middleware.Auth(jwtMgr), paymentHandler.Create)
		payments.Post("/webhook", paymentHandler.Webhook)
	}

	apiKeySvc := apikeymod.NewService(apikeymod.NewRepository(db))
	apiKeyHandler := apikeymod.NewHandler(apiKeySvc)
	keys := v1.Group("/api-keys", middleware.Auth(jwtMgr), middleware.RequireRole(rolemod.Admin))
	keys.Post("/", apiKeyHandler.Create)
	// demo endpoint authenticated by API key instead of a user JWT
	v1.Get("/whoami-apikey", middleware.APIKeyAuth(apiKeySvc), func(c *fiber.Ctx) error {
		return response.OK(c, fiber.Map{"api_key_id": c.Locals("apiKeyID"), "name": c.Locals("apiKeyName")})
	})

	users := v1.Group("/users", middleware.Auth(jwtMgr))
	users.Get("/", userHandler.List)
	users.Post("/", middleware.RequireRole(rolemod.Admin), userHandler.Create)
	users.Get("/search", userHandler.Search)
	users.Get("/me", userHandler.Me)
	users.Get("/:id", middleware.RequireSelfOrAdmin("id"), userHandler.GetByID)
	users.Put("/:id", middleware.RequireSelfOrAdmin("id"), userHandler.Replace)
	users.Patch("/:id", middleware.RequireSelfOrAdmin("id"), userHandler.Patch)
	users.Delete("/:id", middleware.RequireRole(rolemod.Admin), userHandler.Delete)
}

func healthHandler(db *gorm.DB, rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dbStatus := "ok"
		if sqlDB, err := db.DB(); err != nil || sqlDB.PingContext(c.Context()) != nil {
			dbStatus = "down"
		}
		redisStatus := "ok"
		if rdb == nil || rdb.Ping(c.Context()).Err() != nil {
			redisStatus = "down"
		}
		return response.OK(c, fiber.Map{
			"status": "ok",
			"db":     dbStatus,
			"redis":  redisStatus,
		})
	}
}

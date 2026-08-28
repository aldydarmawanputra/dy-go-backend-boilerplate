package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"go-backend-boilerplate/internal/config"
	"go-backend-boilerplate/internal/middleware"
	authmod "go-backend-boilerplate/internal/modules/auth"
	docsmod "go-backend-boilerplate/internal/modules/docs"
	filemod "go-backend-boilerplate/internal/modules/file"
	rolemod "go-backend-boilerplate/internal/modules/role"
	usermod "go-backend-boilerplate/internal/modules/user"
	"go-backend-boilerplate/internal/shared/jwtutil"
	"go-backend-boilerplate/internal/shared/response"
	"go-backend-boilerplate/internal/storage"
)

func registerRoutes(app *fiber.App, cfg *config.Config, db *gorm.DB, rdb *redis.Client, store storage.Storage, jwtMgr *jwtutil.Manager) {
	roleRepo := rolemod.NewRepository(db)

	userRepo := usermod.NewRepository(db)
	userSvc := usermod.NewService(userRepo, roleRepo)
	userHandler := usermod.NewHandler(userSvc)

	refreshStore := authmod.NewRefreshStore(rdb, cfg.RefreshExpireHours)
	authSvc := authmod.NewService(userSvc, userRepo, roleRepo, jwtMgr, refreshStore)
	authHandler := authmod.NewHandler(authSvc)

	app.Get("/health", healthHandler(db, rdb))

	docsHandler := docsmod.NewHandler()
	app.Get("/documentation", func(c *fiber.Ctx) error {
		return c.Redirect("/documentation/en", fiber.StatusFound)
	})
	app.Get("/documentation/en", docsHandler.English)
	app.Get("/documentation/id", docsHandler.Indonesian)

	v1 := app.Group("/api/v1")

	auth := v1.Group("/auth", middleware.RateLimit(cfg.AuthRateLimitMax, cfg.RateLimitWindowSec))
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	if cfg.StorageDriver == "local" || cfg.StorageDriver == "" {
		app.Static("/storage", cfg.StorageLocalPath)
	}
	fileHandler := filemod.NewHandler(store)
	files := v1.Group("/files", middleware.Auth(jwtMgr))
	files.Post("/", fileHandler.Upload)

	users := v1.Group("/users", middleware.Auth(jwtMgr))
	users.Get("/", userHandler.List)
	users.Post("/", middleware.RequireRole(rolemod.Admin), userHandler.Create)
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

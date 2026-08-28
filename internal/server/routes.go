package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"go-backend-boilerplate/internal/middleware"
	authmod "go-backend-boilerplate/internal/modules/auth"
	docsmod "go-backend-boilerplate/internal/modules/docs"
	usermod "go-backend-boilerplate/internal/modules/user"
	"go-backend-boilerplate/internal/shared/jwtutil"
	"go-backend-boilerplate/internal/shared/response"
)

func registerRoutes(app *fiber.App, db *gorm.DB, rdb *redis.Client, jwtMgr *jwtutil.Manager) {
	userRepo := usermod.NewRepository(db)
	userSvc := usermod.NewService(userRepo)
	userHandler := usermod.NewHandler(userSvc)

	authSvc := authmod.NewService(userSvc, userRepo, jwtMgr)
	authHandler := authmod.NewHandler(authSvc)

	app.Get("/health", healthHandler(db, rdb))

	docsHandler := docsmod.NewHandler()
	app.Get("/documentation", func(c *fiber.Ctx) error {
		return c.Redirect("/documentation/en", fiber.StatusFound)
	})
	app.Get("/documentation/en", docsHandler.English)
	app.Get("/documentation/id", docsHandler.Indonesian)

	v1 := app.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	users := v1.Group("/users", middleware.Auth(jwtMgr))
	users.Get("/", userHandler.List)
	users.Post("/", userHandler.Create)
	users.Get("/me", userHandler.Me)
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Replace)
	users.Patch("/:id", userHandler.Patch)
	users.Delete("/:id", userHandler.Delete)
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

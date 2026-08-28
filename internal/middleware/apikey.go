package middleware

import (
	"github.com/gofiber/fiber/v2"

	apikeymod "go-backend-boilerplate/internal/modules/apikey"
	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/response"
)

// APIKeyAuth authenticates a request via the X-API-Key header (machine-to-machine).
func APIKeyAuth(svc apikeymod.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-API-Key")
		if key == "" {
			return response.Error(c, apperror.Unauthorized("missing api key"))
		}
		rec, err := svc.Authenticate(c.Context(), key)
		if err != nil {
			return response.Error(c, apperror.Internal("failed to verify api key"))
		}
		if rec == nil {
			return response.Error(c, apperror.Unauthorized("invalid api key"))
		}
		c.Locals("apiKeyID", rec.ID)
		c.Locals("apiKeyName", rec.Name)
		return c.Next()
	}
}

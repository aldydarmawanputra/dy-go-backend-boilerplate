package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/jwtutil"
	"go-backend-boilerplate/internal/shared/response"
)

func Auth(jwt *jwtutil.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		if header == "" {
			return response.Error(c, apperror.Unauthorized("missing authorization header"))
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return response.Error(c, apperror.Unauthorized("invalid authorization header"))
		}

		claims, err := jwt.Parse(strings.TrimSpace(parts[1]))
		if err != nil {
			return response.Error(c, apperror.Unauthorized("invalid or expired token"))
		}

		c.Locals("userID", claims.Subject)
		c.Locals("roles", claims.Roles)
		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		have, _ := c.Locals("roles").([]string)
		for _, need := range roles {
			for _, h := range have {
				if h == need {
					return c.Next()
				}
			}
		}
		return response.Error(c, apperror.Forbidden("insufficient role"))
	}
}

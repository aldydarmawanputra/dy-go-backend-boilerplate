package middleware

import (
	"github.com/gofiber/fiber/v2"

	"go-backend-boilerplate/internal/shared/i18n"
)

// Locale resolves the request locale (from ?lang= or Accept-Language) and stores
// it in c.Locals("locale") for downstream localization.
func Locale() fiber.Handler {
	return func(c *fiber.Ctx) error {
		loc := c.Query("lang")
		if !i18n.Supported(loc) {
			loc = i18n.Resolve(c.Get(fiber.HeaderAcceptLanguage))
		}
		c.Locals("locale", loc)
		return c.Next()
	}
}

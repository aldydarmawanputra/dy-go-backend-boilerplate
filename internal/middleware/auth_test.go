package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func withRoles(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("roles", roles)
		return c.Next()
	}
}

func ok(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }

func TestRequireRole(t *testing.T) {
	app := fiber.New()
	app.Get("/denied", withRoles([]string{"user"}), RequireRole("admin"), ok)
	app.Get("/allowed", withRoles([]string{"admin"}), RequireRole("admin"), ok)
	app.Get("/no-roles", RequireRole("admin"), ok)

	cases := []struct {
		path string
		want int
	}{
		{"/denied", fiber.StatusForbidden},
		{"/allowed", fiber.StatusOK},
		{"/no-roles", fiber.StatusForbidden},
	}
	for _, tc := range cases {
		resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, tc.path, nil))
		if err != nil {
			t.Fatalf("%s: %v", tc.path, err)
		}
		if resp.StatusCode != tc.want {
			t.Fatalf("%s: status %d, want %d", tc.path, resp.StatusCode, tc.want)
		}
	}
}

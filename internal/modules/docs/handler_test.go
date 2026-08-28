package docs

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestDocsRoutes(t *testing.T) {
	app := fiber.New()
	h := NewHandler()
	app.Get("/documentation/en", h.English)
	app.Get("/documentation/id", h.Indonesian)

	cases := []struct{ path, want string }{
		{"/documentation/en", "# API Documentation"},
		{"/documentation/id", "# Dokumentasi API"},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(fiber.MethodGet, tc.path, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("%s: %v", tc.path, err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("%s: status %d", tc.path, resp.StatusCode)
		}
		if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/markdown") {
			t.Fatalf("%s: content-type %q", tc.path, ct)
		}
		body, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(body), tc.want) {
			t.Fatalf("%s: body missing %q", tc.path, tc.want)
		}
	}
}

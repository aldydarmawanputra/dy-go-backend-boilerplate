package docs

import (
	"embed"

	"github.com/gofiber/fiber/v2"

	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/response"
)

//go:embed templates/*.md
var files embed.FS

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) English(c *fiber.Ctx) error {
	return render(c, "templates/en.md")
}

func (h *Handler) Indonesian(c *fiber.Ctx) error {
	return render(c, "templates/id.md")
}

func render(c *fiber.Ctx, name string) error {
	b, err := files.ReadFile(name)
	if err != nil {
		return response.Error(c, apperror.NotFound("documentation not found"))
	}
	c.Set(fiber.HeaderContentType, "text/markdown; charset=utf-8")
	return c.Send(b)
}

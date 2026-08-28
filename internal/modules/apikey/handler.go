package apikey

import (
	"github.com/gofiber/fiber/v2"

	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/response"
	"go-backend-boilerplate/internal/shared/sanitize"
	"go-backend-boilerplate/internal/shared/validator"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

type createRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}

	plaintext, k, err := h.svc.Create(c.Context(), sanitize.String(req.Name))
	if err != nil {
		return response.Error(c, err)
	}
	// The plaintext key is returned only here; only its hash is stored.
	return response.Created(c, fiber.Map{
		"id":     k.ID,
		"name":   k.Name,
		"prefix": k.Prefix,
		"key":    plaintext,
	})
}

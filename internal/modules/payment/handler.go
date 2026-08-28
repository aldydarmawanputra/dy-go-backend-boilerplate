package payment

import (
	"github.com/gofiber/fiber/v2"

	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/response"
	"go-backend-boilerplate/internal/shared/validator"
)

type Handler struct {
	svc      Service
	currency string
}

func NewHandler(svc Service, currency string) *Handler {
	return &Handler{svc: svc, currency: currency}
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	p, err := h.svc.Create(c.Context(), req, h.currency)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, p)
}

// Webhook is called by the payment provider. In a real integration, verify the
// provider signature here before trusting the payload.
func (h *Handler) Webhook(c *fiber.Ctx) error {
	signature := c.Get("X-Signature")
	if err := h.svc.HandleWebhook(c.Context(), c.Body(), signature); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, fiber.Map{"received": true})
}

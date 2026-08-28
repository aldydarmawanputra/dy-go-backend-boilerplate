package auth

import (
	"github.com/gofiber/fiber/v2"

	"go-backend-boilerplate/internal/modules/user"
	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/response"
	"go-backend-boilerplate/internal/shared/validator"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	u, err := h.svc.Register(c.Context(), req)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, user.ToResponse(u))
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	token, err := h.svc.Login(c.Context(), req)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, TokenResponse{AccessToken: token, TokenType: "Bearer"})
}

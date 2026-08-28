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
	tokens, err := h.svc.Login(c.Context(), req)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, toTokenResponse(tokens))
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	tokens, err := h.svc.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, toTokenResponse(tokens))
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	if err := h.svc.Logout(c.Context(), req.RefreshToken); err != nil {
		return response.Error(c, err)
	}
	return response.NoContent(c)
}

func (h *Handler) VerifyEmail(c *fiber.Ctx) error {
	var req VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	if err := h.svc.VerifyEmail(c.Context(), req.Email, req.Code); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, fiber.Map{"verified": true})
}

func (h *Handler) ResendVerification(c *fiber.Ctx) error {
	var req ResendVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	if err := h.svc.ResendVerification(c.Context(), req.Email); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, fiber.Map{"message": "if the email exists and is unverified, a new code has been sent"})
}

func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	if err := h.svc.ForgotPassword(c.Context(), req.Email); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, fiber.Map{"message": "if the email exists, a reset link has been sent"})
}

func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	if err := h.svc.ResetPassword(c.Context(), req.Token, req.Password); err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, fiber.Map{"message": "password updated"})
}

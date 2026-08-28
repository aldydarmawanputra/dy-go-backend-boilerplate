package user

import (
	"github.com/gofiber/fiber/v2"

	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/pagination"
	"go-backend-boilerplate/internal/shared/response"
	"go-backend-boilerplate/internal/shared/validator"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Me(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(string)
	if userID == "" {
		return response.Error(c, apperror.Unauthorized("missing authentication"))
	}
	u, err := h.svc.GetByID(c.Context(), userID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, ToResponse(u))
}

func (h *Handler) List(c *fiber.Ctx) error {
	keyword := c.Query("q", "")
	p := pagination.Parse(c.Query("limit"), c.Query("offset"))

	users, total, err := h.svc.Search(c.Context(), keyword, p.Limit, p.Offset)
	if err != nil {
		return response.Error(c, err)
	}
	return response.WithMeta(c, ToResponseList(users), pagination.NewMeta(p, len(users), total))
}

func (h *Handler) Search(c *fiber.Ctx) error {
	query := c.Query("q", "")
	p := pagination.Parse(c.Query("limit"), c.Query("offset"))

	users, total, err := h.svc.SearchFullText(c.Context(), query, p.Limit, p.Offset)
	if err != nil {
		return response.Error(c, err)
	}
	return response.WithMeta(c, ToResponseList(users), pagination.NewMeta(p, len(users), total))
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	u, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, ToResponse(u))
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	u, err := h.svc.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, ToResponse(u))
}

func (h *Handler) Replace(c *fiber.Ctx) error {
	var req ReplaceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	u, err := h.svc.Replace(c.Context(), c.Params("id"), req)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, ToResponse(u))
}

func (h *Handler) Patch(c *fiber.Ctx) error {
	var req PatchRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, apperror.BadRequest("invalid request body"))
	}
	if err := validator.Struct(req); err != nil {
		return response.Error(c, err)
	}
	u, err := h.svc.Patch(c.Context(), c.Params("id"), req)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, ToResponse(u))
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Context(), c.Params("id")); err != nil {
		return response.Error(c, err)
	}
	return response.NoContent(c)
}

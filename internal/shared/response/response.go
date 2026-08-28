package response

import (
	"github.com/gofiber/fiber/v2"

	"go-backend-boilerplate/internal/shared/apperror"
)

type Envelope struct {
	Success bool            `json:"success"`
	Data    any             `json:"data,omitempty"`
	Meta    any             `json:"meta,omitempty"`
	Error   *apperror.Error `json:"error,omitempty"`
}

func success(c *fiber.Ctx, status int, data, meta any) error {
	return c.Status(status).JSON(Envelope{Success: true, Data: data, Meta: meta})
}

func OK(c *fiber.Ctx, data any) error { return success(c, fiber.StatusOK, data, nil) }

func Created(c *fiber.Ctx, data any) error { return success(c, fiber.StatusCreated, data, nil) }

func WithMeta(c *fiber.Ctx, data, meta any) error { return success(c, fiber.StatusOK, data, meta) }

func NoContent(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) }

func Error(c *fiber.Ctx, err error) error {
	appErr, ok := apperror.As(err)
	if !ok {
		appErr = apperror.Internal("something went wrong")
	}
	return c.Status(appErr.Status).JSON(Envelope{Success: false, Error: appErr})
}

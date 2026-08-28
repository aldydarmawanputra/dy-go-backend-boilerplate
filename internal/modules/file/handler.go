package file

import (
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/response"
	"go-backend-boilerplate/internal/storage"
)

type Handler struct {
	store storage.Storage
}

func NewHandler(store storage.Storage) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Upload(c *fiber.Ctx) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, apperror.BadRequest("file is required (multipart field 'file')"))
	}

	f, err := fh.Open()
	if err != nil {
		return response.Error(c, apperror.Internal("cannot open uploaded file"))
	}
	defer f.Close()

	ext := strings.ToLower(path.Ext(fh.Filename))
	key := "uploads/" + uuid.NewString() + ext

	url, err := h.store.Save(c.Context(), key, f, fh.Size, fh.Header.Get("Content-Type"))
	if err != nil {
		return response.Error(c, apperror.Internal("failed to store file"))
	}
	return response.Created(c, fiber.Map{"key": key, "url": url, "size": fh.Size})
}

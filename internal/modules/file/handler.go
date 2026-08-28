package file

import (
	"bytes"
	"io"
	"log/slog"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"go-backend-boilerplate/internal/imageproc"
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
	defer func() { _ = f.Close() }()

	data, err := io.ReadAll(f)
	if err != nil {
		return response.Error(c, apperror.Internal("cannot read uploaded file"))
	}

	id := uuid.NewString()
	ext := strings.ToLower(path.Ext(fh.Filename))
	key := "uploads/" + id + ext
	contentType := fh.Header.Get("Content-Type")

	url, err := h.store.Save(c.Context(), key, bytes.NewReader(data), int64(len(data)), contentType)
	if err != nil {
		return response.Error(c, apperror.Internal("failed to store file"))
	}

	out := fiber.Map{"key": key, "url": url, "size": len(data)}

	if imageproc.IsImage(contentType) {
		if thumb, thumbCT, terr := imageproc.Thumbnail(bytes.NewReader(data), 300, 300); terr == nil {
			thumbKey := "uploads/thumbs/" + id + ".jpg"
			if thumbURL, serr := h.store.Save(c.Context(), thumbKey, bytes.NewReader(thumb), int64(len(thumb)), thumbCT); serr == nil {
				out["thumbnail_key"] = thumbKey
				out["thumbnail_url"] = thumbURL
			} else {
				slog.Warn("store thumbnail", "err", serr)
			}
		} else {
			slog.Warn("generate thumbnail", "err", terr)
		}
	}

	return response.Created(c, out)
}

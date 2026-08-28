package storage

import (
	"context"
	"fmt"
	"io"

	"go-backend-boilerplate/internal/config"
)

type Storage interface {
	Save(ctx context.Context, key string, r io.Reader, size int64, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
	URL(key string) string
}

func New(ctx context.Context, cfg *config.Config) (Storage, error) {
	switch cfg.StorageDriver {
	case "local", "":
		return NewLocal(cfg.StorageLocalPath, cfg.StoragePublicBaseURL), nil
	case "r2":
		return NewR2(ctx, cfg)
	case "supabase":
		return NewSupabase(cfg), nil
	default:
		return nil, fmt.Errorf("unknown storage driver %q", cfg.StorageDriver)
	}
}

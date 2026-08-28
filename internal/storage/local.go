package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Local struct {
	basePath      string
	publicBaseURL string
}

func NewLocal(basePath, publicBaseURL string) *Local {
	return &Local{basePath: basePath, publicBaseURL: strings.TrimRight(publicBaseURL, "/")}
}

func (l *Local) safePath(key string) string {
	// Clean against a leading slash so any "../" in key cannot escape basePath.
	clean := filepath.Clean("/" + filepath.FromSlash(key))
	return filepath.Join(l.basePath, clean)
}

func (l *Local) Save(_ context.Context, key string, r io.Reader, _ int64, _ string) (string, error) {
	dest := l.safePath(key)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", err
	}
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}
	return l.URL(key), nil
}

func (l *Local) Delete(_ context.Context, key string) error {
	return os.Remove(l.safePath(key))
}

func (l *Local) URL(key string) string {
	return l.publicBaseURL + "/" + strings.TrimLeft(key, "/")
}

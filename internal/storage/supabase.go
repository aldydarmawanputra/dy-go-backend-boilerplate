package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	appconfig "go-backend-boilerplate/internal/config"
)

type Supabase struct {
	baseURL    string
	serviceKey string
	bucket     string
	publicBase string
	http       *http.Client
}

func NewSupabase(cfg *appconfig.Config) *Supabase {
	return &Supabase{
		baseURL:    strings.TrimRight(cfg.SupabaseURL, "/"),
		serviceKey: cfg.SupabaseServiceKey,
		bucket:     cfg.SupabaseBucket,
		publicBase: strings.TrimRight(cfg.SupabasePublicBaseURL, "/"),
		http:       &http.Client{},
	}
}

func (s *Supabase) objectURL(key string) string {
	return fmt.Sprintf("%s/storage/v1/object/%s/%s", s.baseURL, s.bucket, strings.TrimLeft(key, "/"))
}

func (s *Supabase) Save(ctx context.Context, key string, r io.Reader, _ int64, contentType string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.objectURL(key), r)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.serviceKey)
	req.Header.Set("x-upsert", "true")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase upload failed (%s): %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return s.URL(key), nil
}

func (s *Supabase) Delete(ctx context.Context, key string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, s.objectURL(key), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.serviceKey)

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("supabase delete failed: %s", resp.Status)
	}
	return nil
}

func (s *Supabase) URL(key string) string {
	if s.publicBase != "" {
		return s.publicBase + "/" + strings.TrimLeft(key, "/")
	}
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.baseURL, s.bucket, strings.TrimLeft(key, "/"))
}

package apikey

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"go-backend-boilerplate/internal/shared/apperror"
)

type Service interface {
	// Create returns the plaintext key (shown only once) plus the stored record.
	Create(ctx context.Context, name string) (string, *APIKey, error)
	// Authenticate returns the key record for a valid, non-revoked plaintext key,
	// or (nil, nil) when it is unknown/revoked.
	Authenticate(ctx context.Context, plaintext string) (*APIKey, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, name string) (string, *APIKey, error) {
	raw := make([]byte, 24)
	if _, err := rand.Read(raw); err != nil {
		return "", nil, apperror.Internal("failed to generate key")
	}
	plaintext := "sk_" + hex.EncodeToString(raw)

	k := &APIKey{
		Name:    name,
		KeyHash: hashKey(plaintext),
		Prefix:  plaintext[:11],
	}
	if err := s.repo.Create(ctx, k); err != nil {
		return "", nil, err
	}
	return plaintext, k, nil
}

func (s *service) Authenticate(ctx context.Context, plaintext string) (*APIKey, error) {
	k, err := s.repo.FindByHash(ctx, hashKey(plaintext))
	if err != nil {
		return nil, err
	}
	if k == nil || k.Revoked {
		return nil, nil
	}
	_ = s.repo.TouchLastUsed(ctx, k.ID)
	return k, nil
}

func hashKey(plaintext string) string {
	sum := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum[:])
}

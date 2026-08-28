package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrRefreshUnavailable = errors.New("refresh store unavailable")

type RefreshStore struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRefreshStore(rdb *redis.Client, ttlHours int) *RefreshStore {
	if ttlHours <= 0 {
		ttlHours = 168
	}
	return &RefreshStore{rdb: rdb, ttl: time.Duration(ttlHours) * time.Hour}
}

func (s *RefreshStore) Issue(ctx context.Context, userID string) (string, error) {
	if s.rdb == nil {
		return "", ErrRefreshUnavailable
	}
	token, err := randomToken(32)
	if err != nil {
		return "", err
	}
	if err := s.rdb.Set(ctx, refreshKey(token), userID, s.ttl).Err(); err != nil {
		return "", err
	}
	return token, nil
}

// UserID returns the owner of a refresh token, or ("", nil) if it is unknown or expired.
func (s *RefreshStore) UserID(ctx context.Context, token string) (string, error) {
	if s.rdb == nil {
		return "", ErrRefreshUnavailable
	}
	userID, err := s.rdb.Get(ctx, refreshKey(token)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s *RefreshStore) Revoke(ctx context.Context, token string) error {
	if s.rdb == nil {
		return ErrRefreshUnavailable
	}
	return s.rdb.Del(ctx, refreshKey(token)).Err()
}

func refreshKey(token string) string { return "refresh:" + token }

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

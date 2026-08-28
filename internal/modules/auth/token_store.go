package auth

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrTokenStoreUnavailable = errors.New("token store unavailable")

const (
	nsVerify = "verify"
	nsReset  = "reset"
)

// TokenStore issues and consumes single-use, expiring tokens (email verification,
// password reset) backed by Redis. Consume is atomic (get + delete).
type TokenStore struct {
	rdb *redis.Client
}

func NewTokenStore(rdb *redis.Client) *TokenStore {
	return &TokenStore{rdb: rdb}
}

func (s *TokenStore) Issue(ctx context.Context, namespace, userID string, ttl time.Duration) (string, error) {
	if s.rdb == nil {
		return "", ErrTokenStoreUnavailable
	}
	token, err := randomToken(32)
	if err != nil {
		return "", err
	}
	if err := s.rdb.Set(ctx, namespace+":"+token, userID, ttl).Err(); err != nil {
		return "", err
	}
	return token, nil
}

// Consume returns the user id for a token and deletes it, or ("", nil) if the
// token is unknown or already used/expired.
func (s *TokenStore) Consume(ctx context.Context, namespace, token string) (string, error) {
	if s.rdb == nil {
		return "", ErrTokenStoreUnavailable
	}
	userID, err := s.rdb.GetDel(ctx, namespace+":"+token).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return userID, nil
}

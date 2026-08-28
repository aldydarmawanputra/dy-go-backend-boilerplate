package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

// OTPStore issues and verifies short numeric email-verification codes, backed by
// Redis, with a per-code attempt limit.
type OTPStore struct {
	rdb         *redis.Client
	ttl         time.Duration
	maxAttempts int
}

func NewOTPStore(rdb *redis.Client, expireMinutes, maxAttempts int) *OTPStore {
	if expireMinutes <= 0 {
		expireMinutes = 10
	}
	if maxAttempts <= 0 {
		maxAttempts = 5
	}
	return &OTPStore{
		rdb:         rdb,
		ttl:         time.Duration(expireMinutes) * time.Minute,
		maxAttempts: maxAttempts,
	}
}

func codeKey(userID string) string     { return "verify_code:" + userID }
func attemptsKey(userID string) string { return "verify_attempts:" + userID }

func (s *OTPStore) Issue(ctx context.Context, userID string) (string, error) {
	if s.rdb == nil {
		return "", ErrTokenStoreUnavailable
	}
	code, err := randomCode()
	if err != nil {
		return "", err
	}
	if err := s.rdb.Set(ctx, codeKey(userID), code, s.ttl).Err(); err != nil {
		return "", err
	}
	s.rdb.Del(ctx, attemptsKey(userID))
	return code, nil
}

// Verify checks the code for userID (constant-time), clearing it on success.
// After maxAttempts wrong tries the code is invalidated to blunt brute force.
func (s *OTPStore) Verify(ctx context.Context, userID, code string) (bool, error) {
	if s.rdb == nil {
		return false, ErrTokenStoreUnavailable
	}
	stored, err := s.rdb.Get(ctx, codeKey(userID)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	attempts, err := s.rdb.Incr(ctx, attemptsKey(userID)).Result()
	if err != nil {
		return false, err
	}
	if attempts == 1 {
		s.rdb.Expire(ctx, attemptsKey(userID), s.ttl)
	}
	if int(attempts) > s.maxAttempts {
		s.rdb.Del(ctx, codeKey(userID), attemptsKey(userID))
		return false, nil
	}

	if subtle.ConstantTimeCompare([]byte(stored), []byte(code)) == 1 {
		s.rdb.Del(ctx, codeKey(userID), attemptsKey(userID))
		return true, nil
	}
	return false, nil
}

func randomCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

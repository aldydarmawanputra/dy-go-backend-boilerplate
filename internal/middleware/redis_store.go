package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type redisStorage struct {
	rdb    *redis.Client
	prefix string
}

// RedisStorage adapts a go-redis client to fiber.Storage so middleware such as
// the rate limiter shares state across instances. Returns nil (in-memory
// fallback) when rdb is nil.
func RedisStorage(rdb *redis.Client) fiber.Storage {
	if rdb == nil {
		return nil
	}
	return &redisStorage{rdb: rdb, prefix: "fiber:"}
}

func (s *redisStorage) Get(key string) ([]byte, error) {
	if key == "" {
		return nil, nil
	}
	v, err := s.rdb.Get(context.Background(), s.prefix+key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return v, err
}

func (s *redisStorage) Set(key string, val []byte, exp time.Duration) error {
	if key == "" || len(val) == 0 {
		return nil
	}
	return s.rdb.Set(context.Background(), s.prefix+key, val, exp).Err()
}

func (s *redisStorage) Delete(key string) error {
	if key == "" {
		return nil
	}
	return s.rdb.Del(context.Background(), s.prefix+key).Err()
}

func (s *redisStorage) Reset() error { return nil }

func (s *redisStorage) Close() error { return nil }

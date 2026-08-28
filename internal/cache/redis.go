package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"go-backend-boilerplate/internal/config"
)

func Connect(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return client, err
	}
	return client, nil
}

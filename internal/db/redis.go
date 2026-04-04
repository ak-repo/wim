package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ak-repo/wim/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, cfg config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Redis{Client: client}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}

// Get retrieves a value from the cache based on the provided key.
func (d *Redis) Get(ctx context.Context, key string) (string, error) {
	if d == nil || d.Client == nil {
		return "", fmt.Errorf("redis client is not initialized")
	}
	val, err := d.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// Set sets a key-value pair in the cache with an optional timeout.
func (d *Redis) Set(ctx context.Context, key string, value string, timeout int) error {
	if d == nil || d.Client == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	defaultTimeout := time.Duration(timeout * int(time.Minute))
	err := d.Client.Set(ctx, key, value, defaultTimeout).Err()
	return err
}

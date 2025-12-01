package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis represents a Redis client
type Redis struct {
	client *redis.Client
}

// New creates a new Redis client
func New(client *redis.Client) *Redis {
	return &Redis{}
}

// Close closes the Redis client connection
func (r *Redis) Close() error {
	return r.client.Close()
}

// Set sets a value in Redis with the specified TTL
func (r *Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}

		return err
	}

	return nil
}

// Get retrieves a value from Redis by key
func (r *Redis) Get(ctx context.Context, key string) ([]byte, bool, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return val, true, nil
}

// GetString retrieves a string value from Redis by key
func (r *Redis) GetString(ctx context.Context, key string) (string, bool, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false, nil
		}

		return "", false, err
	}

	return val, true, nil
}

// Del deletes a key from Redis
func (r *Redis) Del(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}

		return err
	}

	return nil
}

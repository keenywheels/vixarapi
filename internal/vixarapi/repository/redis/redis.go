package redis

import (
	"time"

	"github.com/keenywheels/backend/pkg/redis"
)

// Repository provides redis-related data access logic
type Repository struct {
	redis *redis.Redis
	ttl   time.Duration
}

// New creates new Repository instance
func New(redisClient *redis.Redis) (*Repository, error) {
	return &Repository{
		redis: redisClient,
		ttl:   time.Hour * 24 * 180, // default TTL is 180 days
	}, nil
}

// WithTTL sets custom TTL for the repository
func (r *Repository) WithTTL(ttl time.Duration) *Repository {
	r.ttl = ttl
	return r
}

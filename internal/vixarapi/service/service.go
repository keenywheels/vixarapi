package service

import (
	"context"

	"github.com/keenywheels/backend/internal/pkg/client/vk"
	"github.com/keenywheels/backend/internal/vixarapi/models"
	"github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
	"github.com/keenywheels/backend/internal/vixarapi/repository/redis"
)

// IRepository provides interface to communicate with the repository layer
type IRepository interface {
	SearchTokenInfo(context.Context, *postgres.SearchTokenParams) ([]models.TokenInfo, error)
	GetUserByVKID(context.Context, int64) (*models.User, error)
	RegisterVKUser(context.Context, *models.User) (*models.User, error)
}

// IRedisRepository provides interface to communicate with the redis repository layer
type IRedisRepository interface {
	SaveUserSession(ctx context.Context, sessionID string, userInfo *redis.UserInfo) error
	GetUserSession(ctx context.Context, sessionID string) (*redis.UserInfo, error)
	DeleteUserSession(ctx context.Context, userID string) error
	SaveVkTokens(ctx context.Context, key string, tokens *redis.VkTokens) error
	GetVkTokens(ctx context.Context, key string) (*redis.VkTokens, error)
	DeleteVkTokens(ctx context.Context, key string) error
}

// Config holds service configuration
type Config struct {
	SessionSecret string `mapstructure:"session_secret"`
}

// Service provides interest-related business logic
type Service struct {
	repo  IRepository
	redis IRedisRepository
	vk    *vk.Client

	cfg *Config
}

// New creates a new interest service
func New(
	repo IRepository,
	redis IRedisRepository,
	vk *vk.Client,
	cfg *Config,
) *Service {
	return &Service{
		repo:  repo,
		redis: redis,
		vk:    vk,
		cfg:   cfg,
	}
}

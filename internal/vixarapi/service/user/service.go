package user

import (
	"context"

	"github.com/keenywheels/backend/internal/pkg/client/vk"
	"github.com/keenywheels/backend/internal/vixarapi/models"
	"github.com/keenywheels/backend/internal/vixarapi/repository/postgres/search"
	"github.com/keenywheels/backend/internal/vixarapi/repository/postgres/user"
	"github.com/keenywheels/backend/internal/vixarapi/repository/redis/session"
)

// IRepository provides interface to communicate with the postgres repository layer
type IRepository interface {
	GetUserByVKID(context.Context, int64) (*models.User, error)
	RegisterVKUser(context.Context, *models.User) (*models.User, error)
	SaveSearchQuery(context.Context, string, string) (*models.UserQuery, error)
	DeleteSearchQuery(context.Context, string) error
	GetSearchQueries(context.Context, string, uint64, uint64) ([]*models.UserQuery, error)
	AddTokenSub(context.Context, *user.AddTokenSubParams) (string, error)
	GetTokenSubs(context.Context, string, uint64, uint64) ([]*models.UserTokenSub, error)
	DeleteTokenSub(context.Context, string) error
}

// ISearch provides interface to communicate with the search repository layer
type ISearch interface {
	GetLatestToken(ctx context.Context, params *search.GetTokenParams) (*models.Token, error)
}

// ISession provides interface to communicate with the session repository layer
type ISession interface {
	SaveUserSession(ctx context.Context, sessionID string, userInfo *session.UserInfo) error
	GetUserSession(ctx context.Context, sessionID string) (*session.UserInfo, error)
	DeleteUserSession(ctx context.Context, userID string) error
	SaveVkTokens(ctx context.Context, key string, tokens *session.VkTokens) error
	GetVkTokens(ctx context.Context, key string) (*session.VkTokens, error)
	DeleteVkTokens(ctx context.Context, key string) error
}

// Config holds service configuration
type Config struct {
	SessionSecret string `mapstructure:"session_secret"`
}

// Service provides interest-related business logic
type Service struct {
	repo IRepository
	sesh ISession
	srch ISearch
	vk   *vk.Client

	cfg *Config
}

// New creates a new interest service
func New(
	repo IRepository,
	sesh ISession,
	srch ISearch,
	vk *vk.Client,
	cfg *Config,
) *Service {
	return &Service{
		repo: repo,
		sesh: sesh,
		srch: srch,
		vk:   vk,
		cfg:  cfg,
	}
}

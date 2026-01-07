package user

import (
	"context"

	"github.com/keenywheels/backend/internal/vixarapi/delivery/http/cookie"
	service "github.com/keenywheels/backend/internal/vixarapi/service/user"
)

// IService provides user-related service logic
type IService interface {
	HandleVkAuthCallback(context.Context, *service.VkAuthCallbackParams) (*service.VkAuthCallbackResult, error)
	RegisterVkUser(context.Context, *service.RegisterVkUserParams) error
	LogoutUser(context.Context, string) error
	SaveSearchQuery(context.Context, *service.SaveQueryParams) (string, error)
	DeleteSearchQuery(context.Context, string) error
	GetSearchQueries(context.Context, *service.GetSearchQueriesParams) ([]service.Query, error)
	SubscribeToToken(ctx context.Context, params *service.SubscribeToTokenParams) (string, error)
	GetSubscribedTokens(ctx context.Context, userID string, limit, offset uint64) ([]*service.TokenSubInfo, error)
	UnsubscribeFromToken(ctx context.Context, id string) error
}

// Controller contains handlers for endpoints
type Controller struct {
	svc IService
	cm  *cookie.CookieManager
}

// New creates a new controller instance
func New(svc IService, cm *cookie.CookieManager) *Controller {
	return &Controller{
		svc: svc,
		cm:  cm,
	}
}

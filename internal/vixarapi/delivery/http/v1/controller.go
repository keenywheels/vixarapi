package v1

import (
	"context"

	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/vixarapi/service"
)

var _ gen.Handler = (*Controller)(nil)

// IService provides interest-related services
type IService interface {
	SearchTokenInfo(context.Context, *service.SearchTokenInfoParams) ([]service.TokenInfo, error)
	HandleVkAuthCallback(context.Context, *service.VkAuthCallbackParams) (*service.VkAuthCallbackResult, error)
	RegisterVkUser(context.Context, *service.RegisterVkUserParams) error
	LogoutUser(context.Context, string) error
}

// Controller contains handlers for interest-related endpoints
type Controller struct {
	svc IService
}

// New creates a new interest controller
func New(svc IService) *Controller {
	return &Controller{
		svc: svc,
	}
}

package search

import (
	"context"

	service "github.com/keenywheels/backend/internal/vixarapi/service/search"
)

// IService provides search-related service logic
type IService interface {
	SearchTokenInfo(context.Context, *service.SearchTokenInfoParams) ([]service.TokenInfo, error)
}

// Controller contains handlers for endpoints
type Controller struct {
	svc IService
}

// New creates a new controller instance
func New(svc IService) *Controller {
	return &Controller{
		svc: svc,
	}
}

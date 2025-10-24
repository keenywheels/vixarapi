package v1

import (
	"context"

	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/vixarapi/service"
)

var _ gen.Handler = (*Controller)(nil)

// IService provides interest-related services
type IService interface {
	GetAllInterest(context.Context, string) ([]service.Interest, error)
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

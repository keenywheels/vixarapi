package security

import (
	"context"

	gen "github.com/keenywheels/backend/internal/api/v1"
)

var _ gen.SecurityHandler = (*Controller)(nil)

// IService defines the interface for the security service
type IService interface {
	ValidateSession(ctx context.Context, session string) (bool, error)
}

// Controller contains security-related logic
type Controller struct {
	srvc IService
}

// New creates a new security controller
func New(srvc IService) *Controller {
	return &Controller{
		srvc: srvc,
	}
}

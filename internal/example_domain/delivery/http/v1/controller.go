package v1

import (
	gen "github.com/keenywheels/backend/internal/api/v1"
)

var _ gen.Handler = (*Controller)(nil)

// Controller ...
type Controller struct {
}

// NewController ...
func NewController() *Controller {
	return &Controller{}
}

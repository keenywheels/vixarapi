package v1

import (
	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/vixarapi/delivery/http/v1/search"
	"github.com/keenywheels/backend/internal/vixarapi/delivery/http/v1/user"
)

// check that Router implements ogen handler
var _ gen.Handler = (*Router)(nil)

// Router contains handlers for endpoints
type Router struct {
	searchController *search.Controller
	userController   *user.Controller
}

// New creates a new router instance
func New(
	searchController *search.Controller,
	userController *user.Controller,
) *Router {
	return &Router{
		searchController: searchController,
		userController:   userController,
	}
}

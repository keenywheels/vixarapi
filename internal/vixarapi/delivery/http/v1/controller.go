package v1

import (
	"context"

	gen "github.com/keenywheels/backend/internal/api/v1"
	searchSvc "github.com/keenywheels/backend/internal/vixarapi/service/search"
	userSvc "github.com/keenywheels/backend/internal/vixarapi/service/user"
)

var _ gen.Handler = (*Controller)(nil)

// ISearchService provides search-related service logic
type ISearchService interface {
	SearchTokenInfo(context.Context, *searchSvc.SearchTokenInfoParams) ([]searchSvc.TokenInfo, error)
}

// IUserService provides user-related service logic
type IUserService interface {
	HandleVkAuthCallback(context.Context, *userSvc.VkAuthCallbackParams) (*userSvc.VkAuthCallbackResult, error)
	RegisterVkUser(context.Context, *userSvc.RegisterVkUserParams) error
	LogoutUser(context.Context, string) error
	SaveSearchQuery(context.Context, *userSvc.SaveQueryParams) (string, error)
	DeleteSearchQuery(context.Context, string) error
	GetSearchQueries(context.Context, *userSvc.GetSearchQueriesParams) ([]userSvc.Query, error)
}

// Controller contains handlers for endpoints
type Controller struct {
	searchSvc ISearchService
	userSvc   IUserService
}

// New creates a new controller instance
func New(searchSvc ISearchService, userSvc IUserService) *Controller {
	return &Controller{
		searchSvc: searchSvc,
		userSvc:   userSvc,
	}
}

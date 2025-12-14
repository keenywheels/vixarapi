package user

import (
	"context"
	"errors"

	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/vixarapi/delivery/http/security"
	commonService "github.com/keenywheels/backend/internal/vixarapi/service"
	service "github.com/keenywheels/backend/internal/vixarapi/service/user"
	"github.com/keenywheels/backend/pkg/ctxutils"
	"github.com/keenywheels/backend/pkg/httputils"
)

// SaveUserQuery handles saving user search query
func (c *Controller) SaveUserQuery(
	ctx context.Context,
	req *gen.SaveUserQueryRequest,
) (gen.SaveUserQueryRes, error) {
	var (
		op  = "Controller.SaveUserQuery"
		log = ctxutils.GetLogger(ctx)
	)

	// validate request
	if err := req.Validate(); err != nil {
		log.Errorf("[%s] invalid request: %v", op, err)

		return &gen.SaveUserQueryBadRequest{
			Error: httputils.ErrorBadRequest,
		}, nil
	}

	// retrieve user info from context
	userInfo, ok := security.GetUserInfo(ctx)
	if !ok {
		log.Errorf("[%s] missing user info in context", op)

		return &gen.SaveUserQueryInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	// save user query
	id, err := c.svc.SaveSearchQuery(ctx, &service.SaveQueryParams{
		UserID: userInfo.ID,
		Query:  req.Query,
	})
	if err != nil {
		log.Errorf("[%s] failed to save user query: %v", op, err)

		return &gen.SaveUserQueryInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	return &gen.SaveUserQueryResponse{
		ID: id,
	}, nil
}

// DeleteUserSearchQuery handles deleting user search query
func (c *Controller) DeleteUserSearchQuery(
	ctx context.Context,
	params gen.DeleteUserSearchQueryParams,
) (gen.DeleteUserSearchQueryRes, error) {
	var (
		op  = "Controller.DeleteUserSearchQuery"
		log = ctxutils.GetLogger(ctx)
	)

	// delete user search query
	if err := c.svc.DeleteSearchQuery(ctx, params.ID); err != nil {
		switch {
		case errors.Is(err, commonService.ErrNotFound):
			return &gen.DeleteUserSearchQueryNotFound{
				Error: httputils.ErrorNotFound,
			}, nil
		}

		log.Errorf("[%s] failed to delete user search query: %v", op, err)

		return &gen.DeleteUserSearchQueryInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	return &gen.DeleteUserSearchQueryOK{}, nil
}

// GetUserSearchQueries handles retrieving user search queries
func (c *Controller) GetUserSearchQueries(
	ctx context.Context,
	params gen.GetUserSearchQueriesParams,
) (gen.GetUserSearchQueriesRes, error) {
	var (
		op  = "Controller.GetUserSearchQueries"
		log = ctxutils.GetLogger(ctx)
	)

	// retrieve user info from context
	userInfo, ok := security.GetUserInfo(ctx)
	if !ok {
		log.Errorf("[%s] missing user info in context", op)

		return &gen.GetUserSearchQueriesInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	// get user search queries
	queries, err := c.svc.GetSearchQueries(ctx, &service.GetSearchQueriesParams{
		UserID: userInfo.ID,
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		switch {
		case errors.Is(err, commonService.ErrNotFound):
			return &gen.GetUserSearchQueriesNotFound{
				Error: httputils.ErrorNotFound,
			}, nil
		}

		log.Errorf("[%s] failed to get user search queries: %v", op, err)

		return &gen.GetUserSearchQueriesInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	resp := gen.GetUserSearchQueriesOKApplicationJSON(converQueries(queries))

	return &resp, nil
}

// converQueries converts service layer queries to api layer queries
func converQueries(queries []service.Query) []gen.UserSearchQuery {
	resp := make([]gen.UserSearchQuery, 0, len(queries))
	for _, q := range queries {
		resp = append(resp, gen.UserSearchQuery{
			ID:         q.ID,
			Query:      q.Query,
			SearchDate: q.SearchDate,
		})
	}

	return resp
}

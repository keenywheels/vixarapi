package user

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/vixarapi/delivery/http/security"
	commonService "github.com/keenywheels/backend/internal/vixarapi/service"
	service "github.com/keenywheels/backend/internal/vixarapi/service/user"
	"github.com/keenywheels/backend/pkg/ctxutils"
	"github.com/keenywheels/backend/pkg/httputils"
)

const (
	methodDenormalized   = "denormalized"
	methodGlobalMedian   = "global_median"
	methodCategoryMedian = "category_median"
)

// SubscribeUserToToken subscribe the user to tokens update, so they can receive notifications
func (c *Controller) SubscribeUserToToken(
	ctx context.Context,
	req *gen.SubscribeUserToTokenRequest,
) (gen.SubscribeUserToTokenRes, error) {
	var (
		op  = "Controller.SubscribeUserToToken"
		log = ctxutils.GetLogger(ctx)
	)

	// validate request
	if err := req.Validate(); err != nil {
		log.Errorf("[%s] invalid request: %v", op, err)

		return &gen.SubscribeUserToTokenBadRequest{
			Error: httputils.ErrorBadRequest,
		}, nil
	}

	// retrieve user info from context
	userInfo, ok := security.GetUserInfo(ctx)
	if !ok {
		log.Errorf("[%s] missing user info in context", op)

		return &gen.SubscribeUserToTokenInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	// parse method
	method, err := parseMethod(req.Method)
	if err != nil {
		log.Errorf("[%s] invalid method: %v", op, err)

		return &gen.SubscribeUserToTokenBadRequest{
			Error: httputils.ErrorBadRequest,
		}, nil
	}

	// subscribe user to token
	id, err := c.svc.SubscribeToToken(ctx, &service.SubscribeToTokenParams{
		UserID:    userInfo.ID,
		Token:     req.Token,
		Category:  req.Category,
		Method:    method,
		Threshold: req.Threshold,
	})
	if err != nil {
		switch {
		case errors.Is(err, commonService.ErrAlreadyExists):
			return &gen.SubscribeUserToTokenConflict{
				Error: httputils.ErrorConflict,
			}, nil
		}

		log.Errorf("[%s] failed to subscribe user to token: %v", op, err)

		return &gen.SubscribeUserToTokenInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	return &gen.SubscribeUserToTokenResponse{
		ID: id,
	}, nil
}

// GetUserTokenSubs get all user's token subs
func (c *Controller) GetUserTokenSubs(
	ctx context.Context,
	params gen.GetUserTokenSubsParams,
) (gen.GetUserTokenSubsRes, error) {
	var (
		op  = "Controller.GetUserTokenSubs"
		log = ctxutils.GetLogger(ctx)
	)

	// retrieve user info from context
	userInfo, ok := security.GetUserInfo(ctx)
	if !ok {
		log.Errorf("[%s] missing user info in context", op)

		return &gen.GetUserTokenSubsInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	subs, err := c.svc.GetSubscribedTokens(ctx, userInfo.ID, params.Limit.Value, params.Offset.Value)
	if err != nil {
		switch {
		case errors.Is(err, commonService.ErrNotFound):
			return &gen.GetUserTokenSubsNotFound{
				Error: httputils.ErrorNotFound,
			}, nil
		}

		log.Errorf("[%s] failed to get user token subs: %v", op, err)

		return &gen.GetUserTokenSubsInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	resp := gen.GetUserTokenSubsOKApplicationJSON(convertTokenSubs(subs))

	return &resp, nil
}

// DeleteUserTokenSub delete specified token sub from user's subs
func (c *Controller) DeleteUserTokenSub(
	ctx context.Context,
	params gen.DeleteUserTokenSubParams,
) (gen.DeleteUserTokenSubRes, error) {
	var (
		op  = "Controller.DeleteUserTokenSub"
		log = ctxutils.GetLogger(ctx)
	)

	if err := c.svc.UnsubscribeFromToken(ctx, params.ID); err != nil {
		switch {
		case errors.Is(err, commonService.ErrNotFound):
			return &gen.DeleteUserTokenSubNotFound{
				Error: httputils.ErrorNotFound,
			}, nil
		}

		log.Errorf("[%s] failed to unsubscribe user from token: %v", op, err)

		return &gen.DeleteUserTokenSubInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	return &gen.DeleteUserTokenSubOK{}, nil
}

// parseMethod validate method and set default value if wasn't set
func parseMethod(reqMethod gen.OptString) (string, error) {
	// return denormalized method if not set
	if !reqMethod.Set {
		return methodDenormalized, nil
	}

	var (
		method       = strings.ToLower(reqMethod.Value)
		validMethods = []string{methodDenormalized, methodGlobalMedian, methodCategoryMedian}
	)

	if !slices.Contains(validMethods, method) {
		return "", fmt.Errorf("got unexpected method: %s", method)
	}

	return method, nil
}

func convertTokenSubs(subs []*service.TokenSubInfo) []gen.UserTokenSub {
	resp := make([]gen.UserTokenSub, 0, len(subs))
	for _, s := range subs {
		resp = append(resp, gen.UserTokenSub{
			ID:               s.ID,
			Token:            s.Token,
			Category:         s.Category,
			Method:           s.Method,
			CurrentInterest:  s.CurrentInterest,
			PreviousInterest: s.PreviousInterest,
			LastScan:         s.ScanDate,
		})
	}

	return resp
}

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

// VkAuthCallback handles the VK OAuth callback to retrieve tokens and user info.
func (c *Controller) VkAuthCallback(
	ctx context.Context,
	req *gen.VkAuthCallbackRequest,
) (gen.VkAuthCallbackRes, error) {
	var (
		op  = "Controller.VkAuthCallback"
		log = ctxutils.GetLogger(ctx)
	)

	if err := req.Validate(); err != nil {
		log.Errorf("[%s] invalid request: %v", op, err)

		return &gen.VkAuthCallbackBadRequest{
			Error: httputils.ErrorBadRequest,
		}, nil
	}

	res, err := c.svc.HandleVkAuthCallback(ctx, &service.VkAuthCallbackParams{
		Code:         req.Code,
		State:        req.State,
		CodeVerifier: req.CodeVerifier,
		DeviceID:     req.DeviceID,
		RedirectURI:  req.RedirectURI,
	})
	if err != nil {
		log.Errorf("[%s] failed to handle VK auth callback: %v", op, err)

		return &gen.VkAuthCallbackInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	return &gen.VkAuthCallbackResponseHeaders{
		Response: gen.VkAuthCallbackResponse{
			UserExists: res.UserExists,
			Username:   res.Username,
			Email:      res.Email,
			Vkid:       res.VKID,
		},
		SetCookie: gen.OptString{
			Value: c.cm.SessionCookie(res.Session).String(),
			Set:   true,
		},
	}, nil
}

// VkAuthRegister handles VK OAuth registration for new users.
func (c *Controller) VkAuthRegister(
	ctx context.Context,
	req *gen.VkAuthRegisterRequest,
) (gen.VkAuthRegisterRes, error) {
	var (
		op  = "Controller.VkAuthRegister"
		log = ctxutils.GetLogger(ctx)
	)

	if err := req.Validate(); err != nil {
		log.Errorf("[%s] invalid request: %v", op, err)

		return &gen.VkAuthRegisterBadRequest{
			Error: httputils.ErrorBadRequest,
		}, nil
	}

	sessionID, ok := security.GetSessionID(ctx)
	if !ok {
		log.Errorf("[%s] session ID not found in context", op)

		return &gen.VkAuthRegisterInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	if err := c.svc.RegisterVkUser(ctx, &service.RegisterVkUserParams{
		SessionID: sessionID,
		Email:     req.Email,
		Username:  req.Username,
		VKID:      req.Vkid,
	}); err != nil {
		switch {
		case errors.Is(err, commonService.ErrAlreadyExists):
			return &gen.VkAuthRegisterConflict{
				Error: httputils.ErrorConflict,
			}, nil
		}

		log.Errorf("[%s] failed to register VK user: %v", op, err)

		return &gen.VkAuthRegisterInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	return &gen.VkAuthRegisterOK{}, nil
}

// LogoutUser logs out the user by clearing the session
func (c *Controller) LogoutUser(ctx context.Context) (gen.LogoutUserRes, error) {
	var (
		op  = "Controller.LogoutUser"
		log = ctxutils.GetLogger(ctx)
	)

	sessionID, ok := security.GetSessionID(ctx)
	if !ok {
		log.Errorf("[%s] session ID not found in context", op)

		return &gen.LogoutUserInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	if err := c.svc.LogoutUser(ctx, sessionID); err != nil {
		log.Errorf("[%s] failed to logout user: %v", op, err)

		return &gen.LogoutUserInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	return &gen.LogoutUserOK{
		SetCookie: gen.OptString{
			Value: c.cm.SessionCookie("").String(),
			Set:   true,
		},
	}, nil
}

// UserInfo get info of logged-in user
func (c *Controller) UserInfo(ctx context.Context) (gen.UserInfoRes, error) {
	var (
		op  = "Controller.UserInfo"
		log = ctxutils.GetLogger(ctx)
	)

	// get user info from context
	userInfo, ok := security.GetUserInfo(ctx)
	if !ok {
		log.Errorf("[%s] user info not found in context", op)

		return &gen.UserInfoInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	return &gen.UserInfoResponse{
		Username: userInfo.Username,
		Email:    userInfo.Email,
	}, nil
}

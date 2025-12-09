package v1

import (
	"context"
	"errors"
	"fmt"

	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/vixarapi/delivery/http/security"
	"github.com/keenywheels/backend/internal/vixarapi/service"
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

	cookie := fmt.Sprintf("session_id=%s; Path=/; SameSite=Lax; HttpOnly", res.Session)

	return &gen.VkAuthCallbackResponseHeaders{
		Response: gen.VkAuthCallbackResponse{
			UserExists: res.UserExists,
			Username:   res.Username,
			Email:      res.Email,
			Vkid:       res.VKID,
		},
		SetCookie: gen.OptString{
			Value: cookie,
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

	if err := c.svc.RegisterVkUser(ctx, &service.RegisterVkUserParams{
		Email:    req.Email,
		Username: req.Username,
		VKID:     req.Vkid,
	}); err != nil {
		log.Errorf("[%s] failed to register VK user: %v", op, err)

		switch {
		case errors.Is(err, service.ErrAlreadyExists):
			return &gen.VkAuthRegisterConflict{
				Error: httputils.ErrorConflict,
			}, nil
		}

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

	cookie := "session_id=; Path=/; SameSite=Lax; HttpOnly"

	return &gen.LogoutUserOK{
		SetCookie: gen.OptString{
			Value: cookie,
			Set:   true,
		},
	}, nil
}

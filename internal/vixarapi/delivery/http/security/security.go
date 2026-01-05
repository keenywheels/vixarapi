package security

import (
	"context"
	"errors"
	"fmt"
	"strings"

	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/vixarapi/service/user"
	"github.com/keenywheels/backend/pkg/ctxutils"
	"github.com/keenywheels/backend/pkg/logger"
)

// HandleCookieAuth handles cookie-based authentication
func (c *Controller) HandleCookieAuth(
	ctx context.Context,
	operationName gen.OperationName,
	t gen.CookieAuth,
) (context.Context, error) {
	session := t.GetAPIKey()
	if strings.TrimSpace(session) == "" {
		return ctx, fmt.Errorf("failed to validate session: %w", ErrEmptyToken)
	}

	// check if the session is valid
	valid, userInfo, err := c.srvc.ValidateSession(ctx, session)
	if err != nil {
		return ctx, fmt.Errorf("failed to get session: %w", err)
	}

	if !valid {
		return ctx, fmt.Errorf("failed to validate session: %w", ErrInvalidToken)
	}

	// set sessionID and user info in context
	if userInfo == nil {
		return ctx, errors.New("something wrong: got nil user for valid session")
	}

	// update logger with user info fields
	ctxutils.GetLogger(ctx).Add(getFields(userInfo)...)

	// update context with logger, user info and session ID
	ctx = SetSessionID(SetUserInfo(ctx, *userInfo), session)

	return ctx, nil
}

func getFields(userInfo *user.UserSessionInfo) []logger.Field {
	fields := []logger.Field{
		{Key: "user_id", Value: userInfo.ID},
		{Key: "email", Value: userInfo.Email},
	}

	if userInfo.VKID != 0 {
		fields = append(fields, logger.Field{Key: "vkid", Value: userInfo.VKID})
	}

	return fields
}

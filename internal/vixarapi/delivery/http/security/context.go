package security

import (
	"context"

	userSvc "github.com/keenywheels/backend/internal/vixarapi/service/user"
)

// ctxKeySessionID is a type for the session ID context key
type ctxKeySessionID int

// key to get/set session ID in context
const sessionIDKey ctxKeySessionID = 0

// SetSessionID sets the session ID in the context
func SetSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionIDKey, sessionID)
}

// GetSessionID retrieves the session ID from the context
func GetSessionID(ctx context.Context) (string, bool) {
	if sessionID, ok := ctx.Value(sessionIDKey).(string); ok {
		return sessionID, true
	}

	return "", false
}

// ctxKeyUserInfo is a type for the user info context key
type ctxKeyUserInfo int

// key to get/set user info in context
const userInfoKey ctxKeyUserInfo = 0

// SetUserInfo sets the user info in the context
func SetUserInfo(ctx context.Context, userInfo userSvc.UserSessionInfo) context.Context {
	return context.WithValue(ctx, userInfoKey, userInfo)
}

// GetUserInfo retrieves the user info from the context
func GetUserInfo(ctx context.Context) (*userSvc.UserSessionInfo, bool) {
	if userInfo, ok := ctx.Value(userInfoKey).(userSvc.UserSessionInfo); ok {
		return &userInfo, true
	}

	return nil, false
}

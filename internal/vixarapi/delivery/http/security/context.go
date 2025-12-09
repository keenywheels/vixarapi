package security

import "context"

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

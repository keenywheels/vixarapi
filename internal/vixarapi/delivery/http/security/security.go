package security

import (
	"context"
	"errors"
	"fmt"
	"strings"

	gen "github.com/keenywheels/backend/internal/api/v1"
)

var (
	ErrEmptyToken   = errors.New("empty token")
	ErrInvalidToken = errors.New("invalid token")
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
	valid, err := c.srvc.ValidateSession(ctx, session)
	if err != nil {
		return ctx, fmt.Errorf("failed to validate session: %w", err)
	}

	if !valid {
		return ctx, fmt.Errorf("failed to validate session: %w", ErrInvalidToken)
	}

	return ctx, nil
}

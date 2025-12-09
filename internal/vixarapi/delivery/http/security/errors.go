package security

import (
	"errors"

	"github.com/ogen-go/ogen/ogenerrors"
)

var (
	ErrEmptyToken   = errors.New("empty token")
	ErrInvalidToken = errors.New("invalid token")
)

func IsSecurityError(err error) bool {
	return errors.Is(err, ErrEmptyToken) || errors.Is(err, ErrInvalidToken) ||
		errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied)
}

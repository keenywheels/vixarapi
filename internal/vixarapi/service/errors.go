package service

import (
	"errors"
	"fmt"

	"github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
)

// common service layer errors
var (
	ErrNotFound      = errors.New("did not get any data")
	ErrAlreadyExists = errors.New("data already exists")
)

// ParseRepositoryError parses a repository error and returns a corresponding common service error
func ParseRepositoryError(op string, err error) error {
	switch {
	case errors.Is(err, postgres.ErrNotFound):
		return fmt.Errorf("[%s] service error: %w", op, ErrNotFound)
	case errors.Is(err, postgres.ErrAlreadyExists):
		return fmt.Errorf("[%s] service error: %w", op, ErrAlreadyExists)
	}

	return fmt.Errorf("[%s] UNEXPECTED service error: %w", op, err)
}

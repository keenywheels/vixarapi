package service

import (
	"errors"
	"fmt"

	"github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
)

var (
	ErrNotFound      = errors.New("did not get any data")
	ErrAlreadyExists = errors.New("data already exists")
)

// parseRepositoryError parses a repository error and returns a service layer error
func parseRepositoryError(op string, err error) error {
	switch {
	case errors.Is(err, postgres.ErrNotFound):
		return fmt.Errorf("[%s] service error: %w", op, ErrNotFound)
	case errors.Is(err, postgres.ErrAlreadyExists):
		return fmt.Errorf("[%s] service error: %w", op, ErrAlreadyExists)
	}

	return fmt.Errorf("[%s] UNEXPECTED service error: %w", op, err)
}

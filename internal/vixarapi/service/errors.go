package service

import (
	"errors"
	"fmt"

	"github.com/keenywheels/backend/internal/vixarapi/repository"
)

var (
	ErrNotFound = errors.New("did not get any data")
)

// parseServiceError parses a repository error and returns a service layer error
func parseServiceError(op string, err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fmt.Errorf("[%s] service error: %w", op, ErrNotFound)
	}

	return fmt.Errorf("[%s] UNEXPECTED service error: %w", op, err)
}

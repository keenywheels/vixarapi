package repository

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("data not found")
)

// parsePostgresError parses a Postgres error and returns a repo layer error
func parsePostgresError(op string, err error) error {
	return fmt.Errorf("[%s] UNEXPECTED repository error: %w", op, err)
}

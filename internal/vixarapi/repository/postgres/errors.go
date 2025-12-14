package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// common repository layer postgres errors
var (
	ErrNotFound      = errors.New("data not found")
	ErrAlreadyExists = errors.New("data already exists")
)

// ParsePostgresError parses a Postgres error and returns a common repository layer error
func ParsePostgresError(op string, err error) error {
	var pgErr *pgconn.PgError

	switch {
	case errors.As(err, &pgErr):
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("[%s] failed with: %w", op, ErrAlreadyExists)
		}
	case errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("[%s] failed with: %w", op, ErrNotFound)
	}

	return fmt.Errorf("[%s] UNEXPECTED repository error: %w", op, err)
}

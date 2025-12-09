package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound      = errors.New("data not found")
	ErrAlreadyExists = errors.New("data already exists")
)

// parsePostgresError parses a Postgres error and returns a repo layer error
func parsePostgresError(op string, err error) error {
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

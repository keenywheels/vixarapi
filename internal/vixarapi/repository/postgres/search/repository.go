package search

import (
	commonRepo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
	"github.com/keenywheels/backend/pkg/postgres"
)

const (
	searchLimit = 5 * 365 * 10
)

// Tables holds the table definitions
type Tables struct {
	search commonRepo.SearchTokenTable
}

// Repository provides interest-related data access logic
type Repository struct {
	tbls Tables
	db   *postgres.Postgres
}

// New creates new Repository instance
func New(db *postgres.Postgres) *Repository {
	return &Repository{
		tbls: Tables{
			search: commonRepo.NewSearchTokenTable(),
		},
		db: db,
	}
}

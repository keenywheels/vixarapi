package user

import (
	commonRepo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
	"github.com/keenywheels/backend/pkg/postgres"
)

// Tables holds the table definitions
type Tables struct {
	user      commonRepo.UserTable
	userQuery commonRepo.UserQueryTable
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
			user:      commonRepo.NewUserTable(),
			userQuery: commonRepo.NewUserQueryTable(),
		},
		db: db,
	}
}

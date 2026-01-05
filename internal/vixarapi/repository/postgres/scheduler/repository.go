package scheduler

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
	commonRepo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
	"github.com/keenywheels/backend/pkg/postgres"
)

type Tables struct {
	search commonRepo.SearchTokenTable
}

// Repository provides interest-related data access logic
type Repository struct {
	tbls Tables
	db   *postgres.Postgres

	scheduler gocron.Scheduler
}

// New creates new Repository instance
func New(db *postgres.Postgres) (*Repository, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	repo := Repository{
		tbls: Tables{
			search: commonRepo.NewSearchTokenTable(),
		},
		db:        db,
		scheduler: scheduler,
	}

	return &repo, nil
}

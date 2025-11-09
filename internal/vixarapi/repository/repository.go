package repository

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/keenywheels/backend/pkg/logger"
	"github.com/keenywheels/backend/pkg/postgres"
)

const (
	searchLimit = 5 * 365 * 10
)

// SearchTokenFields represents the fields of the search token table
type SearchTokenFields struct {
	TokenName   string
	ScrapeDate  string
	Interest    string
	Sentiment   string
	MaxInterest string
}

// SearchTokenTable represents the structure of the search token table
type SearchTokenTable struct {
	Name   string
	Fields SearchTokenFields
}

// Repository provides interest-related data access logic
type Repository struct {
	tbl SearchTokenTable
	db  *postgres.Postgres

	scheduler gocron.Scheduler
}

// New creates new Repository instance
func New(db *postgres.Postgres) (*Repository, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	repo := Repository{
		tbl: SearchTokenTable{
			Name: "mv_token_search",
			Fields: SearchTokenFields{
				TokenName:   "token_name",
				ScrapeDate:  "scrape_date",
				Interest:    "interest",
				Sentiment:   "sentiment",
				MaxInterest: "max_interest",
			},
		},
		db:        db,
		scheduler: scheduler,
	}

	return &repo, nil
}

// StartScheduler starts the periodic update of the search table
func (r *Repository) StartScheduler(ctx context.Context, logger logger.Logger, cfg *SchedulerConfig) error {
	cfg.fix()

	if err := r.initScheduler(ctx, logger, cfg); err != nil {
		return fmt.Errorf("failed to init scheduler: %w", err)
	}

	r.scheduler.Start()

	return nil
}

// CloseScheduler stops the scheduler
func (r *Repository) CloseScheduler() error {
	return r.scheduler.Shutdown()
}

package scheduler

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/keenywheels/backend/pkg/logger"
)

// StartScheduler starts the periodic update of the search table
func (r *Repository) StartScheduler(ctx context.Context, logger logger.Logger, cfg *Config) error {
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

// initScheduler initializes the scheduler for periodic tasks
func (r *Repository) initScheduler(ctx context.Context, log logger.Logger, cfg *Config) error {
	_, err := r.scheduler.NewJob(
		gocron.CronJob(cfg.RefreshSearchTablePattern, false),
		gocron.NewTask(r.updateSearchTable),
		gocron.WithContext(ctx),
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				log.Errorf("job %s failed: %v", jobName, err)
			}),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to init job: %w", err)
	}

	return nil
}

package search

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/keenywheels/backend/pkg/ctxutils"
)

// StartScheduler starts the periodic update of the search table
func (s *Service) StartScheduler(ctx context.Context, cfg *SchedulerConfig) error {
	cfg.fix()

	if err := s.initScheduler(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init scheduler: %w", err)
	}

	s.scheduler.Start()

	return nil
}

// CloseScheduler stops the scheduler
func (s *Service) CloseScheduler() error {
	return s.scheduler.Shutdown()
}

// initScheduler initializes the scheduler for periodic tasks
func (s *Service) initScheduler(ctx context.Context, cfg *SchedulerConfig) error {
	var (
		log = ctxutils.GetLogger(ctx)
	)

	_, err := s.scheduler.NewJob(
		gocron.CronJob(cfg.RefreshSearchTablePattern, false),
		gocron.NewTask(s.updateSearchTask),
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

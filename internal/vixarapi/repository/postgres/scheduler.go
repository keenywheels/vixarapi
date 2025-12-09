package postgres

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/keenywheels/backend/pkg/logger"
)

const (
	defaultRefreshSearchTablePattern = "0 0 * * *"
)

// SchedulerConfig holds the configuration for the scheduler
type SchedulerConfig struct {
	RefreshSearchTablePattern string `mapstructure:"refresh_search_table_pattern"`
}

// fix validates and sets defaults for SchedulerConfig
func (sc *SchedulerConfig) fix() {
	if sc.RefreshSearchTablePattern == "" {
		sc.RefreshSearchTablePattern = defaultRefreshSearchTablePattern
	}
}

// initScheduler initializes the scheduler for periodic tasks
func (r *Repository) initScheduler(ctx context.Context, log logger.Logger, cfg *SchedulerConfig) error {
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

// updateSearchTable performs the update of the search table
func (r *Repository) updateSearchTable(ctx context.Context) error {
	var (
		op    = "Repository.updateSearchTable"
		query = fmt.Sprintf("REFRESH MATERIALIZED VIEW CONCURRENTLY %s;", r.tbls.search.Name)
	)

	if _, err := r.db.Pool.Exec(ctx, query); err != nil {
		return fmt.Errorf("[%s] failed to refresh materialized view: %w", op, err)
	}

	return nil
}

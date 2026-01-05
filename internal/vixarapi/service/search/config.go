package search

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

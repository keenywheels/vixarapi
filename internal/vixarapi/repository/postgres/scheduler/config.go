package scheduler

const (
	defaultRefreshSearchTablePattern = "0 0 * * *"
)

// Config holds the configuration for the scheduler
type Config struct {
	RefreshSearchTablePattern string `mapstructure:"refresh_search_table_pattern"`
}

// fix validates and sets defaults for SchedulerConfig
func (sc *Config) fix() {
	if sc.RefreshSearchTablePattern == "" {
		sc.RefreshSearchTablePattern = defaultRefreshSearchTablePattern
	}
}

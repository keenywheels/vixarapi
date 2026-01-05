package redis

import "time"

const (
	defaultDialTimeout  = 5 * time.Second
	defaultReadTimeout  = 3 * time.Second
	defaultWriteTimeout = 3 * time.Second
	defaultPoolSize     = 10
	defaultMaxRetries   = 3
	defaultPingTimeout  = 5 * time.Second
)

// Config represents the configuration for Redis
type Config struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	PoolSize     int           `mapstructure:"pool_size"`
	MaxRetries   int           `mapstructure:"max_retries"`
	PingTimeout  time.Duration `mapstructure:"ping_timeout"`
}

// fixConfig sets default values for the Config if they are not provided
func fixConfig(cfg *Config) {
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = defaultDialTimeout
	}

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = defaultReadTimeout
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = defaultWriteTimeout
	}

	if cfg.PoolSize == 0 {
		cfg.PoolSize = defaultPoolSize
	}

	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = defaultMaxRetries
	}

	if cfg.PingTimeout == 0 {
		cfg.PingTimeout = defaultPingTimeout
	}
}

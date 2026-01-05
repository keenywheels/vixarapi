package cookie

import "time"

const (
	defaultSessionExpiration = 30 * 24 * time.Hour // 30 days
	defaultSessionCookieName = "session_id"
)

// BaseInfo holds basic info for cookies
type BaseInfo struct {
	Name       string        `mapstructure:"name"`
	Expiration time.Duration `mapstructure:"expiration"`
}

// Session config for session cookies
type Session struct {
	BaseInfo `mapstructure:",squash"`
}

// Config holds configuration for cookie management
type Config struct {
	Session Session `mapstructure:"session"`
}

// fixConfig sets default values for missing configuration fields
func fixConfig(cfg *Config) {
	if cfg.Session.Name == "" {
		cfg.Session.Name = defaultSessionCookieName
	}

	if cfg.Session.Expiration <= 0 {
		cfg.Session.Expiration = defaultSessionExpiration
	}
}

package redis

// Config represents the configuration for Redis
type Config struct {
	Addr     string
	Password string
	DB       int
}

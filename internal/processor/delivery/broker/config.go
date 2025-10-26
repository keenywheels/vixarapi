package broker

import "time"

// Config holds the configuration for the Broker
type Config struct {
	WorkerCount   int
	MaxRetryCount int
	RetryDelay    time.Duration
}

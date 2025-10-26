package kafka

import "github.com/keenywheels/backend/pkg/logger"

type Option func(*Kafka)

// WithLogger sets a custom logger for the Kafka consumer
func WithLogger(l logger.Logger) Option {
	return func(k *Kafka) {
		k.l = l
	}
}

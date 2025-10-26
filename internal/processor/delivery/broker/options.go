package broker

import "github.com/keenywheels/backend/pkg/logger"

type Option func(*Broker)

// WithLogger sets the logger for the Broker
func WithLogger(l logger.Logger) Option {
	return func(b *Broker) {
		b.l = l
	}
}

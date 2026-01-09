package broker

import (
	"github.com/keenywheels/backend/internal/pkg/producer/kafka"
)

// Topics represents available topics
type Topics struct {
	Notifications string
}

// Broker represents broker instance
type Broker struct {
	topics Topics
	kafka  *kafka.Kafka
}

// New creates new broker instance
func New(kafka *kafka.Kafka, topics Topics) *Broker {
	return &Broker{
		topics: topics,
		kafka:  kafka,
	}
}

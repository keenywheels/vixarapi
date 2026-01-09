package broker

import (
	"github.com/keenywheels/backend/internal/pkg/producer/kafka"
	"github.com/keenywheels/backend/internal/vixarapi/models"
)

// SendNotification put email notification task to kafka
func (b *Broker) SendNotification(event models.Notification) error {
	kafkaMsg := kafka.Message{
		Topic: b.topics.Notifications,
		Value: event,
	}

	return b.kafka.ProduceJSON(kafkaMsg)
}

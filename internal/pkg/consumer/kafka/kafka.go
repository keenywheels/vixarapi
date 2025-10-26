package kafka

import (
	"github.com/IBM/sarama"
	"github.com/keenywheels/backend/pkg/logger"
	"github.com/keenywheels/backend/pkg/logger/zap"
)

// Kafka represents a Kafka consumer
type Kafka struct {
	c sarama.ConsumerGroup
	l logger.Logger
}

// New creates a new Kafka consumer
func New(brokers []string, groupID string, kafkaConfig Config, opts ...Option) (*Kafka, error) {
	k := &Kafka{
		l: zap.New(),
	}

	for _, opt := range opts {
		opt(k)
	}

	// create consumer group
	cfg := sarama.NewConfig()

	cfg.Consumer.Return.Errors = true // custom error handling
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, err
	}

	k.c = consumerGroup

	return k, nil
}

// Close closes the Kafka consumer
func (k *Kafka) Close() error {
	return k.c.Close()
}

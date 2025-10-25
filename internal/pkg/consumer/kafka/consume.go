package kafka

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
)

// StartConsuming starts consuming messages from the specified topics using the provided handler
func (k *Kafka) StartConsuming(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	prefix := "KAFKA CONSUMER"

	for _, topic := range topics {
		if strings.Contains(topic, "__consumer_offsets") {
			return fmt.Errorf("consuming from internal topic %s is not allowed", topic)
		}
	}

	// handle errors
	errCh := k.c.Errors()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err, ok := <-errCh:
				if !ok {
					return
				}

				k.l.Errorf("[%s] got error: %v", prefix, err)
			}
		}
	}()

	// consume loop
	for {
		if ctx.Err() != nil {
			k.l.Infof("[%s] stopping consume because context done: %v", prefix, ctx.Err())
			return ctx.Err()
		}

		// start consuming
		if err := k.c.Consume(ctx, topics, handler); err != nil {
			if ctx.Err() != nil {
				k.l.Infof("[%s] stopping consume because context done: %v", prefix, ctx.Err())
				return ctx.Err()
			}

			k.l.Errorf("[%s] error during consuming: %v", prefix, err)

			return fmt.Errorf("error during consuming: %w", err)
		}
	}
}

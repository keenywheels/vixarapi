package broker

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/keenywheels/backend/pkg/ctxutils"
	"golang.org/x/sync/errgroup"
)

// Setup prepares the broker for message consumption
func (b *Broker) Setup(session sarama.ConsumerGroupSession) error {
	for i := 0; i < b.workerCount; i++ {
		b.wg.Add(1)
		go b.worker(ctxutils.SetLogger(session.Context(), b.l), i)
	}

	return nil
}

// Cleanup cleans up resources after message consumption
func (b *Broker) Cleanup(session sarama.ConsumerGroupSession) error {
	close(b.messageQueue)
	close(b.ackQueue)
	b.wg.Wait()

	return nil
}

// ConsumeClaim get messages from kafka
func (b *Broker) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	gr, ctx := errgroup.WithContext(session.Context())

	gr.Go(func() error {
		for {
			select {
			case msg, ok := <-claim.Messages():
				if !ok {
					return nil
				}

				b.messageQueue <- message{msg: msg, retry: 0} // status doesn't matter in messageQueue

			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	gr.Go(func() error {
		for {
			select {
			case msg, ok := <-b.ackQueue:
				if !ok {
					return nil
				}

				// mark as done if acknowledged or retries exhausted
				if msg.status == statusSuccess || msg.retry == b.maxRetry {
					session.MarkMessage(msg.msg, "")
				} else {
					b.messageQueue <- message{msg: msg.msg, retry: msg.retry + 1}
				}

			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	if err := gr.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

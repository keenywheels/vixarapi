package broker

import (
	"context"
	"fmt"
)

const maxLogMsgLength = 25

// worker processes messages from the message queue
func (b *Broker) worker(ctx context.Context, i int) {
	defer b.wg.Done()

	var (
		prefix = fmt.Sprintf("WORKER %d", i)
	)

	for msg := range b.messageQueue {
		if ctx.Err() != nil {
			return
		}

		b.l.Infof("[%s] got new message from topic %s with id=%s", prefix, msg.msg.Topic, msg.id)
		b.processTasks(ctx, msg)
	}
}

// processTasks processes messages from the message queue
func (b *Broker) processTasks(ctx context.Context, msg message) {
	var (
		op  = "Broker.processTasks"
		err error
	)

	// check topic to choose the right handler
	switch msg.msg.Topic {
	case b.topics.ScraperData:
		err = b.service.TokenizeMessage(ctx, string(msg.msg.Value))
	case b.topics.Notifications:
		err = b.service.NotifyUser(ctx, string(msg.msg.Value))
	default:
		b.l.Warnf("[%s] unknown topic %s for %s", op, msg.msg.Topic, msg.id)
	}

	// parse error
	if err != nil {
		b.l.Errorf("[%s] failed to process message %s from topic %s: %v", op, msg.id, msg.msg.Topic, err)

		// try to retry task
		if msg.retry < b.maxRetry {
			go func() {
				defer func() {
					if rec := recover(); rec != nil {
						b.l.Errorf("[%s] panic while processing message %s: %v", op, msg.id, rec)
					}
				}()

				b.putRetryTask(ctx, msg)
			}()
		} else {
			// all retries exhausted -> mark message as failed and ack it
			b.l.Errorf("[%s] failed to process message %s after %d retries -> ack it", op, msg.id, b.maxRetry)
			msg.status = statusFailed
		}
	} else {
		// no error -> mark message as success
		msg.status = statusSuccess
	}

	b.ackQueue <- msg
}

// putRetryTask puts a message back to the queue for retrying
func (b *Broker) putRetryTask(ctx context.Context, msg message) {
	timeoutCtx, cancel := context.WithTimeout(ctx, b.retryDelay)
	defer cancel()

	select {
	case <-timeoutCtx.Done():
		msg.retry++
		b.messageQueue <- msg
	case <-ctx.Done():
		return
	}
}

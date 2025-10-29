package broker

import (
	"context"
	"fmt"

	"github.com/keenywheels/backend/pkg/ctxutils"
)

const maxLogMsgLength = 25

// worker processes messages from the message queue
func (b *Broker) worker(ctx context.Context, i int) {
	defer b.wg.Done()

	var (
		prefix = fmt.Sprintf("WORKER %d", i)
		log    = ctxutils.GetLogger(ctx)
	)

	for msg := range b.messageQueue {
		if ctx.Err() != nil {
			return
		}

		// TODO: сделать нормальное логгирование
		// TODO: добавить логику с выбором топика
		log.Debugf("[%s] start processing message: partial_msg=%s", prefix, string(msg.msg.Value)[:maxLogMsgLength])

		err := b.service.TokenizeMessage(ctx, string(msg.msg.Value))
		if err != nil {
			b.l.Errorf("[%s] failed to process message: %v", prefix, err)

			// put retry task
			msg.status = statusFailed
			go b.putRetryTask(ctx, msg)

			continue
		}

		// set task as success
		msg.status = statusSuccess
		b.ackQueue <- msg
	}
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

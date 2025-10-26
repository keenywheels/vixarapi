package broker

import (
	"context"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/keenywheels/backend/pkg/logger"
	"github.com/keenywheels/backend/pkg/logger/zap"
)

// default values
const (
	defaultWorkerCount = 5
	defaultMaxRetry    = 3
	defaultRetryDelay  = 5 * time.Second
)

// tasks statuses
const (
	statusSuccess = "success"
	statusFailed  = "failed"
)

// message represents a message in queue
type message struct {
	msg    *sarama.ConsumerMessage
	retry  int
	status string
}

// IService defines the interface for the service layer of processor
type IService interface {
	TokenizeMessage(ctx context.Context, message string) error
}

// Topics holds the topic names
type Topics struct {
	ScraperData string
}

// Broker struct for message broker
type Broker struct {
	l       logger.Logger
	service IService
	topics  Topics

	// settings for message handling
	workerCount  int
	maxRetry     int
	retryDelay   time.Duration
	messageQueue chan message
	ackQueue     chan message
	wg           sync.WaitGroup
}

// New creates a new Broker instance
func New(cfg Config, service IService, topics Topics, opts ...Option) *Broker {
	b := &Broker{
		l:           zap.New(),
		service:     service,
		topics:      topics,
		workerCount: defaultWorkerCount,
		maxRetry:    defaultMaxRetry,
		retryDelay:  defaultRetryDelay,
	}

	// apply options
	for _, opt := range opts {
		opt(b)
	}

	// override with config values if provided
	if cfg.WorkerCount > 0 {
		b.workerCount = cfg.WorkerCount
	}

	if cfg.MaxRetryCount > 0 {
		b.maxRetry = cfg.MaxRetryCount
	}

	if cfg.RetryDelay > 0 {
		b.retryDelay = cfg.RetryDelay
	}

	// initialize queues
	b.messageQueue = make(chan message, b.workerCount*2) // TODO: проверить такой размер буфера по бенчмаркам
	b.ackQueue = make(chan message, b.workerCount*2)

	return b
}

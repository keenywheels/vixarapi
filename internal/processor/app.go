package processor

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/keenywheels/backend/internal/pkg/client/llm"
	"github.com/keenywheels/backend/internal/pkg/consumer/kafka"
	"github.com/keenywheels/backend/internal/processor/delivery/broker"
	"github.com/keenywheels/backend/internal/processor/repository"
	"github.com/keenywheels/backend/internal/processor/service"
	"github.com/keenywheels/backend/pkg/logger"
	"github.com/keenywheels/backend/pkg/logger/zap"
	"github.com/keenywheels/backend/pkg/mailer/smtp"
	"github.com/keenywheels/backend/pkg/postgres"
	"golang.org/x/sync/errgroup"
)

// App represent app environment
type App struct {
	opts *Options

	db *postgres.Postgres

	cfg    *Config
	logger logger.Logger
}

// New creates new application instance with options
func New() *App {
	opts := NewDefaultOpts()
	opts.LoadEnv()
	opts.LoadFlags()

	return &App{
		opts: opts,
	}
}

// Run starts the application
func (app *App) Run() error {
	// read config
	cfg, err := LoadConfig(app.opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load app config: %w", err)
	}

	app.cfg = cfg
	app.initLogger()
	defer func() {
		if err := app.logger.Close(); err != nil {
			log.Printf("failed to close logger: %v", err)
		}
	}()

	// create postgres connection
	db, err := app.getPostgresConn()
	if err != nil {
		return fmt.Errorf("failed to create postgres connection: %w", err)
	}
	defer db.Close()

	// create llm client
	llm := llm.NewClient(&cfg.App.Clients.LLM)

	// create service layer
	mailer := smtp.New(&cfg.App.SMTPCfg)

	repo := repository.New(db)
	service := service.New(repo, llm, mailer)

	// create broker
	brokerOpts := []broker.Option{
		broker.WithLogger(app.logger),
	}

	brokerCfg := broker.Config{
		WorkerCount:   cfg.App.Processor.WorkersCount,
		MaxRetryCount: cfg.App.Processor.MaxRetries,
		RetryDelay:    cfg.App.Processor.RetryDelay,
	}

	b := broker.New(brokerCfg, service, broker.Topics{
		ScraperData:   cfg.KafkaCfg.Topics.ScraperData,
		Notifications: cfg.KafkaCfg.Topics.Notifications,
	}, brokerOpts...)

	// create kafka consumer
	kafkaConsumer, err := kafka.New(
		cfg.KafkaCfg.Brokers,
		cfg.KafkaCfg.GroupID,
		kafka.Config{},
		kafka.WithLogger(app.logger),
	)
	if err != nil {
		return fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	// create errgroup with signal context
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// start consuming
	topics := []string{
		cfg.KafkaCfg.Topics.ScraperData,   // topic with data from scraped sites
		cfg.KafkaCfg.Topics.Notifications, // topic with notifications
	}

	g.Go(func() error {
		app.logger.Infof("starting kafka consumer for topics: %v", topics)

		if err := kafkaConsumer.StartConsuming(ctx, topics, b); err != nil {
			return fmt.Errorf("kafka consumer error: %w", err)
		}

		return nil
	})

	// wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		app.logger.Error("app error: %v", err)

		return err
	}

	return nil
}

// initLogger create new Logger based on config
func (app *App) initLogger() {
	logCfg := app.cfg.App.LoggerCfg
	opts := []zap.Option{}

	// setup opts
	if len(logCfg.LogLevel) != 0 {
		opts = append(opts, zap.LogLvl(logCfg.LogLevel))
	}

	if len(logCfg.Mode) != 0 {
		opts = append(opts, zap.Mode(logCfg.Mode))
	}

	if len(logCfg.Encoding) != 0 {
		opts = append(opts, zap.Encoding(logCfg.Encoding))
	}

	if len(logCfg.LogPath) != 0 {
		opts = append(opts, zap.LogPath(logCfg.LogPath))
	}

	if logCfg.MaxLogSize != 0 {
		opts = append(opts, zap.MaxLogSize(logCfg.MaxLogSize))
	}

	if logCfg.MaxLogBackups != 0 {
		opts = append(opts, zap.MaxLogBackups(logCfg.MaxLogBackups))
	}

	if logCfg.MaxLogAge != 0 {
		opts = append(opts, zap.MaxLogAge(logCfg.MaxLogAge))
	}

	// set logger
	app.logger = zap.New(opts...)
}

// getPostgresConn creates and returns a new Postgres connection
func (app *App) getPostgresConn() (*postgres.Postgres, error) {
	var opts []postgres.Option

	if app.cfg.App.Postgres.MaxPoolSize != 0 {
		opts = append(opts, postgres.MaxPoolSize(app.cfg.App.Postgres.MaxPoolSize))
	}

	if app.cfg.App.Postgres.ConnAttempts != 0 {
		opts = append(opts, postgres.ConnAttempts(app.cfg.App.Postgres.ConnAttempts))
	}

	if app.cfg.App.Postgres.ConnTimeout != 0 {
		opts = append(opts, postgres.ConnTimeout(app.cfg.App.Postgres.ConnTimeout))
	}

	return postgres.New(app.cfg.App.Postgres.DSN(), opts...)
}

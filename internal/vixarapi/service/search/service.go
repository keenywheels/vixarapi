package search

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/keenywheels/backend/internal/vixarapi/models"
	repo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres/search"
)

// IRepository provides interface to communicate with search repository layer
type IRepository interface {
	SearchTokenInfo(context.Context, *repo.SearchTokenParams) ([]models.TokenInfo, error)
	UpdateSearchTable(context.Context) error
	UpdateUserTokenSubs(ctx context.Context, intervalType string, amount int) error
}

// IBroker provides interface to communicate with message broker
type IBroker interface {
	SendNotification(event models.Notification) error
}

// Service provides interest-related business logic
type Service struct {
	r         IRepository
	scheduler gocron.Scheduler
	broker    IBroker
}

// New creates a new interest service
func New(repo IRepository, broker IBroker) (*Service, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &Service{
		r:         repo,
		scheduler: scheduler,
		broker:    broker,
	}, nil
}

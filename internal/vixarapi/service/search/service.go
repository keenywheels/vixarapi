package search

import (
	"context"

	"github.com/keenywheels/backend/internal/vixarapi/models"
	"github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
)

// IRepository provides interface to communicate with the repository layer
type IRepository interface {
	SearchTokenInfo(context.Context, *postgres.SearchTokenParams) ([]models.TokenInfo, error)
}

// Service provides interest-related business logic
type Service struct {
	repo IRepository
}

// New creates a new interest service
func New(
	repo IRepository,
) *Service {
	return &Service{
		repo: repo,
	}
}

package search

import (
	"context"

	"github.com/keenywheels/backend/internal/vixarapi/models"
	repo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres/search"
)

// IRepository provides interface to communicate with search repository layer
type IRepository interface {
	SearchTokenInfo(context.Context, *repo.SearchTokenParams) ([]models.TokenInfo, error)
}

// Service provides interest-related business logic
type Service struct {
	r IRepository
}

// New creates a new interest service
func New(repo IRepository) *Service {
	return &Service{
		r: repo,
	}
}

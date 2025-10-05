package service

import (
	"context"

	"github.com/keenywheels/backend/internal/interest/models"
)

// IRepository provides interface to communicate with the repository layer
type IRepository interface {
	GetAllInterest(context.Context, string) ([]models.Interest, error)
}

// Service provides interest-related business logic
type Service struct {
	repo IRepository
}

// New creates a new interest service
func New(repo IRepository) *Service {
	return &Service{
		repo: repo,
	}
}

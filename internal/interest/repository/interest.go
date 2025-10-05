package repository

import (
	"context"

	"github.com/keenywheels/backend/internal/interest/models"
)

// GetAllInterest retrieves all interest records for the specified token from database
func (r *Repository) GetAllInterest(ctx context.Context, token string) ([]models.Interest, error) {
	mockInterests := []models.Interest{
		{
			Timestamp: 1700000000,
			Features: models.Features{
				Interest: 125,
			},
		},
		{
			Timestamp: 1701209600,
			Features: models.Features{
				Interest: 150,
			},
		},
		{
			Timestamp: 1701814400,
			Features: models.Features{
				Interest: 200,
			},
		},
		{
			Timestamp: 1702419200,
			Features: models.Features{
				Interest: 220,
			},
		},
		{
			Timestamp: 1703024000,
			Features: models.Features{
				Interest: 225,
			},
		},
		{
			Timestamp: 1703628800,
			Features: models.Features{
				Interest: 222,
			},
		},
	}

	return mockInterests, nil
}

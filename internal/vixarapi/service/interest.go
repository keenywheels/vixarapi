package service

import (
	"context"
	"fmt"
)

// Features represents the features of an interest record
type Features struct {
	Interest int
}

// Interest represent an interest record at service layer
type Interest struct {
	Timestamp int64
	Features  Features
}

// GetAllInterest retrieves all interest records for the specified token
func (s *Service) GetAllInterest(ctx context.Context, token string) ([]Interest, error) {
	op := "Service.GetAllInterest"

	interests, err := s.repo.GetAllInterest(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to get interests: %w", op, err)
	}

	return convertToServiceInterest(interests), nil
}

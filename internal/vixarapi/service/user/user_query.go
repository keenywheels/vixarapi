package user

import (
	"context"
	"time"

	"github.com/keenywheels/backend/internal/vixarapi/models"
	"github.com/keenywheels/backend/internal/vixarapi/service"
)

// SaveQueryParams contains parameters for SaveQuery method
type SaveQueryParams struct {
	UserID string
	Query  string
}

// SaveSearchQuery saves user search query
func (s *Service) SaveSearchQuery(ctx context.Context, params *SaveQueryParams) (string, error) {
	op := "Service.SaveSearchQuery"

	res, err := s.repo.SaveSearchQuery(ctx, params.UserID, params.Query)
	if err != nil {
		return "", service.ParseRepositoryError(op, err)
	}

	return res.ID, nil
}

// DeleteSearchQuery deletes user search query by ID
func (s *Service) DeleteSearchQuery(ctx context.Context, id string) error {
	op := "Service.DeleteSearchQuery"

	if err := s.repo.DeleteSearchQuery(ctx, id); err != nil {
		return service.ParseRepositoryError(op, err)
	}

	return nil
}

// Query represents a user search query
type Query struct {
	ID         string
	Query      string
	SearchDate time.Time
}

// GetSearchQueriesParams contains parameters for GetSearchQueries method
type GetSearchQueriesParams struct {
	UserID string
	Limit  uint64
	Offset uint64
}

// GetSearchQueries retrieves user search queries
func (s *Service) GetSearchQueries(ctx context.Context, params *GetSearchQueriesParams) ([]Query, error) {
	op := "Service.GetSearchQueries"

	queries, err := s.repo.GetSearchQueries(ctx, params.UserID, params.Limit, params.Offset)
	if err != nil {
		return nil, service.ParseRepositoryError(op, err)
	}

	return convertQueries(queries), nil
}

// convertQueries converts models.UserQuery slice to Query slice
func convertQueries(queries []*models.UserQuery) []Query {
	res := make([]Query, 0, len(queries))
	for _, q := range queries {
		res = append(res, Query{
			ID:         q.ID,
			Query:      q.Query,
			SearchDate: q.CreatedAt,
		})
	}

	return res
}

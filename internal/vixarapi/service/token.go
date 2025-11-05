package service

import (
	"context"
)

const (
	timedateLayout = "2006-01-02"
)

// Record represent a single record of token data
type Record struct {
	ScrapeDate         string
	Interest           int64
	NormalizedInterest float64
}

// TokenInfo token info in service layer
type TokenInfo struct {
	TokenName string // TODO: подумать над этим полем, мб попробовать привести к нормальной форме?
	Records   []Record
}

// SearchTokenInfo retrieves all interest records for the specified token
func (s *Service) SearchTokenInfo(ctx context.Context, token string) ([]TokenInfo, error) {
	op := "Service.SearchTokenInfo"

	tokensInfo, err := s.repo.SearchTokenInfo(ctx, token)
	if err != nil {
		return nil, parseServiceError(op, err)
	}

	return convertToServiceTokenInfo(tokensInfo), nil
}

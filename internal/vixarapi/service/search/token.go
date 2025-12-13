package search

import (
	"context"
	"time"

	"github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
	"github.com/keenywheels/backend/internal/vixarapi/service"
)

const (
	timedateLayout = "2006-01-02"
)

// Record represent a single record of token data
type Record struct {
	ScrapeDate         string
	Interest           int64
	NormalizedInterest float64
	Sentiment          int16
}

// TokenInfo token info in service layer
type TokenInfo struct {
	TokenName string // TODO: подумать над этим полем, мб попробовать привести к нормальной форме?
	Records   []Record
}

// SearchTokenInfoParams parameters for searching token info
type SearchTokenInfoParams struct {
	Token string
	Start time.Time
	End   time.Time
}

// SearchTokenInfo retrieves all interest records for the specified token
func (s *Service) SearchTokenInfo(ctx context.Context, params *SearchTokenInfoParams) ([]TokenInfo, error) {
	op := "Service.SearchTokenInfo"

	repoParams := &postgres.SearchTokenParams{
		Token: params.Token,
		Start: params.Start,
		End:   params.End,
	}

	tokensInfo, err := s.repo.SearchTokenInfo(ctx, repoParams)
	if err != nil {
		return nil, service.ParseRepositoryError(op, err)
	}

	return convertToServiceTokenInfo(tokensInfo), nil
}

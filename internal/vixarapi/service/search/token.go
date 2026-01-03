package search

import (
	"context"
	"time"

	"github.com/keenywheels/backend/internal/vixarapi/models"
	repo "github.com/keenywheels/backend/internal/vixarapi/repository/postgres/search"
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
	CategoryInterest   float64
	Sentiment          int16
}

// TokenInfo token info in service layer
type TokenInfo struct {
	TokenName string // TODO: подумать над этим полем, мб попробовать привести к нормальной форме?
	Category  string
	Records   []Record
}

// SearchTokenInfoParams parameters for searching token info
type SearchTokenInfoParams struct {
	Token    string
	Category *string
	Start    time.Time
	End      time.Time
}

// SearchTokenInfo retrieves all interest records for the specified token
func (s *Service) SearchTokenInfo(ctx context.Context, params *SearchTokenInfoParams) ([]TokenInfo, error) {
	op := "Service.SearchTokenInfo"

	repoParams := &repo.SearchTokenParams{
		Token:    params.Token,
		Category: params.Category,
		Start:    params.Start,
		End:      params.End,
	}

	tokensInfo, err := s.r.SearchTokenInfo(ctx, repoParams)
	if err != nil {
		return nil, service.ParseRepositoryError(op, err)
	}

	return convertToServiceTokenInfo(tokensInfo), nil
}

// convertToServiceTokenInfo converts repository structs to service layer structs
func convertToServiceTokenInfo(tokens []models.TokenInfo) []TokenInfo {
	resp := make([]TokenInfo, 0, len(tokens))

	for _, t := range tokens {
		records := make([]Record, 0, len(t.Records))
		for _, r := range t.Records {
			records = append(records, Record{
				ScrapeDate:         r.ScrapeDate.Format(timedateLayout),
				Interest:           r.Interest,
				NormalizedInterest: r.GlobalInterest,
				CategoryInterest:   r.CategoryInterest,
				Sentiment:          r.Sentiment,
			})
		}

		resp = append(resp, TokenInfo{
			TokenName: t.TokenName,
			Category:  t.Category,
			Records:   records,
		})
	}

	return resp
}

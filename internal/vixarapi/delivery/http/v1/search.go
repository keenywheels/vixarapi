package v1

import (
	"context"
	"errors"
	"time"

	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/vixarapi/service"
	searchSvc "github.com/keenywheels/backend/internal/vixarapi/service/search"
	"github.com/keenywheels/backend/pkg/ctxutils"
	"github.com/keenywheels/backend/pkg/httputils"
)

// SearchTokenInfo searches for token information based on the provided parameters.
func (c *Controller) SearchTokenInfo(
	ctx context.Context,
	req *gen.SearchTokenInfoRequest,
) (gen.SearchTokenInfoRes, error) {
	var (
		op  = "Controller.SearchTokenInfo"
		log = ctxutils.GetLogger(ctx)
	)

	end := time.Now().UTC()
	if req.End.Set {
		end = req.End.Value.UTC()
	}

	tokensInfo, err := c.searchSvc.SearchTokenInfo(ctx, &searchSvc.SearchTokenInfoParams{
		Token: req.Token,
		Start: req.Start.UTC(),
		End:   end,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return &gen.SearchTokenInfoNotFound{
				Error: httputils.ErrorNotFound,
			}, nil
		}

		log.Errorf("[%s] failed to search token info: %v", op, err)

		return &gen.SearchTokenInfoInternalServerError{
			Error: httputils.ErrorInternalError,
		}, nil
	}

	resp := gen.SearchTokenInfoOKApplicationJSON(convertToSearchTokenInfoResp(tokensInfo))

	return &resp, nil
}

// convertToSearchTokenInfoResp converts service layer structs to api response structs
func convertToSearchTokenInfoResp(tokens []searchSvc.TokenInfo) []gen.TokenInfo {
	resp := make([]gen.TokenInfo, 0, len(tokens))

	for _, t := range tokens {
		records := make([]gen.TokenRecord, 0, len(t.Records))
		for _, r := range t.Records {
			records = append(records, gen.TokenRecord{
				Timestamp: r.ScrapeDate,
				Features: gen.TokenRecordFeatures{
					Interest:           r.Interest,
					InterestNormalized: r.NormalizedInterest,
					Sentiment:          r.Sentiment,
				},
			})
		}

		resp = append(resp, gen.TokenInfo{
			Token:   t.TokenName,
			Records: records,
		})
	}

	return resp
}

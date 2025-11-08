package v1

import (
	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/vixarapi/service"
)

// convertToSearchTokenInfoResp converts service layer structs to api response structs
func convertToSearchTokenInfoResp(tokens []service.TokenInfo) []gen.TokenInfo {
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

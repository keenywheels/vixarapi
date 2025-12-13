package search

import "github.com/keenywheels/backend/internal/vixarapi/models"

// convertToServiceTokenInfo converts repository structs to service layer structs
func convertToServiceTokenInfo(tokens []models.TokenInfo) []TokenInfo {
	resp := make([]TokenInfo, 0, len(tokens))

	for _, t := range tokens {
		records := make([]Record, 0, len(t.Records))
		for _, r := range t.Records {
			records = append(records, Record{
				ScrapeDate:         r.ScrapeDate.Format(timedateLayout),
				Interest:           r.Interest,
				NormalizedInterest: r.NormalizedInterest,
				Sentiment:          r.Sentiment,
			})
		}

		resp = append(resp, TokenInfo{
			TokenName: t.TokenName,
			Records:   records,
		})
	}

	return resp
}

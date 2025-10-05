package v1

import (
	gen "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/interest/service"
)

// convertToGetAllInterestResp converts service layer structs to api response structs
func convertToGetAllInterestResp(interests []service.Interest) []gen.Interest {
	resp := make([]gen.Interest, 0, len(interests))

	for _, interest := range interests {
		resp = append(resp, gen.Interest{
			Timestamp: interest.Timestamp,
			Features: gen.InterestFeatures{
				Interest: interest.Features.Interest,
			},
		})
	}

	return resp
}

package service

import "github.com/keenywheels/backend/internal/vixarapi/models"

// convertToServiceInterest converts repository structs to service layer structs
func convertToServiceInterest(interests []models.Interest) []Interest {
	resp := make([]Interest, 0, len(interests))

	for _, interest := range interests {
		resp = append(resp, Interest{
			Timestamp: interest.Timestamp,
			Features: Features{
				Interest: interest.Features.Interest,
			},
		})
	}

	return resp
}

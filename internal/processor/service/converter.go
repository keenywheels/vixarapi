package service

import (
	"context"
	"strings"
	"time"

	tokenizerbase "github.com/keenywheels/backend/internal/pkg/tokenizer"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/metrics"
	"github.com/keenywheels/backend/internal/processor/models"
	"github.com/keenywheels/backend/pkg/ctxutils"
)

// convertToRepositoryTokens converts a slice of tokens to the format required by the repository layer
func convertToRepositoryTokens(
	ctx context.Context,
	msg *models.ScraperEvent,
	tokens []tokenizerbase.Token,
) []models.TokenData {
	var (
		log    = ctxutils.GetLogger(ctx)
		site   = msg.SiteName
		date   = msg.Date
		result = make([]models.TokenData, 0, len(tokens))
	)

	dateParsed, err := time.Parse(models.ScrapeDataFormat, date)
	if err != nil {
		log.Errorf("failed to parse date %s: %v", date, err)
		dateParsed = time.Now()
	}

	for _, t := range tokens {
		// skip filtered tokens
		if t.IsFiltered() {
			continue
		}

		// get token context
		tokenContext := make([]string, len(t.Context))
		for i, ctxToken := range t.Context {
			tokenContext[i] = ctxToken.Target
		}

		// get metrics
		interest, ok := metrics.Registry[metrics.InterestMetricIndex].Get(t.Target)
		if !ok {
			log.Errorf("failed to get interest metric for token %s", t.Target)
			continue
		}

		interestInt, ok := interest.(int)
		if !ok {
			log.Errorf("interest metric has invalid type for token %s: %T", t.Target, interest)
			continue
		}

		// append to result
		result = append(result, models.TokenData{
			TokenName: t.Target,
			Interest:  interestInt,
			Context:   strings.Join(tokenContext, " "),
			SiteName:  site,
			Date:      dateParsed,
		})
	}

	return result
}

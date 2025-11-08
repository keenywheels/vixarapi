package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/keenywheels/backend/internal/pkg/client/llm"
	tokenizerbase "github.com/keenywheels/backend/internal/pkg/tokenizer"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/metrics"
	"github.com/keenywheels/backend/internal/processor/models"
	"github.com/keenywheels/backend/pkg/ctxutils"
)

// TokenizeMessage processes and tokenizes the given message
func (s *Service) TokenizeMessage(ctx context.Context, message string) error {
	op := "Service.TokenizeMessage"

	var scraperEvent models.ScraperEvent
	if err := json.Unmarshal([]byte(message), &scraperEvent); err != nil {
		return fmt.Errorf("[%s] failed to unmarshal message: %w", op, err)
	}

	// tokenize msg
	tokens := s.tokenizer.Run(tokenizerbase.GetTokens(
		scraperEvent.Msg,
		tokenizerbase.NewTokenConfig(
			tokenizerbase.DefaultTokenSource,
			tokenizerbase.DefaultContextWindow,
		),
	))

	tokensModel, err := s.parseTokens(ctx, &scraperEvent, tokens)
	if err != nil {
		return fmt.Errorf("[%s] failed to parse tokens: %w", op, err)
	}

	// try to insert tokens
	if err := s.repo.InsertTokens(ctx, tokensModel); err != nil {
		return fmt.Errorf("[%s] failed to insert tokens batch: %w", op, err)
	}

	return nil
}

// parseTokens parses tokens and enriches them with features
func (s *Service) parseTokens(
	ctx context.Context,
	msg *models.ScraperEvent,
	tokens []tokenizerbase.Token,
) ([]models.TokenData, error) {
	var (
		log    = ctxutils.GetLogger(ctx)
		site   = msg.SiteName
		date   = msg.Date
		result = make([]models.TokenData, 0, len(tokens))
	)

	dateParsed, err := time.Parse(models.ScrapeDataFormat, date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scrape date %s: %w", date, err)
	}

	for _, t := range tokens {
		// skip filtered tokens
		if t.IsFiltered() {
			continue
		}

		// parse token's context
		tokenContext := make([]string, len(t.Context))
		for i, ctxToken := range t.Context {
			tokenContext[i] = ctxToken.Target
		}

		resp, err := s.llm.SentimentAnalysis(ctx, &llm.SentimentAnalysisRequest{
			Context: strings.Join(tokenContext, " "),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to analyze sentiment for token %s: %w", t.Target, err)
		}

		// get metrics
		interest, ok := metrics.Registry[metrics.InterestMetricIndex].Get(t.Target)
		if !ok {
			log.Errorf("failed to get interest metric for token %s", t.Target)
			continue
		}

		interestInt, ok := interest.(int64)
		if !ok {
			log.Errorf("interest metric has invalid type for token %s: %T", t.Target, interest)
			continue
		}

		// append to result
		result = append(result, models.TokenData{
			TokenName: t.Target,
			Interest:  interestInt,
			Sentiment: resp.Sentiment,
			SiteName:  site,
			Date:      dateParsed,
		})
	}

	return result, nil
}

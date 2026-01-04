package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/keenywheels/backend/internal/pkg/client/llm"
	tokenizerbase "github.com/keenywheels/backend/internal/pkg/tokenizer"
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

	// create tokenizer pipeline
	tokenizer, registry := getTokenizer()

	// tokenize msg
	tokens := tokenizer.Run(tokenizerbase.GetTokens(
		scraperEvent.Msg,
		tokenizerbase.NewTokenConfig(
			tokenizerbase.DefaultTokenSource,
			tokenizerbase.DefaultContextWindow,
		),
	))

	tokensModel, err := s.parseTokens(ctx, &scraperEvent, tokens, registry)
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
	registry metricsRegistry,
) ([]models.TokenData, error) {
	var (
		log      = ctxutils.GetLogger(ctx)
		site     = msg.SiteName
		date     = msg.Date
		category = msg.Category
		result   = make([]models.TokenData, 0, len(tokens))
	)

	dateParsed, err := time.Parse(models.ScrapeDataFormat, date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scrape date %s: %w", date, err)
	}

	uniqRes := make(map[string]int64)
	tokensContext := make(map[string]*strings.Builder)

	for _, t := range tokens {
		// skip filtered tokens
		if t.IsFiltered() {
			continue
		}

		// add token if not exists
		if _, ok := uniqRes[t.Target]; !ok {
			tokensContext[t.Target] = new(strings.Builder) // initialize context builder

			interest, ok := registry[interestMetricKey].Get(t.Target)
			if !ok {
				log.Errorf("failed to get interest metric for token %s", t.Target)
				continue
			}

			interestInt, ok := interest.(int64)
			if !ok {
				log.Errorf("interest metric has invalid type for token %s: %T", t.Target, interest)
				continue
			}

			uniqRes[t.Target] = interestInt
		}

		// append tokens context
		for _, ctxToken := range t.Context {
			tokensContext[t.Target].WriteString(ctxToken.Target + " ")
		}
		tokensContext[t.Target].WriteString("; ")
	}

	// make final result
	for tokenName, interest := range uniqRes {
		// analyze sentiment
		resp, err := s.llm.SentimentAnalysis(ctx, &llm.SentimentAnalysisRequest{
			Context: tokensContext[tokenName].String(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to analyze sentiment for token %s: %w", tokenName, err)
		}

		// append to result
		result = append(result, models.TokenData{
			TokenName: tokenName,
			Interest:  interest,
			Sentiment: resp.Sentiment,
			SiteName:  site,
			Category:  category,
			Date:      dateParsed,
		})
	}

	return result, nil
}

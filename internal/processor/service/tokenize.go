package service

import (
	"context"
	"encoding/json"
	"fmt"

	tokenizerbase "github.com/keenywheels/backend/internal/pkg/tokenizer"
	"github.com/keenywheels/backend/internal/processor/models"
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

	tokensModel := convertToRepositoryTokens(ctx, &scraperEvent, tokens)

	// try to insert tokens
	if err := s.repo.InsertTokens(ctx, tokensModel); err != nil {
		return fmt.Errorf("[%s] failed to insert tokens batch: %w", op, err)
	}

	return nil
}

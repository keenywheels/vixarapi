package service

import (
	"context"

	"github.com/keenywheels/backend/internal/pkg/tokenizer"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/metrics"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/pkg/stemmer"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/stages"
	"github.com/keenywheels/backend/internal/processor/models"
)

// IRepository defines the interface for repository layer interactions
type IRepository interface {
	InsertTokens(ctx context.Context, tokens []models.TokenData) error
}

// Service struct for service layer logic
type Service struct {
	repo      IRepository
	tokenizer *tokenizer.Pipeline
}

// New creates a new instance of Service
func New(repo IRepository) *Service {
	return &Service{
		repo:      repo,
		tokenizer: getTokenizer(),
	}
}

// getTokenizer initializes and returns a tokenizer pipeline
func getTokenizer() *tokenizer.Pipeline {
	stgs := []tokenizer.PipelineStage{
		stages.NewNormalizerStage(),
		stages.NewFilterStage(stages.DefaultTokenMinLength),
		stages.NewStemmerStage(stemmer.DefaultStemmer),
		stages.NewMetricStage(metrics.Registry...),
	}

	return tokenizer.NewPipelineBuilder().AddStages(stgs...).Build()
}

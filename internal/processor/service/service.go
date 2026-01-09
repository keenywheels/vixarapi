package service

import (
	"context"

	"github.com/keenywheels/backend/internal/pkg/client/llm"
	"github.com/keenywheels/backend/internal/pkg/tokenizer"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/metrics"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/pkg/stemmer"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/stages"
	"github.com/keenywheels/backend/internal/processor/models"
	"github.com/keenywheels/backend/pkg/mailer"
)

const (
	interestMetricKey = "interest"
)

// IClientLLM define the intervace for LLM client interactions
type IClientLLM interface {
	SentimentAnalysis(ctx context.Context, req *llm.SentimentAnalysisRequest) (*llm.SentimentAnalysisResponse, error)
}

// IRepository defines the interface for repository layer interactions
type IRepository interface {
	InsertTokens(ctx context.Context, tokens []models.TokenData) error
}

// Service struct for service layer logic
type Service struct {
	repo   IRepository
	llm    IClientLLM
	mailer mailer.Mailer
}

// New creates a new instance of Service
func New(
	repo IRepository,
	llm IClientLLM,
	mailer mailer.Mailer,
) *Service {
	return &Service{
		repo:   repo,
		llm:    llm,
		mailer: mailer,
	}
}

// metricsRegistry is a type alias for a map of metric names to Metric instances
type metricsRegistry map[string]metrics.Metric

// getTokenizer initializes and returns a tokenizer pipeline with metrics registry
func getTokenizer() (*tokenizer.Pipeline, metricsRegistry) {
	interest := metrics.NewInterestMetric()

	registry := metricsRegistry{
		interestMetricKey: interest,
		// add more metrics here if needed
	}

	stgs := []tokenizer.PipelineStage{
		stages.NewNormalizerStage(),
		stages.NewFilterStage(stages.DefaultTokenMinLength),
		stages.NewStemmerStage(stemmer.DefaultStemmer),
		stages.NewMetricStage([]metrics.Metric{
			interest,
		}...),
	}

	return tokenizer.NewPipelineBuilder().AddStages(stgs...).Build(), registry
}

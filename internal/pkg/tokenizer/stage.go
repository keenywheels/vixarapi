package tokenizer

import "github.com/keenywheels/backend/internal/pkg/tokenizer/models"

// PipelineStage defines the interface for a stage in the token processing pipeline.
type PipelineStage interface {
	Execute(tokens []models.Token) []models.Token
	Continue(tokens []models.Token) []models.Token
	SetNext(stage PipelineStage)
}

var _ = PipelineStage(&Stage{})

// Stage represents a processing stage in the token pipeline.
type Stage struct {
	NextStage    PipelineStage
	CallbackFunc func(token *models.Token) error
}

// Execute processes the tokens using the stage's callback function and continues to the next stage.
func (s *Stage) Execute(tokens []models.Token) []models.Token {
	for i := range tokens {
		if err := s.CallbackFunc(&tokens[i]); err != nil {
			continue
		}
	}

	return s.Continue(tokens)
}

// Continue passes the tokens to the next stage in the pipeline, if it exists.
func (s *Stage) Continue(tokens []models.Token) []models.Token {
	if s.NextStage != nil {
		return s.NextStage.Execute(tokens)
	}

	return tokens
}

// SetNext sets the next stage in the pipeline.
func (s *Stage) SetNext(stage PipelineStage) {
	s.NextStage = stage
}

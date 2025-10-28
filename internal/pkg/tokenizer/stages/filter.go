package stages

import (
	"github.com/keenywheels/backend/internal/pkg/tokenizer"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/models"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/pkg/stopwords"
)

const DefaultTokenMinLength = 3

// NewFilterStage creates a new filtering stage that removes tokens
func NewFilterStage(tokenMinLength int) *tokenizer.Stage {
	stage := &tokenizer.Stage{}

	tokenMinLength = getTokenMinLength(tokenMinLength)
	stage.CallbackFunc = func(token *models.Token) error {
		if len([]rune(token.Target)) < tokenMinLength {
			token.Filter()
		}

		if _, isStop := stopwords.All[token.Target]; isStop {
			token.Filter()
		}

		return nil
	}

	return stage
}

// getTokenMinLength returns the minimum token length or the default if invalid
func getTokenMinLength(minLength int) int {
	if minLength <= 0 {
		return DefaultTokenMinLength
	}
	return minLength
}

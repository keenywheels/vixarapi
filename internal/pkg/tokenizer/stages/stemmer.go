package stages

import (
	"fmt"

	"github.com/keenywheels/backend/internal/pkg/tokenizer"
)

// Stemmer defines the interface for stemming tokens.
type Stemmer interface {
	Stem(token string) (string, error)
}

// NewStemmerStage creates a new stemming stage that applies the provided stemmer to tokens.
func NewStemmerStage(stemmer Stemmer) *tokenizer.Stage {
	stage := &tokenizer.Stage{}

	stage.CallbackFunc = func(token *tokenizer.Token) error {
		stemmedToken, err := stemmer.Stem(token.Target)
		if err != nil {
			return fmt.Errorf("stemmer failed: %w", err)
		}

		token.Target = stemmedToken

		return nil
	}

	return stage
}

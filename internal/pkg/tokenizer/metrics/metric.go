package metrics

import "github.com/keenywheels/backend/internal/pkg/tokenizer"

// Metric defines the interface for collecting and retrieving metrics related to tokens.
type Metric interface {
	Collect(token *tokenizer.Token) error
	Get(token string) (any, bool)
}

package metrics

import "github.com/keenywheels/backend/internal/pkg/tokenizer/models"

// Metric defines the interface for collecting and retrieving metrics related to tokens.
type Metric interface {
	Collect(token *models.Token) error
	Get(token string) (any, bool)
}

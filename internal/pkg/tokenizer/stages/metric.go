package stages

import (
	"github.com/keenywheels/backend/internal/pkg/tokenizer"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/metrics"
	"github.com/keenywheels/backend/internal/pkg/tokenizer/models"
)

// NewMetricStage creates a new metric collection stage
func NewMetricStage(metrics ...metrics.Metric) *tokenizer.Stage {
	stage := &tokenizer.Stage{}

	stage.CallbackFunc = func(token *models.Token) error {
		for _, m := range metrics {
			if err := m.Collect(token); err != nil {
				continue
			}
		}
		return nil
	}

	return stage
}

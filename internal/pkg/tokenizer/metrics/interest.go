package metrics

import (
	"sync"

	"github.com/keenywheels/backend/internal/pkg/tokenizer/models"
)

// InterestMetric tracks the number of times tokens target specific entities.
type InterestMetric struct {
	mu     sync.RWMutex
	counts map[string]int
}

// NewInterestMetric creates a new instance of InterestMetric.
func NewInterestMetric() *InterestMetric {
	return &InterestMetric{
		counts: make(map[string]int),
	}
}

// Collect implements collecting interest metrics for a given token.
func (m *InterestMetric) Collect(token *models.Token) error {
	if token.Target == "" {
		return nil
	}

	m.mu.Lock()
	m.counts[token.Target]++
	m.mu.Unlock()

	return nil
}

// Get implements retrieving the interest metric for a given token.
func (m *InterestMetric) Get(token string) (any, bool) {
	m.mu.RLock()
	val, ok := m.counts[token]
	m.mu.RUnlock()

	if !ok {
		return nil, false
	}
	return val, true
}

package metrics

// Indices of available metrics
const (
	InterestMetricIndex = iota
)

// Registry holds all available metrics for token processing.
var Registry = []Metric{
	NewInterestMetric(),
}

package metrics

// Registry holds all available metrics for token processing.
var Registry = []Metric{
	NewInterestMetric(),
}

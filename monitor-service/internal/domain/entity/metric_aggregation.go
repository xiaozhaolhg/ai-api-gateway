package entity

// MetricAggregation represents aggregated metric statistics
type MetricAggregation struct {
	ProviderID   string  `json:"provider_id"`
	Model        string  `json:"model"`
	MetricType   string  `json:"metric_type"`
	AvgValue     float64 `json:"avg_value"`
	MinValue     float64 `json:"min_value"`
	MaxValue     float64 `json:"max_value"`
	SumValue     float64 `json:"sum_value"`
	Count        int64   `json:"count"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
}

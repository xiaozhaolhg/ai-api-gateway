package entity

import "time"

// Metric represents a metric data point
type Metric struct {
	ID           string    `json:"id"`
	ProviderID   string    `json:"provider_id"`
	Model        string    `json:"model"`
	MetricType   string    `json:"metric_type"` // "latency" | "error_rate" | "throughput"
	Value        float64   `json:"value"`
	Timestamp    time.Time `json:"timestamp"`
}

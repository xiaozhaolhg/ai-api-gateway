package entity

import "time"

// AlertRule represents an alert rule
type AlertRule struct {
	ID              string    `json:"id"`
	ProviderID      string    `json:"provider_id"`
	MetricType      string    `json:"metric_type"`
	Threshold       float64   `json:"threshold"`
	Operator        string    `json:"operator"` // ">" | "<" | "="
	WindowMinutes   int       `json:"window_minutes"`
	Severity        string    `json:"severity"` // "info" | "warning" | "critical"
	Enabled         bool      `json:"enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

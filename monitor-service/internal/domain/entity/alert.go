package entity

import "time"

// Alert represents a triggered alert
type Alert struct {
	ID              string    `json:"id"`
	AlertRuleID     string    `json:"alert_rule_id"`
	ProviderID      string    `json:"provider_id"`
	Severity        string    `json:"severity"`
	Message         string    `json:"message"`
	Value           float64   `json:"value"`
	Threshold       float64   `json:"threshold"`
	Timestamp       time.Time `json:"timestamp"`
	Acknowledged    bool      `json:"acknowledged"`
	AcknowledgedAt  *time.Time `json:"acknowledged_at,omitempty"`
}

package entity

import "time"

// RoutingRule represents a routing rule for model-to-provider mapping
type RoutingRule struct {
	ID                 string    `json:"id"`
	ModelPattern       string    `json:"model_pattern"`       // e.g., "gpt-4*" or "claude-*"
	ProviderID         string    `json:"provider_id"`
	Priority           int32     `json:"priority"`
	FallbackProviderID string    `json:"fallback_provider_id"`
	CreatedAt          time.Time `json:"created_at"`
}

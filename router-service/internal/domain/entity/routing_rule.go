package entity

import "time"

// RoutingRule represents a routing rule for model-to-provider mapping
type RoutingRule struct {
	ID                  string    `json:"id" gorm:"primaryKey"`
	UserID              string    `json:"user_id" gorm:"index:idx_routing_rules_user_id"`
	ModelPattern        string    `json:"model_pattern" gorm:"index:idx_routing_rules_model_pattern"`
	ProviderID          string    `json:"provider_id"`
	Priority            int32     `json:"priority"`
	FallbackProviderIDs string    `json:"fallback_provider_ids"` // JSON array stored as TEXT
	FallbackModels      string    `json:"fallback_models"`       // JSON array stored as TEXT
	IsSystemDefault    bool      `json:"is_system_default" gorm:"default:false"`
	CreatedAt           time.Time `json:"created_at"`
}

// RouteResult represents the result of route resolution
type RouteResult struct {
	ProviderID          string   `json:"provider_id"`
	AdapterType         string   `json:"adapter_type"`          // "openai" | "anthropic" | "gemini" | "ollama" | "opencode-zen"
	FallbackProviderIDs []string `json:"fallback_provider_ids"` // ordered list of fallback providers
	FallbackModels      []string `json:"fallback_models"`      // parallel array mapping fallback provider → model
}

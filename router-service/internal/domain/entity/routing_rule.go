package entity

import "time"

// RoutingRule represents a routing rule for model-to-provider mapping
type RoutingRule struct {
	ID                 string    `json:"id"`
	ModelPattern       string    `json:"model_pattern"`       // e.g., "gpt-4*" or "claude-*"
	ProviderID         string    `json:"provider_id"`
	Priority           int32     `json:"priority"`
	FallbackProviderID string    `json:"fallback_provider_id"`
	FallbackModel      string    `json:"fallback_model"`
	CreatedAt          time.Time `json:"created_at"`
}

// RouteResult represents the result of route resolution
type RouteResult struct {
	ProviderID          string   `json:"provider_id"`
	AdapterType         string   `json:"adapter_type"`          // "openai" | "anthropic" | "gemini" | "ollama" | "opencode-zen"
	FallbackProviderIDs []string `json:"fallback_provider_ids"` // ordered list of fallback providers
	FallbackModels      []string `json:"fallback_models"`      // parallel array mapping fallback provider → model
}

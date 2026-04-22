package entity

import "time"

// PricingRule represents a pricing rule for a provider/model
type PricingRule struct {
	ID                string    `json:"id"`
	ProviderID        string    `json:"provider_id"`
	Model             string    `json:"model"`
	PromptPricePer1K  float64   `json:"prompt_price_per_1k"`
	CompletionPricePer1K float64 `json:"completion_price_per_1k"`
	Currency          string    `json:"currency"`
	EffectiveFrom     time.Time `json:"effective_from"`
	EffectiveUntil    *time.Time `json:"effective_until,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

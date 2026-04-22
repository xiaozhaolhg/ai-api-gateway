package entity

import "time"

// UsageRecord represents a usage record for a request
type UsageRecord struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	ProviderID       string    `json:"provider_id"`
	Model            string    `json:"model"`
	PromptTokens     int64     `json:"prompt_tokens"`
	CompletionTokens int64     `json:"completion_tokens"`
	Cost             float64   `json:"cost"`
	Timestamp        time.Time `json:"timestamp"`
}

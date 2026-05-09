package entity

import "time"

// Group represents a user group with authorization configuration
type Group struct {
	ID            string      `json:"id" gorm:"primaryKey"`
	Name          string      `json:"name" gorm:"uniqueIndex;not null"`
	Description   string      `json:"description"`
	ParentGroupID string      `json:"parent_group_id"`
	ModelPatterns Strings     `json:"model_patterns" gorm:"serializer:json"`
	TokenLimit    *TokenLimit `json:"token_limit,omitempty" gorm:"serializer:json"`
	RateLimit     *RateLimit  `json:"rate_limit,omitempty" gorm:"serializer:json"`
	TierID        string      `json:"tier_id,omitempty" gorm:"column:tier_id"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// TokenLimit defines token usage limits for a group
type TokenLimit struct {
	PromptTokens     int64  `json:"prompt_tokens"`
	CompletionTokens int64  `json:"completion_tokens"`
	Period           string `json:"period"` // "daily" | "weekly" | "monthly"
}

// RateLimit defines request rate limits for a group
type RateLimit struct {
	RequestsPerMinute int64 `json:"requests_per_minute"`
	RequestsPerDay    int64 `json:"requests_per_day"`
}

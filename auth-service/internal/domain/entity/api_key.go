package entity

import "time"

// APIKey represents an API key for authentication
type APIKey struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	KeyHash    string     `json:"key_hash"`              // SHA-256 hash of the API key
	Name       string     `json:"name"`
	Scopes     Strings    `json:"scopes" gorm:"serializer:json"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// Strings is a JSON-serializable string slice for GORM
type Strings []string

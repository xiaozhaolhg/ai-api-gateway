package entity

import "time"

// Provider represents an LLM provider
type Provider struct {
	ID        string   `json:"id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`       // "openai" | "anthropic" | "gemini" | "ollama" | "custom"
	BaseURL  string   `json:"base_url"`
	Credentials string `json:"credentials"` // encrypted
	Models   Strings  `json:"models" gorm:"serializer:json"`
	Status   string   `json:"status"`     // "active" | "inactive"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Strings is a JSON-serializable string slice for GORM
type Strings []string

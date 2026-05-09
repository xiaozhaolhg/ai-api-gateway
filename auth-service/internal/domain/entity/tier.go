package entity

import "time"

// Tier represents a permission tier that defines allowed providers and models for groups
type Tier struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsDefault   bool      `json:"is_default"`
	AllowedModels []string  `json:"allowed_models" gorm:"serializer:json"`
	AllowedProviders []string `json:"allowed_providers" gorm:"serializer:json"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
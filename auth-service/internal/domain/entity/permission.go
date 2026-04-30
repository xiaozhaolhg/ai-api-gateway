package entity

import "time"

// Permission represents an authorization rule for a group
type Permission struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	GroupID      string    `json:"group_id" gorm:"index;not null"`
	ResourceType string    `json:"resource_type" gorm:"not null"` // "model" | "provider" | "admin_feature"
	ResourceID   string    `json:"resource_id" gorm:"not null"`   // glob pattern or feature name
	Action       string    `json:"action" gorm:"not null"`        // "access" | "manage" | "view"
	Effect       string    `json:"effect" gorm:"not null"`        // "allow" | "deny"
	Status       string    `json:"status" gorm:"not null"`        // "active" | "revoked"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

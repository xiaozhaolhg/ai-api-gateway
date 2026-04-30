package entity

import "time"

// Budget represents a user's budget limit
type Budget struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Limit          float64    `json:"limit"`
	Period         string     `json:"period"` // "monthly" | "yearly" | "custom"
	SoftCap        float64    `json:"soft_cap"`
	HardCap        float64    `json:"hard_cap"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	AlertThreshold float64    `json:"alert_threshold"` // percentage
	Status         string     `json:"status"`          // "active" | "paused"
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

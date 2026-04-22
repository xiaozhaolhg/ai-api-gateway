package entity

import "time"

// BillingAccount represents a user's billing account
type BillingAccount struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Balance       float64   `json:"balance"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"` // "active" | "suspended" | "closed"
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

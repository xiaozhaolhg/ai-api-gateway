package entity

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`      // "admin" | "user"
	Status    string    `json:"status"`    // "active" | "disabled"
	CreatedAt time.Time `json:"created_at"`
}

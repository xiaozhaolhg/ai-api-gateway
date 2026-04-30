package entity

import "time"

// UserGroupMembership represents a user's membership in a group
type UserGroupMembership struct {
	ID      string    `json:"id" gorm:"primaryKey"`
	UserID  string    `json:"user_id" gorm:"index;not null"`
	GroupID string    `json:"group_id" gorm:"index;not null"`
	AddedAt time.Time `json:"added_at"`
}

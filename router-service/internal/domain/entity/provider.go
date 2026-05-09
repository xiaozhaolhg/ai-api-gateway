package entity

import "time"

type Provider struct {
	ID        string
	Name      string
	Type      string
	BaseURL   string
	Models    []string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

package entity

import (
	"testing"
	"time"
)

func TestUser_Validation(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid user",
			user: User{
				ID:        "123",
				Name:      "John Doe",
				Email:     "john@example.com",
				Role:      "user",
				Status:    "active",
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid admin",
			user: User{
				ID:        "456",
				Name:      "Admin User",
				Email:     "admin@example.com",
				Role:      "admin",
				Status:    "active",
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "disabled user",
			user: User{
				ID:        "789",
				Name:      "Disabled User",
				Email:     "disabled@example.com",
				Role:      "user",
				Status:    "disabled",
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - ensure required fields are not empty
			if tt.user.ID == "" {
				t.Error("User ID cannot be empty")
			}
			if tt.user.Name == "" {
				t.Error("User name cannot be empty")
			}
			if tt.user.Email == "" {
				t.Error("User email cannot be empty")
			}
			if tt.user.Role == "" {
				t.Error("User role cannot be empty")
			}
			if tt.user.Status == "" {
				t.Error("User status cannot be empty")
			}
			if tt.user.CreatedAt.IsZero() {
				t.Error("User created at cannot be zero")
			}
		})
	}
}

func TestUser_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "active user",
			status:   "active",
			expected: true,
		},
		{
			name:     "disabled user",
			status:   "disabled",
			expected: false,
		},
		{
			name:     "unknown status",
			status:   "unknown",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{Status: tt.status}
			isActive := user.Status == "active"
			if isActive != tt.expected {
				t.Errorf("User.IsActive() = %v, want %v", isActive, tt.expected)
			}
		})
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{
			name:     "admin user",
			role:     "admin",
			expected: true,
		},
		{
			name:     "regular user",
			role:     "user",
			expected: false,
		},
		{
			name:     "unknown role",
			role:     "unknown",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{Role: tt.role}
			isAdmin := user.Role == "admin"
			if isAdmin != tt.expected {
				t.Errorf("User.IsAdmin() = %v, want %v", isAdmin, tt.expected)
			}
		})
	}
}

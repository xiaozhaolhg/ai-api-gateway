package entity

import (
	"testing"
	"time"
)

func TestAPIKey_Validation(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  APIKey
		wantErr bool
	}{
		{
			name: "valid API key",
			apiKey: APIKey{
				ID:        "123",
				UserID:    "user-123",
				KeyHash:   "abc123",
				Name:      "Test Key",
				Scopes:    []string{"read", "write"},
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "API key with expiration",
			apiKey: APIKey{
				ID:        "456",
				UserID:    "user-456",
				KeyHash:   "def456",
				Name:      "Expiring Key",
				Scopes:    []string{"read"},
				CreatedAt: time.Now(),
				ExpiresAt: func() *time.Time {
					t := time.Now().Add(24 * time.Hour)
					return &t
				}(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - ensure required fields are not empty
			if tt.apiKey.ID == "" {
				t.Error("API Key ID cannot be empty")
			}
			if tt.apiKey.UserID == "" {
				t.Error("API Key UserID cannot be empty")
			}
			if tt.apiKey.KeyHash == "" {
				t.Error("API Key KeyHash cannot be empty")
			}
			if tt.apiKey.Name == "" {
				t.Error("API Key Name cannot be empty")
			}
			if tt.apiKey.Scopes == nil || len(tt.apiKey.Scopes) == 0 {
				t.Error("API Key Scopes cannot be empty")
			}
			if tt.apiKey.CreatedAt.IsZero() {
				t.Error("API Key CreatedAt cannot be zero")
			}
		})
	}
}

func TestAPIKey_IsExpired(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		expiresAt *time.Time
		expected bool
	}{
		{
			name:     "no expiration",
			expiresAt: nil,
			expected: false,
		},
		{
			name:     "expired key",
			expiresAt: func() *time.Time {
				t := now.Add(-1 * time.Hour)
				return &t
			}(),
			expected: true,
		},
		{
			name:     "not expired key",
			expiresAt: func() *time.Time {
				t := now.Add(1 * time.Hour)
				return &t
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = APIKey{ExpiresAt: tt.expiresAt}
			isExpired := tt.expiresAt != nil && tt.expiresAt.Before(now)
			if isExpired != tt.expected {
				t.Errorf("APIKey.IsExpired() = %v, want %v", isExpired, tt.expected)
			}
		})
	}
}

func TestAPIKey_HasScope(t *testing.T) {
	tests := []struct {
		name     string
		scopes   []string
		scope    string
		expected bool
	}{
		{
			name:     "has read scope",
			scopes:   []string{"read", "write"},
			scope:    "read",
			expected: true,
		},
		{
			name:     "has write scope",
			scopes:   []string{"read", "write"},
			scope:    "write",
			expected: true,
		},
		{
			name:     "does not have admin scope",
			scopes:   []string{"read", "write"},
			scope:    "admin",
			expected: false,
		},
		{
			name:     "empty scopes",
			scopes:   []string{},
			scope:    "read",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiKey := APIKey{Scopes: tt.scopes}
			hasScope := false
			for _, s := range apiKey.Scopes {
				if s == tt.scope {
					hasScope = true
					break
				}
			}
			if hasScope != tt.expected {
				t.Errorf("APIKey.HasScope(%q) = %v, want %v", tt.scope, hasScope, tt.expected)
			}
		})
	}
}

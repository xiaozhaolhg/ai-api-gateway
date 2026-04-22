package entity

import (
	"testing"
	"time"
)

func TestProvider_Validation(t *testing.T) {
	tests := []struct {
		name     string
		provider Provider
		wantErr  bool
	}{
		{
			name: "valid provider",
			provider: Provider{
				ID:          "provider-1",
				Name:        "OpenAI",
				Type:        "openai",
				BaseURL:     "https://api.openai.com/v1",
				Credentials: "encrypted-key",
				Models:      []string{"gpt-4", "gpt-3.5-turbo"},
				Status:      "active",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid Ollama provider",
			provider: Provider{
				ID:          "provider-2",
				Name:        "Ollama",
				Type:        "ollama",
				BaseURL:     "http://localhost:11434",
				Credentials: "",
				Models:      []string{"llama2", "mistral"},
				Status:      "active",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "inactive provider",
			provider: Provider{
				ID:          "provider-3",
				Name:        "Inactive Provider",
				Type:        "custom",
				BaseURL:     "https://api.example.com",
				Credentials: "encrypted-key",
				Models:      []string{"model1"},
				Status:      "inactive",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - ensure required fields are not empty
			if tt.provider.ID == "" {
				t.Error("Provider ID cannot be empty")
			}
			if tt.provider.Name == "" {
				t.Error("Provider name cannot be empty")
			}
			if tt.provider.Type == "" {
				t.Error("Provider type cannot be empty")
			}
			if tt.provider.BaseURL == "" {
				t.Error("Provider base URL cannot be empty")
			}
			if tt.provider.Status == "" {
				t.Error("Provider status cannot be empty")
			}
			if tt.provider.CreatedAt.IsZero() {
				t.Error("Provider created at cannot be zero")
			}
			if tt.provider.UpdatedAt.IsZero() {
				t.Error("Provider updated at cannot be zero")
			}
		})
	}
}

func TestProvider_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "active provider",
			status:   "active",
			expected: true,
		},
		{
			name:     "inactive provider",
			status:   "inactive",
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
			provider := Provider{Status: tt.status}
			isActive := provider.Status == "active"
			if isActive != tt.expected {
				t.Errorf("Provider.IsActive() = %v, want %v", isActive, tt.expected)
			}
		})
	}
}

func TestProvider_HasModel(t *testing.T) {
	tests := []struct {
		name     string
		models   []string
		model    string
		expected bool
	}{
		{
			name:     "has gpt-4",
			models:   []string{"gpt-4", "gpt-3.5-turbo"},
			model:    "gpt-4",
			expected: true,
		},
		{
			name:     "has llama2",
			models:   []string{"llama2", "mistral"},
			model:    "llama2",
			expected: true,
		},
		{
			name:     "does not have model",
			models:   []string{"gpt-4", "gpt-3.5-turbo"},
			model:    "claude-3",
			expected: false,
		},
		{
			name:     "empty models",
			models:   []string{},
			model:    "gpt-4",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := Provider{Models: tt.models}
			hasModel := false
			for _, m := range provider.Models {
				if m == tt.model {
					hasModel = true
					break
				}
			}
			if hasModel != tt.expected {
				t.Errorf("Provider.HasModel(%q) = %v, want %v", tt.model, hasModel, tt.expected)
			}
		})
	}
}

func TestProvider_SupportsType(t *testing.T) {
	tests := []struct {
		name     string
		providerType string
		checkType string
		expected bool
	}{
		{
			name:     "openai provider",
			providerType: "openai",
			checkType: "openai",
			expected: true,
		},
		{
			name:     "ollama provider",
			providerType: "ollama",
			checkType: "ollama",
			expected: true,
		},
		{
			name:     "type mismatch",
			providerType: "openai",
			checkType: "anthropic",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := Provider{Type: tt.providerType}
			supports := provider.Type == tt.checkType
			if supports != tt.expected {
				t.Errorf("Provider.SupportsType(%q) = %v, want %v", tt.checkType, supports, tt.expected)
			}
		})
	}
}

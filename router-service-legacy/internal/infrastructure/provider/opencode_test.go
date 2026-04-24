package provider

import (
	"testing"

	"github.com/ai-api-gateway/router-service-legacy/internal/infrastructure/config"
)

func TestOpenCodeZenFactory_Type(t *testing.T) {
	factory := NewOpenCodeZenFactory()
	if factory.Type() != "opencode_zen" {
		t.Fatalf("expected type 'opencode_zen', got '%s'", factory.Type())
	}
}

func TestOpenCodeZenFactory_Description(t *testing.T) {
	factory := NewOpenCodeZenFactory()
	expected := "OpenCode Zen - Curated AI models for coding agents"
	if factory.Description() != expected {
		t.Fatalf("expected description '%s', got '%s'", expected, factory.Description())
	}
}

func TestOpenCodeZenFactory_Validate(t *testing.T) {
	factory := NewOpenCodeZenFactory()

	// Test valid settings
	settings := config.ProviderSettings{Endpoint: "https://opencode.ai/zen"}
	err := factory.Validate(settings)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test empty endpoint
	settings = config.ProviderSettings{Endpoint: ""}
	err = factory.Validate(settings)
	if err == nil {
		t.Fatal("expected error for empty endpoint")
	}
}

func TestOpenCodeZenFactory_Defaults(t *testing.T) {
	factory := NewOpenCodeZenFactory()
	defaults := factory.Defaults()

	if defaults.Endpoint != "https://opencode.ai/zen" {
		t.Fatalf("expected endpoint 'https://opencode.ai/zen', got '%s'", defaults.Endpoint)
	}
	if defaults.Enabled != false {
		t.Fatal("expected enabled false")
	}
	if defaults.APIKey != "" {
		t.Fatal("expected empty api_key")
	}
}

func TestOpenCodeZenFactory_Create(t *testing.T) {
	factory := NewOpenCodeZenFactory()
	settings := config.ProviderSettings{Endpoint: "https://opencode.ai/zen"}

	provider := factory.Create(settings)
	if provider == nil {
		t.Fatal("expected non-nil provider")
	}

	if provider.Name() != "opencode_zen" {
		t.Fatalf("expected provider name 'opencode_zen', got '%s'", provider.Name())
	}
}

func TestExtractOpenCodeModelName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"opencode_zen:gpt-4", "gpt-4"},
		{"opencode_zen:claude-3", "claude-3"},
		{"gpt-4", "gpt-4"},
		{"opencode:gpt-4", "opencode:gpt-4"}, // Old prefix not stripped
	}

	for _, tt := range tests {
		result := extractOpenCodeModelName(tt.input)
		if result != tt.expected {
			t.Fatalf("expected '%s' for input '%s', got '%s'", tt.expected, tt.input, result)
		}
	}
}

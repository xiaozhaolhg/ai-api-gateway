package provider

import (
	"testing"

	"github.com/ai-api-gateway/router-service/internal/infrastructure/config"
)

func TestOllamaFactory_Type(t *testing.T) {
	factory := NewOllamaFactory()
	if factory.Type() != "ollama" {
		t.Fatalf("expected type 'ollama', got '%s'", factory.Type())
	}
}

func TestOllamaFactory_Description(t *testing.T) {
	factory := NewOllamaFactory()
	expected := "Ollama - Run LLMs locally"
	if factory.Description() != expected {
		t.Fatalf("expected description '%s', got '%s'", expected, factory.Description())
	}
}

func TestOllamaFactory_Validate(t *testing.T) {
	factory := NewOllamaFactory()

	// Test valid settings
	settings := config.ProviderSettings{Endpoint: "http://localhost:11434"}
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

func TestOllamaFactory_Defaults(t *testing.T) {
	factory := NewOllamaFactory()
	defaults := factory.Defaults()

	if defaults.Endpoint != "http://localhost:11434" {
		t.Fatalf("expected endpoint 'http://localhost:11434', got '%s'", defaults.Endpoint)
	}
	if defaults.Enabled != false {
		t.Fatal("expected enabled false")
	}
	if defaults.APIKey != "" {
		t.Fatal("expected empty api_key")
	}
}

func TestOllamaFactory_Create(t *testing.T) {
	factory := NewOllamaFactory()
	settings := config.ProviderSettings{Endpoint: "http://localhost:11434"}

	provider := factory.Create(settings)
	if provider == nil {
		t.Fatal("expected non-nil provider")
	}

	if provider.Name() != "ollama" {
		t.Fatalf("expected provider name 'ollama', got '%s'", provider.Name())
	}
}

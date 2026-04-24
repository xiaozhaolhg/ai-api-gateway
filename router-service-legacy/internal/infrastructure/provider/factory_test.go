package provider

import (
	"fmt"
	"testing"

	"github.com/ai-api-gateway/router-service-legacy/internal/domain/entity"
	"github.com/ai-api-gateway/router-service-legacy/internal/domain/port"
	"github.com/ai-api-gateway/router-service-legacy/internal/infrastructure/config"
)

// mockFactory is a test implementation of ProviderFactory
type mockFactory struct {
	providerType  string
	description   string
	validateFunc  func(config.ProviderSettings) error
	createFunc    func(config.ProviderSettings) (port.Provider, error)
	defaultsFunc  func() config.ProviderSettings
}

func (m *mockFactory) Type() string {
	return m.providerType
}

func (m *mockFactory) Create(settings config.ProviderSettings) (port.Provider, error) {
	if m.createFunc != nil {
		return m.createFunc(settings)
	}
	return nil, nil
}

func (m *mockFactory) Validate(settings config.ProviderSettings) error {
	if m.validateFunc != nil {
		return m.validateFunc(settings)
	}
	return nil
}

func (m *mockFactory) Defaults() config.ProviderSettings {
	if m.defaultsFunc != nil {
		return m.defaultsFunc()
	}
	return config.ProviderSettings{}
}

func (m *mockFactory) Description() string {
	return m.description
}

func TestNewProviderRegistry(t *testing.T) {
	registry := NewProviderRegistry()
	if registry == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestProviderRegistry_Register(t *testing.T) {
	registry := NewProviderRegistry()
	factory := &mockFactory{
		providerType: "test",
		description:  "Test provider",
	}

	err := registry.Register(factory)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test duplicate registration
	err = registry.Register(factory)
	if err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestProviderRegistry_ListTypes(t *testing.T) {
	registry := NewProviderRegistry()
	factory1 := &mockFactory{providerType: "test1"}
	factory2 := &mockFactory{providerType: "test2"}

	registry.Register(factory1)
	registry.Register(factory2)

	types := registry.ListTypes()
	if len(types) != 2 {
		t.Fatalf("expected 2 types, got %d", len(types))
	}
}

func TestProviderRegistry_GetFactory(t *testing.T) {
	registry := NewProviderRegistry()
	factory := &mockFactory{providerType: "test"}
	registry.Register(factory)

	retrieved, err := registry.GetFactory("test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrieved != factory {
		t.Fatal("expected same factory instance")
	}

	_, err = registry.GetFactory("unknown")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestProviderRegistry_GetDefaults(t *testing.T) {
	registry := NewProviderRegistry()
	factory1 := &mockFactory{
		providerType: "test1",
		defaultsFunc: func() config.ProviderSettings {
			return config.ProviderSettings{Endpoint: "http://test1"}
		},
	}
	factory2 := &mockFactory{
		providerType: "test2",
		defaultsFunc: func() config.ProviderSettings {
			return config.ProviderSettings{Endpoint: "http://test2"}
		},
	}

	registry.Register(factory1)
	registry.Register(factory2)

	defaults := registry.GetDefaults()
	if len(defaults) != 2 {
		t.Fatalf("expected 2 defaults, got %d", len(defaults))
	}
	if defaults["test1"].Endpoint != "http://test1" {
		t.Fatal("expected test1 default endpoint")
	}
}

func TestProviderRegistry_ValidateSettings(t *testing.T) {
	registry := NewProviderRegistry()
	factory := &mockFactory{
		providerType: "test",
		validateFunc: func(settings config.ProviderSettings) error {
			if settings.Endpoint == "" {
				return fmt.Errorf("endpoint required")
			}
			return nil
		},
	}
	registry.Register(factory)

	settings := config.ProviderSettings{Endpoint: "http://test"}
	err := registry.ValidateSettings("test", settings)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	settings = config.ProviderSettings{Endpoint: ""}
	err = registry.ValidateSettings("test", settings)
	if err == nil {
		t.Fatal("expected validation error")
	}

	_, err = registry.ValidateSettings("unknown", settings)
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestProviderRegistry_Create(t *testing.T) {
	registry := NewProviderRegistry()
	factory := &mockFactory{
		providerType: "test",
		validateFunc: func(settings config.ProviderSettings) error {
			if settings.Endpoint == "" {
				return fmt.Errorf("endpoint required")
			}
			return nil
		},
		createFunc: func(settings config.ProviderSettings) (port.Provider, error) {
			return &mockProvider{name: "test"}, nil
		},
	}
	registry.Register(factory)

	settings := config.ProviderSettings{Endpoint: "http://test"}
	provider, err := registry.Create("test", settings)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if provider == nil {
		t.Fatal("expected non-nil provider")
	}

	settings = config.ProviderSettings{Endpoint: ""}
	_, err = registry.Create("test", settings)
	if err == nil {
		t.Fatal("expected validation error")
	}

	_, err = registry.Create("unknown", settings)
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

// mockProvider is a test implementation of port.Provider
type mockProvider struct {
	name string
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) ChatCompletion(req port.ChatCompletionRequest) (*port.ChatCompletionResponse, error) {
	return nil, nil
}

func (m *mockProvider) StreamChatCompletion(req port.ChatCompletionRequest) (<-chan port.StreamChunk, error) {
	return nil, nil
}

func (m *mockProvider) ListModels() ([]entity.Model, error) {
	return nil, nil
}

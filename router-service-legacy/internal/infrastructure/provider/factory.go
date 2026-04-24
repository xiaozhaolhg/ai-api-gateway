package provider

import (
	"fmt"
	"sync"

	"github.com/ai-api-gateway/router-service-legacy/internal/domain/port"
	"github.com/ai-api-gateway/router-service-legacy/internal/infrastructure/config"
)

// ProviderFactory defines the interface for creating provider instances
type ProviderFactory interface {
	// Type returns the unique identifier for this provider type
	Type() string
	// Create instantiates a new provider with the given settings
	Create(settings config.ProviderSettings) (port.Provider, error)
	// Validate checks if the settings are valid for this provider type
	Validate(settings config.ProviderSettings) error
	// Defaults returns the default settings for this provider type
	Defaults() config.ProviderSettings
	// Description returns a human-readable description of this provider
	Description() string
}

// ProviderRegistry manages factory registration and provider instantiation
type ProviderRegistry struct {
	factories map[string]ProviderFactory
	mu        sync.RWMutex
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		factories: make(map[string]ProviderFactory),
	}
}

// Register adds a new factory to the registry
func (r *ProviderRegistry) Register(factory ProviderFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	providerType := factory.Type()
	if _, exists := r.factories[providerType]; exists {
		return fmt.Errorf("provider type %s already registered", providerType)
	}

	r.factories[providerType] = factory
	return nil
}

// Create instantiates a provider using the registered factory
func (r *ProviderRegistry) Create(providerType string, settings config.ProviderSettings) (port.Provider, error) {
	r.mu.RLock()
	factory, exists := r.factories[providerType]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown provider type: %s, available types: %v", providerType, r.ListTypes())
	}

	if err := factory.Validate(settings); err != nil {
		return nil, fmt.Errorf("validation failed for provider %s: %w", providerType, err)
	}

	provider, err := factory.Create(settings)
	if err != nil {
		return nil, fmt.Errorf("factory create failed for provider %s: %w", providerType, err)
	}

	return provider, nil
}

// ListTypes returns all registered provider types
func (r *ProviderRegistry) ListTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.factories))
	for t := range r.factories {
		types = append(types, t)
	}
	return types
}

// GetFactory returns the factory for a given provider type
func (r *ProviderRegistry) GetFactory(providerType string) (ProviderFactory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, exists := r.factories[providerType]
	if !exists {
		return nil, fmt.Errorf("unknown provider type: %s", providerType)
	}
	return factory, nil
}

// GetDefaults returns the default settings for all registered providers
func (r *ProviderRegistry) GetDefaults() map[string]config.ProviderSettings {
	r.mu.RLock()
	defer r.mu.RUnlock()

	defaults := make(map[string]config.ProviderSettings)
	for t, factory := range r.factories {
		defaults[t] = factory.Defaults()
	}
	return defaults
}

// ValidateSettings validates settings for a specific provider type
func (r *ProviderRegistry) ValidateSettings(providerType string, settings config.ProviderSettings) error {
	r.mu.RLock()
	factory, exists := r.factories[providerType]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("unknown provider type: %s", providerType)
	}

	return factory.Validate(settings)
}

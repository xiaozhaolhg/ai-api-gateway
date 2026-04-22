package application

import (
	"fmt"
	"sync"

	"github.com/ai-api-gateway/provider-service/internal/domain/port"
)

// AdapterFactory manages provider adapters
type AdapterFactory struct {
	adapters map[string]port.ProviderAdapter
	mu       sync.RWMutex
}

// NewAdapterFactory creates a new adapter factory
func NewAdapterFactory() *AdapterFactory {
	return &AdapterFactory{
		adapters: make(map[string]port.ProviderAdapter),
	}
}

// RegisterAdapter registers a provider adapter
func (f *AdapterFactory) RegisterAdapter(providerType string, adapter port.ProviderAdapter) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.adapters[providerType] = adapter
}

// GetAdapter retrieves a provider adapter by type
func (f *AdapterFactory) GetAdapter(providerType string) (port.ProviderAdapter, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	adapter, ok := f.adapters[providerType]
	if !ok {
		return nil, fmt.Errorf("no adapter registered for provider type: %s", providerType)
	}

	return adapter, nil
}

// ListAdapters returns all registered adapter types
func (f *AdapterFactory) ListAdapters() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	types := make([]string, 0, len(f.adapters))
	for t := range f.adapters {
		types = append(types, t)
	}

	return types
}

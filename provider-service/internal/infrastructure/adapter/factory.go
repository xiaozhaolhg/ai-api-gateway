package adapter

import (
	"github.com/ai-api-gateway/provider-service/internal/application"
	"github.com/ai-api-gateway/provider-service/internal/domain/port"
)

// NewAdapterFactory creates a new adapter factory with all adapters registered
func NewAdapterFactory() *application.AdapterFactory {
	factory := application.NewAdapterFactory()

	// Register all available adapters
	factory.RegisterAdapter("openai", NewOpenAIAdapter())
	factory.RegisterAdapter("anthropic", NewAnthropicAdapter())
	factory.RegisterAdapter("ollama", NewOllamaAdapter())
	factory.RegisterAdapter("opencode-zen", NewOpenCodeZenAdapter())
	factory.RegisterAdapter("gemini", NewGeminiAdapter())

	return factory
}

// GetAdapterByType is a convenience function to get an adapter by type
func GetAdapterByType(providerType string) (port.ProviderAdapter, error) {
	factory := NewAdapterFactory()
	return factory.GetAdapter(providerType)
}

package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ai-api-gateway/router-service-legacy/internal/infrastructure/config"
	"github.com/ai-api-gateway/router-service-legacy/internal/infrastructure/provider"
	"github.com/gin-gonic/gin"
)

// TestFactoryRegistrationAndProviderInstantiation tests factory registration and provider instantiation
func TestFactoryRegistrationAndProviderInstantiation(t *testing.T) {
	registry := provider.NewProviderRegistry()

	// Register factories
	ollamaFactory := provider.NewOllamaFactory()
	if err := registry.Register(ollamaFactory); err != nil {
		t.Fatalf("Failed to register Ollama factory: %v", err)
	}

	opencodeFactory := provider.NewOpenCodeZenFactory()
	if err := registry.Register(opencodeFactory); err != nil {
		t.Fatalf("Failed to register OpenCode Zen factory: %v", err)
	}

	// Test duplicate registration
	err := registry.Register(ollamaFactory)
	if err == nil {
		t.Fatal("Expected error for duplicate registration")
	}

	// Test provider creation
	settings := config.ProviderSettings{
		Enabled:  true,
		Endpoint: "http://localhost:11434",
		APIKey:   "",
	}

	ollamaProvider, err := registry.Create("ollama", settings)
	if err != nil {
		t.Fatalf("Failed to create Ollama provider: %v", err)
	}
	if ollamaProvider == nil {
		t.Fatal("Expected non-nil Ollama provider")
	}
	if ollamaProvider.Name() != "ollama" {
		t.Fatalf("Expected provider name 'ollama', got '%s'", ollamaProvider.Name())
	}

	// Test unknown provider type
	_, err = registry.Create("unknown", settings)
	if err == nil {
		t.Fatal("Expected error for unknown provider type")
	}

	// Test validation
	invalidSettings := config.ProviderSettings{
		Enabled:  true,
		Endpoint: "", // Invalid: empty endpoint
		APIKey:   "",
	}

	_, err = registry.Create("ollama", invalidSettings)
	if err == nil {
		t.Fatal("Expected validation error for empty endpoint")
	}
}

// TestConfigLoadingWithMapBasedStructure tests config loading with new map-based structure
func TestConfigLoadingWithMapBasedStructure(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `
server:
  port: "8080"
provider:
  providers:
    ollama:
      enabled: true
      endpoint: "http://localhost:11434"
      api_key: ""
    opencode_zen:
      enabled: true
      endpoint: "https://opencode.ai/zen"
      api_key: "test-key"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load config
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify map-based structure
	if cfg.Provider.Providers == nil {
		t.Fatal("Expected providers map to be non-nil")
	}

	// Check both providers are present
	if _, ok := cfg.Provider.Providers["ollama"]; !ok {
		t.Fatal("Expected ollama provider in config")
	}
	if _, ok := cfg.Provider.Providers["opencode_zen"]; !ok {
		t.Fatal("Expected opencode_zen provider in config")
	}

	// Test GetEnabledProviders
	enabled := cfg.GetEnabledProviders()
	if len(enabled) != 2 {
		t.Fatalf("Expected 2 enabled providers, got %d", len(enabled))
	}

	// Test empty providers map
	emptyConfigContent := `
server:
  port: "8080"
provider:
  providers: {}
`

	err = os.WriteFile(configPath, []byte(emptyConfigContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write empty config file: %v", err)
	}

	cfg, err = config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load empty config: %v", err)
	}

	enabled = cfg.GetEnabledProviders()
	if len(enabled) != 0 {
		t.Fatalf("Expected 0 enabled providers, got %d", len(enabled))
	}
}

// TestProvidersEndpoint tests the /v1/providers endpoint
func TestProvidersEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test registry
	registry := provider.NewProviderRegistry()
	registry.Register(provider.NewOllamaFactory())
	registry.Register(provider.NewOpenCodeZenFactory())

	// Create test config
	cfg := &config.Config{
		Provider: config.ProviderConfig{
			Providers: map[string]config.ProviderSettings{
				"ollama": {
					Enabled:  true,
					Endpoint: "http://localhost:11434",
					APIKey:   "",
				},
				"opencode_zen": {
					Enabled:  false,
					Endpoint: "https://opencode.ai/zen",
					APIKey:   "test-key",
				},
			},
		},
	}

	// Create a minimal handler for testing
	router := gin.New()

	// Mock providersHandler
	router.GET("/v1/providers", func(c *gin.Context) {
		types := registry.ListTypes()
		defaults := registry.GetDefaults()

		providers := make([]map[string]interface{}, 0, len(types))
		for _, providerType := range types {
			factory, err := registry.GetFactory(providerType)
			if err != nil {
				continue
			}

			providerInfo := map[string]interface{}{
				"type":        factory.Type(),
				"description": factory.Description(),
			}

			if settings, ok := cfg.Provider.Providers[providerType]; ok {
				providerInfo["configured"] = true
				providerInfo["enabled"] = settings.Enabled
				providerInfo["endpoint"] = settings.Endpoint
				providerInfo["has_api_key"] = settings.APIKey != ""
			} else {
				providerInfo["configured"] = false
				providerInfo["enabled"] = false
			}

			if defaultSettings, ok := defaults[providerType]; ok {
				providerInfo["defaults"] = map[string]interface{}{
					"endpoint": defaultSettings.Endpoint,
					"enabled":   defaultSettings.Enabled,
				}
			}

			providers = append(providers, providerInfo)
		}

		c.JSON(http.StatusOK, gin.H{"providers": providers})
	})

	// Test the endpoint
	req, _ := http.NewRequest("GET", "/v1/providers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	providers, ok := response["providers"].([]interface{})
	if !ok {
		t.Fatal("Expected providers to be an array")
	}

	if len(providers) != 2 {
		t.Fatalf("Expected 2 providers, got %d", len(providers))
	}

	// Verify provider structure
	for _, p := range providers {
		providerMap, ok := p.(map[string]interface{})
		if !ok {
			t.Fatal("Expected provider to be a map")
		}

		if _, ok := providerMap["type"]; !ok {
			t.Fatal("Expected provider to have type field")
		}
		if _, ok := providerMap["description"]; !ok {
			t.Fatal("Expected provider to have description field")
		}
		if _, ok := providerMap["configured"]; !ok {
			t.Fatal("Expected provider to have configured field")
		}
		if _, ok := providerMap["defaults"]; !ok {
			t.Fatal("Expected provider to have defaults field")
		}
	}
}

// TestModelRoutingWithNewProviderNames tests model routing with new provider names
func TestModelRoutingWithNewProviderNames(t *testing.T) {
	registry := provider.NewProviderRegistry()
	registry.Register(provider.NewOllamaFactory())
	registry.Register(provider.NewOpenCodeZenFactory())

	// Test model name extraction for different providers
	testCases := []struct {
		model          string
		expectedPrefix string
		expectedName   string
	}{
		{"ollama:llama2", "ollama", "llama2"},
		{"ollama:mistral", "ollama", "mistral"},
		{"opencode_zen:gpt-4", "opencode_zen", "gpt-4"},
		{"opencode_zen:claude-3", "opencode_zen", "claude-3"},
	}

	for _, tc := range testCases {
		t.Run(tc.model, func(t *testing.T) {
			// Extract prefix (first part before colon)
			prefix := ""
			for i, c := range tc.model {
				if c == ':' {
					prefix = tc.model[:i]
					break
				}
			}

			if prefix != tc.expectedPrefix {
				t.Fatalf("Expected prefix '%s', got '%s'", tc.expectedPrefix, prefix)
			}

			// Verify factory exists for this prefix
			_, err := registry.GetFactory(prefix)
			if err != nil {
				t.Fatalf("Expected factory to exist for prefix '%s': %v", prefix, err)
			}
		})
	}

	// Test that old prefix doesn't work
	_, err := registry.GetFactory("opencode")
	if err == nil {
		t.Fatal("Expected error for old prefix 'opencode'")
	}
}

// TestWithBothProvidersEnabled tests with both providers enabled
func TestWithBothProvidersEnabled(t *testing.T) {
	cfg := &config.Config{
		Provider: config.ProviderConfig{
			Providers: map[string]config.ProviderSettings{
				"ollama": {
					Enabled:  true,
					Endpoint: "http://localhost:11434",
					APIKey:   "",
				},
				"opencode_zen": {
					Enabled:  true,
					Endpoint: "https://opencode.ai/zen",
					APIKey:   "test-key",
				},
			},
		},
	}

	enabled := cfg.GetEnabledProviders()
	if len(enabled) != 2 {
		t.Fatalf("Expected 2 enabled providers, got %d", len(enabled))
	}

	if _, ok := enabled["ollama"]; !ok {
		t.Fatal("Expected ollama to be enabled")
	}
	if _, ok := enabled["opencode_zen"]; !ok {
		t.Fatal("Expected opencode_zen to be enabled")
	}
}

// TestWithSingleProviderEnabled tests with single provider enabled
func TestWithSingleProviderEnabled(t *testing.T) {
	cfg := &config.Config{
		Provider: config.ProviderConfig{
			Providers: map[string]config.ProviderSettings{
				"ollama": {
					Enabled:  true,
					Endpoint: "http://localhost:11434",
					APIKey:   "",
				},
				"opencode_zen": {
					Enabled:  false,
					Endpoint: "https://opencode.ai/zen",
					APIKey:   "test-key",
				},
			},
		},
	}

	enabled := cfg.GetEnabledProviders()
	if len(enabled) != 1 {
		t.Fatalf("Expected 1 enabled provider, got %d", len(enabled))
	}

	if _, ok := enabled["ollama"]; !ok {
		t.Fatal("Expected ollama to be enabled")
	}
	if _, ok := enabled["opencode_zen"]; ok {
		t.Fatal("Expected opencode_zen to be disabled")
	}
}

// TestWithNoProvidersEnabled tests with no providers enabled
func TestWithNoProvidersEnabled(t *testing.T) {
	cfg := &config.Config{
		Provider: config.ProviderConfig{
			Providers: map[string]config.ProviderSettings{
				"ollama": {
					Enabled:  false,
					Endpoint: "http://localhost:11434",
					APIKey:   "",
				},
				"opencode_zen": {
					Enabled:  false,
					Endpoint: "https://opencode.ai/zen",
					APIKey:   "test-key",
				},
			},
		},
	}

	enabled := cfg.GetEnabledProviders()
	if len(enabled) != 0 {
		t.Fatalf("Expected 0 enabled providers, got %d", len(enabled))
	}
}

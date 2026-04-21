package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ai-api-gateway/router-service/internal/infrastructure/config"
	"github.com/ai-api-gateway/router-service/internal/infrastructure/provider"
	"github.com/gin-gonic/gin"
)

func TestProvidersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test registry
	registry := provider.NewProviderRegistry()
	registry.Register(provider.NewOllamaFactory())
	registry.Register(provider.NewOpenCodeZenFactory())

	// Create a test config
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

	// Create handler
	h := &Handler{
		registry: registry,
		config:   cfg,
	}

	// Create test router
	router := gin.New()
	router.GET("/v1/providers", h.providersHandler)

	// Create request
	req, _ := http.NewRequest("GET", "/v1/providers", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Check that providers field exists
	providers, ok := response["providers"].([]interface{})
	if !ok {
		t.Fatal("expected providers field to be an array")
	}

	// Check that we have 2 providers
	if len(providers) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(providers))
	}
}

func TestProvidersHandler_NoConfiguredProviders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test registry
	registry := provider.NewProviderRegistry()
	registry.Register(provider.NewOllamaFactory())
	registry.Register(provider.NewOpenCodeZenFactory())

	// Create a test config with no providers
	cfg := &config.Config{
		Provider: config.ProviderConfig{
			Providers: map[string]config.ProviderSettings{},
		},
	}

	// Create handler
	h := &Handler{
		registry: registry,
		config:   cfg,
	}

	// Create test router
	router := gin.New()
	router.GET("/v1/providers", h.providersHandler)

	// Create request
	req, _ := http.NewRequest("GET", "/v1/providers", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Check that providers field exists
	providers, ok := response["providers"].([]interface{})
	if !ok {
		t.Fatal("expected providers field to be an array")
	}

	// Check that we have 2 providers (both unconfigured)
	if len(providers) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(providers))
	}
}

package config

import (
	"os"
	"testing"
)

func TestConfig_GetEnabledProviders(t *testing.T) {
	cfg := &Config{
		Provider: ProviderConfig{
			Providers: map[string]ProviderSettings{
				"ollama": {
					Enabled:  true,
					Endpoint: "http://localhost:11434",
				},
				"opencode_zen": {
					Enabled:  false,
					Endpoint: "https://opencode.ai/zen",
				},
			},
		},
	}

	enabled := cfg.GetEnabledProviders()
	if len(enabled) != 1 {
		t.Fatalf("expected 1 enabled provider, got %d", len(enabled))
	}
	if _, ok := enabled["ollama"]; !ok {
		t.Fatal("expected ollama to be enabled")
	}
	if _, ok := enabled["opencode_zen"]; ok {
		t.Fatal("expected opencode_zen to be disabled")
	}
}

func TestConfig_GetEnabledProviders_Empty(t *testing.T) {
	cfg := &Config{
		Provider: ProviderConfig{
			Providers: map[string]ProviderSettings{},
		},
	}

	enabled := cfg.GetEnabledProviders()
	if len(enabled) != 0 {
		t.Fatalf("expected 0 enabled providers, got %d", len(enabled))
	}
}

func TestConfig_GetEnabledProviders_AllDisabled(t *testing.T) {
	cfg := &Config{
		Provider: ProviderConfig{
			Providers: map[string]ProviderSettings{
				"ollama": {
					Enabled:  false,
					Endpoint: "http://localhost:11434",
				},
				"opencode_zen": {
					Enabled:  false,
					Endpoint: "https://opencode.ai/zen",
				},
			},
		},
	}

	enabled := cfg.GetEnabledProviders()
	if len(enabled) != 0 {
		t.Fatalf("expected 0 enabled providers, got %d", len(enabled))
	}
}

func TestConfig_ResolveEnvVars(t *testing.T) {
	os.Setenv("TEST_ENDPOINT", "http://test.example.com")
	os.Setenv("TEST_API_KEY", "test-key-123")
	defer os.Unsetenv("TEST_ENDPOINT")
	defer os.Unsetenv("TEST_API_KEY")

	cfg := &Config{
		Provider: ProviderConfig{
			Providers: map[string]ProviderSettings{
				"ollama": {
					Enabled:  true,
					Endpoint: "${TEST_ENDPOINT}",
					APIKey:   "${TEST_API_KEY}",
				},
			},
		},
	}

	err := resolveEnvVars(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Provider.Providers["ollama"].Endpoint != "http://test.example.com" {
		t.Fatalf("expected endpoint 'http://test.example.com', got '%s'", cfg.Provider.Providers["ollama"].Endpoint)
	}
	if cfg.Provider.Providers["ollama"].APIKey != "test-key-123" {
		t.Fatalf("expected api_key 'test-key-123', got '%s'", cfg.Provider.Providers["ollama"].APIKey)
	}
}

func TestConfig_ResolveEnvVars_MissingEnv(t *testing.T) {
	cfg := &Config{
		Provider: ProviderConfig{
			Providers: map[string]ProviderSettings{
				"ollama": {
					Enabled:  true,
					Endpoint: "${MISSING_ENDPOINT}",
					APIKey:   "${MISSING_API_KEY}",
				},
			},
		},
	}

	err := resolveEnvVars(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Missing env vars should be replaced with empty string
	if cfg.Provider.Providers["ollama"].Endpoint != "" {
		t.Fatalf("expected empty endpoint, got '%s'", cfg.Provider.Providers["ollama"].Endpoint)
	}
	if cfg.Provider.Providers["ollama"].APIKey != "" {
		t.Fatalf("expected empty api_key, got '%s'", cfg.Provider.Providers["ollama"].APIKey)
	}
}

func TestProviderSettings_BaseURL(t *testing.T) {
	settings := ProviderSettings{Endpoint: "http://localhost:11434/"}
	if settings.BaseURL() != "http://localhost:11434" {
		t.Fatalf("expected 'http://localhost:11434', got '%s'", settings.BaseURL())
	}

	settings = ProviderSettings{Endpoint: "http://localhost:11434"}
	if settings.BaseURL() != "http://localhost:11434" {
		t.Fatalf("expected 'http://localhost:11434', got '%s'", settings.BaseURL())
	}
}

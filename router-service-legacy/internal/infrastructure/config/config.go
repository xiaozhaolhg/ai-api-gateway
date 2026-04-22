package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Provider ProviderConfig `yaml:"provider"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type ProviderConfig struct {
	Providers map[string]ProviderSettings `yaml:"providers"`
}

type ProviderSettings struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
	APIKey   string `yaml:"api_key"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := resolveEnvVars(&cfg); err != nil {
		return nil, fmt.Errorf("failed to resolve env vars: %w", err)
	}

	return &cfg, nil
}

var envVarRegex = regexp.MustCompile(`\$\{([^}]+)\}`)

func resolveEnvVars(cfg *Config) error {
	resolve := func(s string) string {
		return envVarRegex.ReplaceAllStringFunc(s, func(match string) string {
			varName := match[2 : len(match)-1]
			if val := os.Getenv(varName); val != "" {
				return val
			}
			return ""
		})
	}

	for providerType, settings := range cfg.Provider.Providers {
		if settings.Endpoint != "" {
			settings.Endpoint = resolve(settings.Endpoint)
		}
		if settings.APIKey != "" {
			settings.APIKey = resolve(settings.APIKey)
		}
		cfg.Provider.Providers[providerType] = settings
	}

	return nil
}

func (c *Config) GetEnabledProviders() map[string]ProviderSettings {
	enabled := make(map[string]ProviderSettings)
	for providerType, settings := range c.Provider.Providers {
		if settings.Enabled {
			enabled[providerType] = settings
		}
	}
	return enabled
}

func (s ProviderSettings) BaseURL() string {
	return strings.TrimRight(s.Endpoint, "/")
}
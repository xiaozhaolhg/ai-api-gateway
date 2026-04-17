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
	Ollama      ProviderSettings `yaml:"ollama"`
	OpenCodeZen ProviderSettings `yaml:"opencode_zen"`
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

	if cfg.Provider.Ollama.Endpoint != "" {
		cfg.Provider.Ollama.Endpoint = resolve(cfg.Provider.Ollama.Endpoint)
	}
	if cfg.Provider.Ollama.APIKey != "" {
		cfg.Provider.Ollama.APIKey = resolve(cfg.Provider.Ollama.APIKey)
	}
	if cfg.Provider.OpenCodeZen.Endpoint != "" {
		cfg.Provider.OpenCodeZen.Endpoint = resolve(cfg.Provider.OpenCodeZen.Endpoint)
	}
	if cfg.Provider.OpenCodeZen.APIKey != "" {
		cfg.Provider.OpenCodeZen.APIKey = resolve(cfg.Provider.OpenCodeZen.APIKey)
	}

	return nil
}

func (c *Config) GetEnabledProviders() []ProviderSettings {
	var enabled []ProviderSettings
	if c.Provider.Ollama.Enabled {
		enabled = append(enabled, c.Provider.Ollama)
	}
	if c.Provider.OpenCodeZen.Enabled {
		enabled = append(enabled, c.Provider.OpenCodeZen)
	}
	return enabled
}

func (s ProviderSettings) BaseURL() string {
	return strings.TrimRight(s.Endpoint, "/")
}
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the service configuration
type Config struct {
	Server          ServerConfig          `yaml:"server"`
	AuthService     ServiceConfig         `yaml:"auth_service"`
	RouterService   ServiceConfig         `yaml:"router_service"`
	ProviderService ServiceConfig         `yaml:"provider_service"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

// ServiceConfig holds gRPC service configuration
type ServiceConfig struct {
	Address string `yaml:"address"`
}

// Load loads configuration from a YAML file and environment variables
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}
	if host := os.Getenv("HOST"); host != "" {
		cfg.Server.Host = host
	}

	// Set defaults
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}

	return &cfg, nil
}

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the service configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Crypto   CryptoConfig   `yaml:"crypto"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// CryptoConfig holds encryption configuration
type CryptoConfig struct {
	EncryptionKey string `yaml:"encryption_key"`
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
	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		cfg.Database.Path = dbPath
	}
	if key := os.Getenv("ENCRYPTION_KEY"); key != "" {
		cfg.Crypto.EncryptionKey = key
	}

	// Set defaults
	if cfg.Server.Port == "" {
		cfg.Server.Port = "50053"
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Database.Path == "" {
		cfg.Database.Path = "./data/provider.db"
	}

	return &cfg, nil
}

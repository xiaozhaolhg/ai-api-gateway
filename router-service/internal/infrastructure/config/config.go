package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the service configuration
type Config struct {
	Server        ServerConfig   `yaml:"server"`
	Database      DatabaseConfig `yaml:"database"`
	Cache         CacheConfig    `yaml:"cache"`
	ProviderService ServerConfig   `yaml:"provider_service"`
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

// CacheConfig holds cache configuration
type CacheConfig struct {
	Redis RedisConfig `yaml:"redis"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Address     string `yaml:"address"`
	Password    string `yaml:"password"`
	DB          int    `yaml:"db"`
	TTLSeconds  int    `yaml:"ttl_seconds"`
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
	if redisAddr := os.Getenv("REDIS_ADDRESS"); redisAddr != "" {
		cfg.Cache.Redis.Address = redisAddr
	}
	if redisPass := os.Getenv("REDIS_PASSWORD"); redisPass != "" {
		cfg.Cache.Redis.Password = redisPass
	}
	if redisDB := os.Getenv("REDIS_DB"); redisDB != "" {
		// Parse as int
		var db int
		if _, err := fmt.Sscanf(redisDB, "%d", &db); err == nil {
			cfg.Cache.Redis.DB = db
		}
	}
	if redisTTL := os.Getenv("REDIS_TTL_SECONDS"); redisTTL != "" {
		var ttl int
		if _, err := fmt.Sscanf(redisTTL, "%d", &ttl); err == nil {
			cfg.Cache.Redis.TTLSeconds = ttl
		}
	}

	// Set defaults
	if cfg.Server.Port == "" {
		cfg.Server.Port = "50052"
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Database.Path == "" {
		cfg.Database.Path = "./data/router.db"
	}
	if cfg.Cache.Redis.Address == "" {
		cfg.Cache.Redis.Address = "localhost:6379"
	}
	if cfg.Cache.Redis.TTLSeconds == 0 {
		cfg.Cache.Redis.TTLSeconds = 300
	}

	return &cfg, nil
}

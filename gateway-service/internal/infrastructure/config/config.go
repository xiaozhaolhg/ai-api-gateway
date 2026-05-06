package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the service configuration
type Config struct {
	Server          ServerConfig  `yaml:"server"`
	AuthService     ServiceConfig `yaml:"auth_service"`
	RouterService   ServiceConfig `yaml:"router_service"`
	ProviderService ServiceConfig `yaml:"provider_service"`
	BillingService  ServiceConfig `yaml:"billing_service"`
	MonitorService  ServiceConfig `yaml:"monitor_service"`
	Log             LogConfig     `yaml:"log"`
	Timeout         TimeoutConfig `yaml:"timeout"`
	Cache           CacheConfig   `yaml:"cache"`
	GRPC            GRPCConfig    `yaml:"grpc"`
	SSE             SSEConfig     `yaml:"sse"`
	StreamingTokenInterval *int64  `yaml:"streaming_token_interval"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port        string `yaml:"port"`
	Host        string `yaml:"host"`
	MaxBodySize string `yaml:"max_body_size"`
}

// ServiceConfig holds gRPC service configuration
type ServiceConfig struct {
	Address string `yaml:"address"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level          string `yaml:"level"`
	Format         string `yaml:"format"`
	RequestBody    bool   `yaml:"request_body"`
	ResponseBody   bool   `yaml:"response_body"`
	MaskSensitive  bool   `yaml:"mask_sensitive"`
}

// TimeoutConfig holds timeout configuration
type TimeoutConfig struct {
	Auth        string `yaml:"auth"`
	Router      string `yaml:"router"`
	Provider    string `yaml:"provider"`
	Billing     string `yaml:"billing"`
	HealthCheck string `yaml:"health_check"`
}

// CacheConfig holds caching configuration
type CacheConfig struct {
	ModelsTTL string `yaml:"models_ttl"`
}

// GRPCConfig holds gRPC connection configuration
type GRPCConfig struct {
	MaxRetries     int    `yaml:"max_retries"`
	RetryBackoff   string `yaml:"retry_backoff"`
	ConnectTimeout string `yaml:"connect_timeout"`
}

// SSEConfig holds SSE streaming configuration
type SSEConfig struct {
	HeartbeatInterval string `yaml:"heartbeat_interval"`
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
	if addr := os.Getenv("AUTH_SERVICE_ADDRESS"); addr != "" {
		cfg.AuthService.Address = addr
	}
	if addr := os.Getenv("ROUTER_SERVICE_ADDRESS"); addr != "" {
		cfg.RouterService.Address = addr
	}
	if addr := os.Getenv("PROVIDER_SERVICE_ADDRESS"); addr != "" {
		cfg.ProviderService.Address = addr
	}
	if addr := os.Getenv("BILLING_SERVICE_ADDRESS"); addr != "" {
		cfg.BillingService.Address = addr
	}
	if addr := os.Getenv("MONITOR_SERVICE_ADDRESS"); addr != "" {
		cfg.MonitorService.Address = addr
	}
	if envInterval := os.Getenv("STREAMING_TOKEN_INTERVAL"); envInterval != "" {
		val, err := strconv.ParseInt(envInterval, 10, 64)
		if err == nil {
			cfg.StreamingTokenInterval = &val
		}
	}

	// Set defaults
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.MaxBodySize == "" {
		cfg.Server.MaxBodySize = "10MB"
	}
	if cfg.AuthService.Address == "" {
		cfg.AuthService.Address = "localhost:50051"
	}
	if cfg.RouterService.Address == "" {
		cfg.RouterService.Address = "localhost:50052"
	}
	if cfg.ProviderService.Address == "" {
		cfg.ProviderService.Address = "localhost:50053"
	}
	if cfg.BillingService.Address == "" {
		cfg.BillingService.Address = "localhost:50054"
	}
	if cfg.MonitorService.Address == "" {
		cfg.MonitorService.Address = "localhost:50055"
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "json"
	}
	if cfg.Timeout.Auth == "" {
		cfg.Timeout.Auth = "5s"
	}
	if cfg.Timeout.Router == "" {
		cfg.Timeout.Router = "5s"
	}
	if cfg.Timeout.Provider == "" {
		cfg.Timeout.Provider = "30s"
	}
	if cfg.Timeout.Billing == "" {
		cfg.Timeout.Billing = "10s"
	}
	if cfg.Timeout.HealthCheck == "" {
		cfg.Timeout.HealthCheck = "3s"
	}
	if cfg.Cache.ModelsTTL == "" {
		cfg.Cache.ModelsTTL = "5m"
	}
	if cfg.GRPC.MaxRetries == 0 {
		cfg.GRPC.MaxRetries = 3
	}
	if cfg.GRPC.RetryBackoff == "" {
		cfg.GRPC.RetryBackoff = "100ms"
	}
	if cfg.GRPC.ConnectTimeout == "" {
		cfg.GRPC.ConnectTimeout = "5s"
	}
	if cfg.SSE.HeartbeatInterval == "" {
		cfg.SSE.HeartbeatInterval = "15s"
	}
	if cfg.StreamingTokenInterval == nil {
		defaultInterval := int64(1000)
		cfg.StreamingTokenInterval = &defaultInterval
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server port
	if c.Server.Port != "" {
		port, err := strconv.Atoi(c.Server.Port)
		if err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("invalid server port: %s", c.Server.Port)
		}
	}

	// Validate log level
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[strings.ToLower(c.Log.Level)] {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", c.Log.Level)
	}

	// Validate log format
	validLogFormats := map[string]bool{"json": true, "text": true}
	if !validLogFormats[strings.ToLower(c.Log.Format)] {
		return fmt.Errorf("invalid log format: %s (must be json or text)", c.Log.Format)
	}

	// Validate timeout formats
	if err := validateDuration(c.Timeout.Auth, "auth timeout"); err != nil {
		return err
	}
	if err := validateDuration(c.Timeout.Router, "router timeout"); err != nil {
		return err
	}
	if err := validateDuration(c.Timeout.Provider, "provider timeout"); err != nil {
		return err
	}
	if err := validateDuration(c.Timeout.Billing, "billing timeout"); err != nil {
		return err
	}
	if err := validateDuration(c.Timeout.HealthCheck, "health check timeout"); err != nil {
		return err
	}

	// Validate cache TTL
	if err := validateDuration(c.Cache.ModelsTTL, "models cache TTL"); err != nil {
		return err
	}

	// Validate gRPC retry settings
	if c.GRPC.MaxRetries < 0 {
		return fmt.Errorf("grpc max_retries cannot be negative: %d", c.GRPC.MaxRetries)
	}
	if err := validateDuration(c.GRPC.RetryBackoff, "grpc retry_backoff"); err != nil {
		return err
	}
	if err := validateDuration(c.GRPC.ConnectTimeout, "grpc connect_timeout"); err != nil {
		return err
	}

	// Validate SSE heartbeat interval
	if err := validateDuration(c.SSE.HeartbeatInterval, "sse heartbeat_interval"); err != nil {
		return err
	}

	// Validate service addresses are not empty
	if c.AuthService.Address == "" {
		return fmt.Errorf("auth_service address is required")
	}
	if c.RouterService.Address == "" {
		return fmt.Errorf("router_service address is required")
	}
	if c.ProviderService.Address == "" {
		return fmt.Errorf("provider_service address is required")
	}

	return nil
}

// validateDuration validates a duration string
func validateDuration(value, name string) error {
	if value == "" {
		return nil // Empty is allowed (will use defaults)
	}
	_, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("invalid %s: %s", name, value)
	}
	return nil
}

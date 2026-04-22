package entity

import "time"

// ProviderHealthStatus represents the health status of a provider
type ProviderHealthStatus struct {
	ProviderID     string    `json:"provider_id"`
	Status         string    `json:"status"` // "healthy" | "degraded" | "unhealthy"
	LatencyMs      int64     `json:"latency_ms"`
	ErrorRate      float64   `json:"error_rate"`
	LastCheckTime  time.Time `json:"last_check_time"`
	UptimeSeconds  int64     `json:"uptime_seconds"`
}

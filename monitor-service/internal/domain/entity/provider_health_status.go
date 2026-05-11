package entity

// ProviderHealthStatus represents the health status of a provider
type ProviderHealthStatus struct {
	ProviderID  string  `json:"provider_id"`
	Status      string  `json:"status"` // "healthy" | "degraded" | "unhealthy"
	LatencyP50  float64 `json:"latency_p50"`
	LatencyP95  float64 `json:"latency_p95"`
	LatencyP99  float64 `json:"latency_p99"`
	ErrorRate   float64 `json:"error_rate"`
	UptimePct   float64 `json:"uptime_pct"`
	LastCheck   int64   `json:"last_check"`
}

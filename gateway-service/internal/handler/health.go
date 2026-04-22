package handler

import (
	"encoding/json"
	"net/http"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// ServeHTTP handles HTTP requests for /health and /gateway/health
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path

	if path == "/gateway/health" {
		h.gatewayHealth(w, r)
	} else {
		h.health(w, r)
	}
}

// health returns a simple health check
func (h *HealthHandler) health(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// gatewayHealth returns detailed gateway health including dependency status
func (h *HealthHandler) gatewayHealth(w http.ResponseWriter, r *http.Request) {
	// In production, this would check the health of all dependent services
	// For MVP, we'll return a placeholder response

	response := map[string]interface{}{
		"status": "ok",
		"services": map[string]string{
			"auth":     "ok",
			"router":   "ok",
			"provider": "ok",
			"billing":  "ok",
			"monitor":  "ok",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

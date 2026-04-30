package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	authClient     *client.AuthClient
	routerClient   *client.RouterClient
	providerClient *client.ProviderClient
	billingClient  *client.BillingClient
}

// NewHealthHandler creates a new health handler with service clients
func NewHealthHandler(authClient *client.AuthClient, routerClient *client.RouterClient, 
	providerClient *client.ProviderClient, billingClient *client.BillingClient) *HealthHandler {
	return &HealthHandler{
		authClient:     authClient,
		routerClient:   routerClient,
		providerClient: providerClient,
		billingClient:  billingClient,
	}
}

// ServeHTTP handles HTTP requests for /health and /gateway/health (standard HTTP handler)
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path

	if path == "/gateway/health" {
		h.gatewayHealthHTTP(w, r)
	} else {
		h.health(w, r)
	}
}

// GatewayHealth handles Gin requests for /gateway/health (deep health check)
func (h *HealthHandler) GatewayHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	services := make(map[string]interface{})
	anyDegraded := false
	anyUnhealthy := false

	// Check auth service
	authStatus, authLatency := h.checkService(ctx, "auth", func(ctx context.Context) error {
		_, err := h.authClient.ListUsers(ctx, 1, 1)
		return err
	})
	services["auth"] = map[string]interface{}{
		"status":  authStatus,
		"latency": authLatency.String(),
	}
	if authStatus == "degraded" {
		anyDegraded = true
	}
	if authStatus == "unhealthy" {
		anyUnhealthy = true
	}

	// Check router service
	routerStatus, routerLatency := h.checkService(ctx, "router", func(ctx context.Context) error {
		_, err := h.routerClient.ResolveRoute(ctx, "test-model", []string{"test-model"})
		// Ignore "model not found" errors - we're just checking connectivity
		return err
	})
	services["router"] = map[string]interface{}{
		"status":  routerStatus,
		"latency": routerLatency.String(),
	}
	if routerStatus == "degraded" {
		anyDegraded = true
	}
	if routerStatus == "unhealthy" {
		anyUnhealthy = true
	}

	// Check provider service
	providerStatus, providerLatency := h.checkService(ctx, "provider", func(ctx context.Context) error {
		_, err := h.providerClient.ListProviders(ctx, 1, 1)
		return err
	})
	services["provider"] = map[string]interface{}{
		"status":  providerStatus,
		"latency": providerLatency.String(),
	}
	if providerStatus == "degraded" {
		anyDegraded = true
	}
	if providerStatus == "unhealthy" {
		anyUnhealthy = true
	}

	// Check billing service (graceful - billing is optional)
	billingStatus, billingLatency := h.checkService(ctx, "billing", func(ctx context.Context) error {
		_, err := h.billingClient.GetUsage(ctx, "", 1, 1)
		return err
	})
	services["billing"] = map[string]interface{}{
		"status":   billingStatus,
		"latency":  billingLatency.String(),
		"optional": true,
	}
	// Billing being down doesn't affect overall health (it's optional)

	// Determine overall status: healthy (all ok), degraded (some slow), unhealthy (some down)
	overallStatus := "healthy"
	httpStatus := http.StatusOK
	if anyUnhealthy {
		overallStatus = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	} else if anyDegraded {
		overallStatus = "degraded"
	}

	response := map[string]interface{}{
		"status":    overallStatus,
		"gateway":   "healthy",
		"services":  services,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(httpStatus, response)
}

// health returns a simple health check
func (h *HealthHandler) health(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// gatewayHealthHTTP returns detailed gateway health including dependency status (HTTP version)
func (h *HealthHandler) gatewayHealthHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	services := make(map[string]interface{})
	anyDegraded := false
	anyUnhealthy := false

	// Check auth service
	authStatus, authLatency := h.checkService(ctx, "auth", func(ctx context.Context) error {
		_, err := h.authClient.ListUsers(ctx, 1, 1)
		return err
	})
	services["auth"] = map[string]interface{}{
		"status":  authStatus,
		"latency": authLatency.String(),
	}
	if authStatus == "degraded" {
		anyDegraded = true
	}
	if authStatus == "unhealthy" {
		anyUnhealthy = true
	}

	// Check router service
	routerStatus, routerLatency := h.checkService(ctx, "router", func(ctx context.Context) error {
		_, err := h.routerClient.ResolveRoute(ctx, "test-model", []string{"test-model"})
		// Ignore "model not found" errors - we're just checking connectivity
		return err
	})
	services["router"] = map[string]interface{}{
		"status":  routerStatus,
		"latency": routerLatency.String(),
	}
	if routerStatus == "degraded" {
		anyDegraded = true
	}
	if routerStatus == "unhealthy" {
		anyUnhealthy = true
	}

	// Check provider service
	providerStatus, providerLatency := h.checkService(ctx, "provider", func(ctx context.Context) error {
		_, err := h.providerClient.ListProviders(ctx, 1, 1)
		return err
	})
	services["provider"] = map[string]interface{}{
		"status":  providerStatus,
		"latency": providerLatency.String(),
	}
	if providerStatus == "degraded" {
		anyDegraded = true
	}
	if providerStatus == "unhealthy" {
		anyUnhealthy = true
	}

	// Check billing service (graceful - billing is optional)
	billingStatus, billingLatency := h.checkService(ctx, "billing", func(ctx context.Context) error {
		_, err := h.billingClient.GetUsage(ctx, "", 1, 1)
		return err
	})
	services["billing"] = map[string]interface{}{
		"status":   billingStatus,
		"latency":  billingLatency.String(),
		"optional": true,
	}
	// Billing being down doesn't affect overall health (it's optional)

	// Determine overall status: healthy (all ok), degraded (some slow), unhealthy (some down)
	overallStatus := "healthy"
	httpStatus := http.StatusOK
	if anyUnhealthy {
		overallStatus = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	} else if anyDegraded {
		overallStatus = "degraded"
	}

	response := map[string]interface{}{
		"status":    overallStatus,
		"gateway":   "healthy",
		"services":  services,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(response)
}

// checkService performs a health check on a service and returns status and latency
// Status values: "healthy" (ok), "degraded" (slow >500ms), "unhealthy" (error)
func (h *HealthHandler) checkService(ctx context.Context, name string, check func(context.Context) error) (string, time.Duration) {
	start := time.Now()

	// Create a timeout context for this specific check
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := check(checkCtx)
	latency := time.Since(start)

	if err != nil {
		return "unhealthy", latency
	}

	// Degraded if latency > 500ms
	if latency > 500*time.Millisecond {
		return "degraded", latency
	}

	return "healthy", latency
}

// statusToBool returns true if status is healthy or degraded (not down)
func statusToBool(status string) bool {
	return status == "healthy" || status == "degraded"
}

// boolToStatus kept for backward compatibility
func boolToStatus(healthy bool) string {
	if healthy {
		return "healthy"
	}
	return "unhealthy"
}

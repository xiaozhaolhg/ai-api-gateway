package middleware

import (
	"context"
	"io"
	"net/http"

	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// ProxyMiddleware forwards requests to provider-service
type ProxyMiddleware struct {
	providerClient *client.ProviderClient
	billingClient  *client.BillingClient
	monitorClient  *client.MonitorClient
}

// NewProxyMiddleware creates a new proxy middleware
func NewProxyMiddleware(
	providerClient *client.ProviderClient,
	billingClient *client.BillingClient,
	monitorClient *client.MonitorClient,
) *ProxyMiddleware {
	return &ProxyMiddleware{
		providerClient: providerClient,
		billingClient:  billingClient,
		monitorClient:  monitorClient,
	}
}

// Middleware returns the middleware function
func (m *ProxyMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get route info from context (set by RouteMiddleware)
		providerID, ok := r.Context().Value("providerId").(string)
		if !ok {
			http.Error(w, "Provider not resolved", http.StatusInternalServerError)
			return
		}

		// Read request body
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Collect headers
		headers := make(map[string]string)
		for k, v := range r.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}

		// Check if this is a streaming request
		stream := r.URL.Query().Get("stream") == "true"

		if stream {
			m.handleStreamingRequest(w, r, providerID, requestBody, headers)
		} else {
			m.handleNonStreamingRequest(w, r, providerID, requestBody, headers)
		}
	})
}

// handleNonStreamingRequest handles non-streaming requests
func (m *ProxyMiddleware) handleNonStreamingRequest(w http.ResponseWriter, r *http.Request, providerID string, requestBody []byte, headers map[string]string) {
	// Forward request to provider-service
	resp, err := m.providerClient.ForwardRequest(r.Context(), providerID, requestBody, headers)
	if err != nil {
		http.Error(w, "Failed to forward request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ResponseBody)

	// Record usage asynchronously
	go m.recordUsage(r.Context(), providerID, resp.TokenCounts.PromptTokens, resp.TokenCounts.CompletionTokens)
}

// handleStreamingRequest handles streaming requests
func (m *ProxyMiddleware) handleStreamingRequest(w http.ResponseWriter, r *http.Request, providerID string, requestBody []byte, headers map[string]string) {
	// Forward streaming request to provider-service
	stream, err := m.providerClient.StreamRequest(r.Context(), providerID, requestBody, headers)
	if err != nil {
		http.Error(w, "Failed to stream request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Stream chunks
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	for {
		chunk, tokenCounts, done, err := client.ReadStreamChunk(stream)
		if err != nil {
			break
		}

		if done {
			break
		}

		// Write chunk
		w.Write(chunk)
		flusher.Flush()

		// Update token counts
		// In production, you'd accumulate these and record at the end
	}

	// Record usage asynchronously
	// For streaming, token counts are accumulated during the stream
	// This is a simplified version
}

// recordUsage records usage to billing-service asynchronously
func (m *ProxyMiddleware) recordUsage(ctx context.Context, providerID string, promptTokens, completionTokens int64) {
	// Get user ID from context (set by AuthMiddleware)
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return
	}

	// Get model from context (set by RouteMiddleware)
	model, _ := ctx.Value("model").(string)

	// Record usage
	err := m.billingClient.RecordUsage(context.Background(), userID, providerID, model, promptTokens, completionTokens)
	if err != nil {
		// Log error but don't fail the request
		// In production, you'd use proper logging
	}
}

package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ai-api-gateway/gateway-service/internal/client"
)

// ChatCompletionRequest represents the standard OpenAI-style chat completion request
type ChatCompletionRequest struct {
	Model       string                 `json:"model"`
	Messages    []Message              `json:"messages"`
	Stream      bool                   `json:"stream,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Other       map[string]interface{} `json:"-"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ProxyMiddleware forwards requests to provider-service
type ProxyMiddleware struct {
	providerClient *client.ProviderClient
	billingClient  *client.BillingClient
}

// NewProxyMiddleware creates a new proxy middleware
func NewProxyMiddleware(providerClient *client.ProviderClient, billingClient *client.BillingClient) *ProxyMiddleware {
	return &ProxyMiddleware{
		providerClient: providerClient,
		billingClient:  billingClient,
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

// handleNonStreamingRequest handles non-streaming requests with fallback support
func (m *ProxyMiddleware) handleNonStreamingRequest(w http.ResponseWriter, r *http.Request, providerID string, requestBody []byte, headers map[string]string) {
	// Get fallback info from context
	fallbackProviderIDs, _ := r.Context().Value("fallbackProviderIds").([]string)
	fallbackModels, _ := r.Context().Value("fallbackModels").([]string)
	log.Printf("[Fallback] Fallback info: primary=%s, fallbacks=%v, models=%v", providerID, fallbackProviderIDs, fallbackModels)

	// Try primary provider first
	resp, statusCode, err := m.tryNonStreamingProvider(r.Context(), providerID, requestBody, headers)
	if err == nil {
		m.writeNonStreamingResponse(w, r, providerID, resp, requestBody)
		return
	}

	// Check if error is non-retryable (4xx client errors)
	if statusCode >= 400 && statusCode < 500 {
		log.Printf("[Fallback] Primary provider %s returned non-retryable error %d: %v, not attempting fallback", providerID, statusCode, err)
		http.Error(w, fmt.Sprintf("Provider error: %v", err), http.StatusBadGateway)
		return
	}

	// Log fallback
	log.Printf("[Fallback] Primary provider %s failed: %v, attempting fallback", providerID, err)

	// Try fallback providers in order
	for i, fallbackID := range fallbackProviderIDs {
		fallbackModel := ""
		if i < len(fallbackModels) {
			fallbackModel = fallbackModels[i]
		}

		var bodyToSend []byte
		if fallbackModel != "" {
			// Rewrite model in request body
			modifiedBody, err := rewriteModelInRequest(requestBody, fallbackModel)
			if err != nil {
				log.Printf("[Fallback] Failed to rewrite model for %s: %v", fallbackID, err)
				continue
			}
			bodyToSend = modifiedBody
		} else {
			bodyToSend = requestBody
		}

		resp, _, err := m.tryNonStreamingProvider(r.Context(), fallbackID, bodyToSend, headers)
		if err != nil {
			log.Printf("[Fallback] Fallback provider %s failed: %v", fallbackID, err)
			continue
		}

		log.Printf("[Fallback] Successfully fell back to provider %s", fallbackID)
		m.writeNonStreamingResponse(w, r, fallbackID, resp, bodyToSend)
		return
	}

	// All providers failed - return structured error
	m.writeFallbackError(w, "all_providers_failed", "All providers failed")
}

// tryNonStreamingProvider attempts a single non-streaming request to a provider
// Returns response, HTTP status code from provider, and error
func (m *ProxyMiddleware) tryNonStreamingProvider(ctx context.Context, providerID string, requestBody []byte, headers map[string]string) (*client.ForwardRequestResponse, int32, error) {
	resp, err := m.providerClient.ForwardRequest(ctx, providerID, requestBody, headers)
	if err != nil {
		return nil, 0, err
	}
	return resp, resp.StatusCode, nil
}

// writeFallbackError writes a structured error response for fallback failures (W1 fix)
func (m *ProxyMiddleware) writeFallbackError(w http.ResponseWriter, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadGateway)
	errorJSON := fmt.Sprintf(`{"error": {"code": "%s", "message": "%s"}}`, code, message)
	w.Write([]byte(errorJSON))
}

// writeNonStreamingResponse writes the response and records usage
func (m *ProxyMiddleware) writeNonStreamingResponse(w http.ResponseWriter, r *http.Request, providerID string, resp *client.ForwardRequestResponse, requestBody []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp.ResponseBody)

	w.Header().Set("X-Prompt-Tokens", fmt.Sprintf("%d", resp.TokenCounts.PromptTokens))
	w.Header().Set("X-Completion-Tokens", fmt.Sprintf("%d", resp.TokenCounts.CompletionTokens))
	w.Header().Set("X-Total-Tokens", fmt.Sprintf("%d", resp.TokenCounts.TotalTokens))

	req, err := parseChatCompletionRequest(requestBody)
	model := "unknown"
	if err == nil && req.Model != "" {
		model = req.Model
	}
	go m.recordUsage(r.Context(), providerID, model, resp.TokenCounts.PromptTokens, resp.TokenCounts.CompletionTokens)
}

// parseChatCompletionRequest parses the request body into a ChatCompletionRequest
func parseChatCompletionRequest(requestBody []byte) (*ChatCompletionRequest, error) {
	var req ChatCompletionRequest
	if err := json.Unmarshal(requestBody, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// handleStreamingRequest handles streaming requests with fallback support
func (m *ProxyMiddleware) handleStreamingRequest(w http.ResponseWriter, r *http.Request, providerID string, requestBody []byte, headers map[string]string) {
	// Get fallback info from context
	fallbackProviderIDs, _ := r.Context().Value("fallbackProviderIds").([]string)
	fallbackModels, _ := r.Context().Value("fallbackModels").([]string)

	// Try primary provider first
	err := m.tryStreamingProvider(w, r, providerID, requestBody, headers, true)
	if err == nil {
		return
	}

	// Log fallback
	log.Printf("[Fallback] Primary streaming provider %s failed: %v, attempting fallback", providerID, err)

	// Try fallback providers in order
	for i, fallbackID := range fallbackProviderIDs {
		fallbackModel := ""
		if i < len(fallbackModels) {
			fallbackModel = fallbackModels[i]
		}

		var bodyToSend []byte
		if fallbackModel != "" {
			modifiedBody, err := rewriteModelInRequest(requestBody, fallbackModel)
			if err != nil {
				log.Printf("[Fallback] Failed to rewrite model for streaming %s: %v", fallbackID, err)
				continue
			}
			bodyToSend = modifiedBody
		} else {
			bodyToSend = requestBody
		}

		err := m.tryStreamingProvider(w, r, fallbackID, bodyToSend, headers, false)
		if err != nil {
			log.Printf("[Fallback] Fallback streaming provider %s failed: %v", fallbackID, err)
			continue
		}

		log.Printf("[Fallback] Successfully fell back to streaming provider %s", fallbackID)
		return
	}

	// All providers failed - return structured error (W1 fix)
	m.writeFallbackError(w, "all_providers_failed", "All streaming providers failed")
}

// tryStreamingProvider attempts a single streaming request to a provider
// isPrimary indicates if this is the first attempt (sets SSE headers)
func (m *ProxyMiddleware) tryStreamingProvider(w http.ResponseWriter, r *http.Request, providerID string, requestBody []byte, headers map[string]string, isPrimary bool) error {
	stream, err := m.providerClient.StreamRequest(r.Context(), providerID, requestBody, headers)
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	// Set SSE headers only for primary attempt
	if isPrimary {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("X-Accel-Buffering", "no")
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming not supported")
	}

	var totalPromptTokens, totalCompletionTokens int64
	heartbeatTicker := time.NewTicker(15 * time.Second)
	defer heartbeatTicker.Stop()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-heartbeatTicker.C:
				fmt.Fprintf(w, ": ping\n\n")
				flusher.Flush()
			case <-ctx.Done():
				return
			case <-done:
				return
			}
		}
	}()

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			close(done)
			return err
		}

		if len(chunk.ChunkData) > 0 {
			fmt.Fprintf(w, "data: %s\n\n", string(chunk.ChunkData))
			flusher.Flush()
		}

		if chunk.AccumulatedTokens != nil {
			totalPromptTokens += chunk.AccumulatedTokens.PromptTokens
			totalCompletionTokens += chunk.AccumulatedTokens.CompletionTokens
		}

		if chunk.Done {
			break
		}
	}

	close(done)

	finalChunk := map[string]interface{}{
		"prompt_tokens":     totalPromptTokens,
		"completion_tokens": totalCompletionTokens,
		"total_tokens":      totalPromptTokens + totalCompletionTokens,
		"done":              true,
	}
	finalJSON, _ := json.Marshal(finalChunk)
	fmt.Fprintf(w, "data: %s\n\n", string(finalJSON))
	flusher.Flush()

	req, _ := parseChatCompletionRequest(requestBody)
	model := "unknown"
	if req != nil && req.Model != "" {
		model = req.Model
	}
	go m.recordUsage(r.Context(), providerID, model, totalPromptTokens, totalCompletionTokens)
	return nil
}

// rewriteModelInRequest rewrites the model field in a chat completion request body
func rewriteModelInRequest(requestBody []byte, newModel string) ([]byte, error) {
	var req ChatCompletionRequest
	if err := json.Unmarshal(requestBody, &req); err != nil {
		return nil, fmt.Errorf("failed to parse request: %w", err)
	}
	req.Model = newModel
	modified, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	return modified, nil
}

// recordUsage records usage to billing-service asynchronously
func (m *ProxyMiddleware) recordUsage(ctx context.Context, providerID, model string, promptTokens, completionTokens int64) {
	// Get user ID from context (set by AuthMiddleware)
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return
	}

	// Get group ID from context (set by AuthMiddleware)
	groupID, _ := ctx.Value("groupId").(string)

	// Record usage
	err := m.billingClient.RecordUsage(context.Background(), userID, groupID, providerID, model, promptTokens, completionTokens)
	if err != nil {
		// Log error but don't fail the request
		// In production, you'd use proper logging
	}
}

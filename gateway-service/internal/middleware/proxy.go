package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	// Record token counts in response header for client tracking
	w.Header().Set("X-Prompt-Tokens", fmt.Sprintf("%d", resp.TokenCounts.PromptTokens))
	w.Header().Set("X-Completion-Tokens", fmt.Sprintf("%d", resp.TokenCounts.CompletionTokens))
	w.Header().Set("X-Total-Tokens", fmt.Sprintf("%d", resp.TokenCounts.TotalTokens))

	// Parse request and record usage to billing-service asynchronously
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

// handleStreamingRequest handles streaming requests with heartbeat
func (m *ProxyMiddleware) handleStreamingRequest(w http.ResponseWriter, r *http.Request, providerID string, requestBody []byte, headers map[string]string) {
	// Forward streaming request to provider-service
	stream, err := m.providerClient.StreamRequest(r.Context(), providerID, requestBody, headers)
	if err != nil {
		http.Error(w, "Failed to stream request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stream.CloseSend()

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering

	// Stream chunks
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	var totalPromptTokens, totalCompletionTokens int64

	// Create heartbeat ticker (every 15 seconds)
	heartbeatTicker := time.NewTicker(15 * time.Second)
	defer heartbeatTicker.Stop()

	// Create a context that captures client disconnection
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Channel to signal stream completion
	done := make(chan struct{})

	// Start a goroutine to handle heartbeat
	go func() {
		for {
			select {
			case <-heartbeatTicker.C:
				// Send heartbeat comment (SSE comment line)
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
				// Stream completed normally
				break
			}
			// Stream error
			fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
			flusher.Flush()
			break
		}

		// Write SSE chunk
		if len(chunk.ChunkData) > 0 {
			fmt.Fprintf(w, "data: %s\n\n", string(chunk.ChunkData))
			flusher.Flush()
		}

		// Accumulate token counts
		if chunk.AccumulatedTokens != nil {
			totalPromptTokens += chunk.AccumulatedTokens.PromptTokens
			totalCompletionTokens += chunk.AccumulatedTokens.CompletionTokens
		}

		// Check if this is the final chunk
		if chunk.Done {
			break
		}
	}

	// Signal completion to stop heartbeat goroutine
	close(done)

	// Send final SSE message with token counts
	finalChunk := map[string]interface{}{
		"prompt_tokens":     totalPromptTokens,
		"completion_tokens": totalCompletionTokens,
		"total_tokens":      totalPromptTokens + totalCompletionTokens,
		"done":              true,
	}
	finalJSON, _ := json.Marshal(finalChunk)
	fmt.Fprintf(w, "data: %s\n\n", string(finalJSON))
	flusher.Flush()

	// Parse request and record usage to billing-service asynchronously
	req, err := parseChatCompletionRequest(requestBody)
	model := "unknown"
	if err == nil && req.Model != "" {
		model = req.Model
	}
	go m.recordUsage(r.Context(), providerID, model, totalPromptTokens, totalCompletionTokens)
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

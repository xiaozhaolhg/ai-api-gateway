package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/gin-gonic/gin"
)

type ChatCompletionRequest struct {
	Model       string                 `json:"model"`
	Messages    []Message              `json:"messages"`
	Stream      bool                   `json:"stream,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Other       map[string]interface{} `json:"-"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ProxyMiddleware struct {
	providerClient         *client.ProviderClient
	billingClient          *client.BillingClient
	streamingTokenInterval int64
}

func NewProxyMiddleware(providerClient *client.ProviderClient, billingClient *client.BillingClient, streamingTokenInterval int64) *ProxyMiddleware {
	return &ProxyMiddleware{
		providerClient:         providerClient,
		billingClient:          billingClient,
		streamingTokenInterval: streamingTokenInterval,
	}
}

func (m *ProxyMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		providerIDVal, exists := c.Get("providerId")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Provider not resolved"})
			return
		}
		providerID, ok := providerIDVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Provider not resolved"})
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		stream := c.Query("stream") == "true"

		if stream {
			m.handleStreamingRequest(c, providerID, bodyBytes)
		} else {
			m.handleNonStreamingRequest(c, providerID, bodyBytes)
		}
	}
}

// handleNonStreamingRequest handles non-streaming requests with fallback support
func (m *ProxyMiddleware) handleNonStreamingRequest(c *gin.Context, providerID string, requestBody []byte) {
	r := c.Request

	// Get fallback info from context
	fallbackProviderIDs, _ := r.Context().Value("fallbackProviderIds").([]string)
	fallbackModels, _ := r.Context().Value("fallbackModels").([]string)
	log.Printf("[Fallback] Fallback info: primary=%s, fallbacks=%v, models=%v", providerID, fallbackProviderIDs, fallbackModels)

	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	// Try primary provider first
	resp, statusCode, err := m.tryNonStreamingProvider(r.Context(), providerID, requestBody, headers)
	if err == nil {
		m.writeNonStreamingResponse(c, providerID, resp, requestBody)
		return
	}

	// Check if error is non-retryable (4xx client errors)
	if statusCode >= 400 && statusCode < 500 {
		log.Printf("[Fallback] Primary provider %s returned non-retryable error %d: %v, not attempting fallback", providerID, statusCode, err)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Provider error: " + err.Error()})
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
		m.writeNonStreamingResponse(c, fallbackID, resp, bodyToSend)
		return
	}

	// All providers failed - return structured error
	m.writeFallbackError(c, "all_providers_failed", "All providers failed")
}

func parseChatCompletionRequest(requestBody []byte) (*ChatCompletionRequest, error) {
	var req ChatCompletionRequest
	if err := json.Unmarshal(requestBody, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (m *ProxyMiddleware) handleStreamingRequest(c *gin.Context, providerID string, requestBody []byte) {
	r := c.Request
	w := c.Writer

	// Get fallback info from context
	fallbackProviderIDs, _ := r.Context().Value("fallbackProviderIds").([]string)
	fallbackModels, _ := r.Context().Value("fallbackModels").([]string)

	req, parseErr := parseChatCompletionRequest(requestBody)
	model := "unknown"
	if parseErr == nil && req.Model != "" {
		model = req.Model
	}

	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	// Try primary provider first
	err := m.tryStreamingProvider(w, r, providerID, requestBody, headers, model, true)
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

		err := m.tryStreamingProvider(w, r, fallbackID, bodyToSend, headers, model, false)
		if err != nil {
			log.Printf("[Fallback] Fallback streaming provider %s failed: %v", fallbackID, err)
			continue
		}

		log.Printf("[Fallback] Successfully fell back to streaming provider %s", fallbackID)
		return
	}

	// All providers failed - return structured error
	m.writeFallbackError(c, "all_providers_failed", "All streaming providers failed")
}

func (m *ProxyMiddleware) tryStreamingProvider(w http.ResponseWriter, r *http.Request, providerID string, requestBody []byte, headers map[string]string, model string, isPrimary bool) error {
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
	var lastRecordedPromptTokens, lastRecordedCompletionTokens int64

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
			fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
			flusher.Flush()
			break
		}

		if len(chunk.ChunkData) > 0 {
			fmt.Fprintf(w, "data: %s\n\n", string(chunk.ChunkData))
			flusher.Flush()
		}

		if chunk.AccumulatedTokens != nil {
			totalPromptTokens += chunk.AccumulatedTokens.PromptTokens
			totalCompletionTokens += chunk.AccumulatedTokens.CompletionTokens
		}

		if m.streamingTokenInterval > 0 && totalCompletionTokens-lastRecordedCompletionTokens >= m.streamingTokenInterval {
			deltaPrompt := totalPromptTokens - lastRecordedPromptTokens
			deltaCompletion := totalCompletionTokens - lastRecordedCompletionTokens
			go m.recordUsage(r.Context(), providerID, model, deltaPrompt, deltaCompletion)
			lastRecordedPromptTokens = totalPromptTokens
			lastRecordedCompletionTokens = totalCompletionTokens
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

	finalDeltaPrompt := totalPromptTokens - lastRecordedPromptTokens
	finalDeltaCompletion := totalCompletionTokens - lastRecordedCompletionTokens
	if finalDeltaCompletion > 0 {
		go m.recordUsage(r.Context(), providerID, model, finalDeltaPrompt, finalDeltaCompletion)
	}
	return nil
}

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

// tryNonStreamingProvider attempts a single non-streaming request to a provider
// Returns response, HTTP status code from provider, and error
func (m *ProxyMiddleware) tryNonStreamingProvider(ctx context.Context, providerID string, requestBody []byte, headers map[string]string) (*client.ForwardRequestResponse, int32, error) {
	resp, err := m.providerClient.ForwardRequest(ctx, providerID, requestBody, headers)
	if err != nil {
		return nil, 0, err
	}
	return resp, resp.StatusCode, nil
}

// writeFallbackError writes a structured error response for fallback failures
func (m *ProxyMiddleware) writeFallbackError(c *gin.Context, code, message string) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{"code": code, "message": message}})
}

// writeNonStreamingResponse writes the response and records usage
func (m *ProxyMiddleware) writeNonStreamingResponse(c *gin.Context, providerID string, resp *client.ForwardRequestResponse, requestBody []byte) {
	c.Header("Content-Type", "application/json")
	c.Header("X-Prompt-Tokens", fmt.Sprintf("%d", resp.TokenCounts.PromptTokens))
	c.Header("X-Completion-Tokens", fmt.Sprintf("%d", resp.TokenCounts.CompletionTokens))
	c.Header("X-Total-Tokens", fmt.Sprintf("%d", resp.TokenCounts.TotalTokens))
	c.Data(http.StatusOK, "application/json", resp.ResponseBody)

	req, err := parseChatCompletionRequest(requestBody)
	model := "unknown"
	if err == nil && req.Model != "" {
		model = req.Model
	}
	go m.recordUsage(c.Request.Context(), providerID, model, resp.TokenCounts.PromptTokens, resp.TokenCounts.CompletionTokens)
}

func (m *ProxyMiddleware) recordUsage(ctx context.Context, providerID, model string, promptTokens, completionTokens int64) {
	userID, _ := ctx.Value("userId").(string)
	if userID == "" {
		return
	}

	groupID, _ := ctx.Value("groupId").(string)

	err := m.billingClient.RecordUsage(context.Background(), userID, groupID, providerID, model, promptTokens, completionTokens)
	if err != nil {
		// Log error but don't fail the request
	}
}

package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"github.com/ai-api-gateway/provider-service/internal/domain/port"
)

// AnthropicAdapter implements ProviderAdapter for Anthropic API
type AnthropicAdapter struct{}

// NewAnthropicAdapter creates a new Anthropic adapter
func NewAnthropicAdapter() port.ProviderAdapter {
	return &AnthropicAdapter{}
}

// TransformRequest transforms OpenAI format request to Anthropic format
func (a *AnthropicAdapter) TransformRequest(request []byte, headers map[string]string) ([]byte, map[string]string, error) {
	// Parse OpenAI format request
	var openAIReq struct {
		Model    string                 `json:"model"`
		Messages []map[string]interface{} `json:"messages"`
		Stream   bool                   `json:"stream"`
		Temperature float64             `json:"temperature,omitempty"`
		MaxTokens int                   `json:"max_tokens,omitempty"`
	}

	if err := json.Unmarshal(request, &openAIReq); err != nil {
		return nil, nil, fmt.Errorf("invalid OpenAI request format: %w", err)
	}

	if openAIReq.Model == "" {
		return nil, nil, fmt.Errorf("model field is required")
	}

	// Transform to Anthropic format
	anthropicReq := map[string]interface{}{
		"model":         a.convertModelName(openAIReq.Model),
		"max_tokens":    openAIReq.MaxTokens,
		"messages":      a.convertMessages(openAIReq.Messages),
		"stream":        openAIReq.Stream,
	}

	if openAIReq.Temperature > 0 {
		anthropicReq["temperature"] = openAIReq.Temperature
	}

	transformedRequest, err := json.Marshal(anthropicReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal Anthropic request: %w", err)
	}

	// Transform headers
	transformedHeaders := make(map[string]string)
	for k, v := range headers {
		transformedHeaders[k] = v
	}
	transformedHeaders["Content-Type"] = "application/json"
	transformedHeaders["anthropic-version"] = "2023-06-01"

	return transformedRequest, transformedHeaders, nil
}

// TransformResponse transforms Anthropic format response back to OpenAI format.
//
// For non-streaming (isStreaming=false):
//   - Transforms complete Anthropic response to OpenAI format
//   - Extracts token counts from usage data
//
// For streaming (isStreaming=true):
//   - Parses Anthropic SSE chunk format
//   - Transforms content_block_delta events to OpenAI delta format
//   - Detects message_stop event for final chunk
//   - Accumulates token counts during streaming
func (a *AnthropicAdapter) TransformResponse(response []byte, isStreaming bool, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	if !isStreaming {
		return a.transformNonStreamingResponse(response, accumulatedTokens)
	}

	return a.transformStreamingChunk(response, accumulatedTokens)
}

// transformNonStreamingResponse handles non-streaming response transformation
func (a *AnthropicAdapter) transformNonStreamingResponse(response []byte, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	// Parse Anthropic format response
	var anthropicResp struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Model       string  `json:"model"`
		StopReason *string `json:"stop_reason"`
		Usage      struct {
			InputTokens  int64 `json:"input_tokens"`
			OutputTokens int64 `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(response, &anthropicResp); err != nil {
		return nil, accumulatedTokens, false, fmt.Errorf("invalid Anthropic response format: %w", err)
	}

	// Transform to OpenAI format
	openAIResp := map[string]interface{}{
		"id":      anthropicResp.ID,
		"object":  "chat.completion",
		"created": 0, // Anthropic doesn't provide timestamp
		"model":   anthropicResp.Model,
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": a.extractContent(anthropicResp.Content),
				},
				"finish_reason": a.convertStopReason(anthropicResp.StopReason),
			},
		},
		"usage": map[string]interface{}{
			"prompt_tokens":     anthropicResp.Usage.InputTokens,
			"completion_tokens": anthropicResp.Usage.OutputTokens,
			"total_tokens":      anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens,
		},
	}

	transformedResponse, err := json.Marshal(openAIResp)
	if err != nil {
		return nil, accumulatedTokens, false, fmt.Errorf("failed to marshal OpenAI response: %w", err)
	}

	tokenCounts := entity.TokenCounts{
		PromptTokens:      anthropicResp.Usage.InputTokens,
		CompletionTokens:  anthropicResp.Usage.OutputTokens,
		AccumulatedTokens: anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens,
	}

	return transformedResponse, tokenCounts, true, nil
}

// transformStreamingChunk handles streaming SSE chunk transformation
func (a *AnthropicAdapter) transformStreamingChunk(chunk []byte, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	// Check for message_stop event (final chunk indicator in Anthropic)
	if bytes.Contains(chunk, []byte("event: message_stop")) {
		return chunk, accumulatedTokens, true, nil
	}

	// Parse SSE data line
	data, err := a.parseSSEData(chunk)
	if err != nil {
		// If parsing fails, return chunk as-is
		return chunk, accumulatedTokens, false, nil
	}

	// Try to parse the event type
	if eventType, ok := a.extractSSEEventType(chunk); ok {
		switch eventType {
		case "content_block_delta":
			return a.transformContentBlockDelta(data, accumulatedTokens)
		case "message_delta":
			// Message delta contains usage data in some Anthropic versions
			return a.transformMessageDelta(data, accumulatedTokens)
		default:
			// Other events pass through or get transformed
			return chunk, accumulatedTokens, false, nil
		}
	}

	// If no event type, try to parse as JSON
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		return chunk, accumulatedTokens, false, nil
	}

	// Accumulate content length
	contentLength := len(data)
	accumulatedTokens.AccumulatedTokens += int64(contentLength / 4)

	return chunk, accumulatedTokens, false, nil
}

// transformContentBlockDelta transforms Anthropic content_block_delta to OpenAI delta format
func (a *AnthropicAdapter) transformContentBlockDelta(data []byte, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	var delta struct {
		Type  string `json:"type"`
		Index int    `json:"index"`
		Delta struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"delta"`
	}

	if err := json.Unmarshal(data, &delta); err != nil {
		// Return as-is if parsing fails
		return []byte("data: " + string(data) + "\n\n"), accumulatedTokens, false, nil
	}

	// Transform to OpenAI streaming format
	openAIChunk := map[string]interface{}{
		"id":      "chatcmpl-anthropic",
		"object":  "chat.completion.chunk",
		"created": 0,
		"model":   "claude-model",
		"choices": []map[string]interface{}{
			{
				"index": delta.Index,
				"delta": map[string]interface{}{
					"content": delta.Delta.Text,
				},
				"finish_reason": nil,
			},
		},
	}

	transformed, err := json.Marshal(openAIChunk)
	if err != nil {
		return []byte("data: " + string(data) + "\n\n"), accumulatedTokens, false, nil
	}

	// Accumulate token count based on content
	accumulatedTokens.AccumulatedTokens += int64(len(delta.Delta.Text) / 4)

	return []byte("data: " + string(transformed) + "\n\n"), accumulatedTokens, false, nil
}

// transformMessageDelta transforms Anthropic message_delta to OpenAI format with usage
func (a *AnthropicAdapter) transformMessageDelta(data []byte, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	var delta struct {
		Type  string `json:"type"`
		Delta struct {
			StopReason   *string `json:"stop_reason"`
			Usage        struct {
				OutputTokens int64 `json:"output_tokens"`
			} `json:"usage"`
		} `json:"delta"`
		Usage struct {
			OutputTokens int64 `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(data, &delta); err != nil {
		return []byte("data: " + string(data) + "\n\n"), accumulatedTokens, false, nil
	}

	// Update token counts if available
	if delta.Usage.OutputTokens > 0 {
		accumulatedTokens.CompletionTokens = delta.Usage.OutputTokens
		accumulatedTokens.AccumulatedTokens = accumulatedTokens.PromptTokens + accumulatedTokens.CompletionTokens
	}

	// Transform to OpenAI format with finish_reason
	openAIChunk := map[string]interface{}{
		"id":      "chatcmpl-anthropic",
		"object":  "chat.completion.chunk",
		"created": 0,
		"model":   "claude-model",
		"choices": []map[string]interface{}{
			{
				"index":         0,
				"delta":         map[string]interface{}{},
				"finish_reason": a.convertStopReason(delta.Delta.StopReason),
			},
		},
	}

	transformed, err := json.Marshal(openAIChunk)
	if err != nil {
		return []byte("data: " + string(data) + "\n\n"), accumulatedTokens, false, nil
	}

	return []byte("data: " + string(transformed) + "\n\n"), accumulatedTokens, false, nil
}

// extractSSEEventType extracts the event type from SSE format
func (a *AnthropicAdapter) extractSSEEventType(chunk []byte) (string, bool) {
	chunkStr := string(chunk)
	
	// Look for "event: " prefix
	const eventPrefix = "event: "
	if idx := strings.Index(chunkStr, eventPrefix); idx != -1 {
		// Extract event type
		eventStart := idx + len(eventPrefix)
		eventEnd := strings.Index(chunkStr[eventStart:], "\n")
		if eventEnd == -1 {
			eventEnd = len(chunkStr) - eventStart
		}
		return strings.TrimSpace(chunkStr[eventStart : eventStart+eventEnd]), true
	}
	
	return "", false
}

// CountTokens counts tokens in the request/response.
//
// For non-streaming (isStreaming=false):
//   - Returns actual token counts from response usage field if available
//   - Falls back to character-based estimation
//
// For streaming (isStreaming=true):
//   - Returns (0, 0) for intermediate chunks
//   - Returns actual counts for final chunk
//   - May estimate tokens if usage data not available
func (a *AnthropicAdapter) CountTokens(request []byte, response []byte, isStreaming bool) (int64, int64, error) {
	if !isStreaming {
		// Try to extract from response if available
		var resp struct {
			Usage struct {
				InputTokens  int64 `json:"input_tokens"`
				OutputTokens int64 `json:"output_tokens"`
			} `json:"usage"`
		}

		if err := json.Unmarshal(response, &resp); err == nil {
			return resp.Usage.InputTokens, resp.Usage.OutputTokens, nil
		}

		// Fallback to estimation
		promptTokens := int64(len(request) / 4)
		completionTokens := int64(len(response) / 4)
		return promptTokens, completionTokens, nil
	}

	// Streaming: check for final chunk with usage data
	if bytes.Contains(response, []byte("event: message_stop")) {
		// Final chunk - extract any usage data if present
		var resp struct {
			Usage struct {
				OutputTokens int64 `json:"output_tokens"`
			} `json:"usage"`
		}
		if err := json.Unmarshal(response, &resp); err == nil {
			return 0, resp.Usage.OutputTokens, nil
		}
		return 0, 0, nil
	}

	// Intermediate chunk: return 0, 0
	return 0, 0, nil
}

// parseSSEData extracts JSON data from SSE format (data: {...})
func (a *AnthropicAdapter) parseSSEData(chunk []byte) ([]byte, error) {
	chunkStr := string(chunk)

	// Look for "data: " prefix
	const dataPrefix = "data: "
	if idx := strings.Index(chunkStr, dataPrefix); idx != -1 {
		dataStart := idx + len(dataPrefix)
		// Find end of data (next newline or end of string)
		dataEnd := strings.Index(chunkStr[dataStart:], "\n")
		if dataEnd == -1 {
			dataEnd = len(chunkStr) - dataStart
		}
		data := strings.TrimSpace(chunkStr[dataStart : dataStart+dataEnd])
		return []byte(data), nil
	}

	return nil, fmt.Errorf("not a valid SSE data line")
}

// convertModelName converts OpenAI model names to Anthropic model names
func (a *AnthropicAdapter) convertModelName(model string) string {
	// Simple mapping for common models
	modelMap := map[string]string{
		"gpt-4":         "claude-3-opus-20240229",
		"gpt-4-turbo":   "claude-3-sonnet-20240229",
		"gpt-3.5-turbo": "claude-3-haiku-20240307",
	}

	if anthropicModel, ok := modelMap[model]; ok {
		return anthropicModel
	}

	// Return as-is if no mapping found
	return model
}

// convertMessages converts OpenAI messages to Anthropic format
func (a *AnthropicAdapter) convertMessages(messages []map[string]interface{}) []map[string]interface{} {
	anthropicMessages := make([]map[string]interface{}, 0, len(messages))

	for _, msg := range messages {
		role, _ := msg["role"].(string)
		content, _ := msg["content"].(string)

		// Anthropic uses "user" and "assistant" roles
		anthropicRole := role
		if role == "system" {
			// Anthropic doesn't have a system role, convert to user
			anthropicRole = "user"
		}

		anthropicMessages = append(anthropicMessages, map[string]interface{}{
			"role":    anthropicRole,
			"content": content,
		})
	}

	return anthropicMessages
}

// extractContent extracts text content from Anthropic response
func (a *AnthropicAdapter) extractContent(content []struct {
	Type string `json:"type"`
	Text string `json:"text"`
}) string {
	var text string
	for _, c := range content {
		if c.Type == "text" {
			text += c.Text
		}
	}
	return text
}

// convertStopReason converts Anthropic stop reason to OpenAI format
func (a *AnthropicAdapter) convertStopReason(stopReason *string) string {
	if stopReason == nil {
		return "stop"
	}

	reasonMap := map[string]string{
		"end_turn":  "stop",
		"max_tokens": "length",
		"stop_sequence": "stop",
	}

	if openAIReason, ok := reasonMap[*stopReason]; ok {
		return openAIReason
	}

	return *stopReason
}

func (a *AnthropicAdapter) TestConnection(credentials string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("x-api-key", credentials)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid credentials")
	}
	return nil
}

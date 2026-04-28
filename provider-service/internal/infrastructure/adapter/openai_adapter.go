package adapter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"github.com/ai-api-gateway/provider-service/internal/domain/port"
)

// OpenAIAdapter implements ProviderAdapter for OpenAI-compatible APIs
type OpenAIAdapter struct{}

// NewOpenAIAdapter creates a new OpenAI adapter
func NewOpenAIAdapter() port.ProviderAdapter {
	return &OpenAIAdapter{}
}

// TransformRequest transforms the request to OpenAI format
// Since OpenAI is the baseline format, this is mostly a pass-through
func (a *OpenAIAdapter) TransformRequest(request []byte, headers map[string]string) ([]byte, map[string]string, error) {
	// Parse the request to validate it's in OpenAI format
	var req map[string]interface{}
	if err := json.Unmarshal(request, &req); err != nil {
		return nil, nil, fmt.Errorf("invalid OpenAI request format: %w", err)
	}

	// Ensure required fields exist
	if _, ok := req["model"]; !ok {
		return nil, nil, fmt.Errorf("missing required field: model")
	}

	// Transform headers for OpenAI API
	transformedHeaders := make(map[string]string)
	for k, v := range headers {
		transformedHeaders[k] = v
	}

	// Set OpenAI-specific headers
	transformedHeaders["Content-Type"] = "application/json"

	// Return the request as-is (OpenAI format)
	return request, transformedHeaders, nil
}

// TransformResponse transforms the response back to OpenAI format.
//
// For non-streaming (isStreaming=false):
//   - Validates the complete response is in OpenAI format
//   - Returns the response as-is with token counts extracted from usage
//
// For streaming (isStreaming=true):
//   - Parses SSE chunk format (data: {...})
//   - Passes through SSE chunks unchanged
//   - Detects [DONE] marker for final chunk
//   - Extracts usage data from final chunk
func (a *OpenAIAdapter) TransformResponse(response []byte, isStreaming bool, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	if !isStreaming {
		// Non-streaming: validate and return with extracted token counts
		return a.transformNonStreamingResponse(response, accumulatedTokens)
	}

	// Streaming: parse SSE chunk
	return a.transformStreamingChunk(response, accumulatedTokens)
}

// transformNonStreamingResponse handles non-streaming response transformation
func (a *OpenAIAdapter) transformNonStreamingResponse(response []byte, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	// Parse the response to validate it's in OpenAI format
	var resp map[string]interface{}
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, accumulatedTokens, false, fmt.Errorf("invalid OpenAI response format: %w", err)
	}

	// Ensure required fields exist
	if _, ok := resp["id"]; !ok {
		return nil, accumulatedTokens, false, fmt.Errorf("missing required field: id")
	}

	// Extract token counts from response if available
	promptTokens, completionTokens, err := a.extractTokenUsage(response)
	if err != nil {
		// Fallback to accumulated tokens if extraction fails
		promptTokens = accumulatedTokens.PromptTokens
		completionTokens = accumulatedTokens.CompletionTokens
	}

	tokenCounts := entity.TokenCounts{
		PromptTokens:      promptTokens,
		CompletionTokens:  completionTokens,
		AccumulatedTokens: promptTokens + completionTokens,
	}

	// Return the response as-is (OpenAI format), isFinal=true for non-streaming
	return response, tokenCounts, true, nil
}

// transformStreamingChunk handles streaming SSE chunk transformation
func (a *OpenAIAdapter) transformStreamingChunk(chunk []byte, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	// Check for [DONE] marker (final chunk)
	if bytes.Contains(chunk, []byte("[DONE]")) {
		return chunk, accumulatedTokens, true, nil
	}

	// Parse SSE data line (format: "data: {...}")
	data, err := a.parseSSEData(chunk)
	if err != nil {
		// If parsing fails, return chunk as-is and continue
		return chunk, accumulatedTokens, false, nil
	}

	// Try to parse the JSON content
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		// Not valid JSON, pass through as-is
		return chunk, accumulatedTokens, false, nil
	}

	// Check for usage in the chunk (final chunk in some OpenAI responses)
	if usage, ok := event["usage"].(map[string]interface{}); ok {
		promptTokens := int64(0)
		completionTokens := int64(0)

		if pt, ok := usage["prompt_tokens"].(float64); ok {
			promptTokens = int64(pt)
		}
		if ct, ok := usage["completion_tokens"].(float64); ok {
			completionTokens = int64(ct)
		}

		accumulatedTokens.PromptTokens = promptTokens
		accumulatedTokens.CompletionTokens = completionTokens
		accumulatedTokens.AccumulatedTokens = promptTokens + completionTokens
	} else {
		// Accumulate content length for intermediate chunks
		// Estimate tokens from content
		contentLength := len(data)
		accumulatedTokens.AccumulatedTokens += int64(contentLength / 4)
	}

	// For OpenAI, we pass through the SSE chunk unchanged
	// The chunk is already in OpenAI format
	return chunk, accumulatedTokens, false, nil
}

// parseSSEData extracts JSON data from SSE format (data: {...})
func (a *OpenAIAdapter) parseSSEData(chunk []byte) ([]byte, error) {
	// Convert to string for easier processing
	chunkStr := string(chunk)

	// Look for "data: " prefix
	const dataPrefix = "data: "
	if !strings.HasPrefix(chunkStr, dataPrefix) {
		return nil, fmt.Errorf("not a valid SSE data line")
	}

	// Extract the data portion
	data := chunkStr[len(dataPrefix):]

	// Remove trailing newlines
	data = strings.TrimSpace(data)

	return []byte(data), nil
}

// CountTokens counts tokens in the request/response.
//
// For non-streaming (isStreaming=false):
//   - Returns actual token counts from response usage field if available
//   - Falls back to character-based estimation (~4 chars per token)
//
// For streaming (isStreaming=true):
//   - Returns (0, 0) for intermediate chunks
//   - Returns actual counts from usage data for final chunk
//   - May estimate tokens if usage data not available
func (a *OpenAIAdapter) CountTokens(request []byte, response []byte, isStreaming bool) (int64, int64, error) {
	if !isStreaming {
		// Non-streaming: try to extract actual usage, fallback to estimation
		promptTokens, completionTokens, err := a.extractTokenUsage(response)
		if err == nil && (promptTokens > 0 || completionTokens > 0) {
			return promptTokens, completionTokens, nil
		}

		// Fallback to character-based estimation
		promptTokens = int64(len(request) / 4)
		completionTokens = int64(len(response) / 4)
		return promptTokens, completionTokens, nil
	}

	// Streaming: extract from SSE chunk if it's the final one
	if bytes.Contains(response, []byte("[DONE]")) {
		// Final chunk - try to extract usage from previous accumulated data
		// In practice, the caller should pass the final chunk with usage data
		promptTokens, completionTokens, _ := a.extractTokenUsage(response)
		return promptTokens, completionTokens, nil
	}

	// Intermediate chunk: return 0, 0 (accumulated separately during streaming)
	return 0, 0, nil
}

// extractTokenUsage extracts token counts from OpenAI response if available
func (a *OpenAIAdapter) extractTokenUsage(response []byte) (int64, int64, error) {
	var resp struct {
		Usage struct {
			PromptTokens     int64 `json:"prompt_tokens"`
			CompletionTokens int64 `json:"completion_tokens"`
			TotalTokens      int64 `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(response, &resp); err != nil {
		return 0, 0, err
	}

	return resp.Usage.PromptTokens, resp.Usage.CompletionTokens, nil
}

// validateOpenAIRequest validates that the request is in OpenAI format
func (a *OpenAIAdapter) validateOpenAIRequest(req map[string]interface{}) error {
	// Check for required fields
	if _, ok := req["model"]; !ok {
		return fmt.Errorf("missing required field: model")
	}

	// Check for messages or prompt field
	if _, ok := req["messages"]; !ok {
		if _, ok := req["prompt"]; !ok {
			return fmt.Errorf("missing required field: messages or prompt")
		}
	}

	return nil
}

// validateOpenAIResponse validates that the response is in OpenAI format
func (a *OpenAIAdapter) validateOpenAIResponse(resp map[string]interface{}) error {
	// Check for required fields
	if _, ok := resp["id"]; !ok {
		return fmt.Errorf("missing required field: id")
	}

	if _, ok := resp["object"]; !ok {
		return fmt.Errorf("missing required field: object")
	}

	if _, ok := resp["created"]; !ok {
		return fmt.Errorf("missing required field: created")
	}

	return nil
}

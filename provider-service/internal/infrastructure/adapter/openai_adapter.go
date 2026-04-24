package adapter

import (
	"encoding/json"
	"fmt"

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

// TransformResponse transforms the response back to OpenAI format
// Since OpenAI is the baseline format, this is mostly a pass-through
func (a *OpenAIAdapter) TransformResponse(response []byte) ([]byte, error) {
	// Parse the response to validate it's in OpenAI format
	var resp map[string]interface{}
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, fmt.Errorf("invalid OpenAI response format: %w", err)
	}

	// Ensure required fields exist
	if _, ok := resp["id"]; !ok {
		return nil, fmt.Errorf("missing required field: id")
	}

	// Return the response as-is (OpenAI format)
	return response, nil
}

// CountTokens counts tokens in the request/response
// First tries to extract from response usage, falls back to character count estimation
func (a *OpenAIAdapter) CountTokens(request []byte, response []byte) (int64, int64, error) {
	// Try to extract from response usage first
	if len(response) > 0 {
		if promptTokens, completionTokens, err := a.extractTokenUsage(response); err == nil && (promptTokens > 0 || completionTokens > 0) {
			return promptTokens, completionTokens, nil
		}
	}

	// Fallback: ~4 characters per token
	promptTokens := int64(len(request) / 4)
	completionTokens := int64(len(response) / 4)

	return promptTokens, completionTokens, nil
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

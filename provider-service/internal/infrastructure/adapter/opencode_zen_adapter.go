package adapter

import (
	"encoding/json"
	"fmt"

	"github.com/ai-api-gateway/provider-service/internal/domain/port"
)

// OpenCodeZenAdapter implements ProviderAdapter for OpenCode Zen API
type OpenCodeZenAdapter struct{}

// NewOpenCodeZenAdapter creates a new OpenCode Zen adapter
func NewOpenCodeZenAdapter() port.ProviderAdapter {
	return &OpenCodeZenAdapter{}
}

// TransformRequest transforms OpenAI format request to OpenCode Zen format
func (a *OpenCodeZenAdapter) TransformRequest(request []byte, headers map[string]string) ([]byte, map[string]string, error) {
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

	// Transform to OpenCode Zen format
	// OpenCode Zen uses a similar format to OpenAI but with some differences
	zenReq := map[string]interface{}{
		"model":    openAIReq.Model,
		"messages": openAIReq.Messages,
		"stream":   openAIReq.Stream,
	}

	if openAIReq.Temperature > 0 {
		zenReq["temperature"] = openAIReq.Temperature
	}

	if openAIReq.MaxTokens > 0 {
		zenReq["max_tokens"] = openAIReq.MaxTokens
	}

	transformedRequest, err := json.Marshal(zenReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal OpenCode Zen request: %w", err)
	}

	// Transform headers
	transformedHeaders := make(map[string]string)
	for k, v := range headers {
		transformedHeaders[k] = v
	}
	transformedHeaders["Content-Type"] = "application/json"

	return transformedRequest, transformedHeaders, nil
}

// TransformResponse transforms OpenCode Zen format response back to OpenAI format
func (a *OpenCodeZenAdapter) TransformResponse(response []byte) ([]byte, error) {
	// Parse OpenCode Zen format response
	var zenResp struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int64 `json:"prompt_tokens"`
			CompletionTokens int64 `json:"completion_tokens"`
			TotalTokens      int64 `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(response, &zenResp); err != nil {
		return nil, fmt.Errorf("invalid OpenCode Zen response format: %w", err)
	}

	// OpenCode Zen response is already in OpenAI format, so return as-is
	return response, nil
}

// CountTokens counts tokens in the request/response
func (a *OpenCodeZenAdapter) CountTokens(request []byte, response []byte) (int64, int64, error) {
	// Try to extract from response if available
	var resp struct {
		Usage struct {
			PromptTokens     int64 `json:"prompt_tokens"`
			CompletionTokens int64 `json:"completion_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(response, &resp); err == nil {
		return resp.Usage.PromptTokens, resp.Usage.CompletionTokens, nil
	}

	// Fallback to estimation
	promptTokens := int64(len(request) / 4)
	completionTokens := int64(len(response) / 4)

	return promptTokens, completionTokens, nil
}

package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
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

// TransformResponse transforms OpenCode Zen format response back to OpenAI format.
//
// For non-streaming: Pass-through (already OpenAI format) with token extraction
// For streaming: Pass-through with token accumulation
func (a *OpenCodeZenAdapter) TransformResponse(response []byte, isStreaming bool, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	if !isStreaming {
		return a.transformNonStreamingResponse(response, accumulatedTokens)
	}

	// For streaming, pass through and accumulate tokens
	accumulatedTokens.AccumulatedTokens += int64(len(response) / 4)
	return response, accumulatedTokens, false, nil
}

// transformNonStreamingResponse handles non-streaming response transformation
func (a *OpenCodeZenAdapter) transformNonStreamingResponse(response []byte, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	// Parse OpenCode Zen format response to extract token counts
	var zenResp struct {
		Usage struct {
			PromptTokens     int64 `json:"prompt_tokens"`
			CompletionTokens int64 `json:"completion_tokens"`
			TotalTokens      int64 `json:"total_tokens"`
		} `json:"usage"`
	}

	// Try to extract usage data, but don't fail if parsing fails
	json.Unmarshal(response, &zenResp)

	tokenCounts := entity.TokenCounts{
		PromptTokens:      zenResp.Usage.PromptTokens,
		CompletionTokens:  zenResp.Usage.CompletionTokens,
		AccumulatedTokens: zenResp.Usage.TotalTokens,
	}

	// OpenCode Zen response is already in OpenAI format, so return as-is
	return response, tokenCounts, true, nil
}

// CountTokens counts tokens in the request/response.
//
// For non-streaming: Extract from response or estimate
// For streaming: Return 0, 0 for intermediate chunks
func (a *OpenCodeZenAdapter) CountTokens(request []byte, response []byte, isStreaming bool) (int64, int64, error) {
	if !isStreaming {
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

	// Streaming: return 0, 0 (tokens accumulated via TransformResponse)
	return 0, 0, nil
}

func (a *OpenCodeZenAdapter) TestConnection(credentials string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	url := "https://opencode.ai/zen/v1/models"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+credentials)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

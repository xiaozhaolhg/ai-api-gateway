package adapter

import (
	"encoding/json"
	"fmt"

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

// TransformResponse transforms Anthropic format response back to OpenAI format
func (a *AnthropicAdapter) TransformResponse(response []byte) ([]byte, error) {
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
		return nil, fmt.Errorf("invalid Anthropic response format: %w", err)
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
		return nil, fmt.Errorf("failed to marshal OpenAI response: %w", err)
	}

	return transformedResponse, nil
}

// CountTokens counts tokens in the request/response
func (a *AnthropicAdapter) CountTokens(request []byte, response []byte) (int64, int64, error) {
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

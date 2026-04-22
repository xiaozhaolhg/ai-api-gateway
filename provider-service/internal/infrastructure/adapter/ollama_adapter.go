package adapter

import (
	"encoding/json"
	"fmt"

	"github.com/ai-api-gateway/provider-service/internal/domain/port"
)

// OllamaAdapter implements ProviderAdapter for Ollama API
type OllamaAdapter struct{}

// NewOllamaAdapter creates a new Ollama adapter
func NewOllamaAdapter() port.ProviderAdapter {
	return &OllamaAdapter{}
}

// TransformRequest transforms OpenAI format request to Ollama format
func (a *OllamaAdapter) TransformRequest(request []byte, headers map[string]string) ([]byte, map[string]string, error) {
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

	// Transform to Ollama format
	ollamaReq := map[string]interface{}{
		"model":  openAIReq.Model,
		"stream": openAIReq.Stream,
	}

	// Convert messages to prompt (Ollama uses a single prompt field)
	prompt := a.convertMessagesToPrompt(openAIReq.Messages)
	ollamaReq["prompt"] = prompt

	if openAIReq.Temperature > 0 {
		ollamaReq["temperature"] = openAIReq.Temperature
	}

	if openAIReq.MaxTokens > 0 {
		ollamaReq["num_predict"] = openAIReq.MaxTokens
	}

	transformedRequest, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal Ollama request: %w", err)
	}

	// Transform headers
	transformedHeaders := make(map[string]string)
	for k, v := range headers {
		transformedHeaders[k] = v
	}
	transformedHeaders["Content-Type"] = "application/json"

	return transformedRequest, transformedHeaders, nil
}

// TransformResponse transforms Ollama format response back to OpenAI format
func (a *OllamaAdapter) TransformResponse(response []byte) ([]byte, error) {
	// Parse Ollama format response
	var ollamaResp struct {
		Model     string `json:"model"`
		Response  string `json:"response"`
		Done      bool   `json:"done"`
		PromptEvalCount int64 `json:"prompt_eval_count"`
		EvalCount      int64 `json:"eval_count"`
	}

	if err := json.Unmarshal(response, &ollamaResp); err != nil {
		return nil, fmt.Errorf("invalid Ollama response format: %w", err)
	}

	// Transform to OpenAI format
	openAIResp := map[string]interface{}{
		"id":      "ollama-" + ollamaResp.Model,
		"object":  "chat.completion",
		"created": 0, // Ollama doesn't provide timestamp
		"model":   ollamaResp.Model,
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": ollamaResp.Response,
				},
				"finish_reason": a.convertDoneToFinishReason(ollamaResp.Done),
			},
		},
		"usage": map[string]interface{}{
			"prompt_tokens":     ollamaResp.PromptEvalCount,
			"completion_tokens": ollamaResp.EvalCount,
			"total_tokens":      ollamaResp.PromptEvalCount + ollamaResp.EvalCount,
		},
	}

	transformedResponse, err := json.Marshal(openAIResp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OpenAI response: %w", err)
	}

	return transformedResponse, nil
}

// CountTokens counts tokens in the request/response
func (a *OllamaAdapter) CountTokens(request []byte, response []byte) (int64, int64, error) {
	// Try to extract from response if available
	var resp struct {
		PromptEvalCount int64 `json:"prompt_eval_count"`
		EvalCount      int64 `json:"eval_count"`
	}

	if err := json.Unmarshal(response, &resp); err == nil {
		return resp.PromptEvalCount, resp.EvalCount, nil
	}

	// Fallback to estimation
	promptTokens := int64(len(request) / 4)
	completionTokens := int64(len(response) / 4)

	return promptTokens, completionTokens, nil
}

// convertMessagesToPrompt converts OpenAI messages to Ollama prompt format
func (a *OllamaAdapter) convertMessagesToPrompt(messages []map[string]interface{}) string {
	var prompt string

	for _, msg := range messages {
		role, _ := msg["role"].(string)
		content, _ := msg["content"].(string)

		// Ollama uses a simple format with role prefixes
		switch role {
		case "system":
			prompt += "System: " + content + "\n"
		case "user":
			prompt += "User: " + content + "\n"
		case "assistant":
			prompt += "Assistant: " + content + "\n"
		default:
			prompt += role + ": " + content + "\n"
		}
	}

	return prompt
}

// convertDoneToFinishReason converts Ollama done flag to OpenAI finish reason
func (a *OllamaAdapter) convertDoneToFinishReason(done bool) string {
	if done {
		return "stop"
	}
	return "length"
}

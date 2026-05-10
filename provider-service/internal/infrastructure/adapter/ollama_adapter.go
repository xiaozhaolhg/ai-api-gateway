package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
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
	var openAIReq struct {
		Model       string                   `json:"model"`
		Messages    []map[string]interface{} `json:"messages"`
		Stream      bool                     `json:"stream"`
		Temperature float64                  `json:"temperature,omitempty"`
		MaxTokens   int                      `json:"max_tokens,omitempty"`
	}

	if err := json.Unmarshal(request, &openAIReq); err != nil {
		return nil, nil, fmt.Errorf("invalid OpenAI request format: %w", err)
	}

	modelName := openAIReq.Model
	if idx := strings.Index(modelName, ":"); idx != -1 {
		modelName = modelName[idx+1:]
	}

	// Try to match with available Ollama models (check if modelName matches any model name)
	// For now, if model name contains colon, use it directly; otherwise use name:version
	if !strings.Contains(modelName, ":") {
		// Try common version suffixes for common models
		modelName = modelName + ":0.8b"
	}

	ollamaReq := map[string]interface{}{
		"model":   modelName,
		"stream":  openAIReq.Stream,
		"messages": openAIReq.Messages,
	}

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

	transformedHeaders := make(map[string]string)
	for k, v := range headers {
		transformedHeaders[k] = v
	}
	transformedHeaders["Content-Type"] = "application/json"

	return transformedRequest, transformedHeaders, nil
}

// TransformResponse transforms Ollama format response back to OpenAI format.
//
// For non-streaming: Full response transformation with token extraction
// For streaming: Pass-through with token accumulation
func (a *OllamaAdapter) TransformResponse(response []byte, isStreaming bool, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	if !isStreaming {
		return a.transformNonStreamingResponse(response, accumulatedTokens)
	}

	// For streaming, pass through and accumulate tokens
	accumulatedTokens.AccumulatedTokens += int64(len(response) / 4)
	return response, accumulatedTokens, false, nil
}

// transformNonStreamingResponse handles non-streaming response transformation
func (a *OllamaAdapter) transformNonStreamingResponse(response []byte, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	// Parse Ollama format response
	var ollamaResp struct {
		Model   string `json:"model"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Done            bool   `json:"done"`
		DoneReason      string `json:"done_reason"`
		PromptEvalCount int64  `json:"prompt_eval_count"`
		EvalCount       int64  `json:"eval_count"`
	}

	if err := json.Unmarshal(response, &ollamaResp); err != nil {
		return nil, accumulatedTokens, false, fmt.Errorf("invalid Ollama response format: %w", err)
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
					"role":    ollamaResp.Message.Role,
					"content": ollamaResp.Message.Content,
				},
				"finish_reason": a.convertDoneReasonToFinishReason(ollamaResp.DoneReason),
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
		return nil, accumulatedTokens, false, fmt.Errorf("failed to marshal OpenAI response: %w", err)
	}

	tokenCounts := entity.TokenCounts{
		PromptTokens:      ollamaResp.PromptEvalCount,
		CompletionTokens:  ollamaResp.EvalCount,
		AccumulatedTokens: ollamaResp.PromptEvalCount + ollamaResp.EvalCount,
	}

	return transformedResponse, tokenCounts, true, nil
}

// CountTokens counts tokens in the request/response.
//
// For non-streaming: Extract from response or estimate
// For streaming: Return 0, 0 for intermediate chunks
func (a *OllamaAdapter) CountTokens(request []byte, response []byte, isStreaming bool) (int64, int64, error) {
	if !isStreaming {
		// Try to extract from response if available
		var resp struct {
			PromptEvalCount int64 `json:"prompt_eval_count"`
			EvalCount       int64 `json:"eval_count"`
		}

		if err := json.Unmarshal(response, &resp); err == nil {
			return resp.PromptEvalCount, resp.EvalCount, nil
		}

		// Fallback to estimation
		promptTokens := int64(len(request) / 4)
		completionTokens := int64(len(response) / 4)
		return promptTokens, completionTokens, nil
	}

	// Streaming: return 0, 0 (tokens accumulated via TransformResponse)
	return 0, 0, nil
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

// convertDoneReasonToFinishReason converts Ollama done_reason string to OpenAI finish reason
func (a *OllamaAdapter) convertDoneReasonToFinishReason(doneReason string) string {
	switch doneReason {
	case "stop":
		return "stop"
	case "length":
		return "length"
	default:
		return "stop"
	}
}

func (a *OllamaAdapter) TestConnection(credentials string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	baseURL := credentials
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	resp, err := client.Get(baseURL + "/api/tags")
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

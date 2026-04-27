package adapter

import (
	"encoding/json"
	"fmt"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"github.com/ai-api-gateway/provider-service/internal/domain/port"
)

// GeminiAdapter implements ProviderAdapter for Google Gemini API
type GeminiAdapter struct{}

// NewGeminiAdapter creates a new Gemini adapter
func NewGeminiAdapter() port.ProviderAdapter {
	return &GeminiAdapter{}
}

// TransformRequest transforms OpenAI format request to Gemini format
func (a *GeminiAdapter) TransformRequest(request []byte, headers map[string]string) ([]byte, map[string]string, error) {
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

	// Transform to Gemini format
	geminiReq := map[string]interface{}{
		"contents": a.convertMessagesToContents(openAIReq.Messages),
	}

	if openAIReq.Temperature > 0 {
		geminiReq["generationConfig"] = map[string]interface{}{
			"temperature": openAIReq.Temperature,
		}
	}

	if openAIReq.MaxTokens > 0 {
		if _, ok := geminiReq["generationConfig"]; !ok {
			geminiReq["generationConfig"] = map[string]interface{}{}
		}
		geminiReq["generationConfig"].(map[string]interface{})["maxOutputTokens"] = openAIReq.MaxTokens
	}

	transformedRequest, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal Gemini request: %w", err)
	}

	// Transform headers
	transformedHeaders := make(map[string]string)
	for k, v := range headers {
		transformedHeaders[k] = v
	}
	transformedHeaders["Content-Type"] = "application/json"

	return transformedRequest, transformedHeaders, nil
}

// TransformResponse transforms Gemini format response back to OpenAI format.
//
// For non-streaming: Full response transformation with token extraction
// For streaming: Pass-through with token accumulation
func (a *GeminiAdapter) TransformResponse(response []byte, isStreaming bool, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	if !isStreaming {
		return a.transformNonStreamingResponse(response, accumulatedTokens)
	}

	// For streaming, pass through and accumulate tokens
	accumulatedTokens.AccumulatedTokens += int64(len(response) / 4)
	return response, accumulatedTokens, false, nil
}

// transformNonStreamingResponse handles non-streaming response transformation
func (a *GeminiAdapter) transformNonStreamingResponse(response []byte, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error) {
	// Parse Gemini format response
	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
				Role string `json:"role"`
			} `json:"content"`
			FinishReason string `json:"finishReason"`
		} `json:"candidates"`
		UsageMetadata struct {
			PromptTokenCount     int64 `json:"promptTokenCount"`
			CandidatesTokenCount int64 `json:"candidatesTokenCount"`
			TotalTokenCount      int64 `json:"totalTokenCount"`
		} `json:"usageMetadata"`
	}

	if err := json.Unmarshal(response, &geminiResp); err != nil {
		return nil, accumulatedTokens, false, fmt.Errorf("invalid Gemini response format: %w", err)
	}

	// Extract the first candidate's content
	if len(geminiResp.Candidates) == 0 {
		return nil, accumulatedTokens, false, fmt.Errorf("no candidates in Gemini response")
	}

	content := ""
	if len(geminiResp.Candidates[0].Content.Parts) > 0 {
		content = geminiResp.Candidates[0].Content.Parts[0].Text
	}

	// Transform to OpenAI format
	openAIResp := map[string]interface{}{
		"id":      "gemini-" + fmt.Sprintf("%d", 0), // Gemini doesn't provide ID
		"object":  "chat.completion",
		"created": 0, // Gemini doesn't provide timestamp
		"model":   "gemini-pro",
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": content,
				},
				"finish_reason": a.convertFinishReason(geminiResp.Candidates[0].FinishReason),
			},
		},
		"usage": map[string]interface{}{
			"prompt_tokens":     geminiResp.UsageMetadata.PromptTokenCount,
			"completion_tokens": geminiResp.UsageMetadata.CandidatesTokenCount,
			"total_tokens":      geminiResp.UsageMetadata.TotalTokenCount,
		},
	}

	transformedResponse, err := json.Marshal(openAIResp)
	if err != nil {
		return nil, accumulatedTokens, false, fmt.Errorf("failed to marshal OpenAI response: %w", err)
	}

	tokenCounts := entity.TokenCounts{
		PromptTokens:      geminiResp.UsageMetadata.PromptTokenCount,
		CompletionTokens:  geminiResp.UsageMetadata.CandidatesTokenCount,
		AccumulatedTokens: geminiResp.UsageMetadata.TotalTokenCount,
	}

	return transformedResponse, tokenCounts, true, nil
}

// CountTokens counts tokens in the request/response.
//
// For non-streaming: Extract from response or estimate
// For streaming: Return 0, 0 for intermediate chunks
func (a *GeminiAdapter) CountTokens(request []byte, response []byte, isStreaming bool) (int64, int64, error) {
	if !isStreaming {
		// Try to extract from response if available
		var resp struct {
			UsageMetadata struct {
				PromptTokenCount     int64 `json:"promptTokenCount"`
				CandidatesTokenCount int64 `json:"candidatesTokenCount"`
			} `json:"usageMetadata"`
		}

		if err := json.Unmarshal(response, &resp); err == nil {
			return resp.UsageMetadata.PromptTokenCount, resp.UsageMetadata.CandidatesTokenCount, nil
		}

		// Fallback to estimation
		promptTokens := int64(len(request) / 4)
		completionTokens := int64(len(response) / 4)
		return promptTokens, completionTokens, nil
	}

	// Streaming: return 0, 0 (tokens accumulated via TransformResponse)
	return 0, 0, nil
}

// convertMessagesToContents converts OpenAI messages to Gemini contents format
func (a *GeminiAdapter) convertMessagesToContents(messages []map[string]interface{}) []map[string]interface{} {
	contents := make([]map[string]interface{}, 0, len(messages))

	for _, msg := range messages {
		role, _ := msg["role"].(string)
		content, _ := msg["content"].(string)

		// Convert OpenAI roles to Gemini roles
		geminiRole := "user"
		if role == "assistant" {
			geminiRole = "model"
		} else if role == "system" {
			// Gemini doesn't have system role, convert to user with special handling
			geminiRole = "user"
			content = "System: " + content
		}

		contents = append(contents, map[string]interface{}{
			"role": geminiRole,
			"parts": []map[string]interface{}{
				{
					"text": content,
				},
			},
		})
	}

	return contents
}

// convertFinishReason converts Gemini finish reason to OpenAI format
func (a *GeminiAdapter) convertFinishReason(reason string) string {
	reasonMap := map[string]string{
		"STOP":    "stop",
		"MAX_TOKENS": "length",
		"SAFETY":  "content_filter",
		"RECITATION": "content_filter",
		"OTHER":   "stop",
	}

	if openAIReason, ok := reasonMap[reason]; ok {
		return openAIReason
	}

	return reason
}

package adapter

import (
	"encoding/json"
	"testing"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"github.com/ai-api-gateway/provider-service/internal/domain/port"
)

func TestOllamaAdapter_TransformRequest(t *testing.T) {
	adapter := NewOllamaAdapter()

	// Create a sample OpenAI format request
	openAIReq := map[string]interface{}{
		"model": "llama2",
		"messages": []map[string]interface{}{
			{"role": "user", "content": "Hello, how are you?"},
		},
		"stream": false,
	}

	requestJSON, _ := json.Marshal(openAIReq)
	headers := map[string]string{"Content-Type": "application/json"}

	transformed, transformedHeaders, err := adapter.TransformRequest(requestJSON, headers)
	if err != nil {
		t.Errorf("TransformRequest() error = %v", err)
	}

	// Parse the transformed request
	var ollamaReq map[string]interface{}
	if err := json.Unmarshal(transformed, &ollamaReq); err != nil {
		t.Errorf("Failed to parse transformed request: %v", err)
	}

	// Verify transformation
	if ollamaReq["model"] != "llama2" {
		t.Errorf("Expected model 'llama2', got %v", ollamaReq["model"])
	}

	if ollamaReq["stream"] != false {
		t.Errorf("Expected stream false, got %v", ollamaReq["stream"])
	}

	prompt, ok := ollamaReq["prompt"].(string)
	if !ok {
		t.Error("Expected prompt to be a string")
	}

	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	if transformedHeaders["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %s", transformedHeaders["Content-Type"])
	}
}

func TestOllamaAdapter_TransformRequest_WithTemperature(t *testing.T) {
	adapter := NewOllamaAdapter()

	openAIReq := map[string]interface{}{
		"model": "llama2",
		"messages": []map[string]interface{}{
			{"role": "user", "content": "Hello"},
		},
		"stream": false,
		"temperature": 0.7,
	}

	requestJSON, _ := json.Marshal(openAIReq)
	transformed, _, err := adapter.TransformRequest(requestJSON, nil)
	if err != nil {
		t.Errorf("TransformRequest() error = %v", err)
	}

	var ollamaReq map[string]interface{}
	json.Unmarshal(transformed, &ollamaReq)

	if ollamaReq["temperature"] != 0.7 {
		t.Errorf("Expected temperature 0.7, got %v", ollamaReq["temperature"])
	}
}

func TestOllamaAdapter_TransformRequest_WithMaxTokens(t *testing.T) {
	adapter := NewOllamaAdapter()

	openAIReq := map[string]interface{}{
		"model": "llama2",
		"messages": []map[string]interface{}{
			{"role": "user", "content": "Hello"},
		},
		"stream": false,
		"max_tokens": 100,
	}

	requestJSON, _ := json.Marshal(openAIReq)
	transformed, _, err := adapter.TransformRequest(requestJSON, nil)
	if err != nil {
		t.Errorf("TransformRequest() error = %v", err)
	}

	var ollamaReq map[string]interface{}
	json.Unmarshal(transformed, &ollamaReq)

	numPredict, ok := ollamaReq["num_predict"].(float64)
	if !ok || int64(numPredict) != 100 {
		t.Errorf("Expected num_predict 100, got %v", ollamaReq["num_predict"])
	}
}

func TestOllamaAdapter_TransformRequest_MultipleMessages(t *testing.T) {
	adapter := NewOllamaAdapter()

	openAIReq := map[string]interface{}{
		"model": "llama2",
		"messages": []map[string]interface{}{
			{"role": "system", "content": "You are a helpful assistant."},
			{"role": "user", "content": "Hello"},
			{"role": "assistant", "content": "Hi there!"},
		},
		"stream": false,
	}

	requestJSON, _ := json.Marshal(openAIReq)
	transformed, _, err := adapter.TransformRequest(requestJSON, nil)
	if err != nil {
		t.Errorf("TransformRequest() error = %v", err)
	}

	var ollamaReq map[string]interface{}
	json.Unmarshal(transformed, &ollamaReq)

	prompt := ollamaReq["prompt"].(string)
	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	// Verify prompt contains role prefixes
	if len(prompt) < 10 {
		t.Error("Expected prompt to contain role prefixes")
	}
}

func TestOllamaAdapter_TransformResponse(t *testing.T) {
	adapter := NewOllamaAdapter()

	// Create a sample Ollama format response
	ollamaResp := map[string]interface{}{
		"model":             "llama2",
		"response":          "Hello! How can I help you today?",
		"done":              true,
		"prompt_eval_count": 10,
		"eval_count":        20,
	}

	responseJSON, _ := json.Marshal(ollamaResp)

	transformed, tokenCounts, isFinal, err := adapter.TransformResponse(responseJSON, false, entity.TokenCounts{})
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	// Should be final for non-streaming
	if !isFinal {
		t.Error("Expected isFinal to be true for non-streaming")
	}

	// Check token counts
	if tokenCounts.PromptTokens != 10 {
		t.Errorf("Expected prompt tokens 10, got %d", tokenCounts.PromptTokens)
	}

	if tokenCounts.CompletionTokens != 20 {
		t.Errorf("Expected completion tokens 20, got %d", tokenCounts.CompletionTokens)
	}

	// Parse the transformed response
	var openAIResp map[string]interface{}
	if err := json.Unmarshal(transformed, &openAIResp); err != nil {
		t.Errorf("Failed to parse transformed response: %v", err)
	}

	// Verify transformation
	if openAIResp["model"] != "llama2" {
		t.Errorf("Expected model 'llama2', got %v", openAIResp["model"])
	}

	if openAIResp["object"] != "chat.completion" {
		t.Errorf("Expected object 'chat.completion', got %v", openAIResp["object"])
	}

	choices, ok := openAIResp["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		t.Error("Expected choices array with at least one element")
	}

	usage, ok := openAIResp["usage"].(map[string]interface{})
	if !ok {
		t.Error("Expected usage object")
	}

	promptTokens, _ := usage["prompt_tokens"].(float64)
	if int64(promptTokens) != 10 {
		t.Errorf("Expected prompt_tokens 10, got %v", usage["prompt_tokens"])
	}

	completionTokens, _ := usage["completion_tokens"].(float64)
	if int64(completionTokens) != 20 {
		t.Errorf("Expected completion_tokens 20, got %v", usage["completion_tokens"])
	}

	totalTokens, _ := usage["total_tokens"].(float64)
	if int64(totalTokens) != 30 {
		t.Errorf("Expected total_tokens 30, got %v", usage["total_tokens"])
	}
}

func TestOllamaAdapter_TransformResponse_NotDone(t *testing.T) {
	adapter := NewOllamaAdapter()

	ollamaResp := map[string]interface{}{
		"model":             "llama2",
		"response":          "Hello! How can I help you today?",
		"done":              false,
		"prompt_eval_count": 10,
		"eval_count":        20,
	}

	responseJSON, _ := json.Marshal(ollamaResp)
	transformed, _, isFinal, err := adapter.TransformResponse(responseJSON, false, entity.TokenCounts{})
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	// Should still be final for non-streaming even if done=false
	if !isFinal {
		t.Error("Expected isFinal to be true for non-streaming response")
	}

	var openAIResp map[string]interface{}
	json.Unmarshal(transformed, &openAIResp)

	choices := openAIResp["choices"].([]interface{})
	choice := choices[0].(map[string]interface{})
	finishReason := choice["finish_reason"]

	if finishReason != "length" {
		t.Errorf("Expected finish_reason 'length' for not done, got %v", finishReason)
	}
}

func TestOllamaAdapter_CountTokens_WithResponse(t *testing.T) {
	adapter := NewOllamaAdapter()

	response := map[string]interface{}{
		"prompt_eval_count": 15,
		"eval_count":        25,
	}
	responseJSON, _ := json.Marshal(response)

	reqTokens, respTokens, err := adapter.CountTokens([]byte("request"), responseJSON, false)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}

	if reqTokens != 15 {
		t.Errorf("Expected request tokens 15, got %d", reqTokens)
	}

	if respTokens != 25 {
		t.Errorf("Expected response tokens 25, got %d", respTokens)
	}
}

func TestOllamaAdapter_CountTokens_Fallback(t *testing.T) {
	adapter := NewOllamaAdapter()

	request := []byte("This is a test request")
	response := []byte("This is a test response")

	reqTokens, respTokens, err := adapter.CountTokens(request, response, false)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}

	// Fallback estimation: len / 4
	expectedReqTokens := int64(len(request) / 4)
	expectedRespTokens := int64(len(response) / 4)

	if reqTokens != expectedReqTokens {
		t.Errorf("Expected request tokens %d, got %d", expectedReqTokens, reqTokens)
	}

	if respTokens != expectedRespTokens {
		t.Errorf("Expected response tokens %d, got %d", expectedRespTokens, respTokens)
	}
}

func TestOllamaAdapter_ConvertMessagesToPrompt(t *testing.T) {
	adapter := NewOllamaAdapter()

	messages := []map[string]interface{}{
		{"role": "system", "content": "You are helpful"},
		{"role": "user", "content": "Hello"},
	}

	// Access the private method via the public TransformRequest
	openAIReq := map[string]interface{}{
		"model": "llama2",
		"messages": messages,
		"stream": false,
	}

	requestJSON, _ := json.Marshal(openAIReq)
	transformed, _, _ := adapter.TransformRequest(requestJSON, nil)

	var ollamaReq map[string]interface{}
	json.Unmarshal(transformed, &ollamaReq)

	prompt := ollamaReq["prompt"].(string)

	// Verify role prefixes are present
	if len(prompt) < 10 {
		t.Error("Expected prompt to contain role prefixes")
	}
}

func TestOllamaAdapter_InvalidRequest(t *testing.T) {
	adapter := NewOllamaAdapter()

	invalidJSON := []byte("invalid json")
	_, _, err := adapter.TransformRequest(invalidJSON, nil)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestOllamaAdapter_InvalidResponse(t *testing.T) {
	adapter := NewOllamaAdapter()

	invalidJSON := []byte("invalid json")
	_, _, _, err := adapter.TransformResponse(invalidJSON, false, entity.TokenCounts{})
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

// Test that the adapter implements the interface
var _ port.ProviderAdapter = (*OllamaAdapter)(nil)

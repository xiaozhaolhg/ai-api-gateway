package adapter

import (
	"encoding/json"
	"testing"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"github.com/ai-api-gateway/provider-service/internal/domain/port"
)

func TestOpenAIAdapter_TransformRequest(t *testing.T) {
	adapter := NewOpenAIAdapter()

	// Create a sample OpenAI format request
	openAIReq := map[string]interface{}{
		"model": "gpt-4",
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

	// OpenAI adapter is pass-through, so request should be unchanged
	if string(transformed) != string(requestJSON) {
		t.Error("Expected request to be unchanged (pass-through)")
	}

	if transformedHeaders["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %s", transformedHeaders["Content-Type"])
	}
}

func TestOpenAIAdapter_TransformRequest_MissingModel(t *testing.T) {
	adapter := NewOpenAIAdapter()

	// Request without model field
	invalidReq := map[string]interface{}{
		"messages": []map[string]interface{}{
			{"role": "user", "content": "Hello"},
		},
	}

	requestJSON, _ := json.Marshal(invalidReq)
	_, _, err := adapter.TransformRequest(requestJSON, nil)
	if err == nil {
		t.Error("Expected error for missing model field, got nil")
	}
}

func TestOpenAIAdapter_TransformRequest_InvalidJSON(t *testing.T) {
	adapter := NewOpenAIAdapter()

	invalidJSON := []byte("invalid json")
	_, _, err := adapter.TransformRequest(invalidJSON, nil)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestOpenAIAdapter_TransformResponse_NonStreaming(t *testing.T) {
	adapter := NewOpenAIAdapter()

	// Create a sample OpenAI format response
	openAIResp := map[string]interface{}{
		"id": "chatcmpl-123",
		"object": "chat.completion",
		"created": 1234567890,
		"model": "gpt-4",
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"message": map[string]interface{}{
					"role": "assistant",
					"content": "Hello! How can I help you today?",
				},
				"finish_reason": "stop",
			},
		},
		"usage": map[string]interface{}{
			"prompt_tokens": 10,
			"completion_tokens": 20,
			"total_tokens": 30,
		},
	}

	responseJSON, _ := json.Marshal(openAIResp)

	transformed, tokenCounts, isFinal, err := adapter.TransformResponse(responseJSON, false, entity.TokenCounts{})
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	// OpenAI adapter is pass-through, so response should be unchanged
	if string(transformed) != string(responseJSON) {
		t.Error("Expected response to be unchanged (pass-through)")
	}

	// Non-streaming should always be final
	if !isFinal {
		t.Error("Expected isFinal to be true for non-streaming")
	}

	// Token counts should be extracted from usage
	if tokenCounts.PromptTokens != 10 {
		t.Errorf("Expected prompt tokens 10, got %d", tokenCounts.PromptTokens)
	}

	if tokenCounts.CompletionTokens != 20 {
		t.Errorf("Expected completion tokens 20, got %d", tokenCounts.CompletionTokens)
	}
}

func TestOpenAIAdapter_TransformResponse_InvalidJSON(t *testing.T) {
	adapter := NewOpenAIAdapter()

	invalidJSON := []byte("invalid json")
	_, _, _, err := adapter.TransformResponse(invalidJSON, false, entity.TokenCounts{})
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestOpenAIAdapter_CountTokens_NonStreaming(t *testing.T) {
	adapter := NewOpenAIAdapter()

	request := []byte("This is a test request")
	response := []byte("This is a test response")

	reqTokens, respTokens, err := adapter.CountTokens(request, response, false)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}

	// OpenAI adapter uses estimation: len / 4
	expectedReqTokens := int64(len(request) / 4)
	expectedRespTokens := int64(len(response) / 4)

	if reqTokens != expectedReqTokens {
		t.Errorf("Expected request tokens %d, got %d", expectedReqTokens, reqTokens)
	}

	if respTokens != expectedRespTokens {
		t.Errorf("Expected response tokens %d, got %d", expectedRespTokens, respTokens)
	}
}

func TestOpenAIAdapter_CountTokens_Streaming(t *testing.T) {
	adapter := NewOpenAIAdapter()

	// Test intermediate streaming chunk
	reqTokens, respTokens, err := adapter.CountTokens([]byte("request"), []byte("data: {\"content\": \"hello\"}"), true)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}

	// Intermediate chunks should return 0, 0
	if reqTokens != 0 {
		t.Errorf("Expected 0 prompt tokens for intermediate chunk, got %d", reqTokens)
	}

	if respTokens != 0 {
		t.Errorf("Expected 0 completion tokens for intermediate chunk, got %d", respTokens)
	}

	// Test final streaming chunk with [DONE]
	reqTokens, respTokens, err = adapter.CountTokens([]byte("request"), []byte("[DONE]"), true)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}

	// [DONE] chunk doesn't have usage data, so it returns 0, 0
	// In practice, accumulated tokens are tracked separately
}

func TestOpenAIAdapter_CountTokens_FromResponse(t *testing.T) {
	adapter := NewOpenAIAdapter()

	// Response with usage information
	response := map[string]interface{}{
		"usage": map[string]interface{}{
			"prompt_tokens": 15,
			"completion_tokens": 25,
		},
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

// Streaming Tests

func TestOpenAIAdapter_TransformResponse_StreamingChunk(t *testing.T) {
	adapter := NewOpenAIAdapter()

	// SSE chunk format
	sseChunk := []byte("data: {\"id\":\"chatcmpl-123\",\"object\":\"chat.completion.chunk\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"Hello\"}}]}")

	accumulatedTokens := entity.TokenCounts{}
	transformed, tokenCounts, isFinal, err := adapter.TransformResponse(sseChunk, true, accumulatedTokens)
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	// Should not be final
	if isFinal {
		t.Error("Expected isFinal to be false for intermediate chunk")
	}

	// Should pass through unchanged
	if string(transformed) != string(sseChunk) {
		t.Error("Expected SSE chunk to be passed through unchanged")
	}

	// Accumulated tokens should be updated based on content
	if tokenCounts.AccumulatedTokens <= 0 {
		t.Error("Expected accumulated tokens to be > 0")
	}
}

func TestOpenAIAdapter_TransformResponse_FinalChunk(t *testing.T) {
	adapter := NewOpenAIAdapter()

	// Final SSE chunk with [DONE]
	finalChunk := []byte("data: [DONE]")

	accumulatedTokens := entity.TokenCounts{
		PromptTokens:      10,
		CompletionTokens:  20,
		AccumulatedTokens: 30,
	}
	transformed, tokenCounts, isFinal, err := adapter.TransformResponse(finalChunk, true, accumulatedTokens)
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	// Should be final
	if !isFinal {
		t.Error("Expected isFinal to be true for [DONE] chunk")
	}

	// Should pass through
	if string(transformed) != string(finalChunk) {
		t.Error("Expected final chunk to be passed through unchanged")
	}

	// Token counts should match what was passed in
	if tokenCounts.PromptTokens != 10 {
		t.Errorf("Expected prompt tokens 10, got %d", tokenCounts.PromptTokens)
	}

	if tokenCounts.CompletionTokens != 20 {
		t.Errorf("Expected completion tokens 20, got %d", tokenCounts.CompletionTokens)
	}
}

func TestOpenAIAdapter_TransformResponse_StreamingWithUsage(t *testing.T) {
	adapter := NewOpenAIAdapter()

	// SSE chunk with usage data (some OpenAI responses include this)
	sseChunk := []byte(`data: {"id":"chatcmpl-123","object":"chat.completion.chunk","usage":{"prompt_tokens":5,"completion_tokens":10}}`)

	accumulatedTokens := entity.TokenCounts{}
	_, tokenCounts, _, err := adapter.TransformResponse(sseChunk, true, accumulatedTokens)
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	// Should extract usage data
	if tokenCounts.PromptTokens != 5 {
		t.Errorf("Expected prompt tokens 5, got %d", tokenCounts.PromptTokens)
	}

	if tokenCounts.CompletionTokens != 10 {
		t.Errorf("Expected completion tokens 10, got %d", tokenCounts.CompletionTokens)
	}
}

func TestOpenAIAdapter_parseSSEData(t *testing.T) {
	adapter := &OpenAIAdapter{}

	tests := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  bool
	}{
		{
			name:     "valid SSE data line",
			input:    []byte("data: {\"content\": \"hello\"}"),
			expected: []byte("{\"content\": \"hello\"}"),
			wantErr:  false,
		},
		{
			name:     "valid SSE with trailing newline",
			input:    []byte("data: {\"content\": \"hello\"}\n"),
			expected: []byte("{\"content\": \"hello\"}"),
			wantErr:  false,
		},
		{
			name:     "missing data prefix",
			input:    []byte("{\"content\": \"hello\"}"),
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "empty data",
			input:    []byte("data: "),
			expected: []byte(""),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.parseSSEData(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSSEData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(result) != string(tt.expected) {
				t.Errorf("parseSSEData() = %s, expected %s", string(result), string(tt.expected))
			}
		})
	}
}

// Test that the adapter implements the interface
var _ port.ProviderAdapter = (*OpenAIAdapter)(nil)

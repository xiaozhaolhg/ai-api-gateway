package adapter

import (
	"encoding/json"
	"testing"

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

func TestOpenAIAdapter_TransformResponse(t *testing.T) {
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

	transformed, err := adapter.TransformResponse(responseJSON)
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	// OpenAI adapter is pass-through, so response should be unchanged
	if string(transformed) != string(responseJSON) {
		t.Error("Expected response to be unchanged (pass-through)")
	}
}

func TestOpenAIAdapter_TransformResponse_InvalidJSON(t *testing.T) {
	adapter := NewOpenAIAdapter()

	invalidJSON := []byte("invalid json")
	_, err := adapter.TransformResponse(invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestOpenAIAdapter_CountTokens(t *testing.T) {
	adapter := NewOpenAIAdapter()

	request := []byte("This is a test request")
	response := []byte("This is a test response")

	reqTokens, respTokens, err := adapter.CountTokens(request, response)
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

	reqTokens, respTokens, err := adapter.CountTokens([]byte("request"), responseJSON)
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

// Test that the adapter implements the interface
var _ port.ProviderAdapter = (*OpenAIAdapter)(nil)

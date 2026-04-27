package adapter

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"github.com/ai-api-gateway/provider-service/internal/domain/port"
)

func TestAnthropicAdapter_TransformRequest(t *testing.T) {
	adapter := NewAnthropicAdapter()

	// Create an OpenAI format request
	openAIReq := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]interface{}{
			{"role": "system", "content": "You are a helpful assistant."},
			{"role": "user", "content": "Hello, how are you?"},
		},
		"stream":      false,
		"temperature": 0.7,
		"max_tokens":  100,
	}

	requestJSON, _ := json.Marshal(openAIReq)
	headers := map[string]string{"Content-Type": "application/json"}

	transformed, transformedHeaders, err := adapter.TransformRequest(requestJSON, headers)
	if err != nil {
		t.Errorf("TransformRequest() error = %v", err)
	}

	// Parse transformed request
	var anthropicReq map[string]interface{}
	if err := json.Unmarshal(transformed, &anthropicReq); err != nil {
		t.Errorf("Failed to parse transformed request: %v", err)
	}

	// Check model conversion
	if anthropicReq["model"] != "claude-3-opus-20240229" {
		t.Errorf("Expected model 'claude-3-opus-20240229', got %v", anthropicReq["model"])
	}

	// Check messages conversion (system -> user)
	messages, ok := anthropicReq["messages"].([]interface{})
	if !ok || len(messages) != 2 {
		t.Error("Expected 2 messages after conversion")
	}

	// Check headers
	if transformedHeaders["anthropic-version"] != "2023-06-01" {
		t.Errorf("Expected anthropic-version header, got %s", transformedHeaders["anthropic-version"])
	}
}

func TestAnthropicAdapter_TransformRequest_MissingModel(t *testing.T) {
	adapter := NewAnthropicAdapter()

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

func TestAnthropicAdapter_TransformResponse_NonStreaming(t *testing.T) {
	adapter := NewAnthropicAdapter()

	// Create an Anthropic format response
	anthropicResp := map[string]interface{}{
		"id":   "msg_01AbCdEfGhIjKlMn",
		"type": "message",
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "Hello! I'm doing well, thank you for asking. How can I help you today?",
			},
		},
		"model":       "claude-3-opus-20240229",
		"stop_reason": "end_turn",
		"usage": map[string]interface{}{
			"input_tokens":  15,
			"output_tokens": 25,
		},
	}

	responseJSON, _ := json.Marshal(anthropicResp)

	transformed, tokenCounts, isFinal, err := adapter.TransformResponse(responseJSON, false, entity.TokenCounts{})
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	// Parse transformed response
	var openAIResp map[string]interface{}
	if err := json.Unmarshal(transformed, &openAIResp); err != nil {
		t.Errorf("Failed to parse transformed response: %v", err)
	}

	// Check format conversion
	if openAIResp["object"] != "chat.completion" {
		t.Errorf("Expected object 'chat.completion', got %v", openAIResp["object"])
	}

	// Check isFinal
	if !isFinal {
		t.Error("Expected isFinal to be true for non-streaming")
	}

	// Check token counts
	if tokenCounts.PromptTokens != 15 {
		t.Errorf("Expected prompt tokens 15, got %d", tokenCounts.PromptTokens)
	}

	if tokenCounts.CompletionTokens != 25 {
		t.Errorf("Expected completion tokens 25, got %d", tokenCounts.CompletionTokens)
	}
}

func TestAnthropicAdapter_TransformResponse_InvalidJSON(t *testing.T) {
	adapter := NewAnthropicAdapter()

	invalidJSON := []byte("invalid json")
	_, _, _, err := adapter.TransformResponse(invalidJSON, false, entity.TokenCounts{})
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestAnthropicAdapter_CountTokens_NonStreaming(t *testing.T) {
	adapter := NewAnthropicAdapter()

	// Response with usage information
	response := map[string]interface{}{
		"usage": map[string]interface{}{
			"input_tokens":  20,
			"output_tokens": 30,
		},
	}
	responseJSON, _ := json.Marshal(response)

	reqTokens, respTokens, err := adapter.CountTokens([]byte("request"), responseJSON, false)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}

	// Anthropic uses input/output tokens
	if reqTokens != 20 {
		t.Errorf("Expected input tokens 20, got %d", reqTokens)
	}

	if respTokens != 30 {
		t.Errorf("Expected output tokens 30, got %d", respTokens)
	}
}

func TestAnthropicAdapter_CountTokens_Fallback(t *testing.T) {
	adapter := NewAnthropicAdapter()

	// Response without usage information
	response := []byte("This is a test response without usage data")
	request := []byte("This is a test request")

	reqTokens, respTokens, err := adapter.CountTokens(request, response, false)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}

	// Should fall back to estimation
	expectedReqTokens := int64(len(request) / 4)
	expectedRespTokens := int64(len(response) / 4)

	if reqTokens != expectedReqTokens {
		t.Errorf("Expected request tokens %d, got %d", expectedReqTokens, reqTokens)
	}

	if respTokens != expectedRespTokens {
		t.Errorf("Expected response tokens %d, got %d", expectedRespTokens, respTokens)
	}
}

// Streaming Tests

func TestAnthropicAdapter_TransformResponse_StreamingChunk(t *testing.T) {
	adapter := NewAnthropicAdapter()

	// Anthropic SSE chunk format
	sseChunk := []byte("event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"Hello\"}}\n\n")

	accumulatedTokens := entity.TokenCounts{}
	transformed, tokenCounts, isFinal, err := adapter.TransformResponse(sseChunk, true, accumulatedTokens)
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	// Should not be final
	if isFinal {
		t.Error("Expected isFinal to be false for content_block_delta")
	}

	// Should be transformed to OpenAI format
	if !contains(transformed, []byte("chat.completion.chunk")) {
		t.Error("Expected transformed chunk to contain chat.completion.chunk")
	}

	// Accumulated tokens should be updated
	if tokenCounts.AccumulatedTokens <= 0 {
		t.Error("Expected accumulated tokens to be > 0")
	}
}

func TestAnthropicAdapter_TransformResponse_FinalChunk(t *testing.T) {
	adapter := NewAnthropicAdapter()

	// Final SSE chunk with message_stop
	finalChunk := []byte("event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n")

	accumulatedTokens := entity.TokenCounts{
		PromptTokens:      15,
		CompletionTokens:  25,
		AccumulatedTokens: 40,
	}
	_, tokenCounts, isFinal, err := adapter.TransformResponse(finalChunk, true, accumulatedTokens)
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	// Should be final
	if !isFinal {
		t.Error("Expected isFinal to be true for message_stop event")
	}

	// Token counts should be preserved
	if tokenCounts.PromptTokens != 15 {
		t.Errorf("Expected prompt tokens 15, got %d", tokenCounts.PromptTokens)
	}
}

func TestAnthropicAdapter_TransformResponse_MessageDelta(t *testing.T) {
	adapter := NewAnthropicAdapter()

	// Message delta with usage data
	sseChunk := []byte("event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\",\"usage\":{\"output_tokens\":50}},\"usage\":{\"output_tokens\":50}}\n\n")

	accumulatedTokens := entity.TokenCounts{PromptTokens: 20}
	_, tokenCounts, _, err := adapter.TransformResponse(sseChunk, true, accumulatedTokens)
	if err != nil {
		t.Errorf("TransformResponse() error = %v", err)
	}

	// Should update completion tokens from usage
	if tokenCounts.CompletionTokens != 50 {
		t.Errorf("Expected completion tokens 50, got %d", tokenCounts.CompletionTokens)
	}
}

func TestAnthropicAdapter_CountTokens_Streaming(t *testing.T) {
	adapter := NewAnthropicAdapter()

	// Test intermediate streaming chunk
	reqTokens, respTokens, err := adapter.CountTokens([]byte("request"), []byte("data: content"), true)
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

	// Test final streaming chunk with message_stop
	reqTokens, respTokens, err = adapter.CountTokens([]byte("request"), []byte("event: message_stop"), true)
	if err != nil {
		t.Errorf("CountTokens() error = %v", err)
	}

	// Final chunk without usage returns 0, 0
	if reqTokens != 0 || respTokens != 0 {
		t.Errorf("Expected 0 tokens for final chunk without usage, got %d, %d", reqTokens, respTokens)
	}
}

func TestAnthropicAdapter_convertModelName(t *testing.T) {
	adapter := &AnthropicAdapter{}

	tests := []struct {
		input    string
		expected string
	}{
		{"gpt-4", "claude-3-opus-20240229"},
		{"gpt-4-turbo", "claude-3-sonnet-20240229"},
		{"gpt-3.5-turbo", "claude-3-haiku-20240307"},
		{"unknown-model", "unknown-model"}, // passthrough
	}

	for _, tt := range tests {
		result := adapter.convertModelName(tt.input)
		if result != tt.expected {
			t.Errorf("convertModelName(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestAnthropicAdapter_convertMessages(t *testing.T) {
	adapter := &AnthropicAdapter{}

	messages := []map[string]interface{}{
		{"role": "system", "content": "You are helpful."},
		{"role": "user", "content": "Hello"},
		{"role": "assistant", "content": "Hi there!"},
	}

	result := adapter.convertMessages(messages)

	// Check system role converted to user
	if len(result) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(result))
	}

	// System should be converted to user
	firstRole := result[0]["role"].(string)
	if firstRole != "user" {
		t.Errorf("Expected system converted to user, got %s", firstRole)
	}

	// User and assistant should remain
	if result[1]["role"] != "user" {
		t.Error("Expected user role preserved")
	}
	if result[2]["role"] != "assistant" {
		t.Error("Expected assistant role preserved")
	}
}

func TestAnthropicAdapter_convertStopReason(t *testing.T) {
	adapter := &AnthropicAdapter{}

	tests := []struct {
		input    *string
		expected string
	}{
		{nil, "stop"},
		{strPtr("end_turn"), "stop"},
		{strPtr("max_tokens"), "length"},
		{strPtr("stop_sequence"), "stop"},
		{strPtr("unknown"), "unknown"},
	}

	for _, tt := range tests {
		result := adapter.convertStopReason(tt.input)
		if result != tt.expected {
			t.Errorf("convertStopReason(%v) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestAnthropicAdapter_extractContent(t *testing.T) {
	adapter := &AnthropicAdapter{}

	content := []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}{
		{Type: "text", Text: "Hello "},
		{Type: "text", Text: "world!"},
		{Type: "image", Text: "should be ignored"},
	}

	result := adapter.extractContent(content)
	if result != "Hello world!" {
		t.Errorf("Expected 'Hello world!', got '%s'", result)
	}
}

func TestAnthropicAdapter_parseSSEData(t *testing.T) {
	adapter := &AnthropicAdapter{}

	tests := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  bool
	}{
		{
			name:     "valid SSE data line",
			input:    []byte("event: test\ndata: {\"content\": \"hello\"}\n\n"),
			expected: []byte("{\"content\": \"hello\"}"),
			wantErr:  false,
		},
		{
			name:     "missing data prefix",
			input:    []byte("event: test\n{\"content\": \"hello\"}"),
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "data only line",
			input:    []byte("data: {\"test\": true}\n"),
			expected: []byte("{\"test\": true}"),
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

func TestAnthropicAdapter_extractSSEEventType(t *testing.T) {
	adapter := &AnthropicAdapter{}

	tests := []struct {
		input         []byte
		expectedType  string
		expectedFound bool
	}{
		{[]byte("event: content_block_delta\ndata: {}"), "content_block_delta", true},
		{[]byte("event: message_stop\n"), "message_stop", true},
		{[]byte("data: {}"), "", false},
		{[]byte(""), "", false},
	}

	for _, tt := range tests {
		eventType, found := adapter.extractSSEEventType(tt.input)
		if found != tt.expectedFound {
			t.Errorf("extractSSEEventType(%s) found = %v, expected %v", string(tt.input), found, tt.expectedFound)
		}
		if eventType != tt.expectedType {
			t.Errorf("extractSSEEventType(%s) = %s, expected %s", string(tt.input), eventType, tt.expectedType)
		}
	}
}

// Helper functions

func strPtr(s string) *string {
	return &s
}

func contains(haystack, needle []byte) bool {
	return bytes.Contains(haystack, needle)
}

// Test that the adapter implements the interface
var _ port.ProviderAdapter = (*AnthropicAdapter)(nil)

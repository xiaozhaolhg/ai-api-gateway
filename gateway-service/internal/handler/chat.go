package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// ChatHandler handles chat completion requests
type ChatHandler struct{}

// NewChatHandler creates a new chat handler
func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

// ServeHTTP handles HTTP requests for /v1/chat/completions
func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var reqBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if streaming is requested
	stream, _ := reqBody["stream"].(bool)

	if stream {
		h.handleStreamingRequest(w, r, reqBody)
	} else {
		h.handleNonStreamingRequest(w, r, reqBody)
	}
}

// handleNonStreamingRequest handles non-streaming chat completion requests
func (h *ChatHandler) handleNonStreamingRequest(w http.ResponseWriter, r *http.Request, reqBody map[string]interface{}) {
	// Add model to context for route resolution
	ctx := context.WithValue(r.Context(), "model", reqBody["model"])

	// The actual proxying is handled by the ProxyMiddleware
	// This handler just validates the request and passes it through
	// In a real implementation, you'd call the provider-service here

	// For MVP, return a placeholder response
	response := map[string]interface{}{
		"id":      "chatcmpl-placeholder",
		"object":  "chat.completion",
		"created": 0,
		"model":   reqBody["model"],
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": "This is a placeholder response. The actual implementation would proxy to provider-service.",
				},
				"finish_reason": "stop",
			},
		},
		"usage": map[string]interface{}{
			"prompt_tokens":     0,
			"completion_tokens": 0,
			"total_tokens":      0,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleStreamingRequest handles streaming chat completion requests
func (h *ChatHandler) handleStreamingRequest(w http.ResponseWriter, r *http.Request, reqBody map[string]interface{}) {
	// Add model to context for route resolution
	ctx := context.WithValue(r.Context(), "model", reqBody["model"])

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// For MVP, send a single placeholder chunk
	chunk := map[string]interface{}{
		"id":      "chatcmpl-placeholder",
		"object":  "chat.completion.chunk",
		"created": 0,
		"model":   reqBody["model"],
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"delta": map[string]interface{}{
					"role":    "assistant",
					"content": "This is a placeholder streaming response. ",
				},
				"finish_reason": nil,
			},
		},
	}

	chunkJSON, _ := json.Marshal(chunk)
	w.Write([]byte("data: " + string(chunkJSON) + "\n\n"))
	flusher.Flush()

	// Send final chunk
	finalChunk := map[string]interface{}{
		"id":      "chatcmpl-placeholder",
		"object":  "chat.completion.chunk",
		"created": 0,
		"model":   reqBody["model"],
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"delta": map[string]interface{}{},
				"finish_reason": "stop",
			},
		},
	}

	finalChunkJSON, _ := json.Marshal(finalChunk)
	w.Write([]byte("data: " + string(finalChunkJSON) + "\n\n"))
	w.Write([]byte("data: [DONE]\n\n"))
	flusher.Flush()
}

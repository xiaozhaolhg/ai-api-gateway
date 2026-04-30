# Provider Adapter Development Guide

This guide explains how to implement a provider adapter for the AI API Gateway. Provider adapters translate between the gateway's internal format and provider-specific APIs.

## Table of Contents

1. [Overview](#overview)
2. [Interface Definition](#interface-definition)
3. [Step-by-Step Implementation](#step-by-step-implementation)
4. [Testing Guidelines](#testing-guidelines)
5. [OpenAI Adapter Example](#openai-adapter-example)

## Overview

### What is a Provider Adapter?

A provider adapter is a component that:
- Translates gateway requests to provider-specific formats
- Handles provider authentication
- Manages response transformation
- Supports streaming (SSE) responses
- Counts tokens for billing

### Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Gateway   │────▶│    Adapter   │────▶│  Provider   │
│   Service   │◄────│  (You Build) │◄────│    API      │
└─────────────┘     └──────────────┘     └─────────────┘
```

## Interface Definition

### Required Interface

```go
// ProviderAdapter defines the interface for all provider adapters
type ProviderAdapter interface {
    // ForwardRequest handles non-streaming requests
    ForwardRequest(ctx context.Context, req *ForwardRequest) (*ProviderResponse, error)
    
    // StreamRequest handles streaming requests
    StreamRequest(ctx context.Context, req *ForwardRequest) (<-chan StreamChunk, error)
    
    // HealthCheck verifies provider connectivity
    HealthCheck(ctx context.Context) error
}

// ForwardRequest contains the request data
type ForwardRequest struct {
    Model       string
    Messages    []Message
    Parameters  map[string]interface{}
    Headers     map[string]string
}

// ProviderResponse contains the response data
type ProviderResponse struct {
    Content        string
    PromptTokens     int64
    CompletionTokens int64
    FinishReason     string
}

// StreamChunk represents a single streaming chunk
type StreamChunk struct {
    Content          string
    PromptTokens     int64
    CompletionTokens int64
    Done             bool
    Error            error
}
```

## Step-by-Step Implementation

### 1. Create Adapter Structure

```go
package adapters

import (
    "context"
    "net/http"
    "time"
)

// MyProviderAdapter implements the ProviderAdapter interface
type MyProviderAdapter struct {
    baseURL     string
    apiKey      string
    httpClient  *http.Client
    modelMap    map[string]string // Maps gateway model IDs to provider model IDs
}

// NewMyProviderAdapter creates a new adapter instance
func NewMyProviderAdapter(baseURL, apiKey string) *MyProviderAdapter {
    return &MyProviderAdapter{
        baseURL: baseURL,
        apiKey:  apiKey,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        modelMap: map[string]string{
            "myprovider:gpt-4": "gpt-4",
            "myprovider:gpt-3.5": "gpt-3.5-turbo",
        },
    }
}
```

### 2. Implement ForwardRequest

```go
func (a *MyProviderAdapter) ForwardRequest(
    ctx context.Context, 
    req *ForwardRequest,
) (*ProviderResponse, error) {
    // 1. Map model ID
    providerModel := a.modelMap[req.Model]
    if providerModel == "" {
        providerModel = req.Model // Fallback to original
    }
    
    // 2. Build provider-specific request
    providerReq := map[string]interface{}{
        "model": providerModel,
        "messages": a.convertMessages(req.Messages),
    }
    
    // 3. Add parameters
    for k, v := range req.Parameters {
        providerReq[k] = v
    }
    
    // 4. Send request
    resp, err := a.sendRequest(ctx, providerReq)
    if err != nil {
        return nil, fmt.Errorf("provider request failed: %w", err)
    }
    
    // 5. Transform response
    return a.transformResponse(resp), nil
}

func (a *MyProviderAdapter) sendRequest(
    ctx context.Context, 
    body map[string]interface{},
) (map[string]interface{}, error) {
    jsonBody, _ := json.Marshal(body)
    
    httpReq, _ := http.NewRequestWithContext(
        ctx,
        "POST",
        a.baseURL+"/v1/chat/completions",
        bytes.NewReader(jsonBody),
    )
    
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)
    
    resp, err := a.httpClient.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("provider returned %d", resp.StatusCode)
    }
    
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result, nil
}
```

### 3. Implement StreamRequest

```go
func (a *MyProviderAdapter) StreamRequest(
    ctx context.Context, 
    req *ForwardRequest,
) (<-chan StreamChunk, error) {
    // 1. Prepare streaming request
    providerReq := map[string]interface{}{
        "model":    a.modelMap[req.Model],
        "messages": a.convertMessages(req.Messages),
        "stream":   true,
    }
    
    // 2. Send request
    jsonBody, _ := json.Marshal(providerReq)
    httpReq, _ := http.NewRequestWithContext(
        ctx,
        "POST",
        a.baseURL+"/v1/chat/completions",
        bytes.NewReader(jsonBody),
    )
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)
    httpReq.Header.Set("Accept", "text/event-stream")
    
    resp, err := a.httpClient.Do(httpReq)
    if err != nil {
        return nil, err
    }
    
    // 3. Create response channel
    chunkChan := make(chan StreamChunk)
    
    // 4. Start goroutine to process stream
    go func() {
        defer resp.Body.Close()
        defer close(chunkChan)
        
        reader := bufio.NewReader(resp.Body)
        var promptTokens, completionTokens int64
        
        for {
            line, err := reader.ReadString('\n')
            if err != nil {
                if err == io.EOF {
                    break
                }
                chunkChan <- StreamChunk{Error: err}
                return
            }
            
            // Parse SSE format: "data: {...}"
            if !strings.HasPrefix(line, "data: ") {
                continue
            }
            
            data := strings.TrimPrefix(line, "data: ")
            if data == "[DONE]" {
                chunkChan <- StreamChunk{Done: true}
                break
            }
            
            // Parse chunk
            var chunk map[string]interface{}
            if err := json.Unmarshal([]byte(data), &chunk); err != nil {
                continue
            }
            
            // Extract content and tokens
            content := a.extractContent(chunk)
            promptTokens += a.countPromptTokens(chunk)
            completionTokens += a.countCompletionTokens(chunk)
            
            chunkChan <- StreamChunk{
                Content:          content,
                PromptTokens:     promptTokens,
                CompletionTokens: completionTokens,
            }
        }
    }()
    
    return chunkChan, nil
}
```

### 4. Implement HealthCheck

```go
func (a *MyProviderAdapter) HealthCheck(ctx context.Context) error {
    // Simple check: verify we can reach the provider
    req, _ := http.NewRequestWithContext(
        ctx,
        "GET",
        a.baseURL+"/v1/models",
        nil,
    )
    req.Header.Set("Authorization", "Bearer "+a.apiKey)
    
    resp, err := a.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("provider unreachable: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("provider unhealthy: status %d", resp.StatusCode)
    }
    
    return nil
}
```

## Testing Guidelines

### Unit Tests

```go
func TestMyProviderAdapter_ForwardRequest(t *testing.T) {
    // Create mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate request
        if r.Header.Get("Authorization") == "" {
            t.Error("missing authorization header")
        }
        
        // Return mock response
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "choices": []map[string]interface{}{
                {
                    "message": map[string]interface{}{
                        "content": "Hello!",
                    },
                    "finish_reason": "stop",
                },
            },
            "usage": map[string]interface{}{
                "prompt_tokens":     10,
                "completion_tokens": 5,
            },
        })
    }))
    defer server.Close()
    
    // Create adapter
    adapter := NewMyProviderAdapter(server.URL, "test-key")
    
    // Test request
    req := &ForwardRequest{
        Model: "myprovider:gpt-4",
        Messages: []Message{
            {Role: "user", Content: "Hello"},
        },
    }
    
    resp, err := adapter.ForwardRequest(context.Background(), req)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if resp.Content != "Hello!" {
        t.Errorf("expected 'Hello!', got %s", resp.Content)
    }
}
```

### Integration Tests

```go
func TestMyProviderAdapter_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Use real provider API (with test key)
    adapter := NewMyProviderAdapter(
        "https://api.provider.com",
        os.Getenv("PROVIDER_API_KEY"),
    )
    
    // Test health check
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := adapter.HealthCheck(ctx); err != nil {
        t.Fatalf("health check failed: %v", err)
    }
}
```

## OpenAI Adapter Example

Here's a complete, production-ready OpenAI adapter:

```go
package adapters

import (
    "bufio"
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

// OpenAIAdapter implements the ProviderAdapter interface for OpenAI
type OpenAIAdapter struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
    modelMap   map[string]string
}

// NewOpenAIAdapter creates a new OpenAI adapter
func NewOpenAIAdapter(apiKey string) *OpenAIAdapter {
    return &OpenAIAdapter{
        baseURL: "https://api.openai.com",
        apiKey:  apiKey,
        httpClient: &http.Client{
            Timeout: 60 * time.Second,
        },
        modelMap: map[string]string{
            "openai:gpt-4":         "gpt-4",
            "openai:gpt-4-turbo":   "gpt-4-turbo-preview",
            "openai:gpt-3.5":       "gpt-3.5-turbo",
            "openai:gpt-3.5-16k":   "gpt-3.5-turbo-16k",
        },
    }
}

// ForwardRequest implements non-streaming chat completion
func (a *OpenAIAdapter) ForwardRequest(
    ctx context.Context,
    req *ForwardRequest,
) (*ProviderResponse, error) {
    // Map model
    model := a.modelMap[req.Model]
    if model == "" {
        model = req.Model
    }
    
    // Build request body
    body := map[string]interface{}{
        "model":       model,
        "messages":    a.convertMessages(req.Messages),
        "temperature": getFloatParam(req.Parameters, "temperature", 0.7),
        "max_tokens":  getIntParam(req.Parameters, "max_tokens", 150),
    }
    
    // Add optional parameters
    if topP, ok := req.Parameters["top_p"]; ok {
        body["top_p"] = topP
    }
    
    // Send request
    resp, err := a.doRequest(ctx, "/v1/chat/completions", body)
    if err != nil {
        return nil, err
    }
    
    // Parse response
    return a.parseResponse(resp)
}

// StreamRequest implements streaming chat completion
func (a *OpenAIAdapter) StreamRequest(
    ctx context.Context,
    req *ForwardRequest,
) (<-chan StreamChunk, error) {
    model := a.modelMap[req.Model]
    if model == "" {
        model = req.Model
    }
    
    body := map[string]interface{}{
        "model":       model,
        "messages":    a.convertMessages(req.Messages),
        "temperature": getFloatParam(req.Parameters, "temperature", 0.7),
        "max_tokens":  getIntParam(req.Parameters, "max_tokens", 500),
        "stream":      true,
    }
    
    jsonBody, _ := json.Marshal(body)
    httpReq, _ := http.NewRequestWithContext(
        ctx,
        "POST",
        a.baseURL+"/v1/chat/completions",
        bytes.NewReader(jsonBody),
    )
    
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)
    httpReq.Header.Set("Accept", "text/event-stream")
    
    resp, err := a.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    
    if resp.StatusCode != http.StatusOK {
        resp.Body.Close()
        return nil, fmt.Errorf("provider returned %d", resp.StatusCode)
    }
    
    return a.processStream(resp.Body), nil
}

// HealthCheck verifies OpenAI API connectivity
func (a *OpenAIAdapter) HealthCheck(ctx context.Context) error {
    req, _ := http.NewRequestWithContext(
        ctx,
        "GET",
        a.baseURL+"/v1/models",
        nil,
    )
    req.Header.Set("Authorization", "Bearer "+a.apiKey)
    
    resp, err := a.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("openai api unreachable: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("openai api returned %d", resp.StatusCode)
    }
    
    return nil
}

// Helper methods

func (a *OpenAIAdapter) doRequest(
    ctx context.Context,
    path string,
    body map[string]interface{},
) (map[string]interface{}, error) {
    jsonBody, _ := json.Marshal(body)
    
    req, _ := http.NewRequestWithContext(
        ctx,
        "POST",
        a.baseURL+path,
        bytes.NewReader(jsonBody),
    )
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+a.apiKey)
    
    resp, err := a.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("openai returned %d: %s", resp.StatusCode, string(body))
    }
    
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result, nil
}

func (a *OpenAIAdapter) parseResponse(resp map[string]interface{}) (*ProviderResponse, error) {
    choices, ok := resp["choices"].([]interface{})
    if !ok || len(choices) == 0 {
        return nil, fmt.Errorf("no choices in response")
    }
    
    choice, ok := choices[0].(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("invalid choice format")
    }
    
    message, ok := choice["message"].(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("no message in choice")
    }
    
    content, _ := message["content"].(string)
    finishReason, _ := choice["finish_reason"].(string)
    
    // Extract usage
    usage, _ := resp["usage"].(map[string]interface{})
    promptTokens := getInt64FromInterface(usage["prompt_tokens"])
    completionTokens := getInt64FromInterface(usage["completion_tokens"])
    
    return &ProviderResponse{
        Content:          content,
        PromptTokens:     promptTokens,
        CompletionTokens: completionTokens,
        FinishReason:     finishReason,
    }, nil
}

func (a *OpenAIAdapter) processStream(body io.ReadCloser) <-chan StreamChunk {
    chunkChan := make(chan StreamChunk)
    
    go func() {
        defer body.Close()
        defer close(chunkChan)
        
        reader := bufio.NewReader(body)
        var content strings.Builder
        var promptTokens, completionTokens int64
        
        for {
            line, err := reader.ReadString('\n')
            if err != nil {
                if err == io.EOF {
                    break
                }
                chunkChan <- StreamChunk{Error: err}
                return
            }
            
            line = strings.TrimSpace(line)
            if !strings.HasPrefix(line, "data: ") {
                continue
            }
            
            data := strings.TrimPrefix(line, "data: ")
            if data == "[DONE]" {
                break
            }
            
            var chunk map[string]interface{}
            if err := json.Unmarshal([]byte(data), &chunk); err != nil {
                continue
            }
            
            // Extract delta content
            if choices, ok := chunk["choices"].([]interface{}); ok && len(choices) > 0 {
                if choice, ok := choices[0].(map[string]interface{}); ok {
                    if delta, ok := choice["delta"].(map[string]interface{}); ok {
                        if deltaContent, ok := delta["content"].(string); ok {
                            content.WriteString(deltaContent)
                        }
                    }
                    
                    // Check if done
                    finishReason, _ := choice["finish_reason"].(string)
                    if finishReason != "" {
                        chunkChan <- StreamChunk{
                            Content:          content.String(),
                            PromptTokens:     promptTokens,
                            CompletionTokens: completionTokens,
                            Done:             true,
                        }
                        return
                    }
                }
            }
            
            // Estimate tokens (OpenAI doesn't send usage in stream)
            completionTokens = int64(len(content.String())) / 4 // Rough estimate
            
            chunkChan <- StreamChunk{
                Content:          content.String(),
                PromptTokens:     promptTokens,
                CompletionTokens: completionTokens,
            }
        }
        
        // Final chunk
        chunkChan <- StreamChunk{
            Content:          content.String(),
            PromptTokens:     promptTokens,
            CompletionTokens: completionTokens,
            Done:             true,
        }
    }()
    
    return chunkChan
}

func (a *OpenAIAdapter) convertMessages(msgs []Message) []map[string]string {
    result := make([]map[string]string, len(msgs))
    for i, msg := range msgs {
        result[i] = map[string]string{
            "role":    msg.Role,
            "content": msg.Content,
        }
    }
    return result
}

// Utility functions

func getFloatParam(params map[string]interface{}, key string, defaultVal float64) float64 {
    if v, ok := params[key]; ok {
        if f, ok := v.(float64); ok {
            return f
        }
    }
    return defaultVal
}

func getIntParam(params map[string]interface{}, key string, defaultVal int) int {
    if v, ok := params[key]; ok {
        if i, ok := v.(float64); ok {
            return int(i)
        }
    }
    return defaultVal
}

func getInt64FromInterface(v interface{}) int64 {
    if v == nil {
        return 0
    }
    if f, ok := v.(float64); ok {
        return int64(f)
    }
    return 0
}

// Message represents a chat message
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}
```

## Best Practices

1. **Always set timeouts**: Use `context.WithTimeout()` for all external calls
2. **Handle errors gracefully**: Return meaningful error messages
3. **Log requests/responses**: For debugging and monitoring
4. **Validate inputs**: Check model IDs and parameters before sending
5. **Support streaming**: Most modern LLMs support streaming responses
6. **Count tokens accurately**: Essential for billing
7. **Use connection pooling**: Reuse HTTP clients
8. **Implement retries**: With exponential backoff for transient failures

## Registering Your Adapter

Once implemented, register your adapter in the provider-service:

```go
// In provider-service/adapters/registry.go
func init() {
    Register("openai", func(config map[string]string) ProviderAdapter {
        return NewOpenAIAdapter(config["api_key"])
    })
    
    Register("myprovider", func(config map[string]string) ProviderAdapter {
        return NewMyProviderAdapter(config["base_url"], config["api_key"])
    })
}
```

## Next Steps

1. Implement your adapter following this guide
2. Write unit tests for all functions
3. Test against real provider API
4. Submit PR with your adapter
5. Update documentation

For questions, see the [project wiki](https://github.com/ai-api-gateway/wiki) or open an issue.

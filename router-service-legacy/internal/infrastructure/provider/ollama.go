package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	"github.com/ai-api-gateway/router-service/internal/domain/port"
	"github.com/ai-api-gateway/router-service/internal/infrastructure/config"
)

type OllamaFactory struct{}

func NewOllamaFactory() *OllamaFactory {
	return &OllamaFactory{}
}

func (f *OllamaFactory) Type() string {
	return "ollama"
}

func (f *OllamaFactory) Description() string {
	return "Ollama - Run LLMs locally"
}

func (f *OllamaFactory) Validate(settings config.ProviderSettings) error {
	if settings.Endpoint == "" {
		return fmt.Errorf("endpoint is required for Ollama provider")
	}
	return nil
}

func (f *OllamaFactory) Defaults() config.ProviderSettings {
	return config.ProviderSettings{
		Endpoint: "http://localhost:11434",
		Enabled:  false,
		APIKey:   "",
	}
}

func (f *OllamaFactory) Create(settings config.ProviderSettings) (port.Provider, error) {
	return NewOllamaProvider(settings), nil
}

type OllamaProvider struct {
	httpClient *http.Client
	endpoint   string
}

func NewOllamaProvider(settings config.ProviderSettings) *OllamaProvider {
	return &OllamaProvider{
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		endpoint: settings.BaseURL(),
	}
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

func (p *OllamaProvider) ChatCompletion(req port.ChatCompletionRequest) (*port.ChatCompletionResponse, error) {
	ollamaReq := map[string]interface{}{
		"model":    extractModelName(req.Model),
		"messages": transformMessages(req.Messages),
		"stream":   false,
	}

	if req.Temperature > 0 {
		ollamaReq["temperature"] = req.Temperature
	}
	if req.MaxTokens > 0 {
		ollamaReq["options"] = map[string]interface{}{"num_predict": req.MaxTokens}
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", p.endpoint+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var ollamaResp struct {
		Message      MessageContent `json:"message"`
		Done         bool            `json:"done"`
		PromptEvalCount int          `json:"prompt_eval_count"`
		EvalCount    int             `json:"eval_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, err
	}

	return &port.ChatCompletionResponse{
		ID:      generateID(),
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []port.Choice{
			{
				Index: 0,
				Message: port.Message{
					Role:    "assistant",
					Content: ollamaResp.Message.Content,
				},
				FinishReason: "stop",
			},
		},
		Usage: port.Usage{
			PromptTokens:     ollamaResp.PromptEvalCount,
			CompletionTokens:  ollamaResp.EvalCount,
			TotalTokens:       ollamaResp.PromptEvalCount + ollamaResp.EvalCount,
		},
	}, nil
}

func (p *OllamaProvider) StreamChatCompletion(req port.ChatCompletionRequest) (<-chan port.StreamChunk, error) {
	ollamaReq := map[string]interface{}{
		"model":    extractModelName(req.Model),
		"messages": transformMessages(req.Messages),
		"stream":   true,
	}

	if req.Temperature > 0 {
		ollamaReq["temperature"] = req.Temperature
	}
	if req.MaxTokens > 0 {
		ollamaReq["options"] = map[string]interface{}{"num_predict": req.MaxTokens}
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", p.endpoint+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama stream request failed: %w", err)
	}

	ch := make(chan port.StreamChunk, 10)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		decoder := json.NewDecoder(resp.Body)
		for {
			var chunk struct {
				Message MessageContent `json:"message"`
				Done    bool            `json:"done"`
			}
			if err := decoder.Decode(&chunk); err != nil {
				if err == io.EOF {
					break
				}
				return
			}

			ch <- port.StreamChunk{
				ID:      generateID(),
				Created: time.Now().Unix(),
				Model:   req.Model,
				Choices: []port.Choice{
					{
						Delta: &port.Message{
							Role:    "assistant",
							Content: chunk.Message.Content,
						},
					},
				},
				Done: chunk.Done,
			}

			if chunk.Done {
				break
			}
		}
	}()

	return ch, nil
}

func (p *OllamaProvider) ListModels() ([]entity.Model, error) {
	httpReq, err := http.NewRequestWithContext(context.Background(), "GET", p.endpoint+"/api/tags", nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	models := make([]entity.Model, len(result.Models))
	for i, m := range result.Models {
		models[i] = entity.Model{
			ID:       "ollama:" + m.Name,
			Name:     m.Name,
			Provider: "ollama",
		}
	}

	return models, nil
}

type MessageContent struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func transformMessages(msgs []port.Message) []MessageContent {
	result := make([]MessageContent, len(msgs))
	for i, m := range msgs {
		result[i] = MessageContent{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	return result
}

func extractModelName(model string) string {
	if len(model) > 7 && model[:7] == "ollama:" {
		return model[7:]
	}
	return model
}

func generateID() string {
	return fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()/1000)
}
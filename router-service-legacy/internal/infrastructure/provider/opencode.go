package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ai-api-gateway/router-service-legacy/internal/domain/entity"
	"github.com/ai-api-gateway/router-service-legacy/internal/domain/port"
	"github.com/ai-api-gateway/router-service-legacy/internal/infrastructure/config"
)

type OpenCodeZenFactory struct{}

func NewOpenCodeZenFactory() *OpenCodeZenFactory {
	return &OpenCodeZenFactory{}
}

func (f *OpenCodeZenFactory) Type() string {
	return "opencode_zen"
}

func (f *OpenCodeZenFactory) Description() string {
	return "OpenCode Zen - Curated AI models for coding agents"
}

func (f *OpenCodeZenFactory) Validate(settings config.ProviderSettings) error {
	if settings.Endpoint == "" {
		return fmt.Errorf("endpoint is required for OpenCode Zen provider")
	}
	return nil
}

func (f *OpenCodeZenFactory) Defaults() config.ProviderSettings {
	return config.ProviderSettings{
		Endpoint: "https://opencode.ai/zen",
		Enabled:  false,
		APIKey:   "",
	}
}

func (f *OpenCodeZenFactory) Create(settings config.ProviderSettings) (port.Provider, error) {
	return NewOpenCodeZenProvider(settings), nil
}

type OpenCodeZenProvider struct {
	httpClient *http.Client
	endpoint   string
	apiKey     string
}

func NewOpenCodeZenProvider(settings config.ProviderSettings) *OpenCodeZenProvider {
	return &OpenCodeZenProvider{
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		endpoint: settings.BaseURL(),
		apiKey:   settings.APIKey,
	}
}

func (p *OpenCodeZenProvider) Name() string {
	return "opencode_zen"
}

func (p *OpenCodeZenProvider) ChatCompletion(req port.ChatCompletionRequest) (*port.ChatCompletionResponse, error) {
	openaiReq := map[string]interface{}{
		"model":       extractOpenCodeModelName(req.Model),
		"messages":    req.Messages,
		"stream":      false,
		"temperature": req.Temperature,
	}

	if req.MaxTokens > 0 {
		openaiReq["max_tokens"] = req.MaxTokens
	}
	if req.TopP > 0 {
		openaiReq["top_p"] = req.TopP
	}

	body, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", p.endpoint+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("opencode request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("opencode returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var openaiResp struct {
		ID      string `json:"id"`
		Created int64  `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return nil, err
	}

	choices := make([]port.Choice, len(openaiResp.Choices))
	for i, c := range openaiResp.Choices {
		choices[i] = port.Choice{
			Index: i,
			Message: port.Message{
				Role:    c.Message.Role,
				Content: c.Message.Content,
			},
			FinishReason: c.FinishReason,
		}
	}

	return &port.ChatCompletionResponse{
		ID:      openaiResp.ID,
		Created: openaiResp.Created,
		Model:   openaiResp.Model,
		Choices: choices,
		Usage: port.Usage{
			PromptTokens:     openaiResp.Usage.PromptTokens,
			CompletionTokens: openaiResp.Usage.CompletionTokens,
			TotalTokens:      openaiResp.Usage.TotalTokens,
		},
	}, nil
}

func (p *OpenCodeZenProvider) StreamChatCompletion(req port.ChatCompletionRequest) (<-chan port.StreamChunk, error) {
	openaiReq := map[string]interface{}{
		"model":       extractOpenCodeModelName(req.Model),
		"messages":    req.Messages,
		"stream":      true,
		"temperature": req.Temperature,
	}

	if req.MaxTokens > 0 {
		openaiReq["max_tokens"] = req.MaxTokens
	}
	if req.TopP > 0 {
		openaiReq["top_p"] = req.TopP
	}

	body, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", p.endpoint+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("opencode stream request failed: %w", err)
	}

	ch := make(chan port.StreamChunk, 10)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) < 6 || line[:6] != "data: " {
				continue
			}

			if line == "data: [DONE]" {
				ch <- port.StreamChunk{Done: true}
				break
			}

			var delta struct {
				ID      string `json:"id"`
				Created int64  `json:"created"`
				Model   string `json:"model"`
				Choices []struct {
					Delta struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(line[6:]), &delta); err != nil {
				continue
			}

			if len(delta.Choices) > 0 {
				ch <- port.StreamChunk{
					ID:      delta.ID,
					Created: delta.Created,
					Model:   delta.Model,
					Choices: []port.Choice{
						{
							Delta: &port.Message{
								Role:    delta.Choices[0].Delta.Role,
								Content: delta.Choices[0].Delta.Content,
							},
						},
					},
				}
			}
		}
	}()

	return ch, nil
}

func (p *OpenCodeZenProvider) ListModels() ([]entity.Model, error) {
	httpReq, err := http.NewRequestWithContext(context.Background(), "GET", p.endpoint+"/v1/models", nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("opencode returned status %d", resp.StatusCode)
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	models := make([]entity.Model, len(result.Data))
	for i, m := range result.Data {
		models[i] = entity.Model{
			ID:       "opencode_zen:" + m.ID,
			Name:     m.ID,
			Provider: "opencode_zen",
		}
	}

	return models, nil
}

func extractOpenCodeModelName(model string) string {
	if len(model) > 13 && model[:13] == "opencode_zen:" {
		return model[13:]
	}
	return model
}
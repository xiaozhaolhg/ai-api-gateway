package service

import (
	"context"

	"github.com/ai-api-gateway/router-service/internal/domain/port"
)

type ChatCompletionService struct {
	router port.Router
}

func NewChatCompletionService(router port.Router) *ChatCompletionService {
	return &ChatCompletionService{router: router}
}

func (s *ChatCompletionService) HandleChatCompletion(ctx context.Context, req port.ChatCompletionRequest) (*port.ChatCompletionResponse, error) {
	provider, err := s.router.SelectProvider(req.Model)
	if err != nil {
		return nil, err
	}

	return provider.ChatCompletion(req)
}

func (s *ChatCompletionService) HandleStreamChatCompletion(ctx context.Context, req port.ChatCompletionRequest) (<-chan port.StreamChunk, error) {
	provider, err := s.router.SelectProvider(req.Model)
	if err != nil {
		return nil, err
	}

	return provider.StreamChatCompletion(req)
}
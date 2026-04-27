package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"github.com/ai-api-gateway/provider-service/internal/domain/port"
	"github.com/ai-api-gateway/provider-service/internal/infrastructure/crypto"
)

// Service handles provider request forwarding and callback dispatch
type Service struct {
	providerRepo port.ProviderRepository
	adapterFactory *AdapterFactory
	cryptoKey      string // Encryption key for credential decryption
	subscribers    map[string]string // service_name -> gRPC endpoint
	subscribersMu  sync.RWMutex
}

// NewService creates a new application service
func NewService(
	providerRepo port.ProviderRepository,
	adapterFactory *AdapterFactory,
	cryptoKey string,
) *Service {
	return &Service{
		providerRepo:   providerRepo,
		adapterFactory: adapterFactory,
		cryptoKey:      cryptoKey,
		subscribers:    make(map[string]string),
	}
}

// ForwardRequest forwards a non-streaming request to the provider
func (s *Service) ForwardRequest(ctx context.Context, providerID string, requestBody []byte, headers map[string]string) ([]byte, int64, int64, int32, error) {
	// Get provider configuration
	provider, err := s.providerRepo.GetByID(providerID)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to get provider: %w", err)
	}

	if provider.Status != "active" {
		return nil, 0, 0, 0, fmt.Errorf("provider is not active")
	}

	// Decrypt credentials
	decryptedCreds, err := crypto.Decrypt(provider.Credentials, s.cryptoKey)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to decrypt credentials: %w", err)
	}

	// Get adapter for provider type
	adapter, err := s.adapterFactory.GetAdapter(provider.Type)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to get adapter: %w", err)
	}

	// Transform request to provider-specific format
	transformedBody, transformedHeaders, err := adapter.TransformRequest(requestBody, headers)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to transform request: %w", err)
	}

	// Add credentials to headers
	if transformedHeaders == nil {
		transformedHeaders = make(map[string]string)
	}
	transformedHeaders["Authorization"] = "Bearer " + decryptedCreds

	// Make HTTP request to provider
	startTime := time.Now()
	resp, err := s.makeHTTPRequest(ctx, provider.BaseURL, transformedBody, transformedHeaders)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	latencyMs := time.Since(startTime).Milliseconds()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	// Transform response back to OpenAI format (non-streaming)
	transformedResponse, tokenCounts, _, err := adapter.TransformResponse(responseBody, false, entity.TokenCounts{})
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to transform response: %w", err)
	}

	// Use token counts from TransformResponse if available, otherwise fall back to CountTokens
	promptTokens := tokenCounts.PromptTokens
	completionTokens := tokenCounts.CompletionTokens
	
	if promptTokens == 0 && completionTokens == 0 {
		// Fall back to explicit counting if TransformResponse didn't extract tokens
		var err error
		promptTokens, completionTokens, err = adapter.CountTokens(requestBody, transformedResponse, false)
		if err != nil {
			// Log error but don't fail the request
			log.Printf("Failed to count tokens: %v", err)
			promptTokens = 0
			completionTokens = 0
		}
	}

	// Dispatch callbacks asynchronously
	go s.dispatchCallbacks(ctx, providerID, "", "", "", promptTokens, completionTokens, latencyMs, "success", "")

	return transformedResponse, promptTokens, completionTokens, int32(resp.StatusCode), nil
}

// StreamRequest forwards a streaming request to the provider
func (s *Service) StreamRequest(ctx context.Context, providerID string, requestBody []byte, headers map[string]string) (<-chan []byte, <-chan error) {
	chunkChan := make(chan []byte)
	errChan := make(chan error, 1)

	go func() {
		defer close(chunkChan)
		defer close(errChan)

		// Get provider configuration
		provider, err := s.providerRepo.GetByID(providerID)
		if err != nil {
			errChan <- fmt.Errorf("failed to get provider: %w", err)
			return
		}

		if provider.Status != "active" {
			errChan <- fmt.Errorf("provider is not active")
			return
		}

		// Decrypt credentials
		decryptedCreds, err := crypto.Decrypt(provider.Credentials, s.cryptoKey)
		if err != nil {
			errChan <- fmt.Errorf("failed to decrypt credentials: %w", err)
			return
		}

		// Get adapter for provider type
		adapter, err := s.adapterFactory.GetAdapter(provider.Type)
		if err != nil {
			errChan <- fmt.Errorf("failed to get adapter: %w", err)
			return
		}

		// Transform request to provider-specific format
		transformedBody, transformedHeaders, err := adapter.TransformRequest(requestBody, headers)
		if err != nil {
			errChan <- fmt.Errorf("failed to transform request: %w", err)
			return
		}

		// Add credentials to headers
		if transformedHeaders == nil {
			transformedHeaders = make(map[string]string)
		}
		transformedHeaders["Authorization"] = "Bearer " + decryptedCreds

		// Make streaming HTTP request
		startTime := time.Now()
		resp, err := s.makeStreamingHTTPRequest(ctx, provider.BaseURL, transformedBody, transformedHeaders)
		if err != nil {
			errChan <- fmt.Errorf("failed to make request: %w", err)
			return
		}
		defer resp.Body.Close()

		// Stream response chunks
		totalPromptTokens := int64(0)
		totalCompletionTokens := int64(0)

		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				chunk := make([]byte, 1024)
				n, err := resp.Body.Read(chunk)
				if err != nil {
					if err == io.EOF {
						// Stream complete
						latencyMs := time.Since(startTime).Milliseconds()
						go s.dispatchCallbacks(ctx, providerID, "", "", "", totalPromptTokens, totalCompletionTokens, latencyMs, "success", "")
						return
					}
					errChan <- fmt.Errorf("failed to read chunk: %w", err)
					return
				}

				if n > 0 {
					chunkChan <- chunk[:n]
					// Note: Token counting for streaming is more complex and would need SSE parsing
					// For MVP, we'll estimate or skip
				}
			}
		}
	}()

	return chunkChan, errChan
}

// RegisterSubscriber registers a callback subscriber
func (s *Service) RegisterSubscriber(serviceName, callbackEndpoint string) error {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	s.subscribers[serviceName] = callbackEndpoint
	log.Printf("Registered subscriber: %s -> %s", serviceName, callbackEndpoint)
	return nil
}

// UnregisterSubscriber unregisters a callback subscriber
func (s *Service) UnregisterSubscriber(serviceName string) error {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	delete(s.subscribers, serviceName)
	log.Printf("Unregistered subscriber: %s", serviceName)
	return nil
}

// dispatchCallbacks sends ProviderResponseCallback to all registered subscribers
func (s *Service) dispatchCallbacks(ctx context.Context, requestID, userID, groupID, providerID string, promptTokens, completionTokens, latencyMs int64, status, errorCode string) {
	s.subscribersMu.RLock()
	subscribers := make(map[string]string)
	for k, v := range s.subscribers {
		subscribers[k] = v
	}
	s.subscribersMu.RUnlock()

	callback := map[string]interface{}{
		"request_id":         requestID,
		"user_id":           userID,
		"group_id":          groupID,
		"provider_id":       providerID,
		"prompt_tokens":     promptTokens,
		"completion_tokens": completionTokens,
		"latency_ms":        latencyMs,
		"status":            status,
		"error_code":        errorCode,
		"timestamp":         time.Now().Unix(),
	}

	for serviceName, endpoint := range subscribers {
		go func(service, endpoint string) {
			// Fire and forget - log errors but don't block
			if err := s.sendGRPCCallback(ctx, endpoint, callback); err != nil {
				log.Printf("Failed to send callback to %s at %s: %v", service, endpoint, err)
			}
		}(serviceName, endpoint)
	}
}

// sendGRPCCallback sends a callback via gRPC
func (s *Service) sendGRPCCallback(ctx context.Context, endpoint string, callback map[string]interface{}) error {
	// Parse endpoint to get address
	// For MVP, we'll use a simple HTTP POST instead of gRPC for callbacks
	// This simplifies the implementation while maintaining the async fire-and-forget pattern

	callbackJSON, err := json.Marshal(callback)
	if err != nil {
		return fmt.Errorf("failed to marshal callback: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Body = nil // Simplified for MVP

	// For MVP, we'll just log the callback instead of actually sending it
	// In production, this would make an actual HTTP/gRPC call
	log.Printf("Callback to %s: %s", endpoint, string(callbackJSON))

	return nil
}

// makeHTTPRequest makes a non-streaming HTTP request
func (s *Service) makeHTTPRequest(ctx context.Context, url string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	return client.Do(req)
}

// makeStreamingHTTPRequest makes a streaming HTTP request
func (s *Service) makeStreamingHTTPRequest(ctx context.Context, url string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{
		Timeout: 0, // No timeout for streaming
	}

	return client.Do(req)
}

// GetProvider retrieves a provider by ID
func (s *Service) GetProvider(id string) (*entity.Provider, error) {
	return s.providerRepo.GetByID(id)
}

// CreateProvider creates a new provider
func (s *Service) CreateProvider(provider *entity.Provider) error {
	return s.providerRepo.Create(provider)
}

// UpdateProvider updates an existing provider
func (s *Service) UpdateProvider(provider *entity.Provider) error {
	return s.providerRepo.Update(provider)
}

// DeleteProvider deletes a provider
func (s *Service) DeleteProvider(id string) error {
	return s.providerRepo.Delete(id)
}

// ListProviders lists all providers
func (s *Service) ListProviders(page, pageSize int) ([]*entity.Provider, int, error) {
	return s.providerRepo.List(page, pageSize)
}

// GetProviderByType retrieves a provider by type
func (s *Service) GetProviderByType(providerType string) (*entity.Provider, error) {
	return s.providerRepo.GetByType(providerType)
}

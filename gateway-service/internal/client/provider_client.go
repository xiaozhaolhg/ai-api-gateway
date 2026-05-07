package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	providerv1 "github.com/ai-api-gateway/api/gen/provider/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProviderClient is a gRPC client for the provider service with lazy connection.
type ProviderClient struct {
	address string
	client  providerv1.ProviderServiceClient
	conn    *grpc.ClientConn
	mu      sync.RWMutex
}

// NewProviderClient creates a new provider service gRPC client with lazy connection.
func NewProviderClient(address string) (*ProviderClient, error) {
	if address == "" {
		address = "localhost:50053"
	}
	return &ProviderClient{
		address: address,
	}, nil
}

// getClient returns the gRPC client, initializing lazily if needed
func (c *ProviderClient) getClient() (providerv1.ProviderServiceClient, error) {
	c.mu.RLock()
	if c.client != nil {
		defer c.mu.RUnlock()
		return c.client, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if c.client != nil {
		return c.client, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(GRPCInterceptor(DefaultRetryConfig())),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to provider service: %w", err)
	}

	c.conn = conn
	c.client = providerv1.NewProviderServiceClient(conn)
	return c.client, nil
}

// Close closes the gRPC connection.
func (c *ProviderClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ForwardRequest forwards a non-streaming request to the provider service.
func (c *ProviderClient) ForwardRequest(ctx context.Context, providerID string, model string, requestBody []byte, headers map[string]string) (*ForwardRequestResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &providerv1.ForwardRequestRequest{
		ProviderId:  providerID,
		Model:       model,
		RequestBody: requestBody,
		Headers:     headers,
	}

	resp, err := client.ForwardRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %w", err)
	}

	return &ForwardRequestResponse{
		ResponseBody: resp.ResponseBody,
		TokenCounts: TokenCounts{
			PromptTokens:     resp.TokenCounts.PromptTokens,
			CompletionTokens: resp.TokenCounts.CompletionTokens,
			TotalTokens:      resp.TokenCounts.PromptTokens + resp.TokenCounts.CompletionTokens,
		},
		StatusCode: resp.StatusCode,
		Model:      resp.Model,
	}, nil
}

// StreamRequest forwards a streaming request to the provider service.
func (c *ProviderClient) StreamRequest(ctx context.Context, providerID string, model string, requestBody []byte, headers map[string]string) (grpc.ServerStreamingClient[providerv1.ProviderChunk], error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &providerv1.StreamRequestRequest{
		ProviderId:  providerID,
		Model:       model,
		RequestBody: requestBody,
		Headers:     headers,
	}

	stream, err := client.StreamRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to stream request: %w", err)
	}

	return stream, nil
}

// TokenCounts represents token usage statistics.
type TokenCounts struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

// ForwardRequestResponse represents the response from a non-streaming request.
type ForwardRequestResponse struct {
	ResponseBody []byte      `json:"response_body"`
	TokenCounts  TokenCounts `json:"token_counts"`
	StatusCode   int32       `json:"status_code"`
	Model        string      `json:"model"`
}

// Provider represents a provider configuration.
type Provider struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	BaseURL    string   `json:"base_url"`
	Status     string   `json:"status"`
	Models     []string `json:"models"`
	Credentials string   `json:"credentials"`
}

// ListProvidersResponse represents a list of providers.
type ListProvidersResponse struct {
	Providers []*Provider `json:"providers"`
}

// ListProviders lists all providers.
func (c *ProviderClient) ListProviders(ctx context.Context, page, pageSize int32) (*ListProvidersResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &providerv1.ListProvidersRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := client.ListProviders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}

	providers := make([]*Provider, len(resp.Providers))
	for i, p := range resp.Providers {
		providers[i] = &Provider{
			ID:      p.Id,
			Name:    p.Name,
			Type:    p.Type,
			BaseURL: p.BaseUrl,
			Status:  p.Status,
			Models:  p.Models,
		}
	}

	return &ListProvidersResponse{Providers: providers}, nil
}

// GetProvider retrieves a provider by ID.
func (c *ProviderClient) GetProvider(ctx context.Context, id string) (*Provider, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &providerv1.GetProviderRequest{Id: id}
	resp, err := client.GetProvider(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return &Provider{
		ID:     resp.Id,
		Name:   resp.Name,
		Type:   resp.Type,
		Status: resp.Status,
		Models: resp.Models,
	}, nil
}

func (c *ProviderClient) CreateProvider(ctx context.Context, provider *Provider) (*Provider, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &providerv1.CreateProviderRequest{
		Name:        provider.Name,
		Type:        provider.Type,
		BaseUrl:     provider.BaseURL,
		Credentials: provider.Credentials,
		Models:      provider.Models,
		Status:      provider.Status,
	}
	resp, err := client.CreateProvider(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}
	return &Provider{
		ID:          resp.Id,
		Name:        resp.Name,
		Type:        resp.Type,
		BaseURL:     resp.BaseUrl,
		Status:      resp.Status,
		Models:      resp.Models,
		Credentials: resp.Credentials,
	}, nil
}

func (c *ProviderClient) UpdateProvider(ctx context.Context, provider *Provider) (*Provider, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &providerv1.UpdateProviderRequest{
		Id:          provider.ID,
		Name:        provider.Name,
		Type:        provider.Type,
		BaseUrl:     provider.BaseURL,
		Credentials: provider.Credentials,
		Models:      provider.Models,
		Status:      provider.Status,
	}
	resp, err := client.UpdateProvider(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update provider: %w", err)
	}
	return &Provider{
		ID:          resp.Id,
		Name:        resp.Name,
		Type:        resp.Type,
		BaseURL:     resp.BaseUrl,
		Status:      resp.Status,
		Models:      resp.Models,
		Credentials: resp.Credentials,
	}, nil
}

func (c *ProviderClient) DeleteProvider(ctx context.Context, id string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	req := &providerv1.DeleteProviderRequest{Id: id}
	_, err = client.DeleteProvider(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}
	return nil
}

func (c *ProviderClient) HealthCheck(ctx context.Context, id string) (bool, error) {
	client, err := c.getClient()
	if err != nil {
		return false, err
	}

	req := &providerv1.HealthCheckRequest{ProviderId: id}
	resp, err := client.HealthCheck(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to health check provider: %w", err)
	}
	return resp.Healthy, nil
}
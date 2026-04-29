package client

import (
	"context"
	"fmt"

	providerv1 "github.com/ai-api-gateway/api/gen/provider/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProviderClient is a gRPC client for the provider service.
type ProviderClient struct {
	client providerv1.ProviderServiceClient
	conn   *grpc.ClientConn
}

// NewProviderClient creates a new provider service gRPC client.
func NewProviderClient(address string) (*ProviderClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to provider service at %s: %w", address, err)
	}

	return &ProviderClient{
		client: providerv1.NewProviderServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection.
func (c *ProviderClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ForwardRequest forwards a non-streaming request to the provider service.
func (c *ProviderClient) ForwardRequest(ctx context.Context, providerID string, requestBody []byte, headers map[string]string) (*ForwardRequestResponse, error) {
	req := &providerv1.ForwardRequestRequest{
		ProviderId:  providerID,
		RequestBody: requestBody,
		Headers:     headers,
	}

	resp, err := c.client.ForwardRequest(ctx, req)
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
	}, nil
}

// StreamRequest forwards a streaming request to the provider service.
func (c *ProviderClient) StreamRequest(ctx context.Context, providerID string, requestBody []byte, headers map[string]string) (grpc.ServerStreamingClient[providerv1.ProviderChunk], error) {
	req := &providerv1.StreamRequestRequest{
		ProviderId:  providerID,
		RequestBody: requestBody,
		Headers:     headers,
	}

	stream, err := c.client.StreamRequest(ctx, req)
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

// ListProviders lists all providers (stub for MVP).
func (c *ProviderClient) ListProviders(ctx context.Context, page, pageSize int32) (*ListProvidersResponse, error) {
	req := &providerv1.ListProvidersRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.ListProviders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}

	providers := make([]*Provider, len(resp.Providers))
	for i, p := range resp.Providers {
		providers[i] = &Provider{
			ID:     p.Id,
			Name:   p.Name,
			Type:   p.Type,
			Status: p.Status,
			Models: p.Models,
		}
	}

	return &ListProvidersResponse{Providers: providers}, nil
}

// GetProvider retrieves a provider by ID (stub for MVP).
func (c *ProviderClient) GetProvider(ctx context.Context, id string) (*Provider, error) {
	req := &providerv1.GetProviderRequest{Id: id}
	resp, err := c.client.GetProvider(ctx, req)
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
	req := &providerv1.CreateProviderRequest{
		Name:       provider.Name,
		Type:       provider.Type,
		BaseUrl:    provider.BaseURL,
		Credentials: provider.Credentials,
		Models:     provider.Models,
	}
	resp, err := c.client.CreateProvider(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}
	return &Provider{
		ID:         resp.Id,
		Name:       resp.Name,
		Type:       resp.Type,
		BaseURL:    resp.BaseUrl,
		Status:     resp.Status,
		Models:     resp.Models,
		Credentials: resp.Credentials,
	}, nil
}

func (c *ProviderClient) UpdateProvider(ctx context.Context, provider *Provider) (*Provider, error) {
	req := &providerv1.UpdateProviderRequest{
		Id:         provider.ID,
		Name:       provider.Name,
		Type:       provider.Type,
		BaseUrl:    provider.BaseURL,
		Credentials: provider.Credentials,
		Models:     provider.Models,
		Status:     provider.Status,
	}
	resp, err := c.client.UpdateProvider(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update provider: %w", err)
	}
	return &Provider{
		ID:         resp.Id,
		Name:       resp.Name,
		Type:       resp.Type,
		BaseURL:    resp.BaseUrl,
		Status:     resp.Status,
		Models:     resp.Models,
		Credentials: resp.Credentials,
	}, nil
}

func (c *ProviderClient) DeleteProvider(ctx context.Context, id string) error {
	req := &providerv1.DeleteProviderRequest{Id: id}
	_, err := c.client.DeleteProvider(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}
	return nil
}

func (c *ProviderClient) HealthCheck(ctx context.Context, id string) (bool, error) {
	req := &providerv1.HealthCheckRequest{ProviderId: id}
	resp, err := c.client.HealthCheck(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to health check provider: %w", err)
	}
	return resp.Healthy, nil
}
package client

import (
	"context"
	"fmt"
)

// ProviderClient is a stub for provider-service gRPC client
type ProviderClient struct{}

// NewProviderClient creates a new provider service client
func NewProviderClient(address string) (*ProviderClient, error) {
	return &ProviderClient{}, nil
}

// Close is a stub
func (c *ProviderClient) Close() error {
	return nil
}

// ListProviders is a stub
func (c *ProviderClient) ListProviders(ctx context.Context, page, pageSize int32) (*ListProvidersResponse, error) {
	return &ListProvidersResponse{Providers: []*Provider{}}, nil
}

// GetProvider is a stub
func (c *ProviderClient) GetProvider(ctx context.Context, id string) (*Provider, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateProvider is a stub
func (c *ProviderClient) CreateProvider(ctx context.Context, provider *Provider) (*Provider, error) {
	return nil, fmt.Errorf("not implemented")
}

// UpdateProvider is a stub
func (c *ProviderClient) UpdateProvider(ctx context.Context, provider *Provider) (*Provider, error) {
	return nil, fmt.Errorf("not implemented")
}

// DeleteProvider is a stub
func (c *ProviderClient) DeleteProvider(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}

// Provider is a stub type
type Provider struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Status string   `json:"status"`
	Models []string `json:"models"`
}

// ListProvidersResponse is a stub type
type ListProvidersResponse struct {
	Providers []*Provider `json:"providers"`
}
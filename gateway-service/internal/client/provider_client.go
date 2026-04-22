package client

import (
	"context"
	"fmt"
	"io"

	"github.com/ai-api-gateway/api/gen/provider/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProviderClient wraps the provider-service gRPC client
type ProviderClient struct {
	client providerv1.ProviderServiceClient
	conn   *grpc.ClientConn
}

// NewProviderClient creates a new provider service client
func NewProviderClient(address string) (*ProviderClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to provider service: %w", err)
	}

	return &ProviderClient{
		client: providerv1.NewProviderServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close closes the connection
func (c *ProviderClient) Close() error {
	return c.conn.Close()
}

// ForwardRequest forwards a non-streaming request to a provider
func (c *ProviderClient) ForwardRequest(ctx context.Context, providerID string, requestBody []byte, headers map[string]string) (*providerv1.ForwardRequestResponse, error) {
	req := &providerv1.ForwardRequestRequest{
		ProviderId:  providerID,
		RequestBody: requestBody,
		Headers:     headers,
	}

	resp, err := c.client.ForwardRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %w", err)
	}

	return resp, nil
}

// StreamRequest forwards a streaming request to a provider
func (c *ProviderClient) StreamRequest(ctx context.Context, providerID string, requestBody []byte, headers map[string]string) (providerv1.ProviderService_StreamRequestClient, error) {
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

// GetProvider retrieves a provider by ID
func (c *ProviderClient) GetProvider(ctx context.Context, id string) (*providerv1.Provider, error) {
	req := &providerv1.GetProviderRequest{
		Id: id,
	}

	resp, err := c.client.GetProvider(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	return resp, nil
}

// CreateProvider creates a new provider
func (c *ProviderClient) CreateProvider(ctx context.Context, provider *providerv1.Provider) (*providerv1.Provider, error) {
	req := &providerv1.CreateProviderRequest{
		Provider: provider,
	}

	resp, err := c.client.CreateProvider(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	return resp, nil
}

// UpdateProvider updates an existing provider
func (c *ProviderClient) UpdateProvider(ctx context.Context, provider *providerv1.Provider) (*providerv1.Provider, error) {
	req := &providerv1.UpdateProviderRequest{
		Provider: provider,
	}

	resp, err := c.client.UpdateProvider(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update provider: %w", err)
	}

	return resp, nil
}

// DeleteProvider deletes a provider
func (c *ProviderClient) DeleteProvider(ctx context.Context, id string) error {
	req := &providerv1.DeleteProviderRequest{
		Id: id,
	}

	_, err := c.client.DeleteProvider(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete provider: %w", err)
	}

	return nil
}

// ListProviders lists all providers
func (c *ProviderClient) ListProviders(ctx context.Context, page, pageSize int32) (*providerv1.ListProvidersResponse, error) {
	req := &providerv1.ListProvidersRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.ListProviders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}

	return resp, nil
}

// ListModels lists models available from a provider
func (c *ProviderClient) ListModels(ctx context.Context, providerID string) (*providerv1.ListModelsResponse, error) {
	req := &providerv1.ListModelsRequest{
		ProviderId: providerID,
	}

	resp, err := c.client.ListModels(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	return resp, nil
}

// GetProviderByType retrieves a provider by type
func (c *ProviderClient) GetProviderByType(ctx context.Context, providerType string) (*providerv1.Provider, error) {
	req := &providerv1.GetProviderByTypeRequest{
		Type: providerType,
	}

	resp, err := c.client.GetProviderByType(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider by type: %w", err)
	}

	return resp, nil
}

// RegisterSubscriber registers a callback subscriber
func (c *ProviderClient) RegisterSubscriber(ctx context.Context, serviceName, callbackEndpoint string) error {
	req := &providerv1.RegisterSubscriberRequest{
		ServiceName:     serviceName,
		CallbackEndpoint: callbackEndpoint,
	}

	_, err := c.client.RegisterSubscriber(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to register subscriber: %w", err)
	}

	return nil
}

// UnregisterSubscriber unregisters a callback subscriber
func (c *ProviderClient) UnregisterSubscriber(ctx context.Context, serviceName string) error {
	req := &providerv1.UnregisterSubscriberRequest{
		ServiceName: serviceName,
	}

	_, err := c.client.UnregisterSubscriber(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to unregister subscriber: %w", err)
	}

	return nil
}

// ReadStreamChunk reads a chunk from a streaming response
func ReadStreamChunk(stream providerv1.ProviderService_StreamRequestClient) ([]byte, *providerv1.TokenCounts, bool, error) {
	chunk, err := stream.Recv()
	if err == io.EOF {
		return nil, nil, true, nil
	}
	if err != nil {
		return nil, nil, false, err
	}

	return chunk.ChunkData, chunk.AccumulatedTokens, chunk.Done, nil
}

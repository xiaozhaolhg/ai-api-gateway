package client

import (
	"context"
	"fmt"
	"time"

	commonv1 "github.com/ai-api-gateway/api/gen/common/v1"
	routerv1 "github.com/ai-api-gateway/api/gen/router/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RouterClient is a gRPC client for the router service.
type RouterClient struct {
	client routerv1.RouterServiceClient
	conn   *grpc.ClientConn
}

// RouteResolution represents the result of route resolution.
type RouteResolution struct {
	ProviderID          string   `json:"provider_id"`
	AdapterType         string   `json:"adapter_type"`
	FallbackProviderIDs []string `json:"fallback_provider_ids"`
}

// NewRouterClient creates a new router service gRPC client.
func NewRouterClient(address string) (*RouterClient, error) {
	// Create gRPC connection with retry and timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router service at %s: %w", address, err)
	}

	return &RouterClient{
		client: routerv1.NewRouterServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection.
func (c *RouterClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ResolveRoute resolves a model to a provider using the router service.
func (c *RouterClient) ResolveRoute(ctx context.Context, model string, authorizedModels []string) (*RouteResolution, error) {
	req := &routerv1.ResolveRouteRequest{
		Model:            model,
		AuthorizedModels: authorizedModels,
	}

	resp, err := c.client.ResolveRoute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve route for model %s: %w", model, err)
	}

	return &RouteResolution{
		ProviderID:          resp.ProviderId,
		AdapterType:         resp.AdapterType,
		FallbackProviderIDs: resp.FallbackProviderIds,
	}, nil
}

func (c *RouterClient) RefreshRoutingTable(ctx context.Context) error {
	_, err := c.client.RefreshRoutingTable(ctx, &commonv1.Empty{})
	if err != nil {
		return fmt.Errorf("failed to refresh routing table: %w", err)
	}
	return nil
}
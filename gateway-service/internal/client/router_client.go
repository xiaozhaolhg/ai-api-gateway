package client

import (
	"context"
	"fmt"

	"github.com/ai-api-gateway/api/gen/router/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RouterClient wraps the router-service gRPC client
type RouterClient struct {
	client routerv1.RouterServiceClient
	conn   *grpc.ClientConn
}

// NewRouterClient creates a new router service client
func NewRouterClient(address string) (*RouterClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router service: %w", err)
	}

	return &RouterClient{
		client: routerv1.NewRouterServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close closes the connection
func (c *RouterClient) Close() error {
	return c.conn.Close()
}

// ResolveRoute resolves a model name to a provider
func (c *RouterClient) ResolveRoute(ctx context.Context, model string) (*routerv1.RouteResult, error) {
	req := &routerv1.ResolveRouteRequest{
		Model: model,
	}

	resp, err := c.client.ResolveRoute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve route: %w", err)
	}

	return resp, nil
}

// CreateRoutingRule creates a new routing rule
func (c *RouterClient) CreateRoutingRule(ctx context.Context, rule *routerv1.RoutingRule) (*routerv1.RoutingRule, error) {
	req := &routerv1.CreateRoutingRuleRequest{
		Rule: rule,
	}

	resp, err := c.client.CreateRoutingRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create routing rule: %w", err)
	}

	return resp, nil
}

// UpdateRoutingRule updates an existing routing rule
func (c *RouterClient) UpdateRoutingRule(ctx context.Context, rule *routerv1.RoutingRule) (*routerv1.RoutingRule, error) {
	req := &routerv1.UpdateRoutingRuleRequest{
		Rule: rule,
	}

	resp, err := c.client.UpdateRoutingRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update routing rule: %w", err)
	}

	return resp, nil
}

// DeleteRoutingRule deletes a routing rule
func (c *RouterClient) DeleteRoutingRule(ctx context.Context, id string) error {
	req := &routerv1.DeleteRoutingRuleRequest{
		Id: id,
	}

	_, err := c.client.DeleteRoutingRule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete routing rule: %w", err)
	}

	return nil
}

// ListRoutingRules lists all routing rules
func (c *RouterClient) ListRoutingRules(ctx context.Context, page, pageSize int32) (*routerv1.ListRoutingRulesResponse, error) {
	req := &routerv1.ListRoutingRulesRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.ListRoutingRules(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list routing rules: %w", err)
	}

	return resp, nil
}

// RefreshRoutingTable refreshes the routing table cache
func (c *RouterClient) RefreshRoutingTable(ctx context.Context) error {
	req := &routerv1.RefreshRoutingTableRequest{}

	_, err := c.client.RefreshRoutingTable(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to refresh routing table: %w", err)
	}

	return nil
}

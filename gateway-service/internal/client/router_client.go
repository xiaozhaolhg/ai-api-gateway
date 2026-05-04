package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	commonv1 "github.com/ai-api-gateway/api/gen/common/v1"
	routerv1 "github.com/ai-api-gateway/api/gen/router/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RouterClient is a gRPC client for the router service with lazy connection.
type RouterClient struct {
	address string
	client  routerv1.RouterServiceClient
	conn    *grpc.ClientConn
	mu      sync.RWMutex
}

// RouteResolution represents the result of route resolution.
type RouteResolution struct {
	ProviderID          string   `json:"provider_id"`
	AdapterType         string   `json:"adapter_type"`
	FallbackProviderIDs []string `json:"fallback_provider_ids"`
	FallbackModels      []string `json:"fallback_models"`
}

// RoutingRule represents a routing rule for model-to-provider mapping.
type RoutingRule struct {
	ID                 string `json:"id"`
	ModelPattern       string `json:"model_pattern"`
	ProviderID         string `json:"provider_id"`
	Priority           int32  `json:"priority"`
	FallbackProviderID string `json:"fallback_provider_id"`
	FallbackModel      string `json:"fallback_model"`
}

// ListRoutingRulesResponse represents the response from listing routing rules.
type ListRoutingRulesResponse struct {
	Rules []*RoutingRule `json:"rules"`
	Total int32         `json:"total"`
}

// NewRouterClient creates a new router service gRPC client with lazy connection.
func NewRouterClient(address string) (*RouterClient, error) {
	if address == "" {
		address = "localhost:50052"
	}
	return &RouterClient{
		address: address,
	}, nil
}

// getClient returns the gRPC client, initializing lazily if needed
func (c *RouterClient) getClient() (routerv1.RouterServiceClient, error) {
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
		return nil, fmt.Errorf("failed to connect to router service: %w", err)
	}

	c.conn = conn
	c.client = routerv1.NewRouterServiceClient(conn)
	return c.client, nil
}

// Close closes the gRPC connection.
func (c *RouterClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ResolveRoute resolves a model to a provider using the router service.
func (c *RouterClient) ResolveRoute(ctx context.Context, model string, authorizedModels []string) (*RouteResolution, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &routerv1.ResolveRouteRequest{
		Model:            model,
		AuthorizedModels: authorizedModels,
	}

	resp, err := client.ResolveRoute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve route for model %s: %w", model, err)
	}

	return &RouteResolution{
		ProviderID:          resp.ProviderId,
		AdapterType:         resp.AdapterType,
		FallbackProviderIDs: resp.FallbackProviderIds,
		FallbackModels:      resp.FallbackModels,
	}, nil
}

func (c *RouterClient) RefreshRoutingTable(ctx context.Context) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.RefreshRoutingTable(ctx, &commonv1.Empty{})
	if err != nil {
		return fmt.Errorf("failed to refresh routing table: %w", err)
	}
	return nil
}

// ListRoutingRules lists routing rules from the router service.
func (c *RouterClient) ListRoutingRules(ctx context.Context, page, pageSize int32) (*ListRoutingRulesResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &routerv1.GetRoutingRulesRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := client.GetRoutingRules(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list routing rules: %w", err)
	}

	rules := make([]*RoutingRule, len(resp.Rules))
	for i, r := range resp.Rules {
		rules[i] = &RoutingRule{
			ID:                 r.Id,
			ModelPattern:       r.ModelPattern,
			ProviderID:         r.ProviderId,
			Priority:           r.Priority,
			FallbackProviderID: r.FallbackProviderId,
			FallbackModel:      r.FallbackModel,
		}
	}

	return &ListRoutingRulesResponse{
		Rules: rules,
		Total: resp.Total,
	}, nil
}

// CreateRoutingRule creates a new routing rule.
func (c *RouterClient) CreateRoutingRule(ctx context.Context, rule *RoutingRule) (*RoutingRule, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &routerv1.CreateRoutingRuleRequest{
		ModelPattern:       rule.ModelPattern,
		ProviderId:         rule.ProviderID,
		Priority:           rule.Priority,
		FallbackProviderId: rule.FallbackProviderID,
		FallbackModel:      rule.FallbackModel,
	}

	resp, err := client.CreateRoutingRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create routing rule: %w", err)
	}

	return &RoutingRule{
		ID:                 resp.Id,
		ModelPattern:       resp.ModelPattern,
		ProviderID:         resp.ProviderId,
		Priority:           resp.Priority,
		FallbackProviderID: resp.FallbackProviderId,
		FallbackModel:      resp.FallbackModel,
	}, nil
}

// UpdateRoutingRule updates an existing routing rule.
func (c *RouterClient) UpdateRoutingRule(ctx context.Context, rule *RoutingRule) (*RoutingRule, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	req := &routerv1.UpdateRoutingRuleRequest{
		Id:                 rule.ID,
		ModelPattern:       rule.ModelPattern,
		ProviderId:         rule.ProviderID,
		Priority:           rule.Priority,
		FallbackProviderId: rule.FallbackProviderID,
		FallbackModel:      rule.FallbackModel,
	}

	resp, err := client.UpdateRoutingRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update routing rule: %w", err)
	}

	return &RoutingRule{
		ID:                 resp.Id,
		ModelPattern:       resp.ModelPattern,
		ProviderID:         resp.ProviderId,
		Priority:           resp.Priority,
		FallbackProviderID: resp.FallbackProviderId,
		FallbackModel:      resp.FallbackModel,
	}, nil
}

// DeleteRoutingRule deletes a routing rule by ID.
func (c *RouterClient) DeleteRoutingRule(ctx context.Context, id string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.DeleteRoutingRule(ctx, &routerv1.DeleteRoutingRuleRequest{Id: id})
	if err != nil {
		return fmt.Errorf("failed to delete routing rule: %w", err)
	}

	return nil
}
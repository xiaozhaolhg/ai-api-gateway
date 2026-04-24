package client

import (
	"context"
	"fmt"
)

type RouterClient struct{}

func NewRouterClient(address string) (*RouterClient, error) {
	return &RouterClient{}, nil
}

func (c *RouterClient) Close() error {
	return nil
}

func (c *RouterClient) ResolveRoute(ctx context.Context, model string, userID string) (*RouteResolution, error) {
	return nil, fmt.Errorf("not implemented")
}

type RouteResolution struct {
	ProviderID string `json:"provider_id"`
	Endpoint   string `json:"endpoint"`
}
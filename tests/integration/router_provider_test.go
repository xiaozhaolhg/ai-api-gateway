package integration

import (
	"context"
	"testing"
	"time"

	routerv1 "github.com/ai-api-gateway/api/gen/router/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestRouterServiceIntegration tests the router service end-to-end
func TestRouterServiceIntegration(t *testing.T) {
	// This test requires the router service to be running
	// In a full integration test setup, you would:
	// 1. Start router-service in a test container
	// 2. Initialize test data in the database
	// 3. Run test cases
	// 4. Clean up

	t.Skip("Integration test requires running services - skip for MVP")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to router service
	conn, err := grpc.DialContext(ctx, "localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("Failed to connect to router service: %v", err)
	}
	defer conn.Close()

	client := routerv1.NewRouterServiceClient(conn)

	// Test ResolveRoute
	t.Run("ResolveRoute", func(t *testing.T) {
		req := &routerv1.ResolveRouteRequest{
			Model:            "gpt-4",
			AuthorizedModels: []string{"gpt-4", "gpt-3.5-turbo"},
		}

		resp, err := client.ResolveRoute(ctx, req)
		if err != nil {
			t.Fatalf("ResolveRoute failed: %v", err)
		}

		if resp.ProviderId == "" {
			t.Error("Expected non-empty provider ID")
		}

		if resp.AdapterType == "" {
			t.Error("Expected non-empty adapter type")
		}
	})

	// Test CreateRoutingRule
	t.Run("CreateRoutingRule", func(t *testing.T) {
		req := &routerv1.CreateRoutingRuleRequest{
			ModelPattern:       "test-model-*",
			ProviderId:         "test-provider",
			Priority:           10,
			FallbackProviderId: "",
		}

		resp, err := client.CreateRoutingRule(ctx, req)
		if err != nil {
			t.Fatalf("CreateRoutingRule failed: %v", err)
		}

		if resp.Id == "" {
			t.Error("Expected non-empty rule ID")
		}
	})

	// Test GetRoutingRules
	t.Run("GetRoutingRules", func(t *testing.T) {
		req := &routerv1.GetRoutingRulesRequest{
			Page:     1,
			PageSize: 10,
		}

		resp, err := client.GetRoutingRules(ctx, req)
		if err != nil {
			t.Fatalf("GetRoutingRules failed: %v", err)
		}

		if resp.Total < 0 {
			t.Error("Expected non-negative total count")
		}
	})

	// Test RefreshRoutingTable
	t.Run("RefreshRoutingTable", func(t *testing.T) {
		req := &routerv1.Empty{}
		_, err := client.RefreshRoutingTable(ctx, req)
		if err != nil {
			t.Fatalf("RefreshRoutingTable failed: %v", err)
		}
	})
}

// TestGatewayStreamingIntegration tests streaming through gateway
func TestGatewayStreamingIntegration(t *testing.T) {
	// This test requires gateway, router, and provider services to be running
	// In a full integration test setup, you would:
	// 1. Start all services in test containers
	// 2. Send a streaming HTTP request to gateway
	// 3. Verify SSE chunks are received
	// 4. Verify token counts are accumulated correctly
	// 5. Clean up

	t.Skip("Integration test requires running services - skip for MVP")
}

// TestGatewayNonStreamingIntegration tests non-streaming through gateway
func TestGatewayNonStreamingIntegration(t *testing.T) {
	// This test requires gateway, router, and provider services to be running
	// In a full integration test setup, you would:
	// 1. Start all services in test containers
	// 2. Send a non-streaming HTTP request to gateway
	// 3. Verify response is received
	// 4. Verify token counts are reported
	// 5. Clean up

	t.Skip("Integration test requires running services - skip for MVP")
}

package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/ai-api-gateway/router-service/internal/application"
	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	routerv1 "github.com/ai-api-gateway/api/gen/router/v1"
	v1 "github.com/ai-api-gateway/api/gen/common/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mockRoutingRuleRepository is a mock implementation of RoutingRuleRepository
type mockRoutingRuleRepository struct {
	rules []*entity.RoutingRule
}

func (m *mockRoutingRuleRepository) Create(rule *entity.RoutingRule) error {
	if rule.ID == "" {
		rule.ID = "generated-id-" + fmt.Sprintf("%d", len(m.rules)+1)
	}
	m.rules = append(m.rules, rule)
	return nil
}

func (m *mockRoutingRuleRepository) GetByID(id string) (*entity.RoutingRule, error) {
	for _, rule := range m.rules {
		if rule.ID == id {
			return rule, nil
		}
	}
	return nil, nil
}

func (m *mockRoutingRuleRepository) Update(rule *entity.RoutingRule) error {
	for i, r := range m.rules {
		if r.ID == rule.ID {
			m.rules[i] = rule
			return nil
		}
	}
	return nil
}

func (m *mockRoutingRuleRepository) Delete(id string) error {
	for i, r := range m.rules {
		if r.ID == id {
			m.rules = append(m.rules[:i], m.rules[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockRoutingRuleRepository) List(page, pageSize int) ([]*entity.RoutingRule, int, error) {
	return m.rules, len(m.rules), nil
}

// TestHandlerResolveRoute tests the ResolveRoute gRPC method
func TestHandlerResolveRoute(t *testing.T) {
	mockRepo := &mockRoutingRuleRepository{
		rules: []*entity.RoutingRule{
			{
				ID:           "rule1",
				ModelPattern: "gpt-4*",
				ProviderID:   "openai-provider",
				Priority:     10,
			},
		},
	}

	mockCache := &mockCache{}
	service := application.NewService(mockRepo, mockCache)
	handler := NewHandler(service)

	t.Run("SuccessfulResolution", func(t *testing.T) {
		req := &routerv1.ResolveRouteRequest{
			Model:            "gpt-4",
			AuthorizedModels: []string{"gpt-4"},
		}

		resp, err := handler.ResolveRoute(context.Background(), req)
		if err != nil {
			t.Fatalf("ResolveRoute failed: %v", err)
		}

		if resp.ProviderId != "openai-provider" {
			t.Errorf("Expected provider ID 'openai-provider', got '%s'", resp.ProviderId)
		}

		if resp.AdapterType == "" {
			t.Error("Expected non-empty adapter type")
		}
	})

	t.Run("NoMatchingRule", func(t *testing.T) {
		req := &routerv1.ResolveRouteRequest{
			Model:            "unknown-model",
			AuthorizedModels: []string{"unknown-model"},
		}

		_, err := handler.ResolveRoute(context.Background(), req)
		if err == nil {
			t.Error("Expected error for no matching rule")
		}

		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.NotFound {
			t.Errorf("Expected NotFound status, got %v", st.Code())
		}
	})
}

// TestHandlerCreateRoutingRule tests the CreateRoutingRule gRPC method
func TestHandlerCreateRoutingRule(t *testing.T) {
	mockRepo := &mockRoutingRuleRepository{rules: []*entity.RoutingRule{}}
	mockCache := &mockCache{}
	service := application.NewService(mockRepo, mockCache)
	handler := NewHandler(service)

	req := &routerv1.CreateRoutingRuleRequest{
		ModelPattern:       "claude-*",
		ProviderId:         "anthropic-provider",
		Priority:           5,
		FallbackProviderId: "",
	}

	resp, err := handler.CreateRoutingRule(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateRoutingRule failed: %v", err)
	}

	if resp.Id == "" {
		t.Error("Expected non-empty rule ID")
	}

	if resp.ModelPattern != req.ModelPattern {
		t.Errorf("Expected model pattern '%s', got '%s'", req.ModelPattern, resp.ModelPattern)
	}
}

// TestHandlerGetRoutingRules tests the GetRoutingRules gRPC method
func TestHandlerGetRoutingRules(t *testing.T) {
	mockRepo := &mockRoutingRuleRepository{
		rules: []*entity.RoutingRule{
			{ID: "rule1", ModelPattern: "gpt-4*", ProviderID: "openai", Priority: 10},
			{ID: "rule2", ModelPattern: "claude-*", ProviderID: "anthropic", Priority: 5},
		},
	}
	mockCache := &mockCache{}
	service := application.NewService(mockRepo, mockCache)
	handler := NewHandler(service)

	req := &routerv1.GetRoutingRulesRequest{
		Page:     1,
		PageSize: 10,
	}

	resp, err := handler.GetRoutingRules(context.Background(), req)
	if err != nil {
		t.Fatalf("GetRoutingRules failed: %v", err)
	}

	if len(resp.Rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(resp.Rules))
	}

	if resp.Total != 2 {
		t.Errorf("Expected total 2, got %d", resp.Total)
	}
}

// TestHandlerDeleteRoutingRule tests the DeleteRoutingRule gRPC method
func TestHandlerDeleteRoutingRule(t *testing.T) {
	mockRepo := &mockRoutingRuleRepository{
		rules: []*entity.RoutingRule{
			{ID: "rule1", ModelPattern: "gpt-4*", ProviderID: "openai", Priority: 10},
		},
	}
	mockCache := &mockCache{}
	service := application.NewService(mockRepo, mockCache)
	handler := NewHandler(service)

	req := &routerv1.DeleteRoutingRuleRequest{Id: "rule1"}

	_, err := handler.DeleteRoutingRule(context.Background(), req)
	if err != nil {
		t.Fatalf("DeleteRoutingRule failed: %v", err)
	}

	// Verify rule was deleted
	if len(mockRepo.rules) != 0 {
		t.Errorf("Expected 0 rules after deletion, got %d", len(mockRepo.rules))
	}
}

// TestHandlerRefreshRoutingTable tests the RefreshRoutingTable gRPC method
func TestHandlerRefreshRoutingTable(t *testing.T) {
	mockRepo := &mockRoutingRuleRepository{rules: []*entity.RoutingRule{}}
	mockCache := &mockCache{}
	service := application.NewService(mockRepo, mockCache)
	handler := NewHandler(service)

	req := &v1.Empty{}

	_, err := handler.RefreshRoutingTable(context.Background(), req)
	if err != nil {
		t.Fatalf("RefreshRoutingTable failed: %v", err)
	}
}

// TestHandlerResolveFallback tests the ResolveFallback gRPC method (Phase 2+)
func TestHandlerResolveFallback(t *testing.T) {
	mockRepo := &mockRoutingRuleRepository{rules: []*entity.RoutingRule{}}
	mockCache := &mockCache{}
	service := application.NewService(mockRepo, mockCache)
	handler := NewHandler(service)

	req := &routerv1.ResolveFallbackRequest{
		Model:          "gpt-4",
		FailedProviderId: "openai-provider",
	}

	_, err := handler.ResolveFallback(context.Background(), req)
	if err == nil {
		t.Error("Expected Unimplemented error for ResolveFallback")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.Unimplemented {
		t.Errorf("Expected Unimplemented status, got %v", st.Code())
	}
}

// mockCache is a mock implementation of Cache
type mockCache struct{}

func (m *mockCache) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (m *mockCache) Set(ctx context.Context, key string, value string, ttl int) error {
	return nil
}

func (m *mockCache) Delete(ctx context.Context, key string) error {
	return nil
}

func (m *mockCache) ClearPrefix(ctx context.Context, prefix string) error {
	return nil
}

func (m *mockCache) Close() error {
	return nil
}

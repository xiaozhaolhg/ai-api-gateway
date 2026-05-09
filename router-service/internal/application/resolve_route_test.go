package application

import (
	"context"
	"errors"
	"testing"
	"time"

	commonv1 "github.com/ai-api-gateway/api/gen/common/v1"
	providerv1 "github.com/ai-api-gateway/api/gen/provider/v1"
	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	"github.com/ai-api-gateway/router-service/internal/infrastructure/repository"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type mockProviderClient struct {
	findProvidersByModelFunc func(ctx context.Context, in *providerv1.FindProvidersByModelRequest, opts ...grpc.CallOption) (*providerv1.FindProvidersByModelResponse, error)
	healthCheckFunc        func(ctx context.Context, in *providerv1.HealthCheckRequest, opts ...grpc.CallOption) (*providerv1.HealthCheckResponse, error)
}

func (m *mockProviderClient) ForwardRequest(ctx context.Context, in *providerv1.ForwardRequestRequest, opts ...grpc.CallOption) (*providerv1.ForwardRequestResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockProviderClient) StreamRequest(ctx context.Context, in *providerv1.StreamRequestRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[providerv1.ProviderChunk], error) {
	return nil, errors.New("not implemented")
}

func (m *mockProviderClient) GetProvider(ctx context.Context, in *providerv1.GetProviderRequest, opts ...grpc.CallOption) (*providerv1.Provider, error) {
	return nil, errors.New("not implemented")
}

func (m *mockProviderClient) CreateProvider(ctx context.Context, in *providerv1.CreateProviderRequest, opts ...grpc.CallOption) (*providerv1.Provider, error) {
	return nil, errors.New("not implemented")
}

func (m *mockProviderClient) UpdateProvider(ctx context.Context, in *providerv1.UpdateProviderRequest, opts ...grpc.CallOption) (*providerv1.Provider, error) {
	return nil, errors.New("not implemented")
}

func (m *mockProviderClient) DeleteProvider(ctx context.Context, in *providerv1.DeleteProviderRequest, opts ...grpc.CallOption) (*commonv1.Empty, error) {
	return nil, errors.New("not implemented")
}

func (m *mockProviderClient) ListProviders(ctx context.Context, in *providerv1.ListProvidersRequest, opts ...grpc.CallOption) (*providerv1.ListProvidersResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockProviderClient) ListModels(ctx context.Context, in *providerv1.ListModelsRequest, opts ...grpc.CallOption) (*providerv1.ListModelsResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockProviderClient) GetProviderByType(ctx context.Context, in *providerv1.GetProviderByTypeRequest, opts ...grpc.CallOption) (*providerv1.Provider, error) {
	return nil, errors.New("not implemented")
}

func (m *mockProviderClient) HealthCheck(ctx context.Context, in *providerv1.HealthCheckRequest, opts ...grpc.CallOption) (*providerv1.HealthCheckResponse, error) {
	if m.healthCheckFunc != nil {
		return m.healthCheckFunc(ctx, in, opts...)
	}
	return &providerv1.HealthCheckResponse{Healthy: true}, nil
}

func (m *mockProviderClient) FindProvidersByModel(ctx context.Context, in *providerv1.FindProvidersByModelRequest, opts ...grpc.CallOption) (*providerv1.FindProvidersByModelResponse, error) {
	if m.findProvidersByModelFunc != nil {
		return m.findProvidersByModelFunc(ctx, in, opts...)
	}
	return &providerv1.FindProvidersByModelResponse{}, nil
}

func (m *mockProviderClient) RegisterSubscriber(ctx context.Context, in *providerv1.RegisterSubscriberRequest, opts ...grpc.CallOption) (*commonv1.Empty, error) {
	return nil, errors.New("not implemented")
}

func (m *mockProviderClient) UnregisterSubscriber(ctx context.Context, in *providerv1.UnregisterSubscriberRequest, opts ...grpc.CallOption) (*commonv1.Empty, error) {
	return nil, errors.New("not implemented")
}

func setupTestService(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	db.AutoMigrate(&entity.RoutingRule{})
	repo := repository.NewRoutingRuleRepository(db)
	svc, err := NewService(repo, nil, "localhost:50053")
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	return svc, db
}

func setupTestServiceWithMockClient(t *testing.T, mockClient providerv1.ProviderServiceClient) (*Service, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	db.AutoMigrate(&entity.RoutingRule{})
	repo := repository.NewRoutingRuleRepository(db)
	svc := NewServiceWithClient(repo, nil, mockClient)
	return svc, db
}

func TestResolveRoute_UserRuleOverride(t *testing.T) {
	svc, _ := setupTestService(t)

	// Create system rule
	systemRule := &entity.RoutingRule{
		ID:           "system-rule",
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
		IsSystemDefault: true,
	}
	svc.ruleRepo.Create(systemRule)

	// Create user rule that should override
	userRule := &entity.RoutingRule{
		ID:           "user-rule",
		UserID:       "user-123",
		ModelPattern: "ollama:*",
		ProviderID:   "opencode_zen",
	}
	svc.ruleRepo.Create(userRule)

	// Test: user rule should override system rule
	result, err := svc.ResolveRoute(context.Background(), "ollama:llama2", []string{}, "user-123")
	if err != nil {
		t.Fatalf("ResolveRoute failed: %v", err)
	}

	if result.ProviderID != "opencode_zen" {
		t.Errorf("Expected user rule provider 'opencode_zen', got '%s'", result.ProviderID)
	}
}

func TestResolveRoute_SystemRuleFallback(t *testing.T) {
	svc, _ := setupTestService(t)

	// Create only system rule
	systemRule := &entity.RoutingRule{
		ID:           "system-rule",
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
		IsSystemDefault: true,
	}
	svc.ruleRepo.Create(systemRule)

	// Test: with user_id that has no user rule, should fall back to system
	result, err := svc.ResolveRoute(context.Background(), "ollama:llama2", []string{}, "user-456")
	if err != nil {
		t.Fatalf("ResolveRoute failed: %v", err)
	}

	if result.ProviderID != "ollama" {
		t.Errorf("Expected system rule provider 'ollama', got '%s'", result.ProviderID)
	}
}

func TestResolveRoute_NoUserID(t *testing.T) {
	svc, _ := setupTestService(t)

	// Create system rule
	systemRule := &entity.RoutingRule{
		ID:           "system-rule",
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
		IsSystemDefault: true,
	}
	svc.ruleRepo.Create(systemRule)

	// Test: without user_id, should use system rule
	result, err := svc.ResolveRoute(context.Background(), "ollama:llama2", []string{}, "")
	if err != nil {
		t.Fatalf("ResolveRoute failed: %v", err)
	}

	if result.ProviderID != "ollama" {
		t.Errorf("Expected system rule provider 'ollama', got '%s'", result.ProviderID)
	}
}

func TestFindProvidersByModel(t *testing.T) {
	tests := []struct {
		name           string
		model          string
		mockResponse   *providerv1.FindProvidersByModelResponse
		mockError      error
		expectedCount  int
		expectError    bool
	}{
		{
			name:  "successful find single provider",
			model: "llama2",
			mockResponse: &providerv1.FindProvidersByModelResponse{
				Providers: []*providerv1.Provider{
					{Id: "ollama", Name: "Ollama", Type: "ollama", Models: []string{"llama2", "mistral"}, Status: "active", CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
				},
			},
			expectedCount: 1,
		},
		{
			name:  "successful find multiple providers",
			model: "llama2",
			mockResponse: &providerv1.FindProvidersByModelResponse{
				Providers: []*providerv1.Provider{
					{Id: "ollama", Name: "Ollama", Type: "ollama", Models: []string{"llama2"}, Status: "active", CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
					{Id: "opencode_zen", Name: "OpenCode Zen", Type: "opencode_zen", Models: []string{"llama2"}, Status: "active", CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
				},
			},
			expectedCount: 2,
		},
		{
			name:          "no providers found",
			model:         "nonexistent",
			mockResponse:  &providerv1.FindProvidersByModelResponse{Providers: []*providerv1.Provider{}},
			expectedCount: 0,
		},
		{
			name:        "grpc error",
			model:       "llama2",
			mockError:   errors.New("grpc error"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockProviderClient{
				findProvidersByModelFunc: func(ctx context.Context, in *providerv1.FindProvidersByModelRequest, opts ...grpc.CallOption) (*providerv1.FindProvidersByModelResponse, error) {
					if tt.model != in.Model {
						t.Errorf("expected model %s, got %s", tt.model, in.Model)
					}
					return tt.mockResponse, tt.mockError
				},
			}

			svc, _ := setupTestServiceWithMockClient(t, mockClient)
			providers, err := svc.FindProvidersByModel(tt.model)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && len(providers) != tt.expectedCount {
				t.Errorf("expected %d providers, got %d", tt.expectedCount, len(providers))
			}
		})
	}
}

func TestResolveBareModel(t *testing.T) {
	tests := []struct {
		name             string
		model            string
		providers        []*providerv1.Provider
		healthStatuses   map[string]bool
		expectError      bool
		expectedPrimary  string
		expectedFallbacks []string
	}{
		{
			name:  "single healthy provider",
			model: "llama2",
			providers: []*providerv1.Provider{
				{Id: "ollama", Name: "Ollama", Type: "ollama", Models: []string{"llama2"}, Status: "active", CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
			},
			healthStatuses: map[string]bool{"ollama": true},
			expectedPrimary: "ollama",
		},
		{
			name:  "multiple providers selects healthiest",
			model: "llama2",
			providers: []*providerv1.Provider{
				{Id: "ollama", Name: "Ollama", Type: "ollama", Models: []string{"llama2"}, Status: "active", CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
				{Id: "opencode_zen", Name: "OpenCode Zen", Type: "opencode_zen", Models: []string{"llama2"}, Status: "active", CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
			},
			healthStatuses: map[string]bool{"ollama": true, "opencode_zen": true},
			expectedPrimary:  "",
			expectedFallbacks: nil,
		},
		{
			name:  "skips unhealthy providers",
			model: "llama2",
			providers: []*providerv1.Provider{
				{Id: "ollama", Name: "Ollama", Type: "ollama", Models: []string{"llama2"}, Status: "active", CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
				{Id: "opencode_zen", Name: "OpenCode Zen", Type: "opencode_zen", Models: []string{"llama2"}, Status: "active", CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
			},
			healthStatuses: map[string]bool{"ollama": false, "opencode_zen": true},
			expectedPrimary: "opencode_zen",
		},
		{
			name:  "all providers unhealthy",
			model: "llama2",
			providers: []*providerv1.Provider{
				{Id: "ollama", Name: "Ollama", Type: "ollama", Models: []string{"llama2"}, Status: "active", CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
			},
			healthStatuses: map[string]bool{"ollama": false},
			expectError:     true,
		},
		{
			name:           "no providers found",
			model:          "nonexistent",
			providers:       []*providerv1.Provider{},
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockProviderClient{
				findProvidersByModelFunc: func(ctx context.Context, in *providerv1.FindProvidersByModelRequest, opts ...grpc.CallOption) (*providerv1.FindProvidersByModelResponse, error) {
					return &providerv1.FindProvidersByModelResponse{Providers: tt.providers}, nil
				},
				healthCheckFunc: func(ctx context.Context, in *providerv1.HealthCheckRequest, opts ...grpc.CallOption) (*providerv1.HealthCheckResponse, error) {
					healthy := tt.healthStatuses[in.ProviderId]
					return &providerv1.HealthCheckResponse{Healthy: healthy}, nil
				},
			}

			svc, _ := setupTestServiceWithMockClient(t, mockClient)
			result, err := svc.resolveBareModel(context.Background(), tt.model)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && result != nil {
				if tt.expectedPrimary != "" && result.ProviderID != tt.expectedPrimary {
					t.Errorf("expected primary %s, got %s", tt.expectedPrimary, result.ProviderID)
				}
				if tt.expectedPrimary == "" && result.ProviderID == "" {
					t.Error("expected a primary provider to be selected")
				}
				if tt.expectedFallbacks != nil {
					if len(result.FallbackProviderIDs) != len(tt.expectedFallbacks) {
						t.Errorf("expected %d fallbacks, got %d", len(tt.expectedFallbacks), len(result.FallbackProviderIDs))
					}
				}
			}
		})
	}
}

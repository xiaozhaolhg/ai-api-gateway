package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ai-api-gateway/billing-service/internal/application"
	"github.com/ai-api-gateway/billing-service/internal/domain/entity"
)

type MockUsageRepo struct {
	mock.Mock
}

func (m *MockUsageRepo) Create(record *entity.UsageRecord) error {
	return m.Called(record).Error(0)
}

func (m *MockUsageRepo) GetByID(id string) (*entity.UsageRecord, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UsageRecord), args.Error(1)
}

func (m *MockUsageRepo) GetByUserID(userID string, page, pageSize int) ([]*entity.UsageRecord, int, error) {
	args := m.Called(userID, page, pageSize)
	return args.Get(0).([]*entity.UsageRecord), args.Int(1), args.Error(2)
}

func (m *MockUsageRepo) GetAggregation(userID, startDate, endDate, groupBy string) ([]*entity.UsageAggregation, error) {
	args := m.Called(userID, startDate, endDate, groupBy)
	return args.Get(0).([]*entity.UsageAggregation), args.Error(1)
}

type MockPricingRepo struct {
	mock.Mock
}

func (m *MockPricingRepo) Create(rule *entity.PricingRule) error {
	return m.Called(rule).Error(0)
}

func (m *MockPricingRepo) GetByID(id string) (*entity.PricingRule, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.PricingRule), args.Error(1)
}

func (m *MockPricingRepo) GetByProviderAndModel(providerID, model string) (*entity.PricingRule, error) {
	args := m.Called(providerID, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.PricingRule), args.Error(1)
}

func (m *MockPricingRepo) Update(rule *entity.PricingRule) error {
	return m.Called(rule).Error(0)
}

func (m *MockPricingRepo) Delete(id string) error {
	return m.Called(id).Error(0)
}

func (m *MockPricingRepo) List(page, pageSize int) ([]*entity.PricingRule, int, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]*entity.PricingRule), args.Int(1), args.Error(2)
}

type MockAccountRepo struct {
	mock.Mock
}

func (m *MockAccountRepo) Create(account *entity.BillingAccount) error {
	return m.Called(account).Error(0)
}

func (m *MockAccountRepo) GetByID(id string) (*entity.BillingAccount, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BillingAccount), args.Error(1)
}

func (m *MockAccountRepo) GetByUserID(userID string) (*entity.BillingAccount, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BillingAccount), args.Error(1)
}

func (m *MockAccountRepo) Update(account *entity.BillingAccount) error {
	return m.Called(account).Error(0)
}

func (m *MockAccountRepo) Delete(id string) error {
	return m.Called(id).Error(0)
}

type MockBudgetRepo struct {
	mock.Mock
}

func (m *MockBudgetRepo) Create(budget *entity.Budget) error {
	return m.Called(budget).Error(0)
}

func (m *MockBudgetRepo) GetByID(id string) (*entity.Budget, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Budget), args.Error(1)
}

func (m *MockBudgetRepo) GetByUserID(userID string) (*entity.Budget, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Budget), args.Error(1)
}

func (m *MockBudgetRepo) Update(budget *entity.Budget) error {
	return m.Called(budget).Error(0)
}

func (m *MockBudgetRepo) Delete(id string) error {
	return m.Called(id).Error(0)
}

func (m *MockBudgetRepo) List(page, pageSize int) ([]*entity.Budget, int, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]*entity.Budget), args.Int(1), args.Error(2)
}

func TestService_GetUsageAggregation_MultipleRows(t *testing.T) {
	mockUsageRepo := new(MockUsageRepo)
	mockPricingRepo := new(MockPricingRepo)
	mockAccountRepo := new(MockAccountRepo)
	mockBudgetRepo := new(MockBudgetRepo)

	mockUsageRepo.On("GetAggregation", "user-123", "2024-01-01", "2024-01-31", "provider_id").
		Return([]*entity.UsageAggregation{
			{UserID: "user-123", ProviderID: "ollama", Model: "llama2", TotalRequests: 10, TotalPromptTokens: 100, TotalCompletionTokens: 200, TotalCost: 5.50},
			{UserID: "user-123", ProviderID: "openai", Model: "gpt-4", TotalRequests: 5, TotalPromptTokens: 50, TotalCompletionTokens: 150, TotalCost: 12.00},
		}, nil)

	svc := application.NewService(mockUsageRepo, mockPricingRepo, mockAccountRepo, mockBudgetRepo)

	aggs, err := svc.GetUsageAggregation("user-123", "2024-01-01", "2024-01-31", "provider_id")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(aggs))
	assert.Equal(t, "ollama", aggs[0].ProviderID)
	assert.Equal(t, "openai", aggs[1].ProviderID)
	mockUsageRepo.AssertExpectations(t)
}

func TestHandler_GetUsageAggregation(t *testing.T) {
	t.Skip("Handler tests need gRPC test server or full service mock")
}

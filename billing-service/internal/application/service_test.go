package application

import (
	"testing"
	"time"

	"github.com/ai-api-gateway/billing-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUsageRepo struct {
	mock.Mock
}

func (m *MockUsageRepo) Create(record *entity.UsageRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockUsageRepo) GetByID(id string) (*entity.UsageRecord, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UsageRecord), args.Error(1)
}

func (m *MockUsageRepo) GetByUserID(userID string, page, pageSize int, startTime, endTime int64) ([]*entity.UsageRecord, int, error) {
	args := m.Called(userID, page, pageSize, startTime, endTime)
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
	args := m.Called(rule)
	return args.Error(0)
}

func (m *MockPricingRepo) Update(rule *entity.PricingRule) error {
	args := m.Called(rule)
	return args.Error(0)
}

func (m *MockPricingRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPricingRepo) List(page, pageSize int) ([]*entity.PricingRule, int, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]*entity.PricingRule), args.Int(1), args.Error(2)
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

type MockAccountRepo struct {
	mock.Mock
}

func (m *MockAccountRepo) Create(account *entity.BillingAccount) error {
	args := m.Called(account)
	return args.Error(0)
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
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockBudgetRepo struct {
	mock.Mock
}

func (m *MockBudgetRepo) GetByID(id string) (*entity.Budget, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Budget), args.Error(1)
}

func (m *MockBudgetRepo) Create(budget *entity.Budget) error {
	args := m.Called(budget)
	return args.Error(0)
}

func (m *MockBudgetRepo) Update(budget *entity.Budget) error {
	args := m.Called(budget)
	return args.Error(0)
}

func (m *MockBudgetRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBudgetRepo) List(page, pageSize int) ([]*entity.Budget, int, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]*entity.Budget), args.Int(1), args.Error(2)
}

func (m *MockBudgetRepo) GetByUserID(userID string) (*entity.Budget, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Budget), args.Error(1)
}

func TestGetUsageAggregation_MultipleRecordsFromSameRequest(t *testing.T) {
	mockUsageRepo := new(MockUsageRepo)
	mockPricingRepo := new(MockPricingRepo)
	mockAccountRepo := new(MockAccountRepo)
	mockBudgetRepo := new(MockBudgetRepo)

	startDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	mockUsageRepo.On("GetAggregation", "user-123", startDate, endDate, "provider_id").
		Return([]*entity.UsageAggregation{
			{
				UserID:               "user-123",
				GroupID:              "group-1",
				ProviderID:           "ollama",
				Model:                "llama2",
				TotalRequests:        4,
				TotalPromptTokens:     400,
				TotalCompletionTokens: 4000,
				TotalCost:             2.00,
			},
		}, nil)

	svc := NewService(mockUsageRepo, mockPricingRepo, mockAccountRepo, mockBudgetRepo)
	aggs, err := svc.GetUsageAggregation("user-123", startDate, endDate, "provider_id")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(aggs))
	assert.Equal(t, int64(4000), aggs[0].TotalCompletionTokens)
	mockUsageRepo.AssertExpectations(t)
}

func TestGetUsageAggregation_ThreeIntermediateOneFinal(t *testing.T) {
	mockUsageRepo := new(MockUsageRepo)
	mockPricingRepo := new(MockPricingRepo)
	mockAccountRepo := new(MockAccountRepo)
	mockBudgetRepo := new(MockBudgetRepo)

		

	startDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	mockUsageRepo.On("GetAggregation", "user-123", startDate, endDate, "provider_id").
		Return([]*entity.UsageAggregation{
			{
				UserID:               "user-123",
				GroupID:              "group-1",
				ProviderID:           "ollama",
				Model:                "llama2",
				TotalRequests:        4,
				TotalPromptTokens:     400,
				TotalCompletionTokens: 4000,
				TotalCost:             2.00,
			},
		}, nil)

	svc := NewService(mockUsageRepo, mockPricingRepo, mockAccountRepo, mockBudgetRepo)
	aggs, err := svc.GetUsageAggregation("user-123", startDate, endDate, "provider_id")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(aggs))
	assert.Equal(t, int64(4000), aggs[0].TotalCompletionTokens)
	mockUsageRepo.AssertExpectations(t)
}

package port

import (
	"github.com/ai-api-gateway/billing-service/internal/domain/entity"
)

// UsageRecordRepository defines the interface for usage record persistence operations
type UsageRecordRepository interface {
	Create(record *entity.UsageRecord) error
	GetByID(id string) (*entity.UsageRecord, error)
	// GetByUserID retrieves usage records for a user. If startTime or endTime is 0, that bound is not applied.
	GetByUserID(userID string, page, pageSize int, startTime, endTime int64) ([]*entity.UsageRecord, int, error)
	GetAggregation(userID, startDate, endDate, groupBy string) ([]*entity.UsageAggregation, error)
}

// PricingRuleRepository defines the interface for pricing rule persistence operations
type PricingRuleRepository interface {
	Create(rule *entity.PricingRule) error
	GetByID(id string) (*entity.PricingRule, error)
	GetByProviderAndModel(providerID, model string) (*entity.PricingRule, error)
	Update(rule *entity.PricingRule) error
	Delete(id string) error
	List(page, pageSize int) ([]*entity.PricingRule, int, error)
}

// BillingAccountRepository defines the interface for billing account persistence operations
type BillingAccountRepository interface {
	Create(account *entity.BillingAccount) error
	GetByID(id string) (*entity.BillingAccount, error)
	GetByUserID(userID string) (*entity.BillingAccount, error)
	Update(account *entity.BillingAccount) error
	Delete(id string) error
}

// BudgetRepository defines the interface for budget persistence operations
type BudgetRepository interface {
	Create(budget *entity.Budget) error
	GetByID(id string) (*entity.Budget, error)
	GetByUserID(userID string) (*entity.Budget, error)
	Update(budget *entity.Budget) error
	Delete(id string) error
	List(page, pageSize int) ([]*entity.Budget, int, error)
}

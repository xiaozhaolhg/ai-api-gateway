package application

import (
	"fmt"
	"time"

	"github.com/ai-api-gateway/billing-service/internal/domain/entity"
	"github.com/ai-api-gateway/billing-service/internal/domain/port"
)

// Service handles billing logic
type Service struct {
	usageRepo   port.UsageRecordRepository
	pricingRepo port.PricingRuleRepository
	accountRepo port.BillingAccountRepository
	budgetRepo  port.BudgetRepository
}

// NewService creates a new application service
func NewService(
	usageRepo port.UsageRecordRepository,
	pricingRepo port.PricingRuleRepository,
	accountRepo port.BillingAccountRepository,
	budgetRepo port.BudgetRepository,
) *Service {
	return &Service{
		usageRepo:   usageRepo,
		pricingRepo: pricingRepo,
		accountRepo: accountRepo,
		budgetRepo:  budgetRepo,
	}
}

// RecordUsage records usage for a request
func (s *Service) RecordUsage(userID, providerID, model string, promptTokens, completionTokens int64) error {
	// Calculate cost
	cost, err := s.calculateCost(providerID, model, promptTokens, completionTokens)
	if err != nil {
		return fmt.Errorf("failed to calculate cost: %w", err)
	}

	// Create usage record
	record := &entity.UsageRecord{
		UserID:           userID,
		ProviderID:       providerID,
		Model:            model,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		Cost:             cost,
		Timestamp:        time.Now(),
	}

	// Save usage record
	if err := s.usageRepo.Create(record); err != nil {
		return fmt.Errorf("failed to create usage record: %w", err)
	}

	// Deduct from billing account
	if err := s.deductFromAccount(userID, cost); err != nil {
		return fmt.Errorf("failed to deduct from account: %w", err)
	}

	return nil
}

// calculateCost calculates the cost for a request based on pricing rules
func (s *Service) calculateCost(providerID, model string, promptTokens, completionTokens int64) (float64, error) {
	// Get pricing rule for provider/model
	rule, err := s.pricingRepo.GetByProviderAndModel(providerID, model)
	if err != nil {
		// If no specific rule, use default pricing
		return 0.0, nil
	}

	// Calculate cost
	promptCost := float64(promptTokens) / 1000.0 * rule.PromptPricePer1K
	completionCost := float64(completionTokens) / 1000.0 * rule.CompletionPricePer1K
	totalCost := promptCost + completionCost

	return totalCost, nil
}

// deductFromAccount deducts cost from user's billing account
func (s *Service) deductFromAccount(userID string, cost float64) error {
	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get billing account: %w", err)
	}

	// Check balance
	if account.Balance < cost {
		return fmt.Errorf("insufficient balance")
	}

	// Deduct cost
	account.Balance -= cost
	account.UpdatedAt = time.Now()

	if err := s.accountRepo.Update(account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return nil
}

// CheckBudget checks if a user has budget available
func (s *Service) CheckBudget(userID string) (*entity.BudgetStatus, error) {
	// Get budget
	budget, err := s.budgetRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}

	// Get billing account
	account, err := s.accountRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	// Calculate current spend (simplified - in production, you'd aggregate usage records)
	currentSpend := budget.Limit - account.Balance
	remaining := account.Balance

	// Check caps
	softCapExceeded := currentSpend >= budget.SoftCap
	hardCapExceeded := currentSpend >= budget.HardCap

	return &entity.BudgetStatus{
		CurrentSpend:    currentSpend,
		BudgetLimit:     budget.Limit,
		Remaining:       remaining,
		SoftCapExceeded: softCapExceeded,
		HardCapExceeded: hardCapExceeded,
	}, nil
}

// GetUsage retrieves usage records for a user
func (s *Service) GetUsage(userID string, page, pageSize int) ([]*entity.UsageRecord, int, error) {
	return s.usageRepo.GetByUserID(userID, page, pageSize)
}

// GetUsageAggregation retrieves aggregated usage statistics
func (s *Service) GetUsageAggregation(userID, startDate, endDate, groupBy string) ([]*entity.UsageAggregation, error) {
	return s.usageRepo.GetAggregation(userID, startDate, endDate, groupBy)
}

// EstimateCost estimates the cost of a request
func (s *Service) EstimateCost(providerID, model string, promptTokens, completionTokens int64) (float64, error) {
	return s.calculateCost(providerID, model, promptTokens, completionTokens)
}

// CreatePricingRule creates a new pricing rule
func (s *Service) CreatePricingRule(rule *entity.PricingRule) error {
	return s.pricingRepo.Create(rule)
}

// UpdatePricingRule updates an existing pricing rule
func (s *Service) UpdatePricingRule(rule *entity.PricingRule) error {
	return s.pricingRepo.Update(rule)
}

// DeletePricingRule deletes a pricing rule
func (s *Service) DeletePricingRule(id string) error {
	return s.pricingRepo.Delete(id)
}

// ListPricingRules lists all pricing rules with pagination
func (s *Service) ListPricingRules(page, pageSize int) ([]*entity.PricingRule, int, error) {
	return s.pricingRepo.List(page, pageSize)
}

// GetPricingRuleByID retrieves a pricing rule by ID
func (s *Service) GetPricingRuleByID(id string) (*entity.PricingRule, error) {
	return s.pricingRepo.GetByID(id)
}

// GetBillingAccount retrieves a billing account
func (s *Service) GetBillingAccount(userID string) (*entity.BillingAccount, error) {
	return s.accountRepo.GetByUserID(userID)
}

// CreateBillingAccount creates a new billing account
func (s *Service) CreateBillingAccount(account *entity.BillingAccount) error {
	return s.accountRepo.Create(account)
}

// UpdateBillingAccount updates an existing billing account
func (s *Service) UpdateBillingAccount(account *entity.BillingAccount) error {
	return s.accountRepo.Update(account)
}

// CreateBudget creates a new budget
func (s *Service) CreateBudget(budget *entity.Budget) error {
	return s.budgetRepo.Create(budget)
}

// UpdateBudget updates an existing budget
func (s *Service) UpdateBudget(budget *entity.Budget) error {
	return s.budgetRepo.Update(budget)
}

// DeleteBudget deletes a budget
func (s *Service) DeleteBudget(id string) error {
	return s.budgetRepo.Delete(id)
}

// ListBudgets lists all budgets with pagination
func (s *Service) ListBudgets(page, pageSize int) ([]*entity.Budget, int, error) {
	return s.budgetRepo.List(page, pageSize)
}

// GetBudgetByID retrieves a budget by ID
func (s *Service) GetBudgetByID(id string) (*entity.Budget, error) {
	return s.budgetRepo.GetByID(id)
}

// GetBudgetByUserID retrieves a budget by user ID
func (s *Service) GetBudgetByUserID(userID string) (*entity.Budget, error) {
	return s.budgetRepo.GetByUserID(userID)
}

// GetBillingAccountByID retrieves a billing account by ID
func (s *Service) GetBillingAccountByID(id string) (*entity.BillingAccount, error) {
	return s.accountRepo.GetByID(id)
}

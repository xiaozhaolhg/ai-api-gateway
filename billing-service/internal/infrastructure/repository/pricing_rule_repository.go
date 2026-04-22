package repository

import (
	"database/sql"
	"fmt"

	"github.com/ai-api-gateway/billing-service/internal/domain/entity"
	"github.com/google/uuid"
)

// PricingRuleRepository implements the pricing rule repository interface
type PricingRuleRepository struct {
	db *sql.DB
}

// NewPricingRuleRepository creates a new pricing rule repository
func NewPricingRuleRepository(db *sql.DB) *PricingRuleRepository {
	return &PricingRuleRepository{db: db}
}

// Create creates a new pricing rule
func (r *PricingRuleRepository) Create(rule *entity.PricingRule) error {
	rule.ID = uuid.New().String()

	query := `
		INSERT INTO pricing_rules (id, provider_id, model, prompt_price_per_1k, completion_price_per_1k, currency, effective_from, effective_until, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, rule.ID, rule.ProviderID, rule.Model,
		rule.PromptPricePer1K, rule.CompletionPricePer1K, rule.Currency,
		rule.EffectiveFrom, rule.EffectiveUntil, rule.CreatedAt, rule.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create pricing rule: %w", err)
	}

	return nil
}

// GetByID retrieves a pricing rule by ID
func (r *PricingRuleRepository) GetByID(id string) (*entity.PricingRule, error) {
	query := `
		SELECT id, provider_id, model, prompt_price_per_1k, completion_price_per_1k, currency, effective_from, effective_until, created_at, updated_at
		FROM pricing_rules
		WHERE id = ?
	`

	var rule entity.PricingRule
	err := r.db.QueryRow(query, id).Scan(
		&rule.ID, &rule.ProviderID, &rule.Model,
		&rule.PromptPricePer1K, &rule.CompletionPricePer1K, &rule.Currency,
		&rule.EffectiveFrom, &rule.EffectiveUntil, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pricing rule not found")
		}
		return nil, fmt.Errorf("failed to get pricing rule: %w", err)
	}

	return &rule, nil
}

// GetByProviderAndModel retrieves a pricing rule by provider and model
func (r *PricingRuleRepository) GetByProviderAndModel(providerID, model string) (*entity.PricingRule, error) {
	query := `
		SELECT id, provider_id, model, prompt_price_per_1k, completion_price_per_1k, currency, effective_from, effective_until, created_at, updated_at
		FROM pricing_rules
		WHERE provider_id = ? AND model = ?
		AND effective_from <= datetime('now')
		AND (effective_until IS NULL OR effective_until > datetime('now'))
		ORDER BY effective_from DESC
		LIMIT 1
	`

	var rule entity.PricingRule
	err := r.db.QueryRow(query, providerID, model).Scan(
		&rule.ID, &rule.ProviderID, &rule.Model,
		&rule.PromptPricePer1K, &rule.CompletionPricePer1K, &rule.Currency,
		&rule.EffectiveFrom, &rule.EffectiveUntil, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pricing rule not found")
		}
		return nil, fmt.Errorf("failed to get pricing rule: %w", err)
	}

	return &rule, nil
}

// Update updates an existing pricing rule
func (r *PricingRuleRepository) Update(rule *entity.PricingRule) error {
	query := `
		UPDATE pricing_rules
		SET provider_id = ?, model = ?, prompt_price_per_1k = ?, completion_price_per_1k = ?, currency = ?, effective_from = ?, effective_until = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, rule.ProviderID, rule.Model,
		rule.PromptPricePer1K, rule.CompletionPricePer1K, rule.Currency,
		rule.EffectiveFrom, rule.EffectiveUntil, rule.UpdatedAt, rule.ID)
	if err != nil {
		return fmt.Errorf("failed to update pricing rule: %w", err)
	}

	return nil
}

// Delete deletes a pricing rule
func (r *PricingRuleRepository) Delete(id string) error {
	query := `DELETE FROM pricing_rules WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pricing rule: %w", err)
	}

	return nil
}

// List lists all pricing rules
func (r *PricingRuleRepository) List(page, pageSize int) ([]*entity.PricingRule, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM pricing_rules").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pricing rules: %w", err)
	}

	// Get rules
	query := `
		SELECT id, provider_id, model, prompt_price_per_1k, completion_price_per_1k, currency, effective_from, effective_until, created_at, updated_at
		FROM pricing_rules
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query pricing rules: %w", err)
	}
	defer rows.Close()

	var rules []*entity.PricingRule
	for rows.Next() {
		var rule entity.PricingRule
		err := rows.Scan(
			&rule.ID, &rule.ProviderID, &rule.Model,
			&rule.PromptPricePer1K, &rule.CompletionPricePer1K, &rule.Currency,
			&rule.EffectiveFrom, &rule.EffectiveUntil, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pricing rule: %w", err)
		}
		rules = append(rules, &rule)
	}

	return rules, total, nil
}

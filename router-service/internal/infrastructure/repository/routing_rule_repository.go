package repository

import (
	"fmt"
	"strings"

	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	"github.com/ai-api-gateway/router-service/internal/domain/port"
	"gorm.io/gorm"
)

// RoutingRuleRepository implements the RoutingRuleRepository interface using GORM
type RoutingRuleRepository struct {
	db *gorm.DB
}

// NewRoutingRuleRepository creates a new RoutingRuleRepository
func NewRoutingRuleRepository(db *gorm.DB) port.RoutingRuleRepository {
	return &RoutingRuleRepository{db: db}
}

// Create creates a new routing rule
func (r *RoutingRuleRepository) Create(rule *entity.RoutingRule) error {
	return r.db.Create(rule).Error
}

// GetByID retrieves a routing rule by ID
func (r *RoutingRuleRepository) GetByID(id string) (*entity.RoutingRule, error) {
	var rule entity.RoutingRule
	err := r.db.Where("id = ?", id).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// Update updates a routing rule
func (r *RoutingRuleRepository) Update(rule *entity.RoutingRule) error {
	return r.db.Save(rule).Error
}

// UpdateWithOwnership updates a routing rule, verifying ownership for user rules.
func (r *RoutingRuleRepository) UpdateWithOwnership(rule *entity.RoutingRule, requestingUserID string) error {
	var existing entity.RoutingRule
	err := r.db.Where("id = ?", rule.ID).First(&existing).Error
	if err != nil {
		return err
	}

	if existing.UserID != "" && existing.UserID != requestingUserID {
		return fmt.Errorf("not authorized to update this rule")
	}

	return r.db.Save(rule).Error
}

// Delete deletes a routing rule by ID.
func (r *RoutingRuleRepository) Delete(id string) error {
	return r.db.Delete(&entity.RoutingRule{}, "id = ?", id).Error
}

// DeleteWithOwnership deletes a routing rule by ID, verifying ownership for user rules.
func (r *RoutingRuleRepository) DeleteWithOwnership(id string, requestingUserID string) error {
	var rule entity.RoutingRule
	err := r.db.Where("id = ?", id).First(&rule).Error
	if err != nil {
		return err
	}

	if rule.UserID != "" && rule.UserID != requestingUserID {
		return fmt.Errorf("not authorized to delete this rule")
	}

	return r.db.Delete(&entity.RoutingRule{}, "id = ?", id).Error
}

// List retrieves routing rules with pagination
func (r *RoutingRuleRepository) List(page, pageSize int) ([]*entity.RoutingRule, int, error) {
	var rules []*entity.RoutingRule
	var total int64

	if err := r.db.Model(&entity.RoutingRule{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Order("priority ASC").Offset(offset).Limit(pageSize).Find(&rules).Error
	if err != nil {
		return nil, 0, err
	}

	return rules, int(total), nil
}

func (r *RoutingRuleRepository) FindByModel(model string, userID *string) (*entity.RoutingRule, error) {
	if userID != nil && *userID != "" {
		var userRules []entity.RoutingRule
		if err := r.db.Where("user_id = ?", *userID).Order("priority ASC").Find(&userRules).Error; err != nil {
			return nil, err
		}
		for _, rule := range userRules {
			if matchModelPattern(rule.ModelPattern, model) {
				return &rule, nil
			}
		}

		var systemRules []entity.RoutingRule
		if err := r.db.Where("user_id = '' OR user_id IS NULL OR is_system_default = ?", true).Order("priority ASC").Find(&systemRules).Error; err != nil {
			return nil, err
		}
		for _, rule := range systemRules {
			if matchModelPattern(rule.ModelPattern, model) {
				return &rule, nil
			}
		}
		return nil, fmt.Errorf("record not found")
	}

	var rules []entity.RoutingRule
	err := r.db.Where("user_id = '' OR user_id IS NULL OR is_system_default = ?", true).Order("priority ASC").Find(&rules).Error
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if matchModelPattern(rule.ModelPattern, model) {
			return &rule, nil
		}
	}

	return nil, fmt.Errorf("record not found")
}

// matchModelPattern checks if a model matches a pattern (e.g., "ollama:*", "*-gpt4")
func matchModelPattern(pattern, model string) bool {
	if pattern == "*" {
		return true
	}

	if !strings.Contains(pattern, "*") {
		return pattern == model
	}

	// Handle wildcard patterns
	if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(model, suffix)
	}

	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(model, prefix)
	}

	// Handle patterns like "*-gpt4"
	parts := strings.Split(pattern, "*")
	if len(parts) == 2 {
		return strings.HasPrefix(model, parts[0]) && strings.HasSuffix(model, parts[1])
	}

	return false
}

// FindByUserID finds all routing rules for a specific user
func (r *RoutingRuleRepository) FindByUserID(userID string) ([]*entity.RoutingRule, error) {
	var rules []*entity.RoutingRule
	err := r.db.Where("user_id = ?", userID).Order("priority ASC").Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}

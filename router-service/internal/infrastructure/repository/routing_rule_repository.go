package repository

import (
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

// Delete deletes a routing rule by ID
func (r *RoutingRuleRepository) Delete(id string) error {
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

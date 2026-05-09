package repository

import (
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
	"gorm.io/gorm"
)

type TierRepository struct {
	db *gorm.DB
}

func NewTierRepository(db *gorm.DB) port.TierRepository {
	return &TierRepository{db: db}
}

func (r *TierRepository) Create(tier *entity.Tier) error {
	return r.db.Create(tier).Error
}

func (r *TierRepository) GetByID(id string) (*entity.Tier, error) {
	var tier entity.Tier
	err := r.db.Where("id = ?", id).First(&tier).Error
	if err != nil {
		return nil, err
	}
	return &tier, nil
}

func (r *TierRepository) Update(tier *entity.Tier) error {
	return r.db.Save(tier).Error
}

func (r *TierRepository) Delete(id string) error {
	return r.db.Delete(&entity.Tier{}, "id = ?", id).Error
}

func (r *TierRepository) List(page, pageSize int) ([]*entity.Tier, int, error) {
	var tiers []*entity.Tier
	var total int64

	if err := r.db.Model(&entity.Tier{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Find(&tiers).Error
	if err != nil {
		return nil, 0, err
	}

	return tiers, int(total), nil
}

func (r *TierRepository) GetByName(name string) (*entity.Tier, error) {
	var tier entity.Tier
	err := r.db.Where("name = ?", name).First(&tier).Error
	if err != nil {
		return nil, err
	}
	return &tier, nil
}

func (r *TierRepository) GetDefaultTiers() ([]*entity.Tier, error) {
	var tiers []*entity.Tier
	err := r.db.Where("is_default = ?", true).Find(&tiers).Error
	if err != nil {
		return nil, err
	}
	return tiers, nil
}
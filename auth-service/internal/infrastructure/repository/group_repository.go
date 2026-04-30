package repository

import (
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
	"gorm.io/gorm"
)

// GroupRepository implements the GroupRepository interface using GORM
type GroupRepository struct {
	db *gorm.DB
}

// NewGroupRepository creates a new GroupRepository
func NewGroupRepository(db *gorm.DB) port.GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(group *entity.Group) error {
	return r.db.Create(group).Error
}

func (r *GroupRepository) GetByID(id string) (*entity.Group, error) {
	var group entity.Group
	err := r.db.Where("id = ?", id).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepository) Update(group *entity.Group) error {
	return r.db.Save(group).Error
}

func (r *GroupRepository) Delete(id string) error {
	return r.db.Delete(&entity.Group{}, "id = ?", id).Error
}

func (r *GroupRepository) List(page, pageSize int) ([]*entity.Group, int, error) {
	var groups []*entity.Group
	var total int64

	if err := r.db.Model(&entity.Group{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Find(&groups).Error
	if err != nil {
		return nil, 0, err
	}

	return groups, int(total), nil
}

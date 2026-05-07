package repository

import (
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
	"gorm.io/gorm"
)

// PermissionRepository implements the PermissionRepository interface using GORM
type PermissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new PermissionRepository
func NewPermissionRepository(db *gorm.DB) port.PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(permission *entity.Permission) error {
	return r.db.Create(permission).Error
}

func (r *PermissionRepository) GetByID(id string) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) Delete(id string) error {
	return r.db.Delete(&entity.Permission{}, "id = ?", id).Error
}

func (r *PermissionRepository) ListByGroupID(groupID string, page, pageSize int) ([]*entity.Permission, int, error) {
	var permissions []*entity.Permission
	var total int64

	query := r.db.Model(&entity.Permission{})
	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, int(total), nil
}

func (r *PermissionRepository) FindByUserGroups(groupIDs []string, resourceType, resourceID, action string) ([]*entity.Permission, error) {
	var permissions []*entity.Permission

	if len(groupIDs) == 0 {
		return permissions, nil
	}

	query := r.db.Where("group_id IN ?", groupIDs).
		Where("status = ?", "active")

	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}

	// resourceID matching: exact match or glob pattern
	if resourceID != "" {
		query = query.Where("resource_id = ? OR resource_id = ?", resourceID, "*")
	}

	err := query.Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

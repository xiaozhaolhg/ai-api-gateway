package repository

import (
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
	"gorm.io/gorm"
)

// UserGroupRepository implements the UserGroupRepository interface using GORM
type UserGroupRepository struct {
	db *gorm.DB
}

// NewUserGroupRepository creates a new UserGroupRepository
func NewUserGroupRepository(db *gorm.DB) port.UserGroupRepository {
	return &UserGroupRepository{db: db}
}

func (r *UserGroupRepository) Create(membership *entity.UserGroupMembership) error {
	return r.db.Create(membership).Error
}

func (r *UserGroupRepository) Delete(userID, groupID string) error {
	return r.db.Where("user_id = ? AND group_id = ?", userID, groupID).
		Delete(&entity.UserGroupMembership{}).Error
}

func (r *UserGroupRepository) GetByUserID(userID string) ([]*entity.UserGroupMembership, error) {
	var memberships []*entity.UserGroupMembership
	err := r.db.Where("user_id = ?", userID).Find(&memberships).Error
	if err != nil {
		return nil, err
	}
	return memberships, nil
}

func (r *UserGroupRepository) GetByGroupID(groupID string, page, pageSize int) ([]*entity.UserGroupMembership, int, error) {
	var memberships []*entity.UserGroupMembership
	var total int64

	if err := r.db.Model(&entity.UserGroupMembership{}).Where("group_id = ?", groupID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Where("group_id = ?", groupID).Offset(offset).Limit(pageSize).Find(&memberships).Error
	if err != nil {
		return nil, 0, err
	}

	return memberships, int(total), nil
}

func (r *UserGroupRepository) Exists(userID, groupID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.UserGroupMembership{}).
		Where("user_id = ? AND group_id = ?", userID, groupID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserGroupRepository) GetGroupIDsByUserID(userID string) ([]string, error) {
	var memberships []entity.UserGroupMembership
	err := r.db.Where("user_id = ?", userID).Find(&memberships).Error
	if err != nil {
		return nil, err
	}
	groupIDs := make([]string, len(memberships))
	for i, m := range memberships {
		groupIDs[i] = m.GroupID
	}
	return groupIDs, nil
}

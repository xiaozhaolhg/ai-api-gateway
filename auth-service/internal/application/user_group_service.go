package application

import (
	"fmt"
	"time"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
)

// UserGroupService provides user-group membership management logic
type UserGroupService struct {
	userGroupRepo port.UserGroupRepository
}

// NewUserGroupService creates a new UserGroupService
func NewUserGroupService(userGroupRepo port.UserGroupRepository) *UserGroupService {
	return &UserGroupService{userGroupRepo: userGroupRepo}
}

// AddUserToGroup adds a user to a group
func (s *UserGroupService) AddUserToGroup(userID, groupID string) error {
	// Check for duplicate
	exists, err := s.userGroupRepo.Exists(userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to check membership: %w", err)
	}
	if exists {
		return fmt.Errorf("user is already a member of this group")
	}

	membership := &entity.UserGroupMembership{
		ID:      generateID(),
		UserID:  userID,
		GroupID: groupID,
		AddedAt: time.Now(),
	}

	if err := s.userGroupRepo.Create(membership); err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}

	return nil
}

// RemoveUserFromGroup removes a user from a group
func (s *UserGroupService) RemoveUserFromGroup(userID, groupID string) error {
	if err := s.userGroupRepo.Delete(userID, groupID); err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}
	return nil
}

// GetUserGroups returns all group memberships for a user
func (s *UserGroupService) GetUserGroups(userID string) ([]*entity.UserGroupMembership, error) {
	memberships, err := s.userGroupRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}
	return memberships, nil
}

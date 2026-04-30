package application

import (
	"fmt"
	"time"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
)

// PermissionService provides permission management logic
type PermissionService struct {
	permissionRepo port.PermissionRepository
	userGroupRepo  port.UserGroupRepository
}

// NewPermissionService creates a new PermissionService
func NewPermissionService(permissionRepo port.PermissionRepository, userGroupRepo port.UserGroupRepository) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
		userGroupRepo:  userGroupRepo,
	}
}

// GrantPermission creates a new permission for a group
func (s *PermissionService) GrantPermission(groupID, resourceType, resourceID, action, effect string) (*entity.Permission, error) {
	permission := &entity.Permission{
		ID:           generateID(),
		GroupID:      groupID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Action:       action,
		Effect:       effect,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return nil, fmt.Errorf("failed to grant permission: %w", err)
	}

	return permission, nil
}

// RevokePermission deletes a permission
func (s *PermissionService) RevokePermission(id string) error {
	if err := s.permissionRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}
	return nil
}

// ListPermissions lists permissions for a group with pagination
func (s *PermissionService) ListPermissions(groupID string, page, pageSize int) ([]*entity.Permission, int, error) {
	permissions, total, err := s.permissionRepo.ListByGroupID(groupID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list permissions: %w", err)
	}
	return permissions, total, nil
}

// CheckPermission checks if a user has permission for a resource/action
// Deny effects override allow effects for the same resource_type, resource_id, and action
func (s *PermissionService) CheckPermission(userID, resourceType, resourceID, action string) (bool, error) {
	// Resolve user's groups
	memberships, err := s.userGroupRepo.GetByUserID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user groups: %w", err)
	}

	if len(memberships) == 0 {
		return false, nil
	}

	groupIDs := make([]string, len(memberships))
	for i, m := range memberships {
		groupIDs[i] = m.GroupID
	}

	// Find matching permissions
	permissions, err := s.permissionRepo.FindByUserGroups(groupIDs, resourceType, resourceID, action)
	if err != nil {
		return false, fmt.Errorf("failed to find permissions: %w", err)
	}

	hasAllow := false
	for _, p := range permissions {
		if p.Effect == "deny" {
			return false, nil
		}
		if p.Effect == "allow" {
			hasAllow = true
		}
	}

	return hasAllow, nil
}

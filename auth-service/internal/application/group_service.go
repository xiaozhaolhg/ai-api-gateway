package application

import (
	"fmt"
	"time"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
)

// GroupService provides group management logic
type GroupService struct {
	groupRepo port.GroupRepository
}

// NewGroupService creates a new GroupService
func NewGroupService(groupRepo port.GroupRepository) *GroupService {
	return &GroupService{groupRepo: groupRepo}
}

// CreateGroup creates a new group
func (s *GroupService) CreateGroup(name, description, parentGroupID string, modelPatterns []string, tokenLimit *entity.TokenLimit, rateLimit *entity.RateLimit) (*entity.Group, error) {
	group := &entity.Group{
		ID:            generateID(),
		Name:          name,
		Description:   description,
		ParentGroupID: parentGroupID,
		ModelPatterns: modelPatterns,
		TokenLimit:    tokenLimit,
		RateLimit:     rateLimit,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.groupRepo.Create(group); err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return group, nil
}

// UpdateGroup updates an existing group
func (s *GroupService) UpdateGroup(id, name, parentGroupID string) (*entity.Group, error) {
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}

	if name != "" {
		group.Name = name
	}
	if parentGroupID != "" {
		group.ParentGroupID = parentGroupID
	}
	group.UpdatedAt = time.Now()

	if err := s.groupRepo.Update(group); err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	return group, nil
}

// DeleteGroup deletes a group
func (s *GroupService) DeleteGroup(id string) error {
	if err := s.groupRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	return nil
}

// ListGroups lists groups with pagination
func (s *GroupService) ListGroups(page, pageSize int) ([]*entity.Group, int, error) {
	groups, total, err := s.groupRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list groups: %w", err)
	}
	return groups, total, nil
}

func (s *GroupService) GetGroupByID(id string) (*entity.Group, error) {
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}
	return group, nil
}

package application

import (
	"testing"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
)

type mockUGRepo struct {
	memberships []entity.UserGroupMembership
}

func newMockUGRepo() *mockUGRepo {
	return &mockUGRepo{memberships: make([]entity.UserGroupMembership, 0)}
}

func (m *mockUGRepo) GetByUserID(userID string) ([]*entity.UserGroupMembership, error) {
	var result []*entity.UserGroupMembership
	for i := range m.memberships {
		if m.memberships[i].UserID == userID {
			result = append(result, &m.memberships[i])
		}
	}
	return result, nil
}

func (m *mockUGRepo) GetByGroupID(groupID string, page, pageSize int) ([]*entity.UserGroupMembership, int, error) {
	var filtered []*entity.UserGroupMembership
	for i := range m.memberships {
		if m.memberships[i].GroupID == groupID {
			filtered = append(filtered, &m.memberships[i])
		}
	}
	total := len(filtered)
	offset := (page - 1) * pageSize
	if offset >= total {
		return nil, total, nil
	}
	end := offset + pageSize
	if end > total {
		end = total
	}
	return filtered[offset:end], total, nil
}

func (m *mockUGRepo) GetGroupIDsByUserID(userID string) ([]string, error) {
	var ids []string
	for _, mem := range m.memberships {
		if mem.UserID == userID {
			ids = append(ids, mem.GroupID)
		}
	}
	return ids, nil
}

func (m *mockUGRepo) Exists(userID, groupID string) (bool, error) {
	for _, mem := range m.memberships {
		if mem.UserID == userID && mem.GroupID == groupID {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockUGRepo) Create(membership *entity.UserGroupMembership) error {
	m.memberships = append(m.memberships, *membership)
	return nil
}

func (m *mockUGRepo) Delete(userID, groupID string) error {
	for i, mem := range m.memberships {
		if mem.UserID == userID && mem.GroupID == groupID {
			m.memberships = append(m.memberships[:i], m.memberships[i+1:]...)
			return nil
		}
	}
	return nil
}

func TestUserGroupService_AddUserToGroup(t *testing.T) {
	repo := newMockUGRepo()
	svc := NewUserGroupService(repo)

	err := svc.AddUserToGroup("user-1", "group-1")
	if err != nil {
		t.Fatalf("AddUserToGroup() error = %v", err)
	}

	memberships, _ := repo.GetByUserID("user-1")
	if len(memberships) != 1 {
		t.Errorf("Expected 1 membership, got %d", len(memberships))
	}
}

func TestUserGroupService_AddUserToGroup_Duplicate(t *testing.T) {
	repo := newMockUGRepo()
	svc := NewUserGroupService(repo)

	svc.AddUserToGroup("user-1", "group-1")

	err := svc.AddUserToGroup("user-1", "group-1")
	if err == nil {
		t.Error("Expected error for duplicate membership")
	}
}

func TestUserGroupService_RemoveUserFromGroup(t *testing.T) {
	repo := newMockUGRepo()
	svc := NewUserGroupService(repo)

	svc.AddUserToGroup("user-1", "group-1")

	err := svc.RemoveUserFromGroup("user-1", "group-1")
	if err != nil {
		t.Fatalf("RemoveUserFromGroup() error = %v", err)
	}

	memberships, _ := repo.GetByUserID("user-1")
	if len(memberships) != 0 {
		t.Errorf("Expected 0 memberships after removal, got %d", len(memberships))
	}
}

func TestUserGroupService_RemoveUserFromGroup_Idempotent(t *testing.T) {
	repo := newMockUGRepo()
	svc := NewUserGroupService(repo)

	err := svc.RemoveUserFromGroup("user-1", "group-1")
	if err != nil {
		t.Fatalf("RemoveUserFromGroup() for non-existent membership should succeed, got: %v", err)
	}
}

func TestUserGroupService_GetUserGroups(t *testing.T) {
	repo := newMockUGRepo()
	svc := NewUserGroupService(repo)

	svc.AddUserToGroup("user-1", "group-1")
	svc.AddUserToGroup("user-1", "group-2")

	memberships, err := svc.GetUserGroups("user-1")
	if err != nil {
		t.Fatalf("GetUserGroups() error = %v", err)
	}
	if len(memberships) != 2 {
		t.Errorf("Expected 2 memberships, got %d", len(memberships))
	}
}

func TestUserGroupService_GetGroupMembers(t *testing.T) {
	repo := newMockUGRepo()
	svc := NewUserGroupService(repo)

	svc.AddUserToGroup("user-1", "group-1")
	svc.AddUserToGroup("user-2", "group-1")
	svc.AddUserToGroup("user-3", "group-1")

	memberships, total, err := svc.GetGroupMembers("group-1", 1, 10)
	if err != nil {
		t.Fatalf("GetGroupMembers() error = %v", err)
	}
	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(memberships) != 3 {
		t.Errorf("Expected 3 memberships, got %d", len(memberships))
	}
}

func TestUserGroupService_GetGroupMembers_Pagination(t *testing.T) {
	repo := newMockUGRepo()
	svc := NewUserGroupService(repo)

	svc.AddUserToGroup("user-1", "group-1")
	svc.AddUserToGroup("user-2", "group-1")
	svc.AddUserToGroup("user-3", "group-1")

	memberships, total, err := svc.GetGroupMembers("group-1", 1, 2)
	if err != nil {
		t.Fatalf("GetGroupMembers() error = %v", err)
	}
	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(memberships) != 2 {
		t.Errorf("Expected 2 memberships (page size), got %d", len(memberships))
	}
}

func TestUserGroupService_GetGroupIDsByUserID(t *testing.T) {
	repo := newMockUGRepo()

	repo.memberships = append(repo.memberships, entity.UserGroupMembership{UserID: "user-1", GroupID: "group-1"})
	repo.memberships = append(repo.memberships, entity.UserGroupMembership{UserID: "user-1", GroupID: "group-2"})
	repo.memberships = append(repo.memberships, entity.UserGroupMembership{UserID: "user-1", GroupID: "group-3"})

	groupIDs, err := repo.GetGroupIDsByUserID("user-1")
	if err != nil {
		t.Fatalf("GetGroupIDsByUserID() error = %v", err)
	}
	if len(groupIDs) != 3 {
		t.Errorf("Expected 3 group IDs, got %d", len(groupIDs))
	}
}

func TestUserGroupService_GetUserGroups_NoMemberships(t *testing.T) {
	repo := newMockUGRepo()
	svc := NewUserGroupService(repo)

	memberships, err := svc.GetUserGroups("user-no-groups")
	if err != nil {
		t.Fatalf("GetUserGroups() error = %v", err)
	}
	if len(memberships) != 0 {
		t.Errorf("Expected 0 memberships for user with no groups, got %d", len(memberships))
	}
}

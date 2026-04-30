package application

import (
	"testing"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"gorm.io/gorm"
)

type mockPermissionRepo struct {
	permissions map[string]*entity.Permission
}

func newMockPermissionRepo() *mockPermissionRepo {
	return &mockPermissionRepo{permissions: make(map[string]*entity.Permission)}
}

func (m *mockPermissionRepo) Create(permission *entity.Permission) error {
	m.permissions[permission.ID] = permission
	return nil
}

func (m *mockPermissionRepo) GetByID(id string) (*entity.Permission, error) {
	p, ok := m.permissions[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return p, nil
}

func (m *mockPermissionRepo) Delete(id string) error {
	delete(m.permissions, id)
	return nil
}

func (m *mockPermissionRepo) ListByGroupID(groupID string, page, pageSize int) ([]*entity.Permission, int, error) {
	var all []*entity.Permission
	for _, p := range m.permissions {
		if p.GroupID == groupID {
			all = append(all, p)
		}
	}
	total := len(all)
	offset := (page - 1) * pageSize
	if offset >= total {
		return nil, total, nil
	}
	end := offset + pageSize
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

func (m *mockPermissionRepo) FindByUserGroups(groupIDs []string, resourceType, resourceID, action string) ([]*entity.Permission, error) {
	var result []*entity.Permission
	for _, p := range m.permissions {
		if p.Status != "active" {
			continue
		}
		found := false
		for _, gid := range groupIDs {
			if p.GroupID == gid {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		if resourceType != "" && p.ResourceType != resourceType {
			continue
		}
		if action != "" && p.Action != action {
			continue
		}
		if resourceID != "" && p.ResourceID != resourceID && p.ResourceID != "*" {
			continue
		}
		result = append(result, p)
	}
	return result, nil
}

type mockUserGroupRepo struct {
	memberships [] *entity.UserGroupMembership
}

func newMockUserGroupRepo() *mockUserGroupRepo {
	return &mockUserGroupRepo{}
}

func (m *mockUserGroupRepo) Create(membership *entity.UserGroupMembership) error {
	m.memberships = append(m.memberships, membership)
	return nil
}

func (m *mockUserGroupRepo) Delete(userID, groupID string) error {
	for i, mg := range m.memberships {
		if mg.UserID == userID && mg.GroupID == groupID {
			m.memberships = append(m.memberships[:i], m.memberships[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockUserGroupRepo) GetByUserID(userID string) ([]*entity.UserGroupMembership, error) {
	var result []*entity.UserGroupMembership
	for _, mg := range m.memberships {
		if mg.UserID == userID {
			result = append(result, mg)
		}
	}
	return result, nil
}

func (m *mockUserGroupRepo) GetByGroupID(groupID string, page, pageSize int) ([]*entity.UserGroupMembership, int, error) {
	var result []*entity.UserGroupMembership
	for _, mg := range m.memberships {
		if mg.GroupID == groupID {
			result = append(result, mg)
		}
	}
	return result, len(result), nil
}

func (m *mockUserGroupRepo) Exists(userID, groupID string) (bool, error) {
	for _, mg := range m.memberships {
		if mg.UserID == userID && mg.GroupID == groupID {
			return true, nil
		}
	}
	return false, nil
}

func TestPermissionService_GrantPermission(t *testing.T) {
	permRepo := newMockPermissionRepo()
	ugRepo := newMockUserGroupRepo()
	svc := NewPermissionService(permRepo, ugRepo)

	perm, err := svc.GrantPermission("group-1", "model", "gpt-4", "access", "allow")
	if err != nil {
		t.Fatalf("GrantPermission() error = %v", err)
	}
	if perm.GroupID != "group-1" || perm.Effect != "allow" {
		t.Errorf("Unexpected permission: %+v", perm)
	}
	if perm.Status != "active" {
		t.Errorf("Expected status 'active', got %s", perm.Status)
	}
}

func TestPermissionService_RevokePermission(t *testing.T) {
	permRepo := newMockPermissionRepo()
	ugRepo := newMockUserGroupRepo()
	svc := NewPermissionService(permRepo, ugRepo)

	perm, _ := svc.GrantPermission("group-1", "model", "gpt-4", "access", "allow")

	err := svc.RevokePermission(perm.ID)
	if err != nil {
		t.Fatalf("RevokePermission() error = %v", err)
	}

	_, err = permRepo.GetByID(perm.ID)
	if err == nil {
		t.Error("Expected error after revocation")
	}
}

func TestPermissionService_CheckPermission_AllowOnly(t *testing.T) {
	permRepo := newMockPermissionRepo()
	ugRepo := newMockUserGroupRepo()
	svc := NewPermissionService(permRepo, ugRepo)

	// User is in group-1
	ugRepo.Create(&entity.UserGroupMembership{ID: "ug-1", UserID: "user-1", GroupID: "group-1"})

	// Group-1 has allow permission
	svc.GrantPermission("group-1", "model", "gpt-4", "access", "allow")

	allowed, err := svc.CheckPermission("user-1", "model", "gpt-4", "access")
	if err != nil {
		t.Fatalf("CheckPermission() error = %v", err)
	}
	if !allowed {
		t.Error("Expected allowed=true")
	}
}

func TestPermissionService_CheckPermission_DenyOverride(t *testing.T) {
	permRepo := newMockPermissionRepo()
	ugRepo := newMockUserGroupRepo()
	svc := NewPermissionService(permRepo, ugRepo)

	// User is in both group-1 and group-2
	ugRepo.Create(&entity.UserGroupMembership{ID: "ug-1", UserID: "user-1", GroupID: "group-1"})
	ugRepo.Create(&entity.UserGroupMembership{ID: "ug-2", UserID: "user-1", GroupID: "group-2"})

	// group-1 allows
	svc.GrantPermission("group-1", "model", "gpt-4", "access", "allow")
	// group-2 denies
	svc.GrantPermission("group-2", "model", "gpt-4", "access", "deny")

	allowed, err := svc.CheckPermission("user-1", "model", "gpt-4", "access")
	if err != nil {
		t.Fatalf("CheckPermission() error = %v", err)
	}
	if allowed {
		t.Error("Expected allowed=false (deny overrides allow)")
	}
}

func TestPermissionService_CheckPermission_NoGroups(t *testing.T) {
	permRepo := newMockPermissionRepo()
	ugRepo := newMockUserGroupRepo()
	svc := NewPermissionService(permRepo, ugRepo)

	allowed, err := svc.CheckPermission("user-no-groups", "model", "gpt-4", "access")
	if err != nil {
		t.Fatalf("CheckPermission() error = %v", err)
	}
	if allowed {
		t.Error("Expected allowed=false for user with no groups")
	}
}

package application

import (
	"testing"
)

func TestUserGroupService_AddUserToGroup(t *testing.T) {
	repo := newMockUserGroupRepo()
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
	repo := newMockUserGroupRepo()
	svc := NewUserGroupService(repo)

	svc.AddUserToGroup("user-1", "group-1")

	err := svc.AddUserToGroup("user-1", "group-1")
	if err == nil {
		t.Error("Expected error for duplicate membership")
	}
}

func TestUserGroupService_RemoveUserFromGroup(t *testing.T) {
	repo := newMockUserGroupRepo()
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
	repo := newMockUserGroupRepo()
	svc := NewUserGroupService(repo)

	// Removing non-existent membership should succeed (idempotent)
	err := svc.RemoveUserFromGroup("user-1", "group-1")
	if err != nil {
		t.Fatalf("RemoveUserFromGroup() for non-existent membership should succeed, got: %v", err)
	}
}

func TestUserGroupService_GetUserGroups(t *testing.T) {
	repo := newMockUserGroupRepo()
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

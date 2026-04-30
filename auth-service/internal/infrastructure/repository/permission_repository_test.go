package repository

import (
	"testing"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"gorm.io/gorm"
)

func setupPermissionTestDB(t *testing.T) *gorm.DB {
	db := setupTestDB(t)
	err := db.AutoMigrate(&entity.Permission{})
	if err != nil {
		t.Fatalf("Failed to migrate Permission: %v", err)
	}
	return db
}

func TestPermissionRepository_Create(t *testing.T) {
	db := setupPermissionTestDB(t)
	repo := NewPermissionRepository(db)

	perm := &entity.Permission{
		ID:           "perm-1",
		GroupID:      "group-1",
		ResourceType: "model",
		ResourceID:   "gpt-4",
		Action:       "access",
		Effect:       "allow",
		Status:       "active",
	}

	err := repo.Create(perm)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	retrieved, err := repo.GetByID(perm.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if retrieved.ResourceType != "model" || retrieved.Effect != "allow" {
		t.Errorf("Unexpected permission: %+v", retrieved)
	}
}

func TestPermissionRepository_Delete(t *testing.T) {
	db := setupPermissionTestDB(t)
	repo := NewPermissionRepository(db)

	perm := &entity.Permission{
		ID:           "perm-2",
		GroupID:      "group-1",
		ResourceType: "model",
		ResourceID:   "gpt-4",
		Action:       "access",
		Effect:       "allow",
		Status:       "active",
	}
	repo.Create(perm)

	err := repo.Delete(perm.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.GetByID(perm.ID)
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}
}

func TestPermissionRepository_ListByGroupID(t *testing.T) {
	db := setupPermissionTestDB(t)
	repo := NewPermissionRepository(db)

	// Create permissions for two groups
	for i := 0; i < 3; i++ {
		perm := &entity.Permission{
			ID:           "perm-group1-" + string(rune('a'+i)),
			GroupID:      "group-1",
			ResourceType: "model",
			ResourceID:   "model-" + string(rune('a'+i)),
			Action:       "access",
			Effect:       "allow",
			Status:       "active",
		}
		repo.Create(perm)
	}

	perm := &entity.Permission{
		ID:           "perm-group2-1",
		GroupID:      "group-2",
		ResourceType: "model",
		ResourceID:   "gpt-4",
		Action:       "access",
		Effect:       "allow",
		Status:       "active",
	}
	repo.Create(perm)

	perms, total, err := repo.ListByGroupID("group-1", 1, 10)
	if err != nil {
		t.Fatalf("ListByGroupID() error = %v", err)
	}
	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(perms) != 3 {
		t.Errorf("Expected 3 permissions, got %d", len(perms))
	}
}

func TestPermissionRepository_FindByUserGroups(t *testing.T) {
	db := setupPermissionTestDB(t)
	repo := NewPermissionRepository(db)

	// Create allow permission
	perm1 := &entity.Permission{
		ID:           "perm-find-1",
		GroupID:      "group-a",
		ResourceType: "model",
		ResourceID:   "gpt-4",
		Action:       "access",
		Effect:       "allow",
		Status:       "active",
	}
	repo.Create(perm1)

	// Create deny permission
	perm2 := &entity.Permission{
		ID:           "perm-find-2",
		GroupID:      "group-b",
		ResourceType: "model",
		ResourceID:   "gpt-4",
		Action:       "access",
		Effect:       "deny",
		Status:       "active",
	}
	repo.Create(perm2)

	// Create wildcard permission
	perm3 := &entity.Permission{
		ID:           "perm-find-3",
		GroupID:      "group-a",
		ResourceType: "model",
		ResourceID:   "*",
		Action:       "access",
		Effect:       "allow",
		Status:       "active",
	}
	repo.Create(perm3)

	// Find for user in both groups
	perms, err := repo.FindByUserGroups([]string{"group-a", "group-b"}, "model", "gpt-4", "access")
	if err != nil {
		t.Fatalf("FindByUserGroups() error = %v", err)
	}
	if len(perms) < 2 {
		t.Errorf("Expected at least 2 permissions, got %d", len(perms))
	}

	// Empty group IDs
	perms, err = repo.FindByUserGroups([]string{}, "model", "gpt-4", "access")
	if err != nil {
		t.Fatalf("FindByUserGroups() error = %v", err)
	}
	if len(perms) != 0 {
		t.Errorf("Expected 0 permissions for empty groups, got %d", len(perms))
	}
}

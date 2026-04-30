package repository

import (
	"testing"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"gorm.io/gorm"
)

func setupUserGroupTestDB(t *testing.T) *gorm.DB {
	db := setupTestDB(t)
	err := db.AutoMigrate(&entity.UserGroupMembership{})
	if err != nil {
		t.Fatalf("Failed to migrate UserGroupMembership: %v", err)
	}
	return db
}

func TestUserGroupRepository_Create(t *testing.T) {
	db := setupUserGroupTestDB(t)
	repo := NewUserGroupRepository(db)

	membership := &entity.UserGroupMembership{
		ID:      "ug-1",
		UserID:  "user-1",
		GroupID: "group-1",
	}

	err := repo.Create(membership)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	exists, err := repo.Exists("user-1", "group-1")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Expected membership to exist")
	}
}

func TestUserGroupRepository_Delete(t *testing.T) {
	db := setupUserGroupTestDB(t)
	repo := NewUserGroupRepository(db)

	membership := &entity.UserGroupMembership{
		ID:      "ug-2",
		UserID:  "user-2",
		GroupID: "group-2",
	}
	repo.Create(membership)

	err := repo.Delete("user-2", "group-2")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	exists, _ := repo.Exists("user-2", "group-2")
	if exists {
		t.Error("Expected membership to be deleted")
	}
}

func TestUserGroupRepository_GetByUserID(t *testing.T) {
	db := setupUserGroupTestDB(t)
	repo := NewUserGroupRepository(db)

	// User in two groups
	repo.Create(&entity.UserGroupMembership{ID: "ug-3a", UserID: "user-3", GroupID: "group-a"})
	repo.Create(&entity.UserGroupMembership{ID: "ug-3b", UserID: "user-3", GroupID: "group-b"})

	memberships, err := repo.GetByUserID("user-3")
	if err != nil {
		t.Fatalf("GetByUserID() error = %v", err)
	}
	if len(memberships) != 2 {
		t.Errorf("Expected 2 memberships, got %d", len(memberships))
	}
}

func TestUserGroupRepository_GetByGroupID(t *testing.T) {
	db := setupUserGroupTestDB(t)
	repo := NewUserGroupRepository(db)

	// Two users in same group
	repo.Create(&entity.UserGroupMembership{ID: "ug-4a", UserID: "user-4a", GroupID: "group-4"})
	repo.Create(&entity.UserGroupMembership{ID: "ug-4b", UserID: "user-4b", GroupID: "group-4"})

	memberships, total, err := repo.GetByGroupID("group-4", 1, 10)
	if err != nil {
		t.Fatalf("GetByGroupID() error = %v", err)
	}
	if total != 2 {
		t.Errorf("Expected total 2, got %d", total)
	}
	if len(memberships) != 2 {
		t.Errorf("Expected 2 memberships, got %d", len(memberships))
	}
}

func TestUserGroupRepository_Exists(t *testing.T) {
	db := setupUserGroupTestDB(t)
	repo := NewUserGroupRepository(db)

	repo.Create(&entity.UserGroupMembership{ID: "ug-5", UserID: "user-5", GroupID: "group-5"})

	exists, err := repo.Exists("user-5", "group-5")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Expected membership to exist")
	}

	exists, err = repo.Exists("user-5", "nonexistent-group")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Error("Expected membership to not exist")
	}
}

package repository

import (
	"testing"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/infrastructure/migration"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dbPath := ":memory:"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Run migrations
	if err := migration.Migrate(dbPath); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entity.User{
		ID:        "test-user-1",
		Name:      "Test User",
		Email:     "test@example.com",
		Role:      "user",
		Status:    "active",
	}

	err := repo.Create(user)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// Verify user was created
	retrieved, err := repo.GetByID(user.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if retrieved.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, retrieved.Name)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Test non-existent user
	_, err := repo.GetByID("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}

	// Create a user and retrieve it
	user := &entity.User{
		ID:        "test-user-2",
		Name:      "Test User 2",
		Email:     "test2@example.com",
		Role:      "admin",
		Status:    "active",
	}
	repo.Create(user)

	retrieved, err := repo.GetByID(user.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if retrieved.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrieved.ID)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Test non-existent email
	_, err := repo.GetByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("Expected error for non-existent email, got nil")
	}

	// Create a user and retrieve by email
	user := &entity.User{
		ID:        "test-user-3",
		Name:      "Test User 3",
		Email:     "test3@example.com",
		Role:      "user",
		Status:    "active",
	}
	repo.Create(user)

	retrieved, err := repo.GetByEmail(user.Email)
	if err != nil {
		t.Errorf("GetByEmail() error = %v", err)
	}
	if retrieved.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved.Email)
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entity.User{
		ID:        "test-user-4",
		Name:      "Test User 4",
		Email:     "test4@example.com",
		Role:      "user",
		Status:    "active",
	}
	repo.Create(user)

	// Update user
	user.Name = "Updated Name"
	user.Status = "disabled"
	err := repo.Update(user)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	// Verify update
	retrieved, err := repo.GetByID(user.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if retrieved.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", retrieved.Name)
	}
	if retrieved.Status != "disabled" {
		t.Errorf("Expected status 'disabled', got %s", retrieved.Status)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &entity.User{
		ID:        "test-user-5",
		Name:      "Test User 5",
		Email:     "test5@example.com",
		Role:      "user",
		Status:    "active",
	}
	repo.Create(user)

	// Delete user
	err := repo.Delete(user.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(user.ID)
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create multiple users
	for i := 0; i < 5; i++ {
		user := &entity.User{
			ID:     "test-user-list-" + string(rune(i)),
			Name:   "Test User " + string(rune(i)),
			Email:  "testlist" + string(rune(i)) + "@example.com",
			Role:   "user",
			Status: "active",
		}
		repo.Create(user)
	}

	// Test pagination
	users, total, err := repo.List(1, 3)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}
	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}

	// Test second page
	users2, total2, err := repo.List(2, 3)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}
	if total2 != 5 {
		t.Errorf("Expected total 5, got %d", total2)
	}
	if len(users2) != 2 {
		t.Errorf("Expected 2 users on second page, got %d", len(users2))
	}
}

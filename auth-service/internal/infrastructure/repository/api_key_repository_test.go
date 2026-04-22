package repository

import (
	"testing"
	"time"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForAPIKey(t *testing.T) *gorm.DB {
	dbPath := ":memory:"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&entity.User{}, &entity.APIKey{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestAPIKeyRepository_Create(t *testing.T) {
	db := setupTestDBForAPIKey(t)
	repo := NewAPIKeyRepository(db)

	apiKey := &entity.APIKey{
		ID:        "test-key-1",
		UserID:    "user-1",
		KeyHash:   "hash123",
		Name:      "Test Key",
		Scopes:    []string{"read", "write"},
		CreatedAt: time.Now(),
	}

	err := repo.Create(apiKey)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// Verify API key was created
	retrieved, err := repo.GetByID(apiKey.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if retrieved.Name != apiKey.Name {
		t.Errorf("Expected name %s, got %s", apiKey.Name, retrieved.Name)
	}
}

func TestAPIKeyRepository_GetByID(t *testing.T) {
	db := setupTestDBForAPIKey(t)
	repo := NewAPIKeyRepository(db)

	// Test non-existent key
	_, err := repo.GetByID("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent key, got nil")
	}

	// Create an API key and retrieve it
	apiKey := &entity.APIKey{
		ID:        "test-key-2",
		UserID:    "user-2",
		KeyHash:   "hash456",
		Name:      "Test Key 2",
		Scopes:    []string{"read"},
		CreatedAt: time.Now(),
	}
	repo.Create(apiKey)

	retrieved, err := repo.GetByID(apiKey.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if retrieved.ID != apiKey.ID {
		t.Errorf("Expected ID %s, got %s", apiKey.ID, retrieved.ID)
	}
}

func TestAPIKeyRepository_GetByKeyHash(t *testing.T) {
	db := setupTestDBForAPIKey(t)
	repo := NewAPIKeyRepository(db)

	// Test non-existent hash
	_, err := repo.GetByKeyHash("nonexistent-hash")
	if err == nil {
		t.Error("Expected error for non-existent hash, got nil")
	}

	// Create an API key and retrieve by hash
	apiKey := &entity.APIKey{
		ID:        "test-key-3",
		UserID:    "user-3",
		KeyHash:   "hash789",
		Name:      "Test Key 3",
		Scopes:    []string{"read", "write"},
		CreatedAt: time.Now(),
	}
	repo.Create(apiKey)

	retrieved, err := repo.GetByKeyHash(apiKey.KeyHash)
	if err != nil {
		t.Errorf("GetByKeyHash() error = %v", err)
	}
	if retrieved.KeyHash != apiKey.KeyHash {
		t.Errorf("Expected hash %s, got %s", apiKey.KeyHash, retrieved.KeyHash)
	}
}

func TestAPIKeyRepository_GetByUserID(t *testing.T) {
	db := setupTestDBForAPIKey(t)
	repo := NewAPIKeyRepository(db)

	userID := "user-4"

	// Create multiple API keys for the same user
	for i := 0; i < 5; i++ {
		apiKey := &entity.APIKey{
			ID:        "test-key-user-" + string(rune(i)),
			UserID:    userID,
			KeyHash:   "hash" + string(rune(i)),
			Name:      "Test Key " + string(rune(i)),
			Scopes:    []string{"read"},
			CreatedAt: time.Now(),
		}
		repo.Create(apiKey)
	}

	// Create an API key for a different user
	otherKey := &entity.APIKey{
		ID:        "test-key-other",
		UserID:    "other-user",
		KeyHash:   "hash-other",
		Name:      "Other Key",
		Scopes:    []string{"read"},
		CreatedAt: time.Now(),
	}
	repo.Create(otherKey)

	// Test pagination for user's keys
	keys, total, err := repo.GetByUserID(userID, 1, 3)
	if err != nil {
		t.Errorf("GetByUserID() error = %v", err)
	}
	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Verify all keys belong to the user
	for _, key := range keys {
		if key.UserID != userID {
			t.Errorf("Expected key to belong to user %s, got %s", userID, key.UserID)
		}
	}
}

func TestAPIKeyRepository_Delete(t *testing.T) {
	db := setupTestDBForAPIKey(t)
	repo := NewAPIKeyRepository(db)

	apiKey := &entity.APIKey{
		ID:        "test-key-5",
		UserID:    "user-5",
		KeyHash:   "hash-delete",
		Name:      "Test Key 5",
		Scopes:    []string{"read"},
		CreatedAt: time.Now(),
	}
	repo.Create(apiKey)

	// Delete API key
	err := repo.Delete(apiKey.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(apiKey.ID)
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}
}

func TestAPIKeyRepository_Expiration(t *testing.T) {
	db := setupTestDBForAPIKey(t)
	repo := NewAPIKeyRepository(db)

	// Create an expired API key
	expiredTime := time.Now().Add(-1 * time.Hour)
	expiredKey := &entity.APIKey{
		ID:        "test-key-expired",
		UserID:    "user-6",
		KeyHash:   "hash-expired",
		Name:      "Expired Key",
		Scopes:    []string{"read"},
		CreatedAt: time.Now().Add(-2 * time.Hour),
		ExpiresAt: &expiredTime,
	}
	repo.Create(expiredKey)

	// Create a non-expired API key
	futureTime := time.Now().Add(1 * time.Hour)
	validKey := &entity.APIKey{
		ID:        "test-key-valid",
		UserID:    "user-6",
		KeyHash:   "hash-valid",
		Name:      "Valid Key",
		Scopes:    []string{"read"},
		CreatedAt: time.Now(),
		ExpiresAt: &futureTime,
	}
	repo.Create(validKey)

	// Retrieve both keys
	retrievedExpired, err := repo.GetByID(expiredKey.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if retrievedExpired.ExpiresAt == nil {
		t.Error("Expected ExpiresAt to be set for expired key")
	}

	retrievedValid, err := repo.GetByID(validKey.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if retrievedValid.ExpiresAt == nil {
		t.Error("Expected ExpiresAt to be set for valid key")
	}
}

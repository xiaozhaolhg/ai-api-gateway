package repository

import (
	"testing"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"gorm.io/gorm"
)

func setupGroupTestDB(t *testing.T) *gorm.DB {
	db := setupTestDB(t)
	err := db.AutoMigrate(&entity.Group{})
	if err != nil {
		t.Fatalf("Failed to migrate Group: %v", err)
	}
	return db
}

func TestGroupRepository_Create(t *testing.T) {
	db := setupGroupTestDB(t)
	repo := NewGroupRepository(db)

	group := &entity.Group{
		ID:          "group-1",
		Name:        "developers",
		Description: "Developer team",
	}

	err := repo.Create(group)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	retrieved, err := repo.GetByID(group.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if retrieved.Name != "developers" {
		t.Errorf("Expected name 'developers', got %s", retrieved.Name)
	}
}

func TestGroupRepository_CreateWithModelPatternsAndLimits(t *testing.T) {
	db := setupGroupTestDB(t)
	repo := NewGroupRepository(db)

	group := &entity.Group{
		ID:            "group-2",
		Name:          "power-users",
		ModelPatterns: []string{"gpt-4", "claude-*"},
		TokenLimit: &entity.TokenLimit{
			PromptTokens:     100000,
			CompletionTokens: 100000,
			Period:           "daily",
		},
		RateLimit: &entity.RateLimit{
			RequestsPerMinute: 60,
			RequestsPerDay:    10000,
		},
	}

	err := repo.Create(group)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	retrieved, err := repo.GetByID(group.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if len(retrieved.ModelPatterns) != 2 {
		t.Errorf("Expected 2 model patterns, got %d", len(retrieved.ModelPatterns))
	}
	if retrieved.TokenLimit.PromptTokens != 100000 {
		t.Errorf("Expected prompt_tokens 100000, got %d", retrieved.TokenLimit.PromptTokens)
	}
	if retrieved.RateLimit.RequestsPerMinute != 60 {
		t.Errorf("Expected requests_per_minute 60, got %d", retrieved.RateLimit.RequestsPerMinute)
	}
}

func TestGroupRepository_Update(t *testing.T) {
	db := setupGroupTestDB(t)
	repo := NewGroupRepository(db)

	group := &entity.Group{
		ID:   "group-3",
		Name: "developers",
	}
	repo.Create(group)

	group.Name = "senior-devs"
	err := repo.Update(group)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	retrieved, _ := repo.GetByID(group.ID)
	if retrieved.Name != "senior-devs" {
		t.Errorf("Expected name 'senior-devs', got %s", retrieved.Name)
	}
}

func TestGroupRepository_Delete(t *testing.T) {
	db := setupGroupTestDB(t)
	repo := NewGroupRepository(db)

	group := &entity.Group{
		ID:   "group-4",
		Name: "to-delete",
	}
	repo.Create(group)

	err := repo.Delete(group.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.GetByID(group.ID)
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}
}

func TestGroupRepository_List(t *testing.T) {
	db := setupGroupTestDB(t)
	repo := NewGroupRepository(db)

	for i := 0; i < 5; i++ {
		group := &entity.Group{
			ID:   "group-list-" + string(rune('a'+i)),
			Name: "group-" + string(rune('a'+i)),
		}
		repo.Create(group)
	}

	groups, total, err := repo.List(1, 3)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
	if len(groups) != 3 {
		t.Errorf("Expected 3 groups, got %d", len(groups))
	}

	groups2, total2, _ := repo.List(2, 3)
	if total2 != 5 {
		t.Errorf("Expected total 5, got %d", total2)
	}
	if len(groups2) != 2 {
		t.Errorf("Expected 2 groups on second page, got %d", len(groups2))
	}
}

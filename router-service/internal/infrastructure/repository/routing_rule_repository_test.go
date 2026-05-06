package repository

import (
	"testing"

	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*RoutingRuleRepository, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	db.AutoMigrate(&entity.RoutingRule{})
	repo := &RoutingRuleRepository{db: db}
	return repo, db
}

func TestFindByModel_UserRule(t *testing.T) {
	repo, _ := setupTestDB(t)

	userID := "user-123"
	userRule := &entity.RoutingRule{
		ID:           "rule-user",
		UserID:       userID,
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
	}
	repo.Create(userRule)

	systemRule := &entity.RoutingRule{
		ID:           "rule-system",
		ModelPattern: "ollama:*",
		ProviderID:   "opencode_zen",
		IsSystemDefault: true,
	}
	repo.Create(systemRule)

	result, err := repo.FindByModel("ollama:*", &userID)
	if err != nil {
		t.Fatalf("FindByModel failed: %v", err)
	}
	if result.ID != "rule-user" {
		t.Errorf("Expected user rule, got %s", result.ID)
	}
}

func TestFindByModel_SystemRuleFallback(t *testing.T) {
	repo, _ := setupTestDB(t)

	systemRule := &entity.RoutingRule{
		ID:           "rule-system",
		ModelPattern: "ollama:*",
		ProviderID:   "opencode_zen",
		IsSystemDefault: true,
	}
	repo.Create(systemRule)

	userID := "user-123"
	result, err := repo.FindByModel("ollama:*", &userID)
	if err != nil {
		t.Fatalf("FindByModel failed: %v", err)
	}
	if result.ID != "rule-system" {
		t.Errorf("Expected system rule, got %s", result.ID)
	}
}

func TestFindByModel_NoUserID(t *testing.T) {
	repo, _ := setupTestDB(t)

	systemRule := &entity.RoutingRule{
		ID:           "rule-system",
		ModelPattern: "ollama:*",
		ProviderID:   "opencode_zen",
		IsSystemDefault: true,
	}
	repo.Create(systemRule)

	result, err := repo.FindByModel("ollama:*", nil)
	if err != nil {
		t.Fatalf("FindByModel failed: %v", err)
	}
	if result.ID != "rule-system" {
		t.Errorf("Expected system rule, got %s", result.ID)
	}
}

func TestFindByUserID(t *testing.T) {
	repo, _ := setupTestDB(t)

	userID := "user-123"
	rules := []*entity.RoutingRule{
		{ID: "rule-1", UserID: userID, ModelPattern: "ollama:*", ProviderID: "ollama"},
		{ID: "rule-2", UserID: userID, ModelPattern: "openai:*", ProviderID: "openai"},
		{ID: "rule-3", UserID: "other-user", ModelPattern: "anthropic:*", ProviderID: "anthropic"},
	}
	for _, r := range rules {
		repo.Create(r)
	}

	results, err := repo.FindByUserID(userID)
	if err != nil {
		t.Fatalf("FindByUserID failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(results))
	}
}

func TestUpdateWithOwnership_Authorized(t *testing.T) {
	repo, _ := setupTestDB(t)

	rule := &entity.RoutingRule{
		ID:           "rule-1",
		UserID:       "user-123",
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
	}
	repo.Create(rule)

	rule.ProviderID = "opencode_zen"
	err := repo.UpdateWithOwnership(rule, "user-123")
	if err != nil {
		t.Fatalf("UpdateWithOwnership failed: %v", err)
	}
}

func TestUpdateWithOwnership_Unauthorized(t *testing.T) {
	repo, _ := setupTestDB(t)

	rule := &entity.RoutingRule{
		ID:           "rule-1",
		UserID:       "user-123",
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
	}
	repo.Create(rule)

	rule.ProviderID = "opencode_zen"
	err := repo.UpdateWithOwnership(rule, "user-456")
	if err == nil {
		t.Error("Expected error for unauthorized update")
	}
}

func TestDeleteWithOwnership_Authorized(t *testing.T) {
	repo, _ := setupTestDB(t)

	rule := &entity.RoutingRule{
		ID:           "rule-1",
		UserID:       "user-123",
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
	}
	repo.Create(rule)

	err := repo.DeleteWithOwnership("rule-1", "user-123")
	if err != nil {
		t.Fatalf("DeleteWithOwnership failed: %v", err)
	}
}

func TestDeleteWithOwnership_Unauthorized(t *testing.T) {
	repo, _ := setupTestDB(t)

	rule := &entity.RoutingRule{
		ID:           "rule-1",
		UserID:       "user-123",
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
	}
	repo.Create(rule)

	err := repo.DeleteWithOwnership("rule-1", "user-456")
	if err == nil {
		t.Error("Expected error for unauthorized delete")
	}
}

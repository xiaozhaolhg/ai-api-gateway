package application;

import (
	"context"
	"testing"

	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	"github.com/ai-api-gateway/router-service/internal/infrastructure/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestService(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	db.AutoMigrate(&entity.RoutingRule{})
	repo := repository.NewRoutingRuleRepository(db)
	svc := NewService(repo, nil)
	return svc, db
}

func TestResolveRoute_UserRuleOverride(t *testing.T) {
	svc, _ := setupTestService(t)

	// Create system rule
	systemRule := &entity.RoutingRule{
		ID:           "system-rule",
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
		IsSystemDefault: true,
	}
	svc.ruleRepo.Create(systemRule)

	// Create user rule that should override
	userRule := &entity.RoutingRule{
		ID:           "user-rule",
		UserID:       "user-123",
		ModelPattern: "ollama:*",
		ProviderID:   "opencode_zen",
	}
	svc.ruleRepo.Create(userRule)

	// Test: user rule should override system rule
	result, err := svc.ResolveRoute(context.Background(), "ollama:llama2", []string{}, "user-123")
	if err != nil {
		t.Fatalf("ResolveRoute failed: %v", err)
	}

	if result.ProviderID != "opencode_zen" {
		t.Errorf("Expected user rule provider 'opencode_zen', got '%s'", result.ProviderID)
	}
}

func TestResolveRoute_SystemRuleFallback(t *testing.T) {
	svc, _ := setupTestService(t)

	// Create only system rule
	systemRule := &entity.RoutingRule{
		ID:           "system-rule",
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
		IsSystemDefault: true,
	}
	svc.ruleRepo.Create(systemRule)

	// Test: with user_id that has no user rule, should fall back to system
	result, err := svc.ResolveRoute(context.Background(), "ollama:llama2", []string{}, "user-456")
	if err != nil {
		t.Fatalf("ResolveRoute failed: %v", err)
	}

	if result.ProviderID != "ollama" {
		t.Errorf("Expected system rule provider 'ollama', got '%s'", result.ProviderID)
	}
}

func TestResolveRoute_NoUserID(t *testing.T) {
	svc, _ := setupTestService(t)

	// Create system rule
	systemRule := &entity.RoutingRule{
		ID:           "system-rule",
		ModelPattern: "ollama:*",
		ProviderID:   "ollama",
		IsSystemDefault: true,
	}
	svc.ruleRepo.Create(systemRule)

	// Test: without user_id, should use system rule
	result, err := svc.ResolveRoute(context.Background(), "ollama:llama2", []string{}, "")
	if err != nil {
		t.Fatalf("ResolveRoute failed: %v", err)
	}

	if result.ProviderID != "ollama" {
		t.Errorf("Expected system rule provider 'ollama', got '%s'", result.ProviderID)
	}
}

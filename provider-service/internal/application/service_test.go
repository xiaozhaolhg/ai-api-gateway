package application

import (
	"testing"
	"time"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"github.com/ai-api-gateway/provider-service/internal/infrastructure/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&entity.Provider{}); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	return db
}

func TestFindProvidersByModel(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewProviderRepository(db)
	svc := &Service{providerRepo: repo}

	// Create providers with models
	providers := []*entity.Provider{
		{ID: "p1", Name: "Ollama", Type: "ollama", Models: []string{"llama2", "mistral"}, Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "p2", Name: "OpenAI", Type: "openai", Models: []string{"gpt-4"}, Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "p3", Name: "Anthropic", Type: "anthropic", Models: []string{"claude-3", "llama2"}, Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "p4", Name: "NoModels", Type: "custom", Models: []string{}, Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, p := range providers {
		if err := repo.Create(p); err != nil {
			t.Fatalf("Failed to create provider: %v", err)
		}
	}

	// Test 1: Model exists in one provider
	result, err := svc.FindProvidersByModel("gpt-4")
	if err != nil {
		t.Errorf("FindProvidersByModel(gpt-4) error = %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 provider, got %d", len(result))
	}
	if len(result) > 0 && result[0].ID != "p2" {
		t.Errorf("Expected p2, got %s", result[0].ID)
	}

	// Test 2: Model exists in multiple providers
	result, err = svc.FindProvidersByModel("llama2")
	if err != nil {
		t.Errorf("FindProvidersByModel(llama2) error = %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(result))
	}

	// Test 3: Model doesn't exist
	result, err = svc.FindProvidersByModel("non-existent")
	if err != nil {
		t.Errorf("FindProvidersByModel(non-existent) error = %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 providers, got %d", len(result))
	}

	// Test 4: Provider with empty models should not be returned
	result, err = svc.FindProvidersByModel("any-model")
	if err != nil {
		t.Errorf("FindProvidersByModel(any-model) error = %v", err)
	}
	for _, p := range result {
		if p.ID == "p4" {
			t.Error("Provider with empty models should not be returned")
		}
	}
}

// Placeholder callback dispatch tests for provider-service
// Full implementation would test async fire-and-forget, failure isolation

func TestCallbackDispatch_Placeholder(t *testing.T) {
	// TODO: Implement full callback dispatch tests
	// This placeholder ensures the test file exists and compiles
	t.Skip("Callback dispatch tests not yet fully implemented")
}

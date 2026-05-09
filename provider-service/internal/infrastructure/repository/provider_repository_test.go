package repository

import (
	"testing"
	"time"

	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForProvider(t *testing.T) *gorm.DB {
	dbPath := ":memory:"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&entity.Provider{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestProviderRepository_Create(t *testing.T) {
	db := setupTestDBForProvider(t)
	repo := NewProviderRepository(db)

	provider := &entity.Provider{
		ID:          "provider-1",
		Name:        "OpenAI",
		Type:        "openai",
		BaseURL:     "https://api.openai.com/v1",
		Credentials: "encrypted-key",
		Models:      []string{"gpt-4", "gpt-3.5-turbo"},
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(provider)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// Verify provider was created
	retrieved, err := repo.GetByID(provider.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if retrieved.Name != provider.Name {
		t.Errorf("Expected name %s, got %s", provider.Name, retrieved.Name)
	}
}

func TestProviderRepository_GetByID(t *testing.T) {
	db := setupTestDBForProvider(t)
	repo := NewProviderRepository(db)

	// Test non-existent provider
	_, err := repo.GetByID("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent provider, got nil")
	}

	// Create a provider and retrieve it
	provider := &entity.Provider{
		ID:          "provider-2",
		Name:        "Ollama",
		Type:        "ollama",
		BaseURL:     "http://localhost:11434",
		Credentials: "",
		Models:      []string{"llama2", "mistral"},
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Create(provider)

	retrieved, err := repo.GetByID(provider.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if retrieved.ID != provider.ID {
		t.Errorf("Expected ID %s, got %s", provider.ID, retrieved.ID)
	}
}

func TestProviderRepository_GetByType(t *testing.T) {
	db := setupTestDBForProvider(t)
	repo := NewProviderRepository(db)

	// Test non-existent type
	_, err := repo.GetByType("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent type, got nil")
	}

	// Create a provider and retrieve by type
	provider := &entity.Provider{
		ID:          "provider-3",
		Name:        "Anthropic",
		Type:        "anthropic",
		BaseURL:     "https://api.anthropic.com",
		Credentials: "encrypted-key",
		Models:      []string{"claude-3"},
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Create(provider)

	retrieved, err := repo.GetByType(provider.Type)
	if err != nil {
		t.Errorf("GetByType() error = %v", err)
	}
	if retrieved.Type != provider.Type {
		t.Errorf("Expected type %s, got %s", provider.Type, retrieved.Type)
	}
}

func TestProviderRepository_Update(t *testing.T) {
	db := setupTestDBForProvider(t)
	repo := NewProviderRepository(db)

	provider := &entity.Provider{
		ID:          "provider-4",
		Name:        "Test Provider",
		Type:        "custom",
		BaseURL:     "https://api.example.com",
		Credentials: "encrypted-key",
		Models:      []string{"model1"},
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Create(provider)

	// Update provider
	provider.Name = "Updated Provider"
	provider.Status = "inactive"
	provider.UpdatedAt = time.Now()
	err := repo.Update(provider)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	// Verify update
	retrieved, err := repo.GetByID(provider.ID)
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if retrieved.Name != "Updated Provider" {
		t.Errorf("Expected name 'Updated Provider', got %s", retrieved.Name)
	}
	if retrieved.Status != "inactive" {
		t.Errorf("Expected status 'inactive', got %s", retrieved.Status)
	}
}

func TestProviderRepository_Delete(t *testing.T) {
	db := setupTestDBForProvider(t)
	repo := NewProviderRepository(db)

	provider := &entity.Provider{
		ID:          "provider-5",
		Name:        "Test Provider 5",
		Type:        "custom",
		BaseURL:     "https://api.example.com",
		Credentials: "encrypted-key",
		Models:      []string{"model1"},
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Create(provider)

	// Delete provider
	err := repo.Delete(provider.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(provider.ID)
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}
}

func TestProviderRepository_List(t *testing.T) {
	db := setupTestDBForProvider(t)
	repo := NewProviderRepository(db)

	// Create multiple providers
	for i := 0; i < 5; i++ {
		provider := &entity.Provider{
			ID:          "provider-list-" + string(rune('a'+i)),
			Name:        "Test Provider " + string(rune('a'+i)),
			Type:        "custom",
			BaseURL:     "https://api.example.com",
			Credentials: "encrypted-key",
			Models:      []string{"model" + string(rune('a'+i))},
			Status:      "active",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		repo.Create(provider)
	}

	// Test pagination
	providers, total, err := repo.List(1, 3)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}
	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
	if len(providers) != 3 {
		t.Errorf("Expected 3 providers, got %d", len(providers))
	}

	// Test second page
	providers2, total2, err := repo.List(2, 3)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}
	if total2 != 5 {
		t.Errorf("Expected total 5, got %d", total2)
	}
	if len(providers2) != 2 {
		t.Errorf("Expected 2 providers on second page, got %d", len(providers2))
	}
}

func TestProviderRepository_StatusFiltering(t *testing.T) {
	db := setupTestDBForProvider(t)
	repo := NewProviderRepository(db)

	// Create active and inactive providers
	activeProvider := &entity.Provider{
		ID:          "active-1",
		Name:        "Active Provider",
		Type:        "custom",
		BaseURL:     "https://api.example.com",
		Credentials: "encrypted-key",
		Models:      []string{"model1"},
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Create(activeProvider)

	inactiveProvider := &entity.Provider{
		ID:          "inactive-1",
		Name:        "Inactive Provider",
		Type:        "custom",
		BaseURL:     "https://api.example.com",
		Credentials: "encrypted-key",
		Models:      []string{"model2"},
		Status:      "inactive",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Create(inactiveProvider)

	// List all providers
	providers, total, err := repo.List(1, 10)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}
	if total != 2 {
		t.Errorf("Expected total 2, got %d", total)
	}

	// Verify we have both active and inactive
	activeCount := 0
	inactiveCount := 0
	for _, p := range providers {
		if p.Status == "active" {
			activeCount++
		} else if p.Status == "inactive" {
			inactiveCount++
		}
	}

	if activeCount != 1 {
		t.Errorf("Expected 1 active provider, got %d", activeCount)
	}
	if inactiveCount != 1 {
		t.Errorf("Expected 1 inactive provider, got %d", inactiveCount)
	}
}

func TestProviderRepository_FindByModel(t *testing.T) {
	db := setupTestDBForProvider(t)
	repo := NewProviderRepository(db)

	// Create providers with various models
	provider1 := &entity.Provider{
		ID:          "provider-ollama",
		Name:        "Ollama",
		Type:        "ollama",
		BaseURL:     "http://localhost:11434",
		Credentials: "",
		Models:      []string{"llama2", "mistral"},
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Create(provider1)

	provider2 := &entity.Provider{
		ID:          "provider-openai",
		Name:        "OpenAI",
		Type:        "openai",
		BaseURL:     "https://api.openai.com",
		Credentials: "encrypted-key",
		Models:      []string{"gpt-4", "gpt-3.5-turbo"},
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Create(provider2)

	provider3 := &entity.Provider{
		ID:          "provider-anthropic",
		Name:        "Anthropic",
		Type:        "anthropic",
		BaseURL:     "https://api.anthropic.com",
		Credentials: "encrypted-key",
		Models:      []string{"claude-3", "llama2"}, // Also supports llama2
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Create(provider3)

	provider4 := &entity.Provider{
		ID:          "provider-no-models",
		Name:        "No Models Provider",
		Type:        "custom",
		BaseURL:     "https://api.example.com",
		Credentials: "",
		Models:      []string{}, // Empty models
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	repo.Create(provider4)

	// Test 1: Model exists in one provider
	providers, err := repo.FindByModel("gpt-4")
	if err != nil {
		t.Errorf("FindByModel(gpt-4) error = %v", err)
	}
	if len(providers) != 1 {
		t.Errorf("Expected 1 provider for gpt-4, got %d", len(providers))
	}
	if len(providers) > 0 && providers[0].ID != "provider-openai" {
		t.Errorf("Expected provider-openai, got %s", providers[0].ID)
	}

	// Test 2: Model exists in multiple providers
	providers, err = repo.FindByModel("llama2")
	if err != nil {
		t.Errorf("FindByModel(llama2) error = %v", err)
	}
	if len(providers) != 2 {
		t.Errorf("Expected 2 providers for llama2, got %d", len(providers))
	}

	// Test 3: Model doesn't exist
	providers, err = repo.FindByModel("non-existent-model")
	if err != nil {
		t.Errorf("FindByModel(non-existent-model) error = %v", err)
	}
	if len(providers) != 0 {
		t.Errorf("Expected 0 providers for non-existent-model, got %d", len(providers))
	}

	// Test 4: Provider with empty Models field should not be returned
	providers, err = repo.FindByModel("any-model")
	if err != nil {
		t.Errorf("FindByModel(any-model) error = %v", err)
	}
	for _, p := range providers {
		if p.ID == "provider-no-models" {
			t.Error("Provider with empty Models should not be returned")
		}
	}
}

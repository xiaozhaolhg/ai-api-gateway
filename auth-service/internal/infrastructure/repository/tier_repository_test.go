package repository

import (
	"testing"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTierTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&entity.Tier{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestTierRepository_Create(t *testing.T) {
	db := setupTierTestDB(t)
	repo := NewTierRepository(db)

	tier := &entity.Tier{
		ID:          "test-tier-1",
		Name:        "Test Tier",
		Description: "A test tier",
		IsDefault:   false,
		AllowedModels: []string{
			"ollama:llama2",
			"openai:gpt-4",
		},
		AllowedProviders: []string{
			"ollama",
			"openai",
		},
	}

	err := repo.Create(tier)
	require.NoError(t, err)

	// Verify tier was created
	retrieved, err := repo.GetByID(tier.ID)
	require.NoError(t, err)
	require.Equal(t, tier.Name, retrieved.Name)
	require.Equal(t, tier.Description, retrieved.Description)
	require.Equal(t, tier.IsDefault, retrieved.IsDefault)
	require.Equal(t, tier.AllowedModels, retrieved.AllowedModels)
	require.Equal(t, tier.AllowedProviders, retrieved.AllowedProviders)
}

func TestTierRepository_GetByID(t *testing.T) {
	db := setupTierTestDB(t)
	repo := NewTierRepository(db)

	// Test non-existent tier
	_, err := repo.GetByID("non-existent")
	require.Error(t, err)

	// Create a tier and retrieve it
	tier := &entity.Tier{
		ID:          "test-tier-2",
		Name:        "Test Tier 2",
		Description: "Another test tier",
		IsDefault:   true,
		AllowedModels: []string{
			"anthropic:claude-3",
		},
		AllowedProviders: []string{
			"anthropic",
		},
	}
	repo.Create(tier)

	retrieved, err := repo.GetByID(tier.ID)
	require.NoError(t, err)
	require.Equal(t, tier.ID, retrieved.ID)
	require.Equal(t, tier.Name, retrieved.Name)
	require.Equal(t, tier.Description, retrieved.Description)
	require.Equal(t, tier.IsDefault, retrieved.IsDefault)
	require.Equal(t, tier.AllowedModels, retrieved.AllowedModels)
	require.Equal(t, tier.AllowedProviders, retrieved.AllowedProviders)
}

func TestTierRepository_Update(t *testing.T) {
	db := setupTierTestDB(t)
	repo := NewTierRepository(db)

	tier := &entity.Tier{
		ID:          "test-tier-3",
		Name:        "Original Name",
		Description: "Original Description",
		IsDefault:   false,
		AllowedModels: []string{
			"ollama:llama2",
		},
		AllowedProviders: []string{
			"ollama",
		},
	}
	repo.Create(tier)

	// Update tier
	tier.Name = "Updated Name"
	tier.Description = "Updated Description"
	tier.IsDefault = true
	tier.AllowedModels = []string{
		"openai:gpt-4",
		"anthropic:claude-3",
	}
	tier.AllowedProviders = []string{
		"openai",
		"anthropic",
	}

	err := repo.Update(tier)
	require.NoError(t, err)

	// Verify update
	retrieved, err := repo.GetByID(tier.ID)
	require.NoError(t, err)
	require.Equal(t, "Updated Name", retrieved.Name)
	require.Equal(t, "Updated Description", retrieved.Description)
	require.Equal(t, true, retrieved.IsDefault)
	require.Equal(t, []string{"openai:gpt-4", "anthropic:claude-3"}, retrieved.AllowedModels)
	require.Equal(t, []string{"openai", "anthropic"}, retrieved.AllowedProviders)
}

func TestTierRepository_Delete(t *testing.T) {
	db := setupTierTestDB(t)
	repo := NewTierRepository(db)

	tier := &entity.Tier{
		ID:          "test-tier-4",
		Name:        "Test Tier to Delete",
		Description: "This tier will be deleted",
		IsDefault:   false,
		AllowedModels: []string{
			"ollama:llama2",
		},
		AllowedProviders: []string{
			"ollama",
		},
	}
	repo.Create(tier)

	// Delete tier
	err := repo.Delete(tier.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(tier.ID)
	require.Error(t, err)
}

func TestTierRepository_List(t *testing.T) {
	db := setupTierTestDB(t)
	repo := NewTierRepository(db)

	// Create multiple tiers
	for i := 0; i < 5; i++ {
		tier := &entity.Tier{
			ID:     "test-tier-list-" + string(rune('0'+i)),
			Name:   "Test Tier " + string(rune('0'+i)),
			IsDefault: false,
			AllowedModels: []string{
				"ollama:llama2",
			},
			AllowedProviders: []string{
				"ollama",
			},
		}
		repo.Create(tier)
	}

	// Test pagination
	tiers, total, err := repo.List(1, 3)
	require.NoError(t, err)
	require.Equal(t, 5, total)
	require.Len(t, tiers, 3)

	// Test second page
	tiers2, total2, err := repo.List(2, 3)
	require.NoError(t, err)
	require.Equal(t, 5, total2)
	require.Len(t, tiers2, 2)
}

func TestTierRepository_GetByName(t *testing.T) {
	db := setupTierTestDB(t)
	repo := NewTierRepository(db)

	// Test non-existent tier name
	_, err := repo.GetByName("non-existent")
	require.Error(t, err)

	// Create a tier and retrieve by name
	tier := &entity.Tier{
		ID:          "test-tier-name-1",
		Name:        "Unique Tier Name",
		Description: "A tier with unique name",
		IsDefault:   false,
		AllowedModels: []string{
			"ollama:llama2",
		},
		AllowedProviders: []string{
			"ollama",
		},
	}
	repo.Create(tier)

	retrieved, err := repo.GetByName(tier.Name)
	require.NoError(t, err)
	require.Equal(t, tier.ID, retrieved.ID)
	require.Equal(t, tier.Name, retrieved.Name)
	require.Equal(t, tier.Description, retrieved.Description)
	require.Equal(t, tier.IsDefault, retrieved.IsDefault)
	require.Equal(t, tier.AllowedModels, retrieved.AllowedModels)
	require.Equal(t, tier.AllowedProviders, retrieved.AllowedProviders)
}

func TestTierRepository_GetDefaultTiers(t *testing.T) {
	db := setupTierTestDB(t)
	repo := NewTierRepository(db)

	// Create some default tiers and some non-default tiers
	defaultTier1 := &entity.Tier{
		ID:          "default-tier-1",
		Name:        "Basic Tier",
		Description: "Basic access tier",
		IsDefault:   true,
		AllowedModels: []string{
			"ollama:llama2",
		},
		AllowedProviders: []string{
			"ollama",
		},
	}
	defaultTier2 := &entity.Tier{
		ID:          "default-tier-2",
		Name:        "Premium Tier",
		Description: "Premium access tier",
		IsDefault:   true,
		AllowedModels: []string{
			"openai:gpt-4",
			"anthropic:claude-3",
		},
		AllowedProviders: []string{
			"openai",
			"anthropic",
		},
	}
	nonDefaultTier := &entity.Tier{
		ID:          "custom-tier-1",
		Name:        "Custom Tier",
		Description: "Custom access tier",
		IsDefault:   false,
		AllowedModels: []string{
			"gemini:pro",
		},
		AllowedProviders: []string{
			"gemini",
		},
	}

	repo.Create(defaultTier1)
	repo.Create(defaultTier2)
	repo.Create(nonDefaultTier)

	// Get default tiers
	tiers, err := repo.GetDefaultTiers()
	require.NoError(t, err)
	require.Len(t, tiers, 2)

	// Verify we got the default tiers
	tierIDs := map[string]bool{}
	for _, tier := range tiers {
		tierIDs[tier.ID] = true
	}
	require.True(t, tierIDs["default-tier-1"])
	require.True(t, tierIDs["default-tier-2"])
	require.False(t, tierIDs["custom-tier-1"])
}
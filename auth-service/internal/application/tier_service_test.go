package application

import (
	"testing"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTierService_CreateTier(t *testing.T) {
	repo := newMockTierRepoForService()
	groupRepo := newMockGroupRepoForService()
	service := NewTierService(repo, groupRepo)

	tier, err := service.CreateTier("Test Tier", "Test description", false, []string{"ollama:llama2"}, []string{"ollama"})
	require.NoError(t, err)
	require.Equal(t, "Test Tier", tier.Name)
	require.Equal(t, "Test description", tier.Description)
	require.False(t, tier.IsDefault)
}

func TestTierService_GetTier(t *testing.T) {
	repo := newMockTierRepoForService()
	groupRepo := newMockGroupRepoForService()
	service := NewTierService(repo, groupRepo)

	created, _ := service.CreateTier("Test Tier", "Test description", false, []string{"ollama:llama2"}, []string{"ollama"})

	retrieved, err := service.GetTier(created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, retrieved.ID)
}

func TestTierService_UpdateTier(t *testing.T) {
	repo := newMockTierRepoForService()
	groupRepo := newMockGroupRepoForService()
	service := NewTierService(repo, groupRepo)

	tier, _ := service.CreateTier("Original", "Original desc", false, []string{"ollama:llama2"}, []string{"ollama"})

	updated, err := service.UpdateTier(tier.ID, "Updated", "Updated desc", []string{"openai:gpt-4"}, []string{"openai"})
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.Name)
	require.Equal(t, "Updated desc", updated.Description)
}

func TestTierService_UpdateTier_DefaultTier(t *testing.T) {
	repo := newMockTierRepoForService()
	groupRepo := newMockGroupRepoForService()
	service := NewTierService(repo, groupRepo)

	tier, _ := service.CreateTier("Default", "Default desc", true, []string{"ollama:*"}, []string{"ollama"})

	_, err := service.UpdateTier(tier.ID, "Try Update", "", nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "predefined tiers cannot be updated")
}

func TestTierService_DeleteTier(t *testing.T) {
	repo := newMockTierRepoForService()
	groupRepo := newMockGroupRepoForService()
	service := NewTierService(repo, groupRepo)

	tier, _ := service.CreateTier("To Delete", "Delete me", false, []string{"ollama:llama2"}, []string{"ollama"})

	err := service.DeleteTier(tier.ID)
	require.NoError(t, err)

	_, err = service.GetTier(tier.ID)
	require.Error(t, err)
}

func TestTierService_DeleteTier_DefaultTier(t *testing.T) {
	repo := newMockTierRepoForService()
	groupRepo := newMockGroupRepoForService()
	service := NewTierService(repo, groupRepo)

	tier, _ := service.CreateTier("Default Delete", "Default delete me", true, []string{"ollama:*"}, []string{"ollama"})

	err := service.DeleteTier(tier.ID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "predefined tiers cannot be deleted")
}

func TestTierService_DeleteTier_WithGroupReference(t *testing.T) {
	repo := newMockTierRepoForService()
	groupRepo := newMockGroupRepoForService()
	service := NewTierService(repo, groupRepo)

	tier, _ := service.CreateTier("Referenced Tier", "Tier with group reference", false, []string{"ollama:*"}, []string{"ollama"})
	groupRepo.Create(&entity.Group{ID: "group-ref-1", Name: "Reference Group"})

	// Assign tier to group
	err := service.AssignTierToGroup("group-ref-1", tier.ID)
	require.NoError(t, err)

	// Verify tier is assigned
	group, _ := groupRepo.GetByID("group-ref-1")
	require.Equal(t, tier.ID, group.TierID)

	// Attempt to delete tier should fail
	err = service.DeleteTier(tier.ID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot delete tier referenced by groups")

	// Verify tier still exists
	_, err = service.GetTier(tier.ID)
	require.NoError(t, err)
}

func TestTierService_SeedDefaultTiers(t *testing.T) {
	repo := newMockTierRepoForService()
	groupRepo := newMockGroupRepoForService()
	service := NewTierService(repo, groupRepo)

	err := service.SeedDefaultTiers()
	require.NoError(t, err)

	tiers, _, err := service.ListTiers(1, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(tiers), 4)
}

func TestTierService_SeedDefaultTiers_Idempotent(t *testing.T) {
	repo := newMockTierRepoForService()
	groupRepo := newMockGroupRepoForService()
	service := NewTierService(repo, groupRepo)

	err := service.SeedDefaultTiers()
	require.NoError(t, err)

	err = service.SeedDefaultTiers()
	require.NoError(t, err)

	tiers, _, err := service.ListTiers(1, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(tiers), 4)
}

func TestTierService_AssignTierToGroup(t *testing.T) {
	repo := newMockTierRepoForService()
	groupRepo := newMockGroupRepoForService()
	service := NewTierService(repo, groupRepo)

	tier, _ := service.CreateTier("Test Tier", "Test", false, []string{"ollama:*"}, []string{"ollama"})
	groupRepo.Create(&entity.Group{ID: "group-1", Name: "Test Group"})

	err := service.AssignTierToGroup("group-1", tier.ID)
	require.NoError(t, err)

	group, _ := groupRepo.GetByID("group-1")
	require.Equal(t, tier.ID, group.TierID)
}

func TestTierService_RemoveTierFromGroup(t *testing.T) {
	repo := newMockTierRepoForService()
	groupRepo := newMockGroupRepoForService()
	service := NewTierService(repo, groupRepo)

	tier, _ := service.CreateTier("Test Tier", "Test", false, []string{"ollama:*"}, []string{"ollama"})
	group := &entity.Group{ID: "group-1", Name: "Test Group", TierID: tier.ID}
	groupRepo.Create(group)

	err := service.RemoveTierFromGroup("group-1")
	require.NoError(t, err)

	group, _ = groupRepo.GetByID("group-1")
	require.Empty(t, group.TierID)
}

type mockTierRepoForService struct {
	tiers map[string]*entity.Tier
}

func newMockTierRepoForService() *mockTierRepoForService {
	return &mockTierRepoForService{tiers: make(map[string]*entity.Tier)}
}

func (m *mockTierRepoForService) Create(tier *entity.Tier) error {
	m.tiers[tier.ID] = tier
	return nil
}

func (m *mockTierRepoForService) GetByID(id string) (*entity.Tier, error) {
	if tier, ok := m.tiers[id]; ok {
		return tier, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockTierRepoForService) Update(tier *entity.Tier) error {
	m.tiers[tier.ID] = tier
	return nil
}

func (m *mockTierRepoForService) Delete(id string) error {
	delete(m.tiers, id)
	return nil
}

func (m *mockTierRepoForService) List(page, pageSize int) ([]*entity.Tier, int, error) {
	tiers := make([]*entity.Tier, 0, len(m.tiers))
	for _, tier := range m.tiers {
		tiers = append(tiers, tier)
	}
	return tiers, len(tiers), nil
}

func (m *mockTierRepoForService) GetByName(name string) (*entity.Tier, error) {
	for _, tier := range m.tiers {
		if tier.Name == name {
			return tier, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockTierRepoForService) GetDefaultTiers() ([]*entity.Tier, error) {
	var tiers []*entity.Tier
	for _, tier := range m.tiers {
		if tier.IsDefault {
			tiers = append(tiers, tier)
		}
	}
	return tiers, nil
}

type mockGroupRepoForService struct {
	groups map[string]*entity.Group
}

func newMockGroupRepoForService() *mockGroupRepoForService {
	return &mockGroupRepoForService{groups: make(map[string]*entity.Group)}
}

func (m *mockGroupRepoForService) Create(group *entity.Group) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepoForService) GetByID(id string) (*entity.Group, error) {
	if group, ok := m.groups[id]; ok {
		return group, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockGroupRepoForService) Update(group *entity.Group) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockGroupRepoForService) Delete(id string) error {
	delete(m.groups, id)
	return nil
}

func (m *mockGroupRepoForService) List(page, pageSize int) ([]*entity.Group, int, error) {
	groups := make([]*entity.Group, 0, len(m.groups))
	for _, group := range m.groups {
		groups = append(groups, group)
	}
	return groups, len(groups), nil
}
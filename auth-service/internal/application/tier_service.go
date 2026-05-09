package application

import (
	"fmt"
	"log"
	"time"

	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
	"github.com/ai-api-gateway/auth-service/internal/domain/port"
)

type TierService struct {
	tierRepo  port.TierRepository
	groupRepo port.GroupRepository
}

func NewTierService(tierRepo port.TierRepository, groupRepo port.GroupRepository) *TierService {
	return &TierService{tierRepo: tierRepo, groupRepo: groupRepo}
}

func (s *TierService) CreateTier(name, description string, isDefault bool, allowedModels, allowedProviders []string) (*entity.Tier, error) {
	tier := &entity.Tier{
		ID:               generateID(),
		Name:             name,
		Description:      description,
		IsDefault:        isDefault,
		AllowedModels:    allowedModels,
		AllowedProviders: allowedProviders,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.tierRepo.Create(tier); err != nil {
		return nil, fmt.Errorf("failed to create tier: %w", err)
	}

	return tier, nil
}

func (s *TierService) GetTier(id string) (*entity.Tier, error) {
	tier, err := s.tierRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("tier not found: %w", err)
	}
	return tier, nil
}

func (s *TierService) UpdateTier(id, name, description string, allowedModels, allowedProviders []string) (*entity.Tier, error) {
	tier, err := s.tierRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("tier not found: %w", err)
	}

	if tier.IsDefault {
		return nil, fmt.Errorf("predefined tiers cannot be updated")
	}

	if name != "" {
		tier.Name = name
	}
	if description != "" {
		tier.Description = description
	}
	tier.AllowedModels = allowedModels
	tier.AllowedProviders = allowedProviders
	tier.UpdatedAt = time.Now()

	if err := s.tierRepo.Update(tier); err != nil {
		return nil, fmt.Errorf("failed to update tier: %w", err)
	}

	return tier, nil
}

func (s *TierService) DeleteTier(id string) error {
	tier, err := s.tierRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("tier not found: %w", err)
	}

	if tier.IsDefault {
		return fmt.Errorf("predefined tiers cannot be deleted")
	}

	groups, _, err := s.groupRepo.List(1, 1000)
	if err != nil {
		return fmt.Errorf("failed to check group references: %w", err)
	}
	for _, g := range groups {
		if g.TierID == id {
			return fmt.Errorf("cannot delete tier referenced by groups")
		}
	}

	if err := s.tierRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete tier: %w", err)
	}
	return nil
}

func (s *TierService) ListTiers(page, pageSize int) ([]*entity.Tier, int, error) {
	tiers, total, err := s.tierRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tiers: %w", err)
	}
	return tiers, total, nil
}

func (s *TierService) AssignTierToGroup(groupID, tierID string) error {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return fmt.Errorf("group not found: %w", err)
	}

	tier, err := s.tierRepo.GetByID(tierID)
	if err != nil {
		return fmt.Errorf("tier not found: %w", err)
	}

	log.Printf("[DEBUG] AssignTierToGroup: groupID=%s, tierID=%s, tierName=%s", groupID, tierID, tier.Name)
	log.Printf("[DEBUG] Group before: TierID=%s", group.TierID)

	group.TierID = tierID
	group.UpdatedAt = time.Now()

	if err := s.groupRepo.Update(group); err != nil {
		return fmt.Errorf("failed to assign tier to group: %w", err)
	}

	log.Printf("[DEBUG] Group after: TierID=%s", group.TierID)
	return nil
}

func (s *TierService) RemoveTierFromGroup(groupID string) error {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return fmt.Errorf("group not found: %w", err)
	}

	group.TierID = ""
	group.UpdatedAt = time.Now()
	if err := s.groupRepo.Update(group); err != nil {
		return fmt.Errorf("failed to remove tier from group: %w", err)
	}
	return nil
}

func (s *TierService) SeedDefaultTiers() error {
	defaultTiers := []struct {
		name             string
		description      string
		allowedModels    []string
		allowedProviders []string
	}{
		{
			name:             "Basic",
			description:     "Basic access tier with limited models",
			allowedModels:    []string{"ollama:llama2", "ollama:mistral"},
			allowedProviders: []string{"ollama"},
		},
		{
			name:             "Standard",
			description:     "Standard access tier with common models",
			allowedModels:    []string{"ollama:*", "openai:gpt-4", "openai:gpt-4-turbo"},
			allowedProviders: []string{"ollama", "openai"},
		},
		{
			name:             "Premium",
			description:     "Premium access tier with advanced models",
			allowedModels:    []string{"openai:gpt-4*", "anthropic:claude-3", "anthropic:claude-3-sonnet"},
			allowedProviders: []string{"openai", "anthropic"},
		},
		{
			name:             "Enterprise",
			description:     "Enterprise access tier with all models",
			allowedModels:    []string{"*"},
			allowedProviders: []string{"*"},
		},
	}

	for _, dt := range defaultTiers {
		existing, _ := s.tierRepo.GetByName(dt.name)
		if existing != nil {
			continue
		}

		tier := &entity.Tier{
			ID:               generateID(),
			Name:             dt.name,
			Description:      dt.description,
			IsDefault:        true,
			AllowedModels:    dt.allowedModels,
			AllowedProviders: dt.allowedProviders,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := s.tierRepo.Create(tier); err != nil {
			return fmt.Errorf("failed to seed default tier %s: %w", dt.name, err)
		}
	}

	return nil
}
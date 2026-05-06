package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	"github.com/ai-api-gateway/router-service/internal/domain/port"
)

type Service struct {
	ruleRepo port.RoutingRuleRepository
	cache    port.Cache
	ttl      int
}

func NewService(ruleRepo port.RoutingRuleRepository, cache port.Cache) *Service {
	return &Service{
		ruleRepo: ruleRepo,
		cache:    cache,
		ttl:      300,
	}
}

func (s *Service) ResolveRoute(ctx context.Context, model string, authorizedModels []string, userID string) (*entity.RouteResult, error) {
	if s.cache != nil {
		cacheKey := fmt.Sprintf("router:route:%s:%s", model, userID)
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var result entity.RouteResult
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				if s.isAuthorized(result.ProviderID, authorizedModels) {
					return &result, nil
				}
			}
		}
	}

	var rule *entity.RoutingRule
	var err error

	if userID != "" {
		rule, err = s.ruleRepo.FindByModel(model, &userID)
	} else {
		rule, err = s.ruleRepo.FindByModel(model, nil)
	}
	if err != nil {
		for i, c := range model {
			if c == ':' {
				providerType := model[:i]
				return &entity.RouteResult{
					ProviderID:   providerType,
					AdapterType: providerType,
				}, nil
			}
		}
		return nil, fmt.Errorf("no routing rule found for model: %s", model)
	}

	adapterType := s.inferAdapterType(rule.ProviderID)
	fallbackProviderIDs, fallbackModels := s.getFallbackProviders(rule, nil)

	result := &entity.RouteResult{
		ProviderID:          rule.ProviderID,
		AdapterType:         adapterType,
		FallbackProviderIDs: fallbackProviderIDs,
		FallbackModels:      fallbackModels,
	}

	if s.cache != nil {
		cacheKey := fmt.Sprintf("router:route:%s:%s", model, userID)
		resultJSON, _ := json.Marshal(result)
		s.cache.Set(ctx, cacheKey, string(resultJSON), s.ttl)
	}

	return result, nil
}

func (s *Service) isAuthorized(providerID string, authorizedModels []string) bool {
	if len(authorizedModels) == 0 {
		return true
	}
	for _, model := range authorizedModels {
		if strings.Contains(providerID, model) || s.matchPattern(model, providerID) {
			return true
		}
	}
	return false
}

func (s *Service) matchPattern(pattern, model string) bool {
	if pattern == "*" {
		return true
	}

	if !strings.Contains(pattern, "*") {
		return pattern == model
	}

	regexPattern := strings.ReplaceAll(pattern, "*", ".*")
	regexPattern = "^" + regexPattern + "$"

	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(model, prefix)
	}

	if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(model, suffix)
	}

	parts := strings.Split(pattern, "*")
	if len(parts) == 2 {
		return strings.HasPrefix(model, parts[0]) && strings.HasSuffix(model, parts[1])
	}

	return false
}

func (s *Service) inferAdapterType(providerID string) string {
	if strings.Contains(providerID, "openai") {
		return "openai"
	}
	if strings.Contains(providerID, "anthropic") {
		return "anthropic"
	}
	if strings.Contains(providerID, "ollama") {
		return "ollama"
	}
	if strings.Contains(providerID, "opencode-zen") {
		return "opencode-zen"
	}
	return "unknown"
}

func (s *Service) getFallbackProviders(rule *entity.RoutingRule, rules []*entity.RoutingRule) ([]string, []string) {
	var providerIDs []string
	var models []string

	if rule.FallbackProviderIDs != "" {
		json.Unmarshal([]byte(rule.FallbackProviderIDs), &providerIDs)
		json.Unmarshal([]byte(rule.FallbackModels), &models)
	}

	return providerIDs, models
}

// ListRoutingRules returns a paginated list of routing rules.
func (s *Service) ListRoutingRules(page, pageSize int) ([]*entity.RoutingRule, int, error) {
	return s.ruleRepo.List(page, pageSize)
}

// CreateRoutingRule creates a new routing rule.
func (s *Service) CreateRoutingRule(rule *entity.RoutingRule) error {
	return s.ruleRepo.Create(rule)
}

// UpdateRoutingRule updates an existing routing rule.
func (s *Service) UpdateRoutingRule(rule *entity.RoutingRule, requestingUserID string) error {
	return s.ruleRepo.UpdateWithOwnership(rule, requestingUserID)
}

// DeleteRoutingRule deletes a routing rule by ID.
func (s *Service) DeleteRoutingRule(id string, requestingUserID string) error {
	return s.ruleRepo.DeleteWithOwnership(id, requestingUserID)
}

// RefreshRoutingTable invalidates the routing table cache.
func (s *Service) RefreshRoutingTable(ctx context.Context) error {
	if s.cache != nil {
		// Clear all route-related cache keys using the prefix
		return s.cache.ClearPrefix(ctx, "router:")
	}
	return nil
}

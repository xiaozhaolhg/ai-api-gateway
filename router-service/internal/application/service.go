package application

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	"github.com/ai-api-gateway/router-service/internal/domain/port"
)

// Service handles route resolution logic
type Service struct {
	ruleRepo port.RoutingRuleRepository
}

// NewService creates a new application service
func NewService(ruleRepo port.RoutingRuleRepository) *Service {
	return &Service{
		ruleRepo: ruleRepo,
	}
}

// ResolveRoute resolves a model name to a provider based on routing rules
func (s *Service) ResolveRoute(model string) (*entity.RouteResult, error) {
	// Get all routing rules
	rules, _, err := s.ruleRepo.List(0, 1000) // Get all rules
	if err != nil {
		return nil, fmt.Errorf("failed to get routing rules: %w", err)
	}

	// Find matching rules
	var matchingRules []*entity.RoutingRule
	for _, rule := range rules {
		if s.matchPattern(rule.ModelPattern, model) {
			matchingRules = append(matchingRules, rule)
		}
	}

	if len(matchingRules) == 0 {
		return nil, fmt.Errorf("no routing rule found for model: %s", model)
	}

	// Sort by priority (higher priority first)
	sort.Slice(matchingRules, func(i, j int) bool {
		return matchingRules[i].Priority > matchingRules[j].Priority
	})

	// Use the highest priority rule
	rule := matchingRules[0]

	// Determine adapter type based on provider ID
	// In production, this would call provider-service to get provider details
	// For MVP, we'll infer from provider ID or use a default
	adapterType := s.inferAdapterType(rule.ProviderID)

	return &entity.RouteResult{
		ProviderID:           rule.ProviderID,
		AdapterType:          adapterType,
		FallbackProviderIDs:  s.getFallbackProviders(rule, matchingRules),
	}, nil
}

// matchPattern checks if a model name matches a pattern
// Supports wildcard matching with *
func (s *Service) matchPattern(pattern, model string) bool {
	if pattern == "*" {
		return true
	}

	if !strings.Contains(pattern, "*") {
		return pattern == model
	}

	// Convert wildcard pattern to regex
	// Simple implementation: replace * with .*
	regexPattern := strings.ReplaceAll(pattern, "*", ".*")
	regexPattern = "^" + regexPattern + "$"

	// Simple glob matching
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(model, prefix)
	}

	if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(model, suffix)
	}

	// For patterns with * in the middle, do simple matching
	parts := strings.Split(pattern, "*")
	if len(parts) == 2 {
		return strings.HasPrefix(model, parts[0]) && strings.HasSuffix(model, parts[1])
	}

	return false
}

// inferAdapterType infers the adapter type from provider ID
// In production, this would call provider-service to get provider details
func (s *Service) inferAdapterType(providerID string) string {
	// Simple inference based on provider ID naming
	// In production, this would be a lookup in provider-service
	switch {
	case strings.Contains(providerID, "openai"):
		return "openai"
	case strings.Contains(providerID, "anthropic"):
		return "anthropic"
	case strings.Contains(providerID, "gemini"):
		return "gemini"
	case strings.Contains(providerID, "ollama"):
		return "ollama"
	case strings.Contains(providerID, "opencode") || strings.Contains(providerID, "zen"):
		return "opencode-zen"
	default:
		return "openai" // Default to OpenAI-compatible
	}
}

// getFallbackProviders returns fallback provider IDs from the routing rules
func (s *Service) getFallbackProviders(primaryRule *entity.RoutingRule, allRules []*entity.RoutingRule) []string {
	var fallbacks []string

	// Add explicit fallback if configured
	if primaryRule.FallbackProviderID != "" {
		fallbacks = append(fallbacks, primaryRule.FallbackProviderID)
	}

	// Add other matching rules as fallbacks (lower priority)
	for _, rule := range allRules {
		if rule.ID != primaryRule.ID && rule.Priority < primaryRule.Priority {
			fallbacks = append(fallbacks, rule.ProviderID)
		}
	}

	return fallbacks
}

// CreateRoutingRule creates a new routing rule
func (s *Service) CreateRoutingRule(rule *entity.RoutingRule) error {
	return s.ruleRepo.Create(rule)
}

// UpdateRoutingRule updates an existing routing rule
func (s *Service) UpdateRoutingRule(rule *entity.RoutingRule) error {
	return s.ruleRepo.Update(rule)
}

// DeleteRoutingRule deletes a routing rule
func (s *Service) DeleteRoutingRule(id string) error {
	return s.ruleRepo.Delete(id)
}

// ListRoutingRules lists all routing rules
func (s *Service) ListRoutingRules(page, pageSize int) ([]*entity.RoutingRule, int, error) {
	return s.ruleRepo.List(page, pageSize)
}

// RefreshRoutingTable refreshes the routing table cache
// For MVP, this is a no-op since we don't have Redis caching yet
func (s *Service) RefreshRoutingTable() error {
	// In production, this would invalidate Redis cache
	return nil
}

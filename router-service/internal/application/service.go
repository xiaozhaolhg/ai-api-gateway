package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	providerv1 "github.com/ai-api-gateway/api/gen/provider/v1"
	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	"github.com/ai-api-gateway/router-service/internal/domain/port"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	ruleRepo      port.RoutingRuleRepository
	cache         port.Cache
	ttl           int
	providerClient providerv1.ProviderServiceClient
}

func NewService(ruleRepo port.RoutingRuleRepository, cache port.Cache, providerAddr string) (*Service, error) {
	conn, err := grpc.NewClient(providerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to provider-service: %w", err)
	}
	client := providerv1.NewProviderServiceClient(conn)

	return &Service{
		ruleRepo:      ruleRepo,
		cache:         cache,
		ttl:           300,
		providerClient: client,
	}, nil
}

func NewServiceWithClient(ruleRepo port.RoutingRuleRepository, cache port.Cache, providerClient providerv1.ProviderServiceClient) *Service {
	return &Service{
		ruleRepo:      ruleRepo,
		cache:         cache,
		ttl:           300,
		providerClient: providerClient,
	}
}

func (s *Service) ResolveRoute(ctx context.Context, model string, authorizedModels []string, userID string) (*entity.RouteResult, error) {
	log.Printf("[ResolveRoute] Called with model=%s, userID=%s, authorizedModels=%v", model, userID, authorizedModels)
	isBareModel := true
	for _, c := range model {
		if c == ':' {
			isBareModel = false
			break
		}
	}
	log.Printf("[ResolveRoute] isBareModel=%v", isBareModel)

	if isBareModel {
		log.Printf("[ResolveRoute] Calling resolveBareModel")
		return s.resolveBareModel(ctx, model)
	}

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
	if err != nil || rule == nil {
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

func (s *Service) FindProvidersByModel(model string) ([]*entity.Provider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.providerClient.FindProvidersByModel(ctx, &providerv1.FindProvidersByModelRequest{
		Model: model,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call FindProvidersByModel: %w", err)
	}

	providers := make([]*entity.Provider, 0, len(resp.Providers))
	for _, p := range resp.Providers {
		providers = append(providers, &entity.Provider{
			ID:       p.Id,
			Name:     p.Name,
			Type:     p.Type,
			BaseURL:  p.BaseUrl,
			Models:   p.Models,
			Status:   p.Status,
			CreatedAt: time.Unix(p.CreatedAt, 0),
			UpdatedAt: time.Unix(p.UpdatedAt, 0),
		})
	}
	return providers, nil
}

func (s *Service) resolveBareModel(ctx context.Context, bareModel string) (*entity.RouteResult, error) {
	log.Printf("[resolveBareModel] Called with bareModel=%s", bareModel)
	providers, err := s.FindProvidersByModel(bareModel)
	if err != nil {
		log.Printf("[resolveBareModel] FindProvidersByModel error: %v", err)
		return nil, fmt.Errorf("failed to find providers for model %s: %w", bareModel, err)
	}
	log.Printf("[resolveBareModel] Found %d providers for model %s", len(providers), bareModel)
	for i, p := range providers {
		log.Printf("[resolveBareModel] Provider %d: ID=%s, Type=%s, Models=%v", i, p.ID, p.Type, p.Models)
	}
	if len(providers) == 0 {
		return nil, fmt.Errorf("no provider found for model: %s", bareModel)
	}

	type result struct {
		provider *entity.Provider
		healthy   bool
	}

	resultsChan := make(chan result, len(providers))
	for _, p := range providers {
		go func(prov *entity.Provider) {
			log.Printf("[resolveBareModel] Checking health for provider ID=%s", prov.ID)
			healthy, err := s.CheckHealth(ctx, prov.ID)
			log.Printf("[resolveBareModel] Provider %s health: healthy=%v, err=%v", prov.ID, healthy, err)
			resultsChan <- result{prov, healthy}
		}(p)
	}

	var healthyProviders []*entity.Provider
	for i := 0; i < len(providers); i++ {
		r := <-resultsChan
		if r.healthy {
			healthyProviders = append(healthyProviders, r.provider)
		}
	}

	log.Printf("[resolveBareModel] Total healthy providers: %d", len(healthyProviders))
	if len(healthyProviders) == 0 {
		return nil, fmt.Errorf("no healthy provider found for model: %s", bareModel)
	}

	primary := healthyProviders[0]
	var fallbackIDs []string
	var fallbackModels []string
	for i := 1; i < len(healthyProviders); i++ {
		fallbackIDs = append(fallbackIDs, healthyProviders[i].ID)
		fallbackModels = append(fallbackModels, bareModel)
	}

	adapterType := s.inferAdapterType(primary.Type)
	log.Printf("[resolveBareModel] primary.ID=%s, primary.Type=%s, bareModel=%s", primary.ID, primary.Type, bareModel)
	return &entity.RouteResult{
		ProviderID:          primary.ID,
		AdapterType:         adapterType,
		FallbackProviderIDs: fallbackIDs,
		FallbackModels:      fallbackModels,
	}, nil
}

func (s *Service) CheckHealth(ctx context.Context, providerID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := s.providerClient.HealthCheck(ctx, &providerv1.HealthCheckRequest{
		ProviderId: providerID,
	})
	if err != nil {
		return false, nil
	}
	return resp.Healthy, nil
}

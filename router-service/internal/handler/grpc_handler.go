package handler

import (
	"context"
	"encoding/json"

	routerv1 "github.com/ai-api-gateway/api/gen/router/v1"
	v1 "github.com/ai-api-gateway/api/gen/common/v1"
	"github.com/ai-api-gateway/router-service/internal/application"
	"github.com/ai-api-gateway/router-service/internal/domain/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

// Handler implements the RouterService gRPC interface.
type Handler struct {
	routerv1.UnimplementedRouterServiceServer
	service *application.Service
}

// NewHandler creates a new gRPC handler.
func NewHandler(service *application.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ResolveRoute resolves a model name to a provider.
func (h *Handler) ResolveRoute(ctx context.Context, req *routerv1.ResolveRouteRequest) (*routerv1.RouteResult, error) {
	result, err := h.service.ResolveRoute(ctx, req.Model, req.AuthorizedModels, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "no route found for model %s: %v", req.Model, err)
	}

	return &routerv1.RouteResult{
		ProviderId:         result.ProviderID,
		AdapterType:        result.AdapterType,
		FallbackProviderIds: result.FallbackProviderIDs,
		FallbackModels:      result.FallbackModels,
	}, nil
}

// GetRoutingRules lists all routing rules.
func (h *Handler) GetRoutingRules(ctx context.Context, req *routerv1.GetRoutingRulesRequest) (*routerv1.ListRoutingRulesResponse, error) {
	rules, total, err := h.service.ListRoutingRules(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get routing rules: %v", err)
	}

	protoRules := make([]*routerv1.RoutingRule, len(rules))
	for i, rule := range rules {
		ids, models := unmarshalFallback(rule.FallbackProviderIDs, rule.FallbackModels)
		protoRules[i] = &routerv1.RoutingRule{
			Id:                 rule.ID,
			UserId:              rule.UserID,
			ModelPattern:        rule.ModelPattern,
			ProviderId:          rule.ProviderID,
			Priority:            rule.Priority,
			FallbackProviderIds: ids,
			FallbackModels:      models,
			IsSystemDefault:     rule.IsSystemDefault,
		}
	}

	return &routerv1.ListRoutingRulesResponse{
		Rules: protoRules,
		Total: int32(total),
	}, nil
}

// CreateRoutingRule creates a new routing rule.
func (h *Handler) CreateRoutingRule(ctx context.Context, req *routerv1.CreateRoutingRuleRequest) (*routerv1.RoutingRule, error) {
	idsJSON, modelsJSON := marshalFallback(req.FallbackProviderIds, req.FallbackModels)

	rule := &entity.RoutingRule{
		ID:                  uuid.New().String(),
		UserID:              req.UserId,
		ModelPattern:        req.ModelPattern,
		ProviderID:          req.ProviderId,
		Priority:            req.Priority,
		FallbackProviderIDs: idsJSON,
		FallbackModels:      modelsJSON,
		IsSystemDefault:     req.IsSystemDefault,
	}

	if err := h.service.CreateRoutingRule(rule); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create routing rule: %v", err)
	}

	ids, models := unmarshalFallback(rule.FallbackProviderIDs, rule.FallbackModels)
	return &routerv1.RoutingRule{
		Id:                 rule.ID,
		UserId:              rule.UserID,
		ModelPattern:        rule.ModelPattern,
		ProviderId:          rule.ProviderID,
		Priority:            rule.Priority,
		FallbackProviderIds: ids,
		FallbackModels:      models,
		IsSystemDefault:     rule.IsSystemDefault,
	}, nil
}

// UpdateRoutingRule updates an existing routing rule.
func (h *Handler) UpdateRoutingRule(ctx context.Context, req *routerv1.UpdateRoutingRuleRequest) (*routerv1.RoutingRule, error) {
	idsJSON, modelsJSON := marshalFallback(req.FallbackProviderIds, req.FallbackModels)

	rule := &entity.RoutingRule{
		ID:                  req.Id,
		UserID:              req.UserId,
		ModelPattern:        req.ModelPattern,
		ProviderID:          req.ProviderId,
		Priority:            req.Priority,
		FallbackProviderIDs: idsJSON,
		FallbackModels:      modelsJSON,
		IsSystemDefault:     req.IsSystemDefault,
	}

	if err := h.service.UpdateRoutingRule(rule, req.RequestingUserId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update routing rule: %v", err)
	}

	ids, models := unmarshalFallback(rule.FallbackProviderIDs, rule.FallbackModels)
	return &routerv1.RoutingRule{
		Id:                 rule.ID,
		UserId:              rule.UserID,
		ModelPattern:        rule.ModelPattern,
		ProviderId:          rule.ProviderID,
		Priority:            rule.Priority,
		FallbackProviderIds: ids,
		FallbackModels:      models,
		IsSystemDefault:     rule.IsSystemDefault,
	}, nil
}

// DeleteRoutingRule deletes a routing rule by ID.
func (h *Handler) DeleteRoutingRule(ctx context.Context, req *routerv1.DeleteRoutingRuleRequest) (*v1.Empty, error) {
	if err := h.service.DeleteRoutingRule(req.Id, req.RequestingUserId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete routing rule: %v", err)
	}
	return &v1.Empty{}, nil
}

// RefreshRoutingTable invalidates the routing table cache.
func (h *Handler) RefreshRoutingTable(ctx context.Context, req *v1.Empty) (*v1.Empty, error) {
	if err := h.service.RefreshRoutingTable(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to refresh routing table: %v", err)
	}
	return &v1.Empty{}, nil
}

// ResolveFallback resolves a fallback route (Phase 2+ feature).
func (h *Handler) ResolveFallback(ctx context.Context, req *routerv1.ResolveFallbackRequest) (*routerv1.RouteResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "ResolveFallback is not implemented in Phase 1")
}

// Shutdown gracefully shuts down the handler.
func (h *Handler) Shutdown(ctx context.Context) error {
	// No resources to clean up in the handler itself
	return nil
}

func marshalFallback(protoIDs, protoModels []string) (string, string) {
	idsJSON, _ := json.Marshal(protoIDs)
	modelsJSON, _ := json.Marshal(protoModels)
	return string(idsJSON), string(modelsJSON)
}

func unmarshalFallback(entityIDs, entityModels string) ([]string, []string) {
	var ids []string
	var models []string
	json.Unmarshal([]byte(entityIDs), &ids)
	json.Unmarshal([]byte(entityModels), &models)
	return ids, models
}

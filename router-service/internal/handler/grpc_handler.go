package handler

import (
	"context"

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
	// Call application service with authorized models filtering
	routeResult, err := h.service.ResolveRoute(ctx, req.Model, req.AuthorizedModels)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "no route found for model %s: %v", req.Model, err)
	}

	// Convert domain entity to proto message
	return &routerv1.RouteResult{
		ProviderId:          routeResult.ProviderID,
		AdapterType:         routeResult.AdapterType,
		FallbackProviderIds: routeResult.FallbackProviderIDs,
		FallbackModels:      routeResult.FallbackModels,
	}, nil
}

// GetRoutingRules lists all routing rules.
func (h *Handler) GetRoutingRules(ctx context.Context, req *routerv1.GetRoutingRulesRequest) (*routerv1.ListRoutingRulesResponse, error) {
	rules, total, err := h.service.ListRoutingRules(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get routing rules: %v", err)
	}

	// Convert domain entities to proto messages
	protoRules := make([]*routerv1.RoutingRule, len(rules))
	for i, rule := range rules {
		protoRules[i] = &routerv1.RoutingRule{
			Id:                 rule.ID,
			ModelPattern:       rule.ModelPattern,
			ProviderId:         rule.ProviderID,
			Priority:           rule.Priority,
			FallbackProviderId: rule.FallbackProviderID,
			FallbackModel:      rule.FallbackModel,
		}
	}

	return &routerv1.ListRoutingRulesResponse{
		Rules: protoRules,
		Total: int32(total),
	}, nil
}

// CreateRoutingRule creates a new routing rule.
func (h *Handler) CreateRoutingRule(ctx context.Context, req *routerv1.CreateRoutingRuleRequest) (*routerv1.RoutingRule, error) {
	rule := &entity.RoutingRule{
		ID:                 uuid.New().String(), // Always generate UUID for new rules
		ModelPattern:       req.ModelPattern,
		ProviderID:         req.ProviderId,
		Priority:           req.Priority,
		FallbackProviderID: req.FallbackProviderId,
		FallbackModel:      req.FallbackModel,
	}

	if err := h.service.CreateRoutingRule(rule); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create routing rule: %v", err)
	}

	return &routerv1.RoutingRule{
		Id:                 rule.ID,
		ModelPattern:       rule.ModelPattern,
		ProviderId:         rule.ProviderID,
		Priority:           rule.Priority,
		FallbackProviderId: rule.FallbackProviderID,
		FallbackModel:      rule.FallbackModel,
	}, nil
}

// UpdateRoutingRule updates an existing routing rule.
func (h *Handler) UpdateRoutingRule(ctx context.Context, req *routerv1.UpdateRoutingRuleRequest) (*routerv1.RoutingRule, error) {
	rule := &entity.RoutingRule{
		ID:                 req.Id,
		ModelPattern:       req.ModelPattern,
		ProviderID:         req.ProviderId,
		Priority:           req.Priority,
		FallbackProviderID: req.FallbackProviderId,
		FallbackModel:      req.FallbackModel,
	}

	if err := h.service.UpdateRoutingRule(rule); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update routing rule: %v", err)
	}

	return &routerv1.RoutingRule{
		Id:                 rule.ID,
		ModelPattern:       rule.ModelPattern,
		ProviderId:         rule.ProviderID,
		Priority:           rule.Priority,
		FallbackProviderId: rule.FallbackProviderID,
		FallbackModel:      rule.FallbackModel,
	}, nil
}

// DeleteRoutingRule deletes a routing rule by ID.
func (h *Handler) DeleteRoutingRule(ctx context.Context, req *routerv1.DeleteRoutingRuleRequest) (*v1.Empty, error) {
	if err := h.service.DeleteRoutingRule(req.Id); err != nil {
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

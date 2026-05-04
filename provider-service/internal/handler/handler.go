package handler

import (
	"context"
	"log"

	commonv1 "github.com/ai-api-gateway/api/gen/common/v1"
	providerv1 "github.com/ai-api-gateway/api/gen/provider/v1"
	"github.com/ai-api-gateway/provider-service/internal/application"
	"github.com/ai-api-gateway/provider-service/internal/domain/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	providerv1.UnimplementedProviderServiceServer
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func maskCredentials(provider *providerv1.Provider) {
	if provider != nil {
		provider.Credentials = "***"
	}
}

func toProto(p *entity.Provider) *providerv1.Provider {
	if p == nil {
		return nil
	}
	return &providerv1.Provider{
		Id:        p.ID,
		Name:       p.Name,
		Type:       p.Type,
		BaseUrl:    p.BaseURL,
		Credentials: p.Credentials,
		Models:     p.Models,
		Status:     p.Status,
		CreatedAt:  p.CreatedAt.Unix(),
		UpdatedAt:  p.UpdatedAt.Unix(),
	}
}

func toEntity(p *providerv1.Provider) *entity.Provider {
	if p == nil {
		return nil
	}
	return &entity.Provider{
		ID:        p.Id,
		Name:       p.Name,
		Type:       p.Type,
		BaseURL:    p.BaseUrl,
		Credentials: p.Credentials,
		Models:     p.Models,
		Status:     p.Status,
	}
}

func (h *Handler) GetProvider(ctx context.Context, req *providerv1.GetProviderRequest) (*providerv1.Provider, error) {
	provider, err := h.service.GetProvider(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "provider not found: %v", err)
	}
	proto := toProto(provider)
	maskCredentials(proto)
	return proto, nil
}

func (h *Handler) CreateProvider(ctx context.Context, req *providerv1.CreateProviderRequest) (*providerv1.Provider, error) {
	log.Printf("[CreateProvider] Received: Name=%s, Type=%s, BaseUrl=%s, Status=%s", req.Name, req.Type, req.BaseUrl, req.Status)
	provider := toEntity(&providerv1.Provider{
		Id:         req.Name, // Use Name as ID for MVP (e.g., "Ollama", "OpenCode Zen")
		Name:       req.Name,
		Type:       req.Type,
		BaseUrl:    req.BaseUrl,
		Credentials: req.Credentials,
		Models:     req.Models,
		Status:     req.Status,
	})
	if err := h.service.CreateProvider(provider); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create provider: %v", err)
	}
	proto := toProto(provider)
	maskCredentials(proto)
	return proto, nil
}

func (h *Handler) UpdateProvider(ctx context.Context, req *providerv1.UpdateProviderRequest) (*providerv1.Provider, error) {
	provider := toEntity(&providerv1.Provider{
		Id:         req.Id,
		Name:        req.Name,
		Type:        req.Type,
		BaseUrl:     req.BaseUrl,
		Credentials:  req.Credentials,
		Models:      req.Models,
		Status:      req.Status,
	})
	updated, err := h.service.UpdateProvider(provider)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update provider: %v", err)
	}
	proto := toProto(updated)
	maskCredentials(proto)
	return proto, nil
}

func (h *Handler) DeleteProvider(ctx context.Context, req *providerv1.DeleteProviderRequest) (*commonv1.Empty, error) {
	if err := h.service.DeleteProvider(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete provider: %v", err)
	}
	return &commonv1.Empty{}, nil
}

func (h *Handler) ListProviders(ctx context.Context, req *providerv1.ListProvidersRequest) (*providerv1.ListProvidersResponse, error) {
	providers, total, err := h.service.ListProviders(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list providers: %v", err)
	}
	protoProviders := make([]*providerv1.Provider, 0, len(providers))
	for _, p := range providers {
		proto := toProto(p)
		maskCredentials(proto)
		protoProviders = append(protoProviders, proto)
	}
	return &providerv1.ListProvidersResponse{
		Providers: protoProviders,
		Total:     int32(total),
	}, nil
}

func (h *Handler) ListModels(ctx context.Context, req *providerv1.ListModelsRequest) (*providerv1.ListModelsResponse, error) {
	provider, err := h.service.GetProvider(req.ProviderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "provider not found: %v", err)
	}
	return &providerv1.ListModelsResponse{
		Models: provider.Models,
	}, nil
}

func (h *Handler) GetProviderByType(ctx context.Context, req *providerv1.GetProviderByTypeRequest) (*providerv1.Provider, error) {
	provider, err := h.service.GetProviderByType(req.Type)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "provider not found: %v", err)
	}
	proto := toProto(provider)
	maskCredentials(proto)
	return proto, nil
}

func (h *Handler) HealthCheck(ctx context.Context, req *providerv1.HealthCheckRequest) (*providerv1.HealthCheckResponse, error) {
	healthy, err := h.service.HealthCheck(req.ProviderId)
	if err != nil {
		return &providerv1.HealthCheckResponse{Healthy: false, Error: err.Error()}, nil
	}
	return &providerv1.HealthCheckResponse{Healthy: healthy}, nil
}

func (h *Handler) ForwardRequest(ctx context.Context, req *providerv1.ForwardRequestRequest) (*providerv1.ForwardRequestResponse, error) {
	respBody, promptTokens, completionTokens, statusCode, err := h.service.ForwardRequest(ctx, req.ProviderId, req.RequestBody, req.Headers)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "forward request failed: %v", err)
	}
	return &providerv1.ForwardRequestResponse{
		ResponseBody: respBody,
		TokenCounts: &commonv1.TokenCounts{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
		},
		StatusCode: statusCode,
	}, nil
}

func (h *Handler) StreamRequest(req *providerv1.StreamRequestRequest, stream providerv1.ProviderService_StreamRequestServer) error {
	ctx := stream.Context()
	chunkChan, errChan := h.service.StreamRequest(ctx, req.ProviderId, req.RequestBody, req.Headers)

	for {
		select {
		case chunk, ok := <-chunkChan:
			if !ok {
				return nil
			}
			if err := stream.Send(&providerv1.ProviderChunk{ChunkData: chunk}); err != nil {
				return status.Errorf(codes.Internal, "failed to send chunk: %v", err)
			}
		case err := <-errChan:
			if err != nil {
				return status.Errorf(codes.Internal, "stream request failed: %v", err)
			}
			return nil
		}
	}
}

func (h *Handler) RegisterSubscriber(ctx context.Context, req *providerv1.RegisterSubscriberRequest) (*commonv1.Empty, error) {
	if err := h.service.RegisterSubscriber(req.ServiceName, req.CallbackEndpoint); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register subscriber: %v", err)
	}
	return &commonv1.Empty{}, nil
}

func (h *Handler) UnregisterSubscriber(ctx context.Context, req *providerv1.UnregisterSubscriberRequest) (*commonv1.Empty, error) {
	if err := h.service.UnregisterSubscriber(req.ServiceName); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unregister subscriber: %v", err)
	}
	return &commonv1.Empty{}, nil
}

func (h *Handler) Shutdown(ctx context.Context) error {
	return nil
}

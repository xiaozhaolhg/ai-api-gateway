package handler

import (
	"context"
	"strconv"
)

// ProviderService defines the interface for provider operations
type ProviderService interface {
	ListProviders(ctx context.Context, page, pageSize int) (*ListProvidersResp, error)
	CreateProvider(ctx context.Context, provider *Provider) (*Provider, error)
	UpdateProvider(ctx context.Context, provider *Provider) (*Provider, error)
	DeleteProvider(ctx context.Context, id string) error
}

type Provider struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Status string   `json:"status"`
}

type ListProvidersResp struct {
	Providers []Provider `json:"providers"`
}

type AdminProvidersHandler struct {
	svc ProviderService
}

func NewAdminProvidersHandler(svc ProviderService) *AdminProvidersHandler {
	return &AdminProvidersHandler{svc: svc}
}

func (h *AdminProvidersHandler) ListProviders(page, pageSize int) (*ListProvidersResp, error) {
	return h.svc.ListProviders(context.Background(), page, pageSize)
}

func (h *AdminProvidersHandler) CreateProvider(provider *Provider) (*Provider, error) {
	return h.svc.CreateProvider(context.Background(), provider)
}

func (h *AdminProvidersHandler) UpdateProvider(provider *Provider) (*Provider, error) {
	return h.svc.UpdateProvider(context.Background(), provider)
}

func (h *AdminProvidersHandler) DeleteProvider(id string) error {
	return h.svc.DeleteProvider(context.Background(), id)
}

func parsePageParams(r interface{ Get(string) string }) (int, int) {
	page, _ := strconv.Atoi(r.Get("page"))
	pageSize, _ := strconv.Atoi(r.Get("page_size"))
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	return page, pageSize
}
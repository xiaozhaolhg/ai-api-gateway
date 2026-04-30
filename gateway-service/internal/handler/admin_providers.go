package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/ai-api-gateway/gateway-service/internal/errors"
	"github.com/ai-api-gateway/gateway-service/internal/middleware"
)

type ProviderService interface {
	ListProviders(ctx context.Context, page, pageSize int32) (*client.ListProvidersResponse, error)
	CreateProvider(ctx context.Context, provider *client.Provider) (*client.Provider, error)
	UpdateProvider(ctx context.Context, provider *client.Provider) (*client.Provider, error)
	DeleteProvider(ctx context.Context, id string) error
	HealthCheck(ctx context.Context, id string) (bool, error)
}

type AdminProvidersHandler struct {
	providerSvc *client.ProviderClient
	routerSvc   *client.RouterClient
}

// NewAdminProvidersHandler creates a new AdminProvidersHandler with lazy connection clients
func NewAdminProvidersHandler(providerAddr, routerAddr string) *AdminProvidersHandler {
	// Lazy connection - clients are created without connecting
	// Connection happens on first request
	providerSvc, _ := client.NewProviderClient(providerAddr)
	routerSvc, _ := client.NewRouterClient(routerAddr)

	return &AdminProvidersHandler{
		providerSvc: providerSvc,
		routerSvc:   routerSvc,
	}
}

func (h *AdminProvidersHandler) ListProviders(c *gin.Context) {
	page, pageSize := parsePageParams(c)
	resp, err := h.providerSvc.ListProviders(c.Request.Context(), page, pageSize)
	if err != nil {
		middleware.HandleGRPCError(c, err, "provider service")
		return
	}
	for i := range resp.Providers {
		resp.Providers[i].Credentials = "***"
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AdminProvidersHandler) CreateProvider(c *gin.Context) {
	var provider client.Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		middleware.AbortWithError(c, errors.New(errors.ErrBadRequest, "invalid request"))
		return
	}
	created, err := h.providerSvc.CreateProvider(c.Request.Context(), &provider)
	if err != nil {
		middleware.HandleGRPCError(c, err, "provider service")
		return
	}
	if err := h.routerSvc.RefreshRoutingTable(c.Request.Context()); err != nil {
		// Log but don't fail - provider was created successfully
		c.JSON(http.StatusOK, created)
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *AdminProvidersHandler) UpdateProvider(c *gin.Context) {
	id := c.Param("id")
	var provider client.Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		middleware.AbortWithError(c, errors.New(errors.ErrBadRequest, "invalid request"))
		return
	}
	provider.ID = id
	updated, err := h.providerSvc.UpdateProvider(c.Request.Context(), &provider)
	if err != nil {
		middleware.HandleGRPCError(c, err, "provider service")
		return
	}
	if err := h.routerSvc.RefreshRoutingTable(c.Request.Context()); err != nil {
		// Log but don't fail - provider was updated successfully
		c.JSON(http.StatusOK, updated)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *AdminProvidersHandler) DeleteProvider(c *gin.Context) {
	id := c.Param("id")
	if err := h.providerSvc.DeleteProvider(c.Request.Context(), id); err != nil {
		middleware.HandleGRPCError(c, err, "provider service")
		return
	}
	if err := h.routerSvc.RefreshRoutingTable(c.Request.Context()); err != nil {
		// Log but don't fail - provider was deleted successfully
		c.JSON(http.StatusOK, gin.H{"message": "provider deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "provider deleted"})
}

func (h *AdminProvidersHandler) HealthCheck(c *gin.Context) {
	id := c.Param("id")
	healthy, err := h.providerSvc.HealthCheck(c.Request.Context(), id)
	if err != nil {
		middleware.HandleGRPCError(c, err, "provider service")
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "healthy": healthy})
}

func parsePageParams(c *gin.Context) (int32, int32) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return int32(page), int32(pageSize)
}

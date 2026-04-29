package handler

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
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
	routerSvc  *client.RouterClient
}

func NewAdminProvidersHandler() *AdminProvidersHandler {
	providerAddr := os.Getenv("PROVIDER_SERVICE_ADDRESS")
	if providerAddr == "" {
		providerAddr = "localhost:50053"
	}
	providerSvc, err := client.NewProviderClient(providerAddr)
	if err != nil {
		panic(err)
	}
	routerAddr := os.Getenv("ROUTER_SERVICE_ADDRESS")
	if routerAddr == "" {
		routerAddr = "localhost:50052"
	}
	routerSvc, err := client.NewRouterClient(routerAddr)
	if err != nil {
		panic(err)
	}
	return &AdminProvidersHandler{
		providerSvc: providerSvc,
		routerSvc:  routerSvc,
	}
}

func (h *AdminProvidersHandler) ListProviders(c *gin.Context) {
	page, pageSize := parsePageParams(c)
	resp, err := h.providerSvc.ListProviders(context.Background(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	created, err := h.providerSvc.CreateProvider(context.Background(), &provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.routerSvc.RefreshRoutingTable(context.Background()); err != nil {
		c.JSON(http.StatusOK, created)
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *AdminProvidersHandler) UpdateProvider(c *gin.Context) {
	id := c.Param("id")
	var provider client.Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	provider.ID = id
	updated, err := h.providerSvc.UpdateProvider(context.Background(), &provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.routerSvc.RefreshRoutingTable(context.Background()); err != nil {
		c.JSON(http.StatusOK, updated)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *AdminProvidersHandler) DeleteProvider(c *gin.Context) {
	id := c.Param("id")
	if err := h.providerSvc.DeleteProvider(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.routerSvc.RefreshRoutingTable(context.Background()); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "provider deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "provider deleted"})
}

func (h *AdminProvidersHandler) HealthCheck(c *gin.Context) {
	id := c.Param("id")
	healthy, err := h.providerSvc.HealthCheck(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "healthy": healthy})
}

func parsePageParams(c *gin.Context) (int32, int32) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return int32(page), int32(pageSize)
}

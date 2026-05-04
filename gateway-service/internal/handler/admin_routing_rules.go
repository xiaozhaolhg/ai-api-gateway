package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/ai-api-gateway/gateway-service/internal/errors"
	"github.com/ai-api-gateway/gateway-service/internal/middleware"
)

type AdminRoutingRulesHandler struct {
	routerSvc *client.RouterClient
}

// NewAdminRoutingRulesHandler creates a new AdminRoutingRulesHandler with lazy connection
func NewAdminRoutingRulesHandler(routerAddr string) *AdminRoutingRulesHandler {
	routerSvc, _ := client.NewRouterClient(routerAddr)
	return &AdminRoutingRulesHandler{
		routerSvc: routerSvc,
	}
}

func (h *AdminRoutingRulesHandler) ListRoutingRules(c *gin.Context) {
	page, pageSize := parsePageParams(c)
	resp, err := h.routerSvc.ListRoutingRules(c.Request.Context(), page, pageSize)
	if err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}
	c.JSON(http.StatusOK, resp.Rules)
}

func (h *AdminRoutingRulesHandler) CreateRoutingRule(c *gin.Context) {
	var rule client.RoutingRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		middleware.AbortWithError(c, errors.New(errors.ErrBadRequest, "invalid request"))
		return
	}
	created, err := h.routerSvc.CreateRoutingRule(c.Request.Context(), &rule)
	if err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}
	if err := h.routerSvc.RefreshRoutingTable(c.Request.Context()); err != nil {
		// Log but don't fail - rule was created successfully
		c.JSON(http.StatusCreated, created)
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *AdminRoutingRulesHandler) UpdateRoutingRule(c *gin.Context) {
	id := c.Param("id")
	var rule client.RoutingRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		middleware.AbortWithError(c, errors.New(errors.ErrBadRequest, "invalid request"))
		return
	}
	rule.ID = id
	updated, err := h.routerSvc.UpdateRoutingRule(c.Request.Context(), &rule)
	if err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}
	if err := h.routerSvc.RefreshRoutingTable(c.Request.Context()); err != nil {
		// Log but don't fail - rule was updated successfully
		c.JSON(http.StatusOK, updated)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *AdminRoutingRulesHandler) DeleteRoutingRule(c *gin.Context) {
	id := c.Param("id")
	if err := h.routerSvc.DeleteRoutingRule(c.Request.Context(), id); err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}
	if err := h.routerSvc.RefreshRoutingTable(c.Request.Context()); err != nil {
		// Log but don't fail - rule was deleted successfully
		c.JSON(http.StatusOK, gin.H{"message": "routing rule deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "routing rule deleted"})
}

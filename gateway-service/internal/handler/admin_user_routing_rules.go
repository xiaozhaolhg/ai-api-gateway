package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/ai-api-gateway/gateway-service/internal/errors"
	"github.com/ai-api-gateway/gateway-service/internal/middleware"
)

type AdminUserRoutingRulesHandler struct {
	routerSvc *client.RouterClient
}

// NewAdminUserRoutingRulesHandler creates a new AdminUserRoutingRulesHandler
func NewAdminUserRoutingRulesHandler(routerAddr string) *AdminUserRoutingRulesHandler {
	routerSvc, _ := client.NewRouterClient(routerAddr)
	return &AdminUserRoutingRulesHandler{
		routerSvc: routerSvc,
	}
}

func (h *AdminUserRoutingRulesHandler) ListUserRoutingRules(c *gin.Context) {
	userID := c.Param("userId")

	resp, err := h.routerSvc.ListRoutingRules(c.Request.Context(), 0, 1000)
	if err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}

	rules := filterRulesByUser(resp.Rules, userID)
	c.JSON(http.StatusOK, rules)
}

func (h *AdminUserRoutingRulesHandler) CreateUserRoutingRule(c *gin.Context) {
	userID := c.Param("userId")

	var rule client.RoutingRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		middleware.AbortWithError(c, errors.New(errors.ErrBadRequest, "invalid request"))
		return
	}
	rule.UserID = userID

	created, err := h.routerSvc.CreateRoutingRule(c.Request.Context(), &rule)
	if err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}

	if err := h.routerSvc.RefreshRoutingTable(c.Request.Context()); err != nil {
		c.JSON(http.StatusCreated, created)
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *AdminUserRoutingRulesHandler) GetUserRoutingRule(c *gin.Context) {
	userID := c.Param("userId")
	id := c.Param("id")

	resp, err := h.routerSvc.ListRoutingRules(c.Request.Context(), 0, 1000)
	if err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}

	rule := findRuleByID(resp.Rules, id)
	if rule == nil || rule.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Routing rule not found"})
		return
	}

	c.JSON(http.StatusOK, rule)
}

func (h *AdminUserRoutingRulesHandler) UpdateUserRoutingRule(c *gin.Context) {
	userID := c.Param("userId")
	id := c.Param("id")

	var rule client.RoutingRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		middleware.AbortWithError(c, errors.New(errors.ErrBadRequest, "invalid request"))
		return
	}
	rule.ID = id
	rule.UserID = userID

	updated, err := h.routerSvc.UpdateRoutingRule(c.Request.Context(), &rule, "")
	if err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}

	if err := h.routerSvc.RefreshRoutingTable(c.Request.Context()); err != nil {
		c.JSON(http.StatusOK, updated)
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *AdminUserRoutingRulesHandler) DeleteUserRoutingRule(c *gin.Context) {
	userID := c.Param("userId")
	id := c.Param("id")

	// First verify the rule belongs to the user
	resp, err := h.routerSvc.ListRoutingRules(c.Request.Context(), 0, 1000)
	if err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}

	rule := findRuleByID(resp.Rules, id)
	if rule == nil || rule.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Routing rule not found"})
		return
	}

	if err := h.routerSvc.DeleteRoutingRule(c.Request.Context(), id, ""); err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}

	if err := h.routerSvc.RefreshRoutingTable(c.Request.Context()); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Routing rule deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Routing rule deleted"})
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
	"github.com/ai-api-gateway/gateway-service/internal/errors"
	"github.com/ai-api-gateway/gateway-service/internal/middleware"
)

type UserRoutingRulesHandler struct {
	routerSvc *client.RouterClient
}

// NewUserRoutingRulesHandler creates a new UserRoutingRulesHandler
func NewUserRoutingRulesHandler(routerAddr string) *UserRoutingRulesHandler {
	routerSvc, _ := client.NewRouterClient(routerAddr)
	return &UserRoutingRulesHandler{
		routerSvc: routerSvc,
	}
}

func (h *UserRoutingRulesHandler) ListRoutingRules(c *gin.Context) {
	userID, _ := c.Get("userId")

	resp, err := h.routerSvc.ListRoutingRules(c.Request.Context(), 0, 1000)
	if err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}

	rules := filterRulesByUser(resp.Rules, userID.(string))
	c.JSON(http.StatusOK, rules)
}

func (h *UserRoutingRulesHandler) CreateRoutingRule(c *gin.Context) {
	userID, _ := c.Get("userId")

	var rule client.RoutingRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		middleware.AbortWithError(c, errors.New(errors.ErrBadRequest, "invalid request"))
		return
	}

	rule.UserID = userID.(string)

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

func (h *UserRoutingRulesHandler) GetRoutingRule(c *gin.Context) {
	userID, _ := c.Get("userId")
	id := c.Param("id")

	resp, err := h.routerSvc.ListRoutingRules(c.Request.Context(), 0, 1000)
	if err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}

	rule := findRuleByID(resp.Rules, id)
	if rule == nil || rule.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, rule)
}

func (h *UserRoutingRulesHandler) UpdateRoutingRule(c *gin.Context) {
	userID, _ := c.Get("userId")
	id := c.Param("id")

	var rule client.RoutingRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		middleware.AbortWithError(c, errors.New(errors.ErrBadRequest, "invalid request"))
		return
	}

	rule.ID = id
	rule.UserID = userID.(string)

	updated, err := h.routerSvc.UpdateRoutingRule(c.Request.Context(), &rule, userID.(string))
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

func (h *UserRoutingRulesHandler) DeleteRoutingRule(c *gin.Context) {
	userID, _ := c.Get("userId")
	id := c.Param("id")

	err := h.routerSvc.DeleteRoutingRule(c.Request.Context(), id, userID.(string))
	if err != nil {
		middleware.HandleGRPCError(c, err, "router service")
		return
	}

	if err := h.routerSvc.RefreshRoutingTable(c.Request.Context()); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Routing rule deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Routing rule deleted"})
}

func filterRulesByUser(rules []*client.RoutingRule, userID string) []*client.RoutingRule {
	filtered := make([]*client.RoutingRule, 0)
	for _, r := range rules {
		if r.UserID == userID {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func findRuleByID(rules []*client.RoutingRule, id string) *client.RoutingRule {
	for _, r := range rules {
		if r.ID == id {
			return r
		}
	}
	return nil
}

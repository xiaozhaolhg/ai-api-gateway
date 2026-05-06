package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ai-api-gateway/gateway-service/internal/client"
	monitorv1 "github.com/ai-api-gateway/api/gen/monitor/v1"
)

// AdminAlertsHandler handles alert and alert rule management endpoints
type AdminAlertsHandler struct {
	monitorClient *client.MonitorClient
}

// NewAdminAlertsHandler creates a new admin alerts handler
func NewAdminAlertsHandler(monitorClient *client.MonitorClient) *AdminAlertsHandler {
	// Use monitorv1 to avoid unused import error
	_ = monitorv1.AlertRule{}
	
	return &AdminAlertsHandler{
		monitorClient: monitorClient,
	}
}

// ListAlertRules returns all alert rules
func (h *AdminAlertsHandler) ListAlertRules(c *gin.Context) {
	page := int32(1)
	pageSize := int32(10)
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = int32(p)
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 {
		pageSize = int32(ps)
	}

	if h.monitorClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "monitor service unavailable"})
		return
	}

	resp, err := h.monitorClient.ListAlertRules(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert proto to UI-compatible format
	alertRules := make([]gin.H, len(resp.Rules))
	for i, rule := range resp.Rules {
		alertRules[i] = gin.H{
			"id":         rule.GetId(),
			"name":       rule.GetMetricType() + " Alert", // Generate name from metric type
			"metric":     rule.GetMetricType(),
			"condition":  rule.GetCondition(),
			"threshold":  rule.GetThreshold(),
			"channel":    rule.GetChannel(),
			"status":     rule.GetStatus(),
			"created_at": time.Now().Format(time.RFC3339), // TODO: Add created_at to proto
			"updated_at": time.Now().Format(time.RFC3339), // TODO: Add updated_at to proto
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"alert_rules": alertRules,
		"total":       resp.GetTotal(),
	})
}

// CreateAlertRule creates a new alert rule
func (h *AdminAlertsHandler) CreateAlertRule(c *gin.Context) {
	var req struct {
		MetricType    string  `json:"metric_type" binding:"required"`
		Condition    string  `json:"condition" binding:"required"`
		Threshold    float64 `json:"threshold" binding:"required"`
		Channel      string  `json:"channel" binding:"required"`
		ChannelConfig string  `json:"channel_config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if h.monitorClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "monitor service unavailable"})
		return
	}

	resp, err := h.monitorClient.CreateAlertRule(
		c.Request.Context(),
		req.MetricType,
		req.Condition,
		req.Threshold,
		req.Channel,
		req.ChannelConfig,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         resp.GetId(),
		"metric_type": resp.GetMetricType(),
		"condition":  resp.GetCondition(),
		"threshold":  resp.GetThreshold(),
		"channel":    resp.GetChannel(),
		"status":     resp.GetStatus(),
		"created_at": time.Now().Format(time.RFC3339), // TODO: Add created_at to proto
		"updated_at": time.Now().Format(time.RFC3339), // TODO: Add updated_at to proto
	})
}

// UpdateAlertRule updates an existing alert rule
func (h *AdminAlertsHandler) UpdateAlertRule(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		MetricType    string  `json:"metric_type"`
		Condition    string  `json:"condition"`
		Threshold    float64 `json:"threshold"`
		Channel      string  `json:"channel"`
		ChannelConfig string  `json:"channel_config"`
		Status       string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if h.monitorClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "monitor service unavailable"})
		return
	}

	resp, err := h.monitorClient.UpdateAlertRule(
		c.Request.Context(),
		id,
		req.MetricType,
		req.Condition,
		req.Threshold,
		req.Channel,
		req.ChannelConfig,
		req.Status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         resp.GetId(),
		"metric_type": resp.GetMetricType(),
		"condition":  resp.GetCondition(),
		"threshold":  resp.GetThreshold(),
		"channel":    resp.GetChannel(),
		"status":     resp.GetStatus(),
		"created_at": time.Now().Format(time.RFC3339), // TODO: Add created_at to proto
		"updated_at": time.Now().Format(time.RFC3339), // TODO: Add updated_at to proto
	})
}

// DeleteAlertRule deletes an alert rule
func (h *AdminAlertsHandler) DeleteAlertRule(c *gin.Context) {
	id := c.Param("id")

	if h.monitorClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "monitor service unavailable"})
		return
	}

	if err := h.monitorClient.DeleteAlertRule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert rule deleted"})
}

// ListAlerts returns all alerts
func (h *AdminAlertsHandler) ListAlerts(c *gin.Context) {
	ruleID := c.Query("rule_id")
	status := c.Query("status")
	page := int32(1)
	pageSize := int32(10)
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = int32(p)
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 {
		pageSize = int32(ps)
	}

	if h.monitorClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "monitor service unavailable"})
		return
	}

	resp, err := h.monitorClient.GetAlerts(c.Request.Context(), ruleID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert proto to UI-compatible format
	alerts := make([]gin.H, len(resp.Alerts))
	for i, alert := range resp.Alerts {
		alertData := gin.H{
			"id":           alert.GetId(),
			"rule_id":     alert.GetRuleId(),
			"severity":     "warning", // TODO: Add severity to proto
			"status":       alert.GetStatus(),
			"triggered_at": alert.GetTriggeredAt(),
			"description":  "Alert triggered", // TODO: Add description to proto
		}
		if alert.GetAcknowledgedAt() > 0 {
			alertData["acknowledged_at"] = alert.GetAcknowledgedAt()
		}
		alerts[i] = alertData
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"total":   resp.GetTotal(),
	})
}

// AcknowledgeAlert acknowledges an alert
func (h *AdminAlertsHandler) AcknowledgeAlert(c *gin.Context) {
	id := c.Param("id")

	if h.monitorClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "monitor service unavailable"})
		return
	}

	resp, err := h.monitorClient.AcknowledgeAlert(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             resp.GetId(),
		"rule_id":         resp.GetRuleId(),
		"severity":       "warning", // TODO: Add severity to proto
		"status":         resp.GetStatus(),
		"triggered_at":    resp.GetTriggeredAt(),
		"description":     "Alert acknowledged", // TODO: Add description to proto
		"acknowledged_at": resp.GetAcknowledgedAt(),
	})
}

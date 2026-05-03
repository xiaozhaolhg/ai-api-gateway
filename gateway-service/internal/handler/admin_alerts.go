package handler

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// AlertRule represents an alert rule configuration
type AlertRule struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Metric    string    `json:"metric"`
	Condition string    `json:"condition"`
	Threshold float64   `json:"threshold"`
	Channel   string    `json:"channel"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Alert represents a triggered alert
type Alert struct {
	ID            string    `json:"id"`
	RuleID        string    `json:"rule_id"`
	Severity      string    `json:"severity"`
	Status        string    `json:"status"`
	TriggeredAt   time.Time `json:"triggered_at"`
	Description   string    `json:"description"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
}

// AdminAlertsHandler handles alert and alert rule management endpoints
type AdminAlertsHandler struct {
	mu         sync.RWMutex
	alertRules []AlertRule
	alerts     []Alert
}

// NewAdminAlertsHandler creates a new admin alerts handler
func NewAdminAlertsHandler() *AdminAlertsHandler {
	h := &AdminAlertsHandler{
		alertRules: make([]AlertRule, 0),
		alerts:     make([]Alert, 0),
	}
	
	// Initialize with some default alert rules
	h.initializeDefaultData()
	
	return h
}

func (h *AdminAlertsHandler) initializeDefaultData() {
	now := time.Now()
	
	// Default alert rules
	h.alertRules = []AlertRule{
		{
			ID:        "1",
			Name:      "Budget Alert",
			Metric:    "budget_usage",
			Condition: "greater_than",
			Threshold: 80,
			Channel:   "email",
			Status:    "active",
			CreatedAt: now.Add(-30 * 24 * time.Hour),
			UpdatedAt: now.Add(-30 * 24 * time.Hour),
		},
		{
			ID:        "2",
			Name:      "Error Rate Alert",
			Metric:    "error_rate",
			Condition: "greater_than",
			Threshold: 5,
			Channel:   "slack",
			Status:    "active",
			CreatedAt: now.Add(-30 * 24 * time.Hour),
			UpdatedAt: now.Add(-30 * 24 * time.Hour),
		},
		{
			ID:        "3",
			Name:      "Latency Alert",
			Metric:    "latency",
			Condition: "greater_than",
			Threshold: 2000,
			Channel:   "email",
			Status:    "active",
			CreatedAt: now.Add(-30 * 24 * time.Hour),
			UpdatedAt: now.Add(-30 * 24 * time.Hour),
		},
	}
	
	// Default alerts
	ackTime := now.Add(-24 * time.Hour)
	h.alerts = []Alert{
		{
			ID:          "1",
			RuleID:      "1",
			Severity:    "warning",
			Status:      "triggered",
			TriggeredAt: now.Add(-48 * time.Hour),
			Description: "Budget usage exceeded 80% for Development Budget",
		},
		{
			ID:             "2",
			RuleID:         "2",
			Severity:       "critical",
			Status:         "acknowledged",
			TriggeredAt:    now.Add(-72 * time.Hour),
			Description:    "Error rate exceeded 5% for provider OpenCode Zen",
			AcknowledgedAt: &ackTime,
		},
	}
}

// ListAlertRules returns all alert rules
func (h *AdminAlertsHandler) ListAlertRules(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	// Convert to response format with ISO timestamps
	rules := make([]gin.H, len(h.alertRules))
	for i, rule := range h.alertRules {
		rules[i] = gin.H{
			"id":         rule.ID,
			"name":       rule.Name,
			"metric":     rule.Metric,
			"condition":  rule.Condition,
			"threshold":  rule.Threshold,
			"channel":    rule.Channel,
			"status":     rule.Status,
			"created_at": rule.CreatedAt.Format(time.RFC3339),
			"updated_at": rule.UpdatedAt.Format(time.RFC3339),
		}
	}
	
	c.JSON(http.StatusOK, rules)
}

// CreateAlertRule creates a new alert rule
func (h *AdminAlertsHandler) CreateAlertRule(c *gin.Context) {
	var req struct {
		Name      string  `json:"name" binding:"required"`
		Metric    string  `json:"metric" binding:"required"`
		Condition string  `json:"condition" binding:"required"`
		Threshold float64 `json:"threshold" binding:"required"`
		Channel   string  `json:"channel" binding:"required"`
		Status    string  `json:"status"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	
	if req.Status == "" {
		req.Status = "active"
	}
	
	h.mu.Lock()
	defer h.mu.Unlock()
	
	now := time.Now()
	rule := AlertRule{
		ID:        generateID(),
		Name:      req.Name,
		Metric:    req.Metric,
		Condition: req.Condition,
		Threshold: req.Threshold,
		Channel:   req.Channel,
		Status:    req.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	h.alertRules = append(h.alertRules, rule)
	
	c.JSON(http.StatusCreated, gin.H{
		"id":         rule.ID,
		"name":       rule.Name,
		"metric":     rule.Metric,
		"condition":  rule.Condition,
		"threshold":  rule.Threshold,
		"channel":    rule.Channel,
		"status":     rule.Status,
		"created_at": rule.CreatedAt.Format(time.RFC3339),
		"updated_at": rule.UpdatedAt.Format(time.RFC3339),
	})
}

// UpdateAlertRule updates an existing alert rule
func (h *AdminAlertsHandler) UpdateAlertRule(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Name      string  `json:"name"`
		Metric    string  `json:"metric"`
		Condition string  `json:"condition"`
		Threshold float64 `json:"threshold"`
		Channel   string  `json:"channel"`
		Status    string  `json:"status"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	
	h.mu.Lock()
	defer h.mu.Unlock()
	
	for i, rule := range h.alertRules {
		if rule.ID == id {
			if req.Name != "" {
				rule.Name = req.Name
			}
			if req.Metric != "" {
				rule.Metric = req.Metric
			}
			if req.Condition != "" {
				rule.Condition = req.Condition
			}
			if req.Threshold != 0 {
				rule.Threshold = req.Threshold
			}
			if req.Channel != "" {
				rule.Channel = req.Channel
			}
			if req.Status != "" {
				rule.Status = req.Status
			}
			rule.UpdatedAt = time.Now()
			h.alertRules[i] = rule
			
			c.JSON(http.StatusOK, gin.H{
				"id":         rule.ID,
				"name":       rule.Name,
				"metric":     rule.Metric,
				"condition":  rule.Condition,
				"threshold":  rule.Threshold,
				"channel":    rule.Channel,
				"status":     rule.Status,
				"created_at": rule.CreatedAt.Format(time.RFC3339),
				"updated_at": rule.UpdatedAt.Format(time.RFC3339),
			})
			return
		}
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "alert rule not found"})
}

// DeleteAlertRule deletes an alert rule
func (h *AdminAlertsHandler) DeleteAlertRule(c *gin.Context) {
	id := c.Param("id")
	
	h.mu.Lock()
	defer h.mu.Unlock()
	
	for i, rule := range h.alertRules {
		if rule.ID == id {
			h.alertRules = append(h.alertRules[:i], h.alertRules[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "alert rule deleted"})
			return
		}
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "alert rule not found"})
}

// ListAlerts returns all alerts
func (h *AdminAlertsHandler) ListAlerts(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	// Convert to response format with ISO timestamps
	alerts := make([]gin.H, len(h.alerts))
	for i, alert := range h.alerts {
		alertData := gin.H{
			"id":           alert.ID,
			"rule_id":      alert.RuleID,
			"severity":     alert.Severity,
			"status":       alert.Status,
			"triggered_at": alert.TriggeredAt.Format(time.RFC3339),
			"description":  alert.Description,
		}
		if alert.AcknowledgedAt != nil {
			alertData["acknowledged_at"] = alert.AcknowledgedAt.Format(time.RFC3339)
		}
		alerts[i] = alertData
	}
	
	c.JSON(http.StatusOK, alerts)
}

// AcknowledgeAlert acknowledges an alert
func (h *AdminAlertsHandler) AcknowledgeAlert(c *gin.Context) {
	id := c.Param("id")
	
	h.mu.Lock()
	defer h.mu.Unlock()
	
	now := time.Now()
	
	for i, alert := range h.alerts {
		if alert.ID == id {
			h.alerts[i].Status = "acknowledged"
			h.alerts[i].AcknowledgedAt = &now
			c.JSON(http.StatusOK, gin.H{"message": "alert acknowledged"})
			return
		}
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
}

// generateID creates a simple unique ID
var idCounter int
var idMutex sync.Mutex

func generateID() string {
	idMutex.Lock()
	defer idMutex.Unlock()
	idCounter++
	return "rule-" + strconv.Itoa(idCounter) + "-" + strconv.FormatInt(time.Now().Unix(), 10)
}

package application

import (
	"fmt"
	"time"

	"github.com/ai-api-gateway/monitor-service/internal/domain/entity"
	"github.com/ai-api-gateway/monitor-service/internal/domain/port"
)

// Service handles monitoring logic
type Service struct {
	metricRepo     port.MetricRepository
	alertRuleRepo  port.AlertRuleRepository
	alertRepo      port.AlertRepository
	healthRepo     port.ProviderHealthRepository
}

// NewService creates a new application service
func NewService(
	metricRepo port.MetricRepository,
	alertRuleRepo port.AlertRuleRepository,
	alertRepo port.AlertRepository,
	healthRepo port.ProviderHealthRepository,
) *Service {
	return &Service{
		metricRepo:    metricRepo,
		alertRuleRepo: alertRuleRepo,
		alertRepo:     alertRepo,
		healthRepo:    healthRepo,
	}
}

// RecordMetric records a metric
func (s *Service) RecordMetric(metric *entity.Metric) error {
	if err := s.metricRepo.Create(metric); err != nil {
		return fmt.Errorf("failed to record metric: %w", err)
	}

	// Evaluate alert rules
	go s.evaluateAlerts(metric)

	return nil
}

// evaluateAlerts evaluates alert rules against a metric
func (s *Service) evaluateAlerts(metric *entity.Metric) {
	rules, err := s.alertRuleRepo.GetEnabledRules()
	if err != nil {
		return
	}

	for _, rule := range rules {
		if rule.ProviderID != metric.ProviderID || rule.MetricType != metric.MetricType {
			continue
		}

		if s.shouldTriggerAlert(rule, metric.Value) {
			alert := &entity.Alert{
				AlertRuleID: rule.ID,
				ProviderID:  metric.ProviderID,
				Severity:    rule.Severity,
				Message:     fmt.Sprintf("Metric %s exceeded threshold: %.2f %s %.2f", rule.MetricType, metric.Value, rule.Operator, rule.Threshold),
				Value:       metric.Value,
				Threshold:   rule.Threshold,
				Timestamp:   time.Now(),
			}

			s.alertRepo.Create(alert)
		}
	}
}

// shouldTriggerAlert checks if an alert should be triggered
func (s *Service) shouldTriggerAlert(rule *entity.AlertRule, value float64) bool {
	switch rule.Operator {
	case ">":
		return value > rule.Threshold
	case "<":
		return value < rule.Threshold
	case "=":
		return value == rule.Threshold
	default:
		return false
	}
}

// GetMetrics retrieves metrics for a provider
func (s *Service) GetMetrics(providerID string, page, pageSize int) ([]*entity.Metric, int, error) {
	return s.metricRepo.GetByProviderID(providerID, page, pageSize)
}

// GetMetricAggregation retrieves aggregated metrics
func (s *Service) GetMetricAggregation(providerID, metricType, startDate, endDate string) (*entity.MetricAggregation, error) {
	return s.metricRepo.GetAggregation(providerID, metricType, startDate, endDate)
}

// ReportProviderHealth reports provider health status
func (s *Service) ReportProviderHealth(providerID string, status string, latencyMs int64) error {
	healthStatus := &entity.ProviderHealthStatus{
		ProviderID: providerID,
		Status:     status,
		LatencyP50: float64(latencyMs),
		LastCheck:  time.Now().Unix(),
	}

	return s.healthRepo.Upsert(healthStatus)
}

// GetProviderHealth retrieves provider health status
func (s *Service) GetProviderHealth(providerID string) (*entity.ProviderHealthStatus, error) {
	return s.healthRepo.GetByProviderID(providerID)
}

// CreateAlertRule creates a new alert rule
func (s *Service) CreateAlertRule(rule *entity.AlertRule) error {
	return s.alertRuleRepo.Create(rule)
}

// UpdateAlertRule updates an existing alert rule
func (s *Service) UpdateAlertRule(rule *entity.AlertRule) error {
	return s.alertRuleRepo.Update(rule)
}

// DeleteAlertRule deletes an alert rule
func (s *Service) DeleteAlertRule(id string) error {
	return s.alertRuleRepo.Delete(id)
}

// GetAlerts retrieves alerts for a provider
func (s *Service) GetAlerts(providerID string, page, pageSize int) ([]*entity.Alert, int, error) {
	return s.alertRepo.GetByProviderID(providerID, page, pageSize)
}

// AcknowledgeAlert acknowledges an alert
func (s *Service) AcknowledgeAlert(alertID string) error {
	return s.alertRepo.Acknowledge(alertID)
}

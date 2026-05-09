package handler

import (
	"context"
	"time"

	commonv1 "github.com/ai-api-gateway/api/gen/common/v1"
	monitorv1 "github.com/ai-api-gateway/api/gen/monitor/v1"
	"github.com/ai-api-gateway/monitor-service/internal/application"
	"github.com/ai-api-gateway/monitor-service/internal/domain/entity"
)

// Handler implements the MonitorService gRPC interface
type Handler struct {
	monitorv1.UnimplementedMonitorServiceServer
	service *application.Service
}

// NewHandler creates a new monitor service handler
func NewHandler(service *application.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// OnProviderResponse handles provider response callback
func (h *Handler) OnProviderResponse(ctx context.Context, req *commonv1.ProviderResponseCallback) (*commonv1.Empty, error) {
	// Record latency metric
	metric := &entity.Metric{
		ProviderID: req.ProviderId,
		Model:      req.Model,
		MetricType: "latency",
		Value:      float64(req.LatencyMs),
		Timestamp:  time.Now(),
	}

	err := h.service.RecordMetric(metric)
	if err != nil {
		return nil, err
	}

	return &commonv1.Empty{}, nil
}

// RecordMetric records a metric
func (h *Handler) RecordMetric(ctx context.Context, req *monitorv1.RecordMetricRequest) (*commonv1.Empty, error) {
	metric := &entity.Metric{
		ProviderID: req.Labels["provider"],
		Model:      req.Labels["model"],
		MetricType: req.MetricType,
		Value:      req.Value,
		Timestamp:  time.Unix(req.Timestamp, 0),
	}

	err := h.service.RecordMetric(metric)
	if err != nil {
		return nil, err
	}

	return &commonv1.Empty{}, nil
}

// GetMetrics retrieves metrics for a provider
func (h *Handler) GetMetrics(ctx context.Context, req *monitorv1.GetMetricsRequest) (*monitorv1.ListMetricsResponse, error) {
	providerID := req.Labels["provider"]
	if providerID == "" {
		return &monitorv1.ListMetricsResponse{
			Metrics: []*monitorv1.Metric{},
			Total:   0,
		}, nil
	}

	metrics, total, err := h.service.GetMetrics(providerID, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	// Convert domain entities to proto messages
	protoMetrics := make([]*monitorv1.Metric, len(metrics))
	for i, metric := range metrics {
		protoMetrics[i] = &monitorv1.Metric{
			Id:        metric.ID,
			Type:      metric.MetricType,
			Labels:    map[string]string{"provider": metric.ProviderID, "model": metric.Model},
			Value:     metric.Value,
			Timestamp: metric.Timestamp.Unix(),
		}
	}

	return &monitorv1.ListMetricsResponse{
		Metrics: protoMetrics,
		Total:   int32(total),
	}, nil
}

// GetMetricAggregation retrieves aggregated metrics
func (h *Handler) GetMetricAggregation(ctx context.Context, req *monitorv1.GetMetricAggregationRequest) (*monitorv1.MetricAggregationResponse, error) {
	providerID := req.Labels["provider"]
	if providerID == "" {
		return &monitorv1.MetricAggregationResponse{
			Aggregation: &monitorv1.MetricAggregation{},
		}, nil
	}

	agg, err := h.service.GetMetricAggregation(providerID, req.MetricType, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}

	return &monitorv1.MetricAggregationResponse{
		Aggregation: &monitorv1.MetricAggregation{
			MetricType: agg.MetricType,
			Labels:     map[string]string{"provider": agg.ProviderID, "model": agg.Model},
			Avg:        agg.Avg,
			Min:        agg.Min,
			Max:        agg.Max,
			P50:        agg.P50,
			P95:        agg.P95,
			P99:        agg.P99,
			Count:      agg.Count,
		},
	}, nil
}

// GetProviderHealth retrieves provider health status
func (h *Handler) GetProviderHealth(ctx context.Context, req *monitorv1.GetProviderHealthRequest) (*monitorv1.ProviderHealthStatus, error) {
	status, err := h.service.GetProviderHealth(req.ProviderId)
	if err != nil {
		return nil, err
	}

	return &monitorv1.ProviderHealthStatus{
		ProviderId:    status.ProviderID,
		Status:        status.Status,
		LatencyMs:     status.LatencyMs,
		ErrorRate:     status.ErrorRate,
		LastCheckTime: status.LastCheckTime.Unix(),
		UptimeSeconds: status.UptimeSeconds,
	}, nil
}

// ReportProviderHealth reports provider health status
func (h *Handler) ReportProviderHealth(ctx context.Context, req *monitorv1.ReportProviderHealthRequest) (*commonv1.Empty, error) {
	err := h.service.ReportProviderHealth(req.ProviderId, req.Status, int64(req.Latency))
	if err != nil {
		return nil, err
	}

	return &commonv1.Empty{}, nil
}

// CreateAlertRule creates a new alert rule
func (h *Handler) CreateAlertRule(ctx context.Context, req *monitorv1.CreateAlertRuleRequest) (*monitorv1.AlertRule, error) {
	err := h.service.CreateAlertRule(req.Rule)
	if err != nil {
		return nil, err
	}

	return req.Rule, nil
}

// UpdateAlertRule updates an existing alert rule
func (h *Handler) UpdateAlertRule(ctx context.Context, req *monitorv1.UpdateAlertRuleRequest) (*monitorv1.AlertRule, error) {
	err := h.service.UpdateAlertRule(req.Rule)
	if err != nil {
		return nil, err
	}

	return req.Rule, nil
}

// DeleteAlertRule deletes an alert rule
func (h *Handler) DeleteAlertRule(ctx context.Context, req *monitorv1.DeleteAlertRuleRequest) (*commonv1.Empty, error) {
	err := h.service.DeleteAlertRule(req.Id)
	if err != nil {
		return nil, err
	}

	return &commonv1.Empty{}, nil
}

// GetAlerts retrieves alerts for a provider
func (h *Handler) GetAlerts(ctx context.Context, req *monitorv1.GetAlertsRequest) (*monitorv1.ListAlertsResponse, error) {
	alerts, total, err := h.service.GetAlerts(req.ProviderId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	// Convert domain entities to proto messages
	protoAlerts := make([]*monitorv1.Alert, len(alerts))
	for i, alert := range alerts {
		var acknowledgedAt int64
		if alert.AcknowledgedAt != nil {
			acknowledgedAt = alert.AcknowledgedAt.Unix()
		}

		protoAlerts[i] = &monitorv1.Alert{
			Id:             alert.ID,
			AlertRuleId:    alert.AlertRuleID,
			ProviderId:     alert.ProviderID,
			Severity:       alert.Severity,
			Message:        alert.Message,
			Value:          alert.Value,
			Threshold:      alert.Threshold,
			Timestamp:      alert.Timestamp.Unix(),
			Acknowledged:   alert.Acknowledged,
			AcknowledgedAt: acknowledgedAt,
		}
	}

	return &monitorv1.ListAlertsResponse{
		Alerts: protoAlerts,
		Total:  int32(total),
	}, nil
}

// AcknowledgeAlert acknowledges an alert
func (h *Handler) AcknowledgeAlert(ctx context.Context, req *monitorv1.AcknowledgeAlertRequest) (*monitorv1.Alert, error) {
	err := h.service.AcknowledgeAlert(req.AlertId)
	if err != nil {
		return nil, err
	}

	// Return the acknowledged alert
	// For MVP, we'll return a placeholder
	return &monitorv1.Alert{}, nil
}

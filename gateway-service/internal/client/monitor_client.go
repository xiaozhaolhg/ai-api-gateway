package client

import (
	"context"
	"fmt"

	"github.com/ai-api-gateway/api/gen/monitor/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// MonitorClient wraps the monitor-service gRPC client
type MonitorClient struct {
	client monitorv1.MonitorServiceClient
	conn   *grpc.ClientConn
}

// NewMonitorClient creates a new monitor service client
func NewMonitorClient(address string) (*MonitorClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to monitor service: %w", err)
	}

	return &MonitorClient{
		client: monitorv1.NewMonitorServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close closes the connection
func (c *MonitorClient) Close() error {
	return c.conn.Close()
}

// OnProviderResponse handles provider response callback
func (c *MonitorClient) OnProviderResponse(ctx context.Context, callback *monitorv1.ProviderResponseCallback) error {
	req := &monitorv1.OnProviderResponseRequest{
		Callback: callback,
	}

	_, err := c.client.OnProviderResponse(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to handle provider response callback: %w", err)
	}

	return nil
}

// RecordMetric records a metric
func (c *MonitorClient) RecordMetric(ctx context.Context, metric *monitorv1.Metric) error {
	req := &monitorv1.RecordMetricRequest{
		Metric: metric,
	}

	_, err := c.client.RecordMetric(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to record metric: %w", err)
	}

	return nil
}

// GetMetrics retrieves metrics
func (c *MonitorClient) GetMetrics(ctx context.Context, providerID string, page, pageSize int32) (*monitorv1.ListMetricsResponse, error) {
	req := &monitorv1.GetMetricsRequest{
		ProviderId: providerID,
		Page:       page,
		PageSize:   pageSize,
	}

	resp, err := c.client.GetMetrics(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return resp, nil
}

// GetMetricAggregation retrieves aggregated metrics
func (c *MonitorClient) GetMetricAggregation(ctx context.Context, providerID string, startDate, endDate string) (*monitorv1.MetricAggregationResponse, error) {
	req := &monitorv1.GetMetricAggregationRequest{
		ProviderId: providerID,
		StartDate:  startDate,
		EndDate:   endDate,
	}

	resp, err := c.client.GetMetricAggregation(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric aggregation: %w", err)
	}

	return resp, nil
}

// GetProviderHealth retrieves provider health status
func (c *MonitorClient) GetProviderHealth(ctx context.Context, providerID string) (*monitorv1.ProviderHealthStatus, error) {
	req := &monitorv1.GetProviderHealthRequest{
		ProviderId: providerID,
	}

	resp, err := c.client.GetProviderHealth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider health: %w", err)
	}

	return resp, nil
}

// ReportProviderHealth reports provider health status
func (c *MonitorClient) ReportProviderHealth(ctx context.Context, providerID string, status string, latencyMs int64) error {
	req := &monitorv1.ReportProviderHealthRequest{
		ProviderId: providerID,
		Status:     status,
		LatencyMs:  latencyMs,
	}

	_, err := c.client.ReportProviderHealth(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to report provider health: %w", err)
	}

	return nil
}

// CreateAlertRule creates a new alert rule
func (c *MonitorClient) CreateAlertRule(ctx context.Context, rule *monitorv1.AlertRule) (*monitorv1.AlertRule, error) {
	req := &monitorv1.CreateAlertRuleRequest{
		Rule: rule,
	}

	resp, err := c.client.CreateAlertRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert rule: %w", err)
	}

	return resp, nil
}

// UpdateAlertRule updates an existing alert rule
func (c *MonitorClient) UpdateAlertRule(ctx context.Context, rule *monitorv1.AlertRule) (*monitorv1.AlertRule, error) {
	req := &monitorv1.UpdateAlertRuleRequest{
		Rule: rule,
	}

	resp, err := c.client.UpdateAlertRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update alert rule: %w", err)
	}

	return resp, nil
}

// DeleteAlertRule deletes an alert rule
func (c *MonitorClient) DeleteAlertRule(ctx context.Context, id string) error {
	req := &monitorv1.DeleteAlertRuleRequest{
		Id: id,
	}

	_, err := c.client.DeleteAlertRule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete alert rule: %w", err)
	}

	return nil
}

// GetAlerts retrieves alerts
func (c *MonitorClient) GetAlerts(ctx context.Context, providerID string, page, pageSize int32) (*monitorv1.ListAlertsResponse, error) {
	req := &monitorv1.GetAlertsRequest{
		ProviderId: providerID,
		Page:       page,
		PageSize:   pageSize,
	}

	resp, err := c.client.GetAlerts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}

	return resp, nil
}

// AcknowledgeAlert acknowledges an alert
func (c *MonitorClient) AcknowledgeAlert(ctx context.Context, alertID string) (*monitorv1.Alert, error) {
	req := &monitorv1.AcknowledgeAlertRequest{
		AlertId: alertID,
	}

	resp, err := c.client.AcknowledgeAlert(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to acknowledge alert: %w", err)
	}

	return resp, nil
}

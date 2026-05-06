package client

import (
	"context"
	"fmt"

	monitorv1 "github.com/ai-api-gateway/api/gen/monitor/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// MonitorClient wraps monitor-service gRPC client
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

// Close closes connection
func (c *MonitorClient) Close() error {
	return c.conn.Close()
}

// ListAlertRules retrieves alert rules
func (c *MonitorClient) ListAlertRules(ctx context.Context, page, pageSize int32) (*monitorv1.ListAlertRulesResponse, error) {
	req := &monitorv1.ListAlertRulesRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.ListAlertRules(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list alert rules: %w", err)
	}

	return resp, nil
}

// CreateAlertRule creates a new alert rule
func (c *MonitorClient) CreateAlertRule(ctx context.Context, metricType, condition string, threshold float64, channel, channelConfig string) (*monitorv1.AlertRule, error) {
	req := &monitorv1.CreateAlertRuleRequest{
		MetricType:   metricType,
		Condition:    condition,
		Threshold:    threshold,
		Channel:      channel,
		ChannelConfig: channelConfig,
	}

	resp, err := c.client.CreateAlertRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert rule: %w", err)
	}

	return resp, nil
}

// UpdateAlertRule updates an existing alert rule
func (c *MonitorClient) UpdateAlertRule(ctx context.Context, id, metricType, condition string, threshold float64, channel, channelConfig, status string) (*monitorv1.AlertRule, error) {
	req := &monitorv1.UpdateAlertRuleRequest{
		Id:           id,
		MetricType:   metricType,
		Condition:    condition,
		Threshold:    threshold,
		Channel:      channel,
		ChannelConfig: channelConfig,
		Status:       status,
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
func (c *MonitorClient) GetAlerts(ctx context.Context, ruleID, status string, page, pageSize int32) (*monitorv1.ListAlertsResponse, error) {
	req := &monitorv1.GetAlertsRequest{
		RuleId: ruleID,
		Status:  status,
		Page:    page,
		PageSize: pageSize,
	}

	resp, err := c.client.GetAlerts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}

	return resp, nil
}

// AcknowledgeAlert acknowledges an alert
func (c *MonitorClient) AcknowledgeAlert(ctx context.Context, id string) (*monitorv1.Alert, error) {
	req := &monitorv1.AcknowledgeAlertRequest{
		Id: id,
	}

	resp, err := c.client.AcknowledgeAlert(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to acknowledge alert: %w", err)
	}

	return resp, nil
}
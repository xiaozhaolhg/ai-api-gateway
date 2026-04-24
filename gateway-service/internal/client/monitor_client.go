package client

import (
	"context"
	"fmt"
	"time"
)

type MonitorClient struct{}

func NewMonitorClient(address string) (*MonitorClient, error) {
	return &MonitorClient{}, nil
}

func (c *MonitorClient) Close() error {
	return nil
}

func (c *MonitorClient) PublishMetric(ctx context.Context, metric *Metric) error {
	return fmt.Errorf("not implemented")
}

func (c *MonitorClient) GetMetrics(ctx context.Context, service string, since time.Time) ([]*Metric, error) {
	return nil, fmt.Errorf("not implemented")
}

type Metric struct {
	Service   string         `json:"service"`
	Name      string         `json:"name"`
	Value     float64        `json:"value"`
	Timestamp time.Time     `json:"timestamp"`
	Labels   map[string]string `json:"labels"`
}
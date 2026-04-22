package port

import (
	"github.com/ai-api-gateway/monitor-service/internal/domain/entity"
)

// MetricRepository defines the interface for metric persistence operations
type MetricRepository interface {
	Create(metric *entity.Metric) error
	GetByProviderID(providerID string, page, pageSize int) ([]*entity.Metric, int, error)
	GetAggregation(providerID, metricType, startDate, endDate string) (*entity.MetricAggregation, error)
}

// AlertRuleRepository defines the interface for alert rule persistence operations
type AlertRuleRepository interface {
	Create(rule *entity.AlertRule) error
	GetByID(id string) (*entity.AlertRule, error)
	Update(rule *entity.AlertRule) error
	Delete(id string) error
	List(page, pageSize int) ([]*entity.AlertRule, int, error)
	GetEnabledRules() ([]*entity.AlertRule, error)
}

// AlertRepository defines the interface for alert persistence operations
type AlertRepository interface {
	Create(alert *entity.Alert) error
	GetByID(id string) (*entity.Alert, error)
	GetByProviderID(providerID string, page, pageSize int) ([]*entity.Alert, int, error)
	Acknowledge(id string) error
}

// ProviderHealthRepository defines the interface for provider health persistence operations
type ProviderHealthRepository interface {
	Upsert(status *entity.ProviderHealthStatus) error
	GetByProviderID(providerID string) (*entity.ProviderHealthStatus, error)
	List(page, pageSize int) ([]*entity.ProviderHealthStatus, int, error)
}

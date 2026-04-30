# monitor-service

## Purpose

Observability domain — metrics collection, alerting, provider health monitoring.

## Service Responsibility

- **Role**: Metrics, health monitoring, alerting
- **Owned Entities**: Metric, AlertRule, Alert, ProviderHealthStatus
- **Data Layer**: monitor-db (SQLite/PostgreSQL)

## Dependencies

### Calls To

| Service | Methods | Purpose |
|---|---|---|
| (none) | — | Does not call other internal services |

### Called By

| Service | Methods | Purpose |
|---|---|---|
| provider-service | `OnProviderResponse` | Latency/error data from callbacks |
| provider-service | `ReportProviderHealth` | Periodic health probes |
| gateway-service | `RecordMetric` | Direct metric recording |
| gateway-service | `GetMetrics`, `GetProviderHealth` | Query operations |
| gateway-service | `CreateAlertRule`, `GetAlerts` | Alert management |

### Data Dependencies

- **Database**: monitor-db (Metric, AlertRule, Alert, ProviderHealthStatus)
- **Cache**: Redis (recent metrics, health status)

## Key Design

### Metrics Collection

Two paths:
1. **Provider callback** (automatic): provider-service dispatches OnProviderResponse with latency, error data
2. **Direct call** (any service): Any service calls RecordMetric for custom metrics

### Provider Health

- Real-time data from provider callbacks
- Periodic probes from provider-service via ReportProviderHealth
- Aggregate: latency_p50/p95/p99, error_rate, uptime_pct

### Alerting

- Threshold-based rules (e.g., latency > 5s, error_rate > 5%)
- Channels: email, webhook, Slack
- States: firing → acknowledged → resolved

### Key Operations

- **OnProviderResponse**: Receive callback data from provider-service
- **RecordMetric**: Custom metric recording
- **GetMetrics/MetricAggregation**: Query metrics
- **GetProviderHealth/ReportProviderHealth**: Health tracking
- **CreateAlertRule/UpdateAlertRule/DeleteAlertRule**: Alert rule CRUD
- **GetAlerts/AcknowledgeAlert**: Alert lifecycle

## Requirements

### Requirement: Metrics collection
The monitor-service SHALL collect and store metrics from provider callbacks and direct service calls.

#### Scenario: Provider callback metrics
- **WHEN** provider-service dispatches OnProviderResponse callback
- **THEN** the monitor-service SHALL store latency and error data as metrics

#### Scenario: Direct metric recording
- **WHEN** any service calls RecordMetric RPC
- **THEN** the monitor-service SHALL store the custom metric with timestamp
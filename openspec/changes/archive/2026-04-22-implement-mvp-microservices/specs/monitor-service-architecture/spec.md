## MODIFIED Requirements

### Requirement: DDD four-layer architecture
The monitor-service SHALL implement four-layer Clean Architecture: Domain, Application, Infrastructure, and Handler with dependency direction from outer to inner layers.

#### Scenario: Domain layer has no external dependencies
- **WHEN** the domain layer is imported
- **THEN** it SHALL NOT import any code from application, infrastructure, or handler layers

### Requirement: Metric entity and repository
The monitor-service SHALL own the Metric entity with fields: id, type, labels (map), value, timestamp. It SHALL provide a MetricRepository interface.

#### Scenario: Record metric from provider callback
- **WHEN** an OnProviderResponse callback is received
- **THEN** the service SHALL extract latency, error data and create Metric records

#### Scenario: Record custom metric
- **WHEN** a RecordMetric request is received
- **THEN** the service SHALL persist the Metric to SQLite

### Requirement: ProviderHealthStatus entity and repository
The monitor-service SHALL own the ProviderHealthStatus entity with fields: provider_id, latency_p50, latency_p95, latency_p99, error_rate, uptime_pct, last_check, status. It SHALL provide a ProviderHealthRepository interface.

#### Scenario: Update provider health from callbacks
- **WHEN** OnProviderResponse callbacks accumulate for a provider
- **THEN** the service SHALL update the provider's health status with aggregated latency and error rate

#### Scenario: Get provider health
- **WHEN** a GetProviderHealth request is received
- **THEN** the service SHALL return the current ProviderHealthStatus for the provider

### Requirement: AlertRule and Alert entities
The monitor-service SHALL own AlertRule (id, metric_type, condition, threshold, channel, channel_config, status) and Alert (id, rule_id, triggered_at, status, acknowledged_at, resolved_at) entities with repositories.

#### Scenario: Create alert rule
- **WHEN** a CreateAlertRule request is received
- **THEN** the service SHALL persist the rule and evaluate it against incoming metrics

#### Scenario: Alert triggered
- **WHEN** a metric value exceeds an alert rule threshold
- **THEN** the service SHALL create an Alert with status "firing"

### Requirement: gRPC server implementation
The monitor-service SHALL implement the MonitorService gRPC server as defined in the api/ module proto definitions.

#### Scenario: All proto RPCs are implemented
- **WHEN** the monitor-service starts
- **THEN** it SHALL register all RPCs defined in monitor.proto: OnProviderResponse, RecordMetric, GetMetrics, GetMetricAggregation, GetProviderHealth, ListProviderHealth, ReportProviderHealth, and CRUD for alert rules and alert lifecycle

### Requirement: SQLite persistence
The monitor-service SHALL use SQLite as its database with tables for metrics, provider_health_statuses, alert_rules, and alerts, managed via migrations.

#### Scenario: Database initialized on startup
- **WHEN** the monitor-service starts
- **THEN** it SHALL create the SQLite database file and run migrations if the database is new

## MODIFIED Requirements

### Requirement: Unit tests for domain layer
The monitor-service SHALL have unit tests for domain layer components with no external dependencies.

#### Scenario: Metric entity tests
- **WHEN** unit tests run for Metric entity
- **THEN** field validation and type checking SHALL be verified

#### Scenario: Alert rule evaluation tests
- **WHEN** unit tests run for alert rule evaluation
- **THEN** threshold comparison (gt, lt, eq) SHALL be verified

### Requirement: Repository integration tests
The monitor-service SHALL have integration tests that verify SQLite repository implementations.

#### Scenario: MetricRepository integration test
- **WHEN** integration tests run against SQLite MetricRepository
- **THEN** metric recording and filtered queries SHALL work correctly

#### Scenario: ProviderHealthRepository integration test
- **WHEN** integration tests run against SQLite ProviderHealthRepository
- **THEN** health status updates and queries SHALL work correctly

### Requirement: gRPC server tests
The monitor-service SHALL have tests that verify gRPC server behavior.

#### Scenario: OnProviderResponse gRPC test
- **WHEN** a gRPC OnProviderResponse callback is received
- **THEN** the service SHALL create Metric records and update ProviderHealthStatus

#### Scenario: GetMetrics gRPC test
- **WHEN** a gRPC GetMetrics request is sent with filters
- **THEN** the service SHALL return matching metrics

### Requirement: Test Coverage
The monitor-service domain and application layers SHALL maintain at least 70% code coverage.

#### Scenario: Coverage check
- **WHEN** `go test -cover ./internal/domain/... ./internal/application/...` is run
- **THEN** coverage SHALL be at least 70%

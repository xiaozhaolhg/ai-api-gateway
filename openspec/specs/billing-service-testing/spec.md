## Purpose

Testing specifications for the billing-service microservice with gRPC interface and SQLite persistence.
## Requirements
### Requirement: Unit tests for domain layer
The billing-service SHALL have unit tests for domain layer components with no external dependencies.

#### Scenario: UsageRecord entity tests
- **WHEN** unit tests run for UsageRecord entity
- **THEN** field validation and cost calculation SHALL be verified

#### Scenario: Budget check logic tests
- **WHEN** unit tests run for budget checking
- **THEN** soft cap and hard cap threshold detection SHALL be verified

### Requirement: Repository integration tests
The billing-service SHALL have integration tests that verify SQLite repository implementations.

#### Scenario: UsageRecordRepository integration test
- **WHEN** integration tests run against SQLite UsageRecordRepository
- **THEN** CRUD operations and filtered queries SHALL work correctly

#### Scenario: PricingRuleRepository integration test
- **WHEN** integration tests run against SQLite PricingRuleRepository
- **THEN** pricing rule CRUD and cost calculation SHALL work correctly

### Requirement: gRPC server tests
The billing-service SHALL have tests that verify gRPC server behavior.

#### Scenario: OnProviderResponse gRPC test
- **WHEN** a gRPC OnProviderResponse callback is received
- **THEN** the service SHALL create a UsageRecord with correct token counts and cost

#### Scenario: GetUsage gRPC test
- **WHEN** a gRPC GetUsage request is sent with filters
- **THEN** the service SHALL return matching usage records

### Requirement: Test Coverage
The billing-service domain and application layers SHALL maintain at least 70% code coverage.

#### Scenario: Coverage check
- **WHEN** `go test -cover ./internal/domain/... ./internal/application/...` is run
- **THEN** coverage SHALL be at least 70%


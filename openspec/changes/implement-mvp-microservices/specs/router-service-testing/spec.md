## MODIFIED Requirements

### Requirement: Unit Tests for Domain Layer
The router-service SHALL have unit tests for domain layer components that have no external dependencies.

#### Scenario: Routing logic tests
- **WHEN** unit tests run for the Router component
- **THEN** model pattern matching (exact, wildcard) SHALL have test coverage

#### Scenario: RoutingRule entity tests
- **WHEN** unit tests run for RoutingRule entity
- **THEN** priority ordering and pattern matching SHALL be verified

### Requirement: Integration Tests for Repository
The router-service SHALL have integration tests that verify SQLite repository implementations.

#### Scenario: RoutingRuleRepository integration test
- **WHEN** integration tests run against SQLite RoutingRuleRepository
- **THEN** CRUD operations and pattern matching queries SHALL work correctly

### Requirement: gRPC server tests
The router-service SHALL have tests that verify gRPC server behavior using in-process gRPC connections.

#### Scenario: ResolveRoute gRPC test
- **WHEN** a gRPC ResolveRoute request is sent with a known model
- **THEN** the response SHALL contain correct provider_id and adapter_type

#### Scenario: ResolveRoute not found
- **WHEN** a gRPC ResolveRoute request is sent with an unknown model
- **THEN** the response SHALL return a NOT_FOUND gRPC error

### Requirement: Test Coverage
The router-service domain and application layers SHALL maintain at least 70% code coverage.

#### Scenario: Coverage check
- **WHEN** `go test -cover ./internal/domain/... ./internal/application/...` is run
- **THEN** coverage SHALL be at least 70%

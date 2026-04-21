## ADDED Requirements

### Requirement: Unit tests for domain layer
The auth-service SHALL have unit tests for domain layer components with no external dependencies.

#### Scenario: User entity tests
- **WHEN** unit tests run for User entity
- **THEN** field validation and status transitions SHALL be verified

#### Scenario: APIKey hashing tests
- **WHEN** unit tests run for API key creation
- **THEN** key generation, hashing, and verification SHALL be tested

### Requirement: Repository integration tests
The auth-service SHALL have integration tests that verify SQLite repository implementations.

#### Scenario: UserRepository integration test
- **WHEN** integration tests run against SQLite UserRepository
- **THEN** CRUD operations SHALL work correctly with a test database

#### Scenario: APIKeyRepository integration test
- **WHEN** integration tests run against SQLite APIKeyRepository
- **THEN** key creation, hash lookup, and deletion SHALL work correctly

### Requirement: gRPC server tests
The auth-service SHALL have tests that verify gRPC server behavior using in-process gRPC connections.

#### Scenario: ValidateAPIKey gRPC test
- **WHEN** a gRPC ValidateAPIKey request is sent with a valid key
- **THEN** the response SHALL contain correct UserIdentity

#### Scenario: CreateUser gRPC test
- **WHEN** a gRPC CreateUser request is sent
- **THEN** the user SHALL be persisted and returned with generated id

### Requirement: Test coverage threshold
The auth-service domain and application layers SHALL maintain at least 70% code coverage.

#### Scenario: Coverage check
- **WHEN** `go test -cover ./internal/domain/... ./internal/application/...` is run
- **THEN** coverage SHALL be at least 70%

## MODIFIED Requirements

### Requirement: Unit tests for domain layer
The provider-service SHALL have unit tests for domain layer components with no external dependencies.

#### Scenario: Provider entity tests
- **WHEN** unit tests run for Provider entity
- **THEN** field validation and status transitions SHALL be verified

#### Scenario: Adapter interface tests
- **WHEN** unit tests run for ProviderAdapter interface
- **THEN** mock implementations SHALL verify request/response transformation contracts

### Requirement: Adapter integration tests
The provider-service SHALL have integration tests that verify provider adapters work with their respective APIs.

#### Scenario: Ollama adapter integration test
- **WHEN** integration tests run against OllamaProvider
- **THEN** if Ollama is available, requests SHALL succeed and responses SHALL be in OpenAI format

#### Scenario: OpenAI-compatible adapter test
- **WHEN** integration tests run against an OpenAI-compatible endpoint
- **THEN** requests SHALL succeed if the endpoint is available and API key is configured

### Requirement: gRPC server tests
The provider-service SHALL have tests that verify gRPC server behavior.

#### Scenario: ForwardRequest gRPC test
- **WHEN** a gRPC ForwardRequest is sent with a valid provider_id
- **THEN** the service SHALL forward the request and return the response

#### Scenario: CreateProvider gRPC test
- **WHEN** a gRPC CreateProvider is sent
- **THEN** the provider SHALL be persisted with encrypted credentials

### Requirement: Callback dispatch tests
The provider-service SHALL have tests that verify callback dispatch to subscribers.

#### Scenario: Callback dispatched after response
- **WHEN** a provider response completes
- **THEN** the service SHALL dispatch OnProviderResponse to all registered subscribers

#### Scenario: Callback failure does not block
- **WHEN** a subscriber callback fails
- **THEN** the service SHALL log the error and NOT block the response to the caller

### Requirement: Test Coverage
The provider-service domain and application layers SHALL maintain at least 70% code coverage.

#### Scenario: Coverage check
- **WHEN** `go test -cover ./internal/domain/... ./internal/application/...` is run
- **THEN** coverage SHALL be at least 70%

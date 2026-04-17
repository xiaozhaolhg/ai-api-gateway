## ADDED Requirements

### Requirement: Unit Tests for Domain Layer
The router-service SHALL have unit tests for domain layer components that have no external dependencies.

#### Scenario: Provider interface tests
- **WHEN** unit tests run for the Provider interface
- **THEN** mock implementations SHALL be used to verify routing logic works correctly

#### Scenario: Router logic tests
- **WHEN** unit tests run for the Router component
- **THEN** each routing strategy (by model prefix) SHALL have test coverage

### Requirement: Integration Tests for Provider Adapters
The router-service SHALL have integration tests that verify provider adapters work with their respective APIs.

#### Scenario: Ollama adapter integration test
- **WHEN** integration tests run against OllamaProvider
- **THEN** if Ollama is available at the configured endpoint, requests SHALL succeed

#### Scenario: OpenCode Zen adapter integration test
- **WHEN** integration tests run against OpenCodeZenProvider
- **THEN** if the API key is configured, requests SHALL succeed against the Zen API

### Requirement: End-to-End API Tests
The router-service SHALL have e2e tests that verify the complete request/response cycle through the HTTP API.

#### Scenario: Non-streaming e2e test
- **WHEN** a POST request is made to /v1/chat/completions with stream: false
- **THEN** the response SHALL be in valid OpenAI format with all required fields (id, object, created, model, choices)

#### Scenario: Streaming e2e test
- **WHEN** a POST request is made to /v1/chat/completions with stream: true
- **THEN** the response SHALL be streamed SSE with proper formatting and end with data: [DONE]

### Requirement: Model Listing Tests
The router-service SHALL have tests for the /v1/models endpoint.

#### Scenario: Models list returns configured providers
- **WHEN** GET /v1/models is called
- **THEN** the response SHALL include models from all enabled providers

### Requirement: Test Fixtures
The router-service SHALL provide test fixtures for common request/response patterns.

#### Scenario: Test fixtures available
- **WHEN** tests need sample requests
- **THEN** fixture files SHALL be available in the testdata directory

### Requirement: Test Coverage
The router-service SHALL maintain adequate test coverage for the domain and application layers.

#### Scenario: Coverage threshold
- **WHEN** go test -cover is run
- **THEN** the domain and application layers SHALL achieve at least 70% code coverage
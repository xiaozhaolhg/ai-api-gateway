## Purpose

Testing specifications for the gateway-service microservice with HTTP interface and middleware pipeline.
## Requirements
### Requirement: Unit tests for middleware pipeline
The gateway-service SHALL have unit tests for each middleware component using mock gRPC clients.

#### Scenario: Auth middleware test
- **WHEN** unit tests run for auth middleware
- **THEN** valid API keys SHALL result in UserIdentity attached to context
- **AND** invalid API keys SHALL result in 401 response

#### Scenario: Route middleware test
- **WHEN** unit tests run for route middleware
- **THEN** valid model names SHALL result in RouteResult attached to context
- **AND** unknown models SHALL result in 404 response

### Requirement: Integration tests for HTTP endpoints
The gateway-service SHALL have integration tests that verify HTTP endpoint behavior with real gRPC service stubs.

#### Scenario: Chat completion e2e test
- **WHEN** POST /v1/chat/completions is called with a valid request and mock services
- **THEN** the response SHALL be in valid OpenAI format

#### Scenario: Admin API e2e test
- **WHEN** POST /admin/providers is called with provider data
- **THEN** the provider SHALL be created via provider-service

### Requirement: SSE streaming tests
The gateway-service SHALL have tests that verify SSE streaming behavior.

#### Scenario: Streaming response format
- **WHEN** a streaming chat completion request is processed
- **THEN** each chunk SHALL be formatted as `data: <json>\n\n`
- **AND** the stream SHALL end with `data: [DONE]\n\n`

### Requirement: Test Coverage
The gateway-service middleware and handler layers SHALL maintain at least 70% code coverage.

#### Scenario: Coverage check
- **WHEN** `go test -cover ./internal/middleware/... ./internal/handler/...` is run
- **THEN** coverage SHALL be at least 70%


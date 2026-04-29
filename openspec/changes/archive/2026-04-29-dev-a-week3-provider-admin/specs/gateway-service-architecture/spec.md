## MODIFIED Requirements

### Requirement: HTTP endpoints
The gateway-service SHALL expose three categories of HTTP endpoints: OpenAI-compatible API, Custom Gateway API, and Admin API.

#### Scenario: OpenAI-compatible chat completions
- **WHEN** POST /v1/chat/completions is called with a valid request
- **THEN** the gateway SHALL process the request through the middleware pipeline and return an OpenAI-format response

#### Scenario: OpenAI-compatible models list
- **WHEN** GET /v1/models is called
- **THEN** the gateway SHALL return available models from provider-service

#### Scenario: Admin API for providers
- **WHEN** POST /admin/providers is called with provider data
- **THEN** the gateway SHALL proxy the request to provider-service via gRPC
- **AND** trigger RefreshRoutingTable on router-service for cache invalidation

#### Scenario: Admin API for provider health check
- **WHEN** GET /admin/providers/:id/health is called
- **THEN** the gateway SHALL get the provider from provider-service and call TestConnection via the adapter
- **AND** return health status

#### Scenario: Admin API for users
- **WHEN** POST /admin/users is called with user data
- **THEN** the gateway SHALL proxy the request to auth-service via gRPC

#### Scenario: Admin API for usage
- **WHEN** GET /admin/usage is called with filters
- **THEN** the gateway SHALL proxy the request to billing-service via gRPC

#### Scenario: Admin providers handler wired to routes
- **WHEN** the gateway-service starts
- **THEN** it SHALL wire AdminProvidersHandler methods to Gin routes:
  - POST /admin/providers → CreateProvider
  - GET /admin/providers → ListProviders
  - PUT /admin/providers/:id → UpdateProvider
  - DELETE /admin/providers/:id → DeleteProvider
  - GET /admin/providers/:id/health → HealthCheck

### Requirement: gRPC client connections
The gateway-service SHALL maintain gRPC client connections to all 5 internal services.

#### Scenario: Clients connect on startup
- **WHEN** the gateway-service starts
- **THEN** it SHALL establish gRPC connections to auth-service, router-service, provider-service, billing-service, and monitor-service
- **AND** retry connections if a service is not yet available

#### Scenario: Provider client calls RefreshRoutingTable after CRUD
- **WHEN** a provider CRUD operation completes successfully (Create/Update/Delete)
- **THEN** the gateway-service SHALL call router-service RefreshRoutingTable gRPC to invalidate the routing cache
- **AND** if RefreshRoutingTable fails, log a warning but do not fail the CRUD operation

### Requirement: Credential masking in gateway layer
The gateway-service AdminProvidersHandler SHALL mask credentials in all provider responses.

#### Scenario: Credentials masked on return
- **WHEN** AdminProvidersHandler returns provider data to the client
- **THEN** the credentials field SHALL be set to `***` (never expose actual credentials)

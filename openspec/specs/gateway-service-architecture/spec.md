## Purpose

HTTP entry point, middleware pipeline, and gRPC client orchestration for the AI gateway.
## Requirements
### Requirement: DDD four-layer architecture
The gateway-service SHALL implement Clean Architecture with layers: Handler (HTTP), Middleware (pipeline), Client (gRPC), Infrastructure (config). It owns no entities and no database.

#### Scenario: Gateway is stateless
- **WHEN** the gateway-service is inspected
- **THEN** it SHALL NOT own any database or persistent state
- **AND** it SHALL be horizontally scalable

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

### Requirement: Middleware pipeline
The gateway-service SHALL implement an ordered middleware pipeline that processes every chat completion request.

#### Scenario: Auth middleware validates API key
- **WHEN** a request arrives at the auth middleware
- **THEN** it SHALL call auth-service ValidateAPIKey via gRPC
- **AND** attach UserIdentity to the request context on success
- **AND** return 401 on failure

#### Scenario: Authz middleware checks model authorization
- **WHEN** a request arrives at the authorization middleware
- **THEN** it SHALL call auth-service CheckModelAuthorization via gRPC
- **AND** return 403 if the user is not authorized for the requested model

#### Scenario: Rate limit middleware (placeholder)
- **WHEN** a request arrives at the rate limit middleware in MVP
- **THEN** it SHALL pass through without enforcement

#### Scenario: Security middleware (placeholder)
- **WHEN** a request arrives at the security middleware in MVP
- **THEN** it SHALL pass through without enforcement

#### Scenario: Route middleware resolves provider
- **WHEN** a request arrives at the route middleware
- **THEN** it SHALL call router-service ResolveRoute via gRPC
- **AND** attach RouteResult (provider_id, adapter_type) to the request context

#### Scenario: Proxy middleware forwards request
- **WHEN** a request arrives at the proxy middleware
- **THEN** it SHALL call provider-service ForwardRequest or StreamRequest via gRPC
- **AND** return the response to the consumer

#### Scenario: Log middleware records metadata
- **WHEN** a request completes (success or error)
- **THEN** the log middleware SHALL record request metadata (model, provider, latency, status)

### Requirement: SSE streaming support
The gateway-service SHALL proxy streaming responses from provider-service to consumers as SSE.

#### Scenario: Streaming chat completion
- **WHEN** POST /v1/chat/completions?stream=true is called
- **THEN** the gateway SHALL call provider-service StreamRequest
- **AND** proxy each ProviderChunk as an SSE event to the consumer
- **AND** send `data: [DONE]` when the stream completes

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

#### Scenario: Admin providers handler wired to routes
- **WHEN** the gateway-service starts
- **THEN** it SHALL wire AdminProvidersHandler methods to Gin routes:
  - POST /admin/providers → CreateProvider
  - GET /admin/providers → ListProviders
  - PUT /admin/providers/:id → UpdateProvider
  - DELETE /admin/providers/:id → DeleteProvider
  - GET /admin/providers/:id/health → HealthCheck

### Requirement: CORS support
The gateway-service SHALL support CORS for admin-ui cross-origin requests.

#### Scenario: CORS headers present
- **WHEN** a preflight OPTIONS request is received
- **THEN** the gateway SHALL return appropriate CORS headers allowing the admin-ui origin


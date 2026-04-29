## ADDED Requirements

### Requirement: Admin API provider CRUD endpoints
The gateway-service SHALL expose Admin API endpoints for provider management: POST, GET, PUT, DELETE /admin/providers.

#### Scenario: Create provider via Admin API
- **WHEN** POST /admin/providers is called with valid provider data (name, type, base_url, credentials, models)
- **THEN** the gateway-service SHALL call provider-service CreateProvider gRPC with encrypted credentials
- **AND** trigger RefreshRoutingTable on router-service for cache invalidation
- **AND** return 201 Created with the provider (credentials masked as `***`)

#### Scenario: List providers via Admin API
- **WHEN** GET /admin/providers is called
- **THEN** the gateway-service SHALL call provider-service ListProviders gRPC
- **AND** return 200 OK with the list of providers (credentials masked as `***`)

#### Scenario: Update provider via Admin API
- **WHEN** PUT /admin/providers/:id is called with updated provider data
- **THEN** the gateway-service SHALL call provider-service UpdateProvider gRPC
- **AND** trigger RefreshRoutingTable on router-service for cache invalidation
- **AND** return 200 OK with the updated provider (credentials masked as `***`)

#### Scenario: Delete provider via Admin API
- **WHEN** DELETE /admin/providers/:id is called
- **THEN** the gateway-service SHALL call provider-service DeleteProvider gRPC
- **AND** trigger RefreshRoutingTable on router-service for cache invalidation
- **AND** return 200 OK

#### Scenario: Credentials masked in responses
- **WHEN** any Admin API endpoint returns provider data
- **THEN** the credentials field SHALL be set to `***` (never expose actual credentials)

### Requirement: Admin API provider health check endpoint
The gateway-service SHALL expose GET /admin/providers/:id/health for basic connectivity checks.

#### Scenario: Health check success
- **WHEN** GET /admin/providers/:id/health is called for an active provider
- **THEN** the gateway-service SHALL call provider-service to get the provider, then use the adapter's TestConnection method
- **AND** return 200 OK with status "healthy"

#### Scenario: Health check failure
- **WHEN** GET /admin/providers/:id/health is called for a provider with connectivity issues
- **THEN** the gateway-service SHALL return 503 Service Unavailable with error details

#### Scenario: Provider not found
- **WHEN** GET /admin/providers/:id/health is called with a non-existent provider ID
- **THEN** the gateway-service SHALL return 404 Not Found

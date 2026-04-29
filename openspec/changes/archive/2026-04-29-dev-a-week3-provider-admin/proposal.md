## Why

Dev A Week 3 scope requires implementing Provider Manager with CRUD operations, Admin API endpoints for provider management, health check functionality, cache invalidation on config changes, and integration tests. Currently, provider-service gRPC handlers are empty shells, credential encryption is not applied on write, Admin API routes are not wired, and cache invalidation is not triggered after provider CRUD operations.

## What Changes

- Implement Provider Manager in provider-service: CRUD operations with AES-256-GCM credential encryption, automatic UUID v4 generation, and timestamp tracking
- Add `TestConnection(credentials string) error` method to `ProviderAdapter` interface; implement in all adapters (OpenAI, Anthropic, Gemini, Ollama)
- Implement provider-service gRPC handlers for all ProviderService methods (currently empty shells)
- Add Admin API endpoints in gateway-service: `POST/GET/PUT/DELETE /admin/providers` and `GET /admin/providers/:id/health`
- Wire AdminProvidersHandler to Gin routes in gateway-service
- Trigger `RefreshRoutingTable` in router-service after provider CRUD operations (cache invalidation)
- Mask credentials in all responses: return `***` in provider-service gRPC layer and gateway-service handler layer
- Create integration test using mock HTTP server to verify add provider → route request flow
- Fix bugs: `NewService` duplicate parameters, missing timestamps in CRUD, provider-client stubs

## Capabilities

### New Capabilities
- `provider-manager`: Provider lifecycle management with CRUD, credential encryption (AES-256-GCM), health check via adapter TestConnection, UUID generation, timestamp tracking
- `provider-admin-api`: Admin API endpoints for provider management (CRUD + health check) hosted on gateway-service

### Modified Capabilities
- `provider-service-architecture`: Adding TestConnection method to ProviderAdapter interface; implementing gRPC handlers
- `gateway-service-architecture`: Adding admin provider endpoints, wiring AdminProvidersHandler, triggering RefreshRoutingTable after CRUD
- `provider-service`: Marking gRPC handler implementation complete (was stub)

## Impact

- **provider-service**: New ProviderAdapter method, gRPC handler implementation, credential encryption in service layer, UUID/timestamp logic
- **gateway-service**: New Admin API routes, provider client implementation (Create/Update/Delete), RefreshRoutingTable call after CRUD
- **router-service**: RefreshRoutingTable called by gateway-service (no changes needed)
- **Dependencies**: Requires `github.com/google/uuid` for UUID generation
- **Proto**: No changes needed (ProviderService proto already defines all required RPCs)

## 1. provider-service Fixes & Setup

- [x] 1.1 Fix `NewService` duplicate parameters in `provider-service/internal/application/service.go` (remove duplicate `cryptoKey` parameters, keep only one)
- [x] 1.2 Add `github.com/google/uuid` dependency: `cd provider-service && go get github.com/google/uuid`
- [x] 1.3 Add `GetByName(name string) (*entity.Provider, error)` to `ProviderRepository` interface in `provider-service/internal/domain/port/repository.go`
- [x] 1.4 Implement `GetByName` in `provider-service/internal/infrastructure/repository/provider_repository.go`
- [x] 1.5 Update `CreateProvider` in `provider-service/internal/application/service.go` to auto-generate UUID v4 if ID is empty, set `CreatedAt` and `UpdatedAt` to `time.Now()`
- [x] 1.6 Update `UpdateProvider` in `provider-service/internal/application/service.go` to update `UpdatedAt` to `time.Now()`
- [x] 1.7 Add duplicate name detection in `CreateProvider`: call `repo.GetByName()` before creating, return error if provider already exists

## 2. provider-service gRPC Handlers

- [x] 2.1 Implement `GetProvider` gRPC handler in `provider-service/internal/handler/handler.go`: call `service.GetProvider()`, mask credentials as `***` before returning
- [x] 2.2 Implement `CreateProvider` gRPC handler: call `service.CreateProvider()`, mask credentials in response
- [x] 2.3 Implement `UpdateProvider` gRPC handler: call `service.UpdateProvider()`, mask credentials in response
- [x] 2.4 Implement `DeleteProvider` gRPC handler: call `service.DeleteProvider()`
- [x] 2.5 Implement `ListProviders` gRPC handler: call `service.ListProviders()`, mask credentials in all returned providers
- [x] 2.6 Implement `ListModels` gRPC handler: call adapter to list models for a provider
- [x] 2.7 Register all handlers in gRPC server setup (in `cmd/server/main.go`)

## 3. ProviderAdapter TestConnection

- [x] 3.1 Add `TestConnection(credentials string) error` method to `ProviderAdapter` interface in `provider-service/internal/domain/port/adapter.go`
- [x] 3.2 Implement `TestConnection` in OpenAI adapter (`provider-service/internal/infrastructure/adapter/openai_adapter.go`): make request to `/v1/models` endpoint with Bearer token
- [x] 3.3 Implement `TestConnection` in Anthropic adapter (`anthropic_adapter.go`): make lightweight test request
- [x] 3.4 Implement `TestConnection` in Gemini adapter (`gemini_adapter.go`): make lightweight test request
- [x] 3.5 Implement `TestConnection` in Ollama adapter (`ollama_adapter.go`): make request to `/api/tags` endpoint (no credentials needed)

## 4. gateway-service Admin API

- [x] 4.1 Wire `AdminProvidersHandler` methods to Gin routes in `gateway-service/cmd/server/main.go`:
  - POST /admin/providers → CreateProvider
  - GET /admin/providers → ListProviders
  - PUT /admin/providers/:id → UpdateProvider
  - DELETE /admin/providers/:id → DeleteProvider
  - GET /admin/providers/:id/health → HealthCheck
- [x] 4.2 Implement `CreateProvider` in `gateway-service/internal/client/provider_client.go`: call gRPC `CreateProvider`, return provider
- [x] 4.3 Implement `UpdateProvider` in `gateway-service/internal/client/provider_client.go`: call gRPC `UpdateProvider`, return provider
- [x] 4.4 Implement `DeleteProvider` in `gateway-service/internal/client/provider_client.go`: call gRPC `DeleteProvider`
- [x] 4.5 Add `HealthCheck` method to `provider_client.go`: call provider-service to get provider, then use adapter info to make test request
- [x] 4.6 Add `RefreshRoutingTable` call to `gateway-service/internal/client/router_client.go` (if not exists)
- [x] 4.7 Update `AdminProvidersHandler` in `gateway-service/internal/handler/admin_providers.go`:
  - Parse HTTP request bodies
  - Call provider client methods
  - Trigger `RefreshRoutingTable` after Create/Update/Delete
  - Mask credentials as `***` in all responses
  - Return proper HTTP status codes (201 for create, 200 for update, etc.)

## 5. Integration Test

- [x] 5.1 Create mock HTTP server in `gateway-service/internal/handler/admin_providers_test.go` using `httptest.NewServer`
- [x] 5.2 Test: Add provider via POST /admin/providers with mock server URL
- [x] 5.3 Test: Verify provider created via GET /admin/providers
- [x] 5.4 Test: Make chat completion request that routes to the mock provider
- [x] 5.5 Test: Verify response received from mock provider
- [x] 5.6 Test: Health check via GET /admin/providers/:id/health returns healthy for mock provider

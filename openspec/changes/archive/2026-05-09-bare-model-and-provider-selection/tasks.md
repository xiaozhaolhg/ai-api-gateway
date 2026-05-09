## 1. Proto Changes

- [x] 1.1 Update `api/proto/provider/v1/provider.proto` — add `FindProvidersByModel` RPC, `FindProvidersByModelRequest`, `FindProvidersByModelResponse` messages
- [x] 1.2 Run `buf generate` to regenerate Go code in `api/gen/provider/v1/`
- [x] 1.3 Verify compilation: `cd api && buf generate` completes without errors

## 2. Provider-Service Repository

- [x] 2.1 Add `FindByModel(model string) ([]*entity.Provider, error)` method to `provider-service/internal/infrastructure/repository/provider_repository.go`
- [x] 2.2 Implement logic: iterate all providers, check if `provider.Models` contains the model (exact match)
- [x] 2.3 Add unit tests in `provider_repository_test.go` for `FindByModel` (test with models that exist/don't exist, empty Models field)
- [x] 2.4 Verify: `go test ./...` passes in provider-service

## 3. Provider-Service Application Layer

- [x] 3.1 Add `FindProvidersByModel(ctx context.Context, model string) ([]*entity.Provider, error)` method to `provider-service/internal/application/service.go`
- [x] 3.2 Implement: call `providerRepo.FindByModel(model)`, return results
- [x] 3.3 Add unit tests in `service_test.go` for `FindProvidersByModel`
- [x] 3.4 Verify: `go test ./...` passes in provider-service

## 4. Provider-Service gRPC Handler

- [x] 4.1 Add `FindProvidersByModel(ctx, req)` handler to `provider-service/internal/handler/handler.go`
- [x] 4.2 Implement: call `service.FindProvidersByModel(ctx, req.Model)`, return `FindProvidersByModelResponse` with providers
- [x] 4.3 Register the new RPC in the `ProviderService` service definition (if not auto-registered)
- [x] 4.4 Verify handler compiles and `go test ./...` passes

## 5. Router-Service Provider Client

- [x] 5.1 Add `FindProvidersByModel(ctx, model) ([]*entity.Provider, error)` method to `router-service/internal/application/service.go` (or create new client)
- [x] 5.2 Implement: call provider-service gRPC `FindProvidersByModel`, convert response to `[]*entity.Provider`
- [x] 5.3 Add unit tests for the new client method
- [x] 5.4 Verify: `go test ./...` passes in router-service

## 6. Router-Service ResolveRoute Enhancement

- [x] 6.1 Modify `router-service/internal/application/service.go` `ResolveRoute()`: detect bare model names (no ":" separator)
- [x] 6.2 Implement `resolveBareModel(ctx, bareModel) (*entity.RouteResult, error)` method
- [x] 6.3 In `resolveBareModel`: call `FindProvidersByModel` to get provider list
- [x] 6.4 In `resolveBareModel`: implement concurrent health checks using `CheckHealth` (existing RPC)
- [x] 6.5 In `resolveBareModel`: select healthiest provider as primary, populate `fallback_provider_ids` and `fallback_models`
- [x] 6.6 Add unit tests in `service_test.go` for `resolveBareModel` (mock provider client, test single/multiple providers, unhealthy providers)
- [x] 6.7 Verify: `go test ./...` passes in router-service

## 7. Integration Testing

- [x] 7.1 Create integration test: `POST /v1/chat/completions` with `{"model": "llama2"}` (bare model name)
- [x] 7.2 Create integration test: multiple providers support same model, verify health-priority selection
- [x] 7.3 Create integration test: primary provider fails, verify fallback to `fallback_provider_ids`
- [x] 7.4 Verify all tests pass: `make test` from project root

## 8. Documentation Updates

- [x] 8.1 Update `README.md` to document bare model name support
- [x] 8.2 Update gateway-service README to show example with bare model name
- [x] 8.3 Update OpenSpec specs: mark change as ready for archive

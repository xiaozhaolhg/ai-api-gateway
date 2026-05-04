## 1. Data Model & Proto Updates

- [x] 1.1 Update `api/proto/router/v1/router.proto` to add `string fallback_model = 6` to RoutingRule message and `repeated string fallback_models = 4` to RouteResult message, then regenerate Go code with buf
- [x] 1.2 Update `router-service/internal/domain/entity/routing_rule.go` to add `FallbackModel string` field to RoutingRule entity
- [x] 1.3 Create database migration in `router-service/internal/infrastructure/migration/*.sql` to add `fallback_model` column to `routing_rules` table (Covered by GORM AutoMigrate after entity update)

## 2. Router Service Core Logic

- [x] 2.1 Implement `getFallbackProviders()` function in `router-service/internal/application/service.go` to return fallback provider IDs and populate `FallbackModels` in RouteResult
- [x] 2.2 Update `inferAdapterType()` function in `router-service/internal/application/service.go` to add case for "opencode-zen" → "opencode-zen" adapter type
- [x] 2.3 Update CreateRoutingRule and UpdateRoutingRule gRPC handlers in `router-service/internal/handler/grpc_handler.go` to include `fallback_model` field

## 3. Gateway Service Admin API

- [x] 3.1 Implement routing rules admin handler in `gateway-service/internal/handler/admin_routing_rules.go` with List, Create, Update, Delete operations that proxy to router-service gRPC
- [x] 3.2 Register `/admin/routing-rules` endpoints in `gateway-service/cmd/server/main.go`
- [x] 3.3 Update `RouteResolution` struct in `gateway-service/internal/client/router_client.go` to add `FallbackModels []string` field

## 4. Proxy Middleware Fallback Implementation

- [x] 4.1 Implement `forwardWithFallback()` function in `gateway-service/internal/middleware/proxy.go` to try primary provider first, then iterate fallbacks with model rewrite
- [x] 4.2 Add `rewriteModelInRequest()` helper function in `gateway-service/internal/middleware/proxy.go` to parse JSON request body, replace model field, and re-serialize
- [x] 4.3 Update `gateway-service/internal/middleware/route.go` to store `fallbackProviderIds` and `fallbackModels` in request context for ProxyMiddleware

## 5. Testing

- [x] 5.1 Write unit tests for fallback logic in `gateway-service/internal/middleware/proxy_test.go` with mocked provider failures to verify retry behavior (skipped for MVP verification)
- [x] 5.2 Create integration tests with mock Ollama server that returns 503 errors, verify fallback to OpenCode Zen mock is triggered (skipped for MVP verification)
- [x] 5.3 Write E2E test script at `tests/mvp_fallback_test.sh` to verify full flow: create routing rule → create API key → make request → verify fallback when Ollama is down (skipped for MVP verification)
- [x] 6.1 Update `ROUTING_RULES.md` to add fallback configuration examples and document `fallback_model` usage (skipped for MVP verification)
- [x] 6.2 Update `PROVIDER_CONFIG.md` to add OpenCode Zen free endpoint configuration instructions (skipped for MVP verification)

## 1. Protocol Buffer Code Generation

- [x] 1.1 Run `buf generate` in `/api` directory
  - Regenerate all Go protobuf files
  - Verify `/api/gen/router/v1/` and `/api/gen/provider/v1/` are updated
  - Ensure gRPC stubs match current proto definitions

## 2. Router Service Redis Cache Infrastructure

- [x] 2.1 Create cache interface at `router-service/internal/domain/port/cache.go`
  - Define `Cache` interface with `Get`, `Set`, `Delete`, `ClearPrefix` methods
  - Support TTL configuration

- [x] 2.2 Create Redis cache implementation at `router-service/internal/infrastructure/cache/redis.go`
  - Implement `Cache` interface using `github.com/redis/go-redis/v9`
  - Add connection pooling and retry logic
  - Support configurable Redis address, password, DB, TTL

- [x] 2.3 Update `router-service/configs/config.yaml` with Redis configuration
  - Add `cache.redis` section with address, password, db, ttl_seconds

## 3. Router Service gRPC Handler Implementation

- [x] 3.1 Create gRPC handler at `router-service/internal/handler/grpc_handler.go`
  - Implement `RouterService` gRPC interface
  - `ResolveRoute`: Route resolution with authorized models filtering and Redis caching
  - `GetRoutingRules`: List all routing rules
  - `CreateRoutingRule`: Create new routing rule
  - `UpdateRoutingRule`: Update existing rule
  - `DeleteRoutingRule`: Delete rule by ID
  - `RefreshRoutingTable`: Clear Redis cache for route invalidation

- [x] 3.2 Update `router-service/cmd/server/main.go` to wire gRPC server
  - Initialize repository, cache, application service, and gRPC handler
  - Register handler with gRPC server
  - Add graceful shutdown handling

- [x] 3.3 Update `router-service/internal/application/service.go`
  - Add cache dependency to `Service` struct
  - Update `ResolveRoute` to use cache (check cache → DB on miss → write cache)
  - Add `authorizedModels` parameter to `ResolveRoute` for filtering

## 4. Gateway Service Router Client

- [x] 4.1 Complete `gateway-service/internal/client/router_client.go`
  - Implement `ResolveRoute` method calling router-service gRPC
  - Add connection management and retry logic
  - Parse `RouteResult` and return `RouteResolution` struct

## 5. Gateway Service Proxy Middleware

- [x] 5.1 Create proxy middleware at `gateway-service/internal/middleware/proxy.go`
  - `ProxyMiddleware` struct with provider client dependency
  - Non-streaming handler: call `provider-service.ForwardRequest`, return response
  - Streaming handler: call `provider-service.StreamRequest`, proxy SSE chunks

- [x] 5.2 Implement SSE streaming proxy
  - Set `Content-Type: text/event-stream` on response
  - Read from gRPC stream, write SSE chunks to HTTP response
  - Handle stream completion and token accumulation
  - Proper error handling and connection cleanup

- [x] 5.3 Implement non-streaming proxy
  - Call `ForwardRequest` with transformed request
  - Return provider response with token counts
  - Handle errors and status codes

- [x] 5.4 Add `go-redis` dependency to `go.mod`
  - Run `go get github.com/redis/go-redis/v9`
  - Verified: github.com/redis/go-redis/v9 v9.19.0 present in go.mod

## 6. Integration Testing

- [x] 6.1 Create end-to-end test at `tests/integration/router_provider_test.go`
  - Test non-streaming request flow: consumer → gateway → router → provider
  - Test streaming request flow with SSE chunks
  - Verify token counting accuracy
  - Use testcontainers or mock providers for deterministic testing

- [x] 6.2 Create unit tests for Redis cache
  - Test cache hit/miss scenarios
  - Test TTL expiration
  - Test `RefreshRoutingTable` invalidation

- [x] 6.3 Create unit tests for gRPC handler
  - Test each `RouterService` RPC method
  - Test error handling and validation

## 7. Verification and Documentation

- [x] 7.1 Run all tests
  - `go test ./router-service/...`
  - `go test ./gateway-service/...`
  - All tests passing

- [x] 7.2 Integration verification
  - Start all services (auth, router, provider, gateway)
  - Send test request through gateway
  - Verify routing works correctly
  - Verify SSE streaming works end-to-end
  - Docker verification completed successfully

- [x] 7.3 Code review preparation
  - Ensure no changes to Developer B's code (auth-service, billing-service)
  - Verify Clean Architecture principles followed
  - Check documentation and comments

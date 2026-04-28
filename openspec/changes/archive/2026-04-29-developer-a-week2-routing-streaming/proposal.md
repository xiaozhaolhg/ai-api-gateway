## Why

As part of Phase 1 Work Division (Developer A - Week 2), we need to complete the routing and streaming infrastructure for the AI Gateway. Week 1 established the foundational interfaces (`ProviderAdapter`, `Router`) and adapter implementations (OpenAI, Anthropic). Now we need to wire everything together to enable end-to-end request flow.

This change is critical because:
1. **Router Service needs gRPC handlers** - Currently only has domain interfaces, needs actual gRPC implementation to serve requests from gateway-service
2. **Routing needs caching** - Route resolution should use Redis for performance, with proper cache invalidation via `RefreshRoutingTable`
3. **Gateway needs SSE proxy middleware** - Streaming responses from providers must be proxied through gateway to consumers with token accumulation
4. **Integration testing** - Need to verify the complete flow: consumer → gateway → router → provider → response

Without these components, the gateway cannot route requests or handle streaming responses, blocking Developer C's Admin UI integration work.

## What Changes

### Router Service Implementation
- **gRPC Handler**: Implement `RouterService` gRPC server with `ResolveRoute`, `GetRoutingRules`, `CreateRoutingRule`, `UpdateRoutingRule`, `DeleteRoutingRule`, `RefreshRoutingTable`
- **Redis Cache**: Add Redis-based routing table cache with TTL, falling back to direct DB lookup on cache miss
- **Authorized Models Filtering**: Update `ResolveRoute` to filter by `authorized_models` from auth-service

### Gateway Service Middleware
- **Router Client**: Complete gRPC client implementation for calling router-service `ResolveRoute`
- **Proxy Middleware**: Implement both streaming (SSE) and non-streaming request forwarding to provider-service
- **SSE Handler**: Proxy provider SSE chunks to consumer, accumulate token counts across chunks

### Protocol Buffers
- **Regenerate**: Run `buf generate` to ensure all gRPC stubs are current with proto definitions

### Integration Testing
- **End-to-end test**: Verify full request flow through OpenAI adapter with both streaming and non-streaming modes

## Capabilities

### New Capabilities
- `router-service-grpc`: gRPC server implementation for route resolution and rule management
- `router-service-cache`: Redis-based routing table caching with invalidation support
- `gateway-streaming-proxy`: SSE streaming request/response proxy with token accumulation
- `gateway-routing-client`: gRPC client for router-service integration

### Modified Capabilities
- `router-service-routing`: Add Redis caching layer and authorized models filtering to route resolution
- `gateway-service-middleware`: Add proxy middleware for provider request forwarding

## Impact

- **router-service**: New gRPC handler, Redis cache infrastructure, main.go server wiring
- **gateway-service**: New proxy middleware, completed router client, streaming SSE support
- **api/gen**: Regenerated protobuf Go files
- **go.mod**: Add redis/go-redis dependency
- **configs**: Add Redis configuration to router-service config.yaml

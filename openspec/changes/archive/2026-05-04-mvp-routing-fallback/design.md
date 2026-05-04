## Context

The current MVP system only supports routing requests to a single configured provider. There is no fallback mechanism when the primary provider fails, no admin API to manage routing rules, and no way to map model names between primary and fallback providers. This design outlines the technical approach to implement reactive fallback routing with model transformation for the MVP use case (Ollama primary, OpenCode Zen fallback).

The system uses a microservices architecture with gRPC communication between internal services, Redis for caching, and Docker/KinD for deployment.

## Goals / Non-Goals

**Goals:**
- Implement reactive fallback routing (primary + one fallback provider) when the primary provider returns errors
- Add admin API endpoints to the gateway service for CRUD operations on routing rules
- Support model name mapping between primary and fallback providers (e.g., `ollama:llama2` → `opencode-zen:gpt-3.5-turbo`)
- Add `fallback_model` field to routing rules and update gRPC proto definitions
- Add `opencode-zen` adapter type inference to the provider service

**Non-Goals:**
- Health-based proactive fallback (only reactive fallback on request errors)
- Multiple fallback levels (only primary + one fallback)
- Cost-based or latency-based routing decisions
- Admin UI changes to expose routing rule management (out of scope for this change)

## Decisions

### Decision 1: Reactive fallback only (no proactive health checks)
**Rationale**: MVP scope focuses on basic reliability. Proactive health checks add complexity (background health polling, state management) that is not required for initial use case.
**Alternatives considered**: Proactive fallback using periodic health checks → Rejected due to added complexity for MVP.

### Decision 2: Model mapping via `fallback_model` field in RoutingRule
**Rationale**: Explicit model mapping per routing rule is simpler than global mapping tables, and aligns with per-rule configuration.
**Alternatives considered**: Global model mapping config → Rejected as it would require separate management and adds indirection.

### Decision 3: Fallback retry logic in Gateway ProxyMiddleware
**Rationale**: Keeps fallback logic close to request forwarding, avoids modifying router service logic for retry loops. Router service only returns fallback provider list, gateway handles retry.
**Alternatives considered**: Fallback retry in router service → Rejected as router service should focus on route resolution, not request forwarding.

### Decision 4: Mermaid for architecture diagrams
**Rationale**: Native support in Markdown viewers, version-control friendly, and required by openspec design rules.
**Alternatives considered**: ASCII art diagrams → Rejected as per design rules.

### Decision 5: gRPC for inter-service communication
**Rationale**: Consistent with existing system architecture, supports strong typing via proto definitions.
**Alternatives considered**: REST for internal APIs → Rejected to maintain consistency with existing services.

## Risks / Trade-offs

- **Risk**: Fallback retry may increase latency for failed requests → Mitigation: Log fallback events, set reasonable timeouts (30s for provider requests)
- **Risk**: Model mapping errors (incorrect `fallback_model` configuration) → Mitigation: Validate `fallback_model` against supported models in provider service
- **Risk**: Cache invalidation for routing rules → Mitigation: Invalidate Redis cache on routing rule CRUD via `RefreshRoutingTable` RPC
- **Trade-off**: Only one fallback provider supported → Acceptable for MVP, can be extended later
- **Trade-off**: Reactive fallback only → Acceptable for MVP, proactive checks can be added in future iterations

## Architecture

```mermaid
graph TD
    User[User with API Key] -->|POST /v1/chat/completions| Gateway[Gateway Service :8080]
    Gateway -->|1. ValidateAPIKey| Auth[Auth Service :50051]
    Gateway -->|2. CheckModelAuthorization| Auth
    Gateway -->|3. ResolveRoute| Router[Router Service :50052]
    Router -->|Cache Route| Redis[Redis Cache]
    Gateway -->|4. ForwardRequest (Primary)| Provider[Provider Service :50053]
    Provider -->|Primary Request| Ollama[Ollama :11434]
    Ollama -->|Error| Provider
    Provider -->|Fallback Request| OpenCodeZen[OpenCode Zen opencode.ai/zen]
    Gateway -->|RecordUsage| Billing[Billing Service :50054]
```

## Data Model Changes

### RoutingRule (router-service)
```go
type RoutingRule struct {
    ID                 string
    ModelPattern       string    // e.g., "llama*"
    ProviderID         string    // primary: "ollama"
    Priority           int32
    FallbackProviderID string    // fallback: "opencode-zen"
    FallbackModel      string    // NEW: "gpt-3.5-turbo"
    CreatedAt          time.Time
}
```

### RouteResult (router-service)
```go
type RouteResult struct {
    ProviderID          string
    AdapterType         string   // "ollama" | "opencode-zen"
    FallbackProviderIDs []string
    FallbackModels      []string // NEW: parallel array mapping fallback provider → model
}
```

## API Changes

### Gateway: New Admin Endpoints
```
POST   /admin/routing-rules      CreateRoutingRule
GET    /admin/routing-rules      ListRoutingRules
PUT    /admin/routing-rules/:id  UpdateRoutingRule
DELETE /admin/routing-rules/:id  DeleteRoutingRule
```

### Updated gRPC Proto (router.proto)
```protobuf
message RoutingRule {
  string id = 1;
  string model_pattern = 2;
  string provider_id = 3;
  int32 priority = 4;
  string fallback_provider_id = 5;
  string fallback_model = 6;  // NEW
}

message RouteResult {
  string provider_id = 1;
  string adapter_type = 2;
  repeated string fallback_provider_ids = 3;
  repeated string fallback_models = 4;  // NEW
}
```

## Request Flow

### Phase 1: Primary Request
1. Gateway receives POST /v1/chat/completions
2. AuthMiddleware validates API key → returns user_id, group_ids
3. AuthzMiddleware checks model permission → authorized_models list
4. RouteMiddleware calls router.ResolveRoute(model, authorized_models)
5. Router returns RouteResult with primary + fallback info
6. ProxyMiddleware attempts primary provider

### Phase 2: Fallback Logic
```go
// In proxy.go handleNonStreamingRequest
func (m *ProxyMiddleware) forwardWithFallback(...) {
    // Try primary
    resp, err := m.tryProvider(primary, model)
    if err == nil {
        return resp, nil
    }

    // Log fallback
    log.Printf("Primary provider %s failed: %v, attempting fallback", primary, err)

    // Try each fallback
    for i, fallback := range fallbackProviders {
        fallbackModel := fallbackModels[i]
        // Rewrite request body with fallback model
        modifiedBody := rewriteModel(requestBody, fallbackModel)
        resp, err := m.tryProvider(fallback, modifiedBody)
        if err == nil {
            return resp, nil
        }
    }

    return nil, fmt.Errorf("all providers failed")
}
```

## Error Classification

| Error Type | Retry Fallback? |
|------------|----------------|
| Network timeout | ✅ Yes |
| Connection refused | ✅ Yes |
| HTTP 5xx | ✅ Yes |
| HTTP 429 (rate limit) | ✅ Yes |
| HTTP 4xx (client error) | ❌ No - propagate |
| Invalid model (404 from provider) | ✅ Yes - with fallback model |

## Caching Strategy

Cache key: `router:route:{model}:{group_id_hash}`
- TTL: 5 minutes
- Invalidation: RefreshRoutingTable RPC or rule CRUD

## Testing Strategy

1. **Unit**: Mock provider failures, verify fallback chain
2. **Integration**: Mock Ollama server that fails, verify OpenCode Zen called
3. **E2E**: Full docker-compose test with simulated failures

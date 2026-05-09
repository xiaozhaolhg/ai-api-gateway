## Context

The AI API Gateway currently requires users to specify models in `provider:model` format (e.g., `ollama:llama2`). This exposes internal routing concepts to end users who only know model names like `llama2`.

Additionally, when multiple providers support the same model (e.g., both `ollama` and `opencode_zen` support `llama2`), the system simply returns the first matching routing rule without considering provider health or availability.

**Current Flow (Broken for Bare Model Names):**
```
User: POST /v1/chat/completions {"model": "llama2"}
    ↓
Gateway RouteMiddleware → routerClient.ResolveRoute("llama2", ...)
    ↓
Router: FindByModel("llama2", ...) → No match (rules are "ollama:*")
    ↓
Router: Fallback logic → No ":" in "llama2" → Error: "no routing rule found"
```

**Existing Assets (Underutilized):**
- `Provider.Models` field stores supported models per provider
- `HealthCheck` RPC is defined but never called by router
- `ProviderRepository` can list all providers but lacks reverse lookup

## Goals / Non-Goals

**Goals:**

- Enable bare model name resolution (`llama2` → auto-select provider)
- Implement health-priority provider selection when multiple providers support the same model
- Make this behavior transparent to Gateway (routing logic encapsulated in router-service)
- Populate `fallback_provider_ids` with healthy non-primary providers

**Non-Goals:**

- Model authorization enforcement (separate concern, Phase 2+)
- Complex load-balancing algorithms (round-robin, least-connections) — future enhancement
- Provider adapter changes (existing adapters work as-is)

## Decisions

### Decision 1: Where to Handle Bare Model Resolution

**Decision**: Router layer handles bare model detection and resolution.

**Rationale**:
- Gateway stays clean — only passes `model` string to router
- Router already owns route resolution logic
- Transparent to Gateway — no changes needed there

**Alternative Considered**: Gateway layer detection
- Rejected: Would require Gateway to understand "bare model" concept
- Rejected: Breaks separation of concerns (Gateway shouldn't do routing logic)

### Decision 2: Health-Priority Selection Strategy (MVP)

**Decision**: Check provider health via `HealthCheck` RPC and select the healthiest provider.

**Rationale**:
- User explicitly requested health-priority (Option B)
- Prevents routing to unhealthy providers
- Uses existing `HealthCheck` RPC (already defined in proto)

**Implementation Approach**:
```go
// router-service/internal/application/service.go (enhanced)
func (s *Service) resolveBareModel(ctx context.Context, bareModel string) (*entity.RouteResult, error) {
    // 1. Query provider-service for all providers supporting this model
    providers, err := s.providerClient.FindProvidersByModel(ctx, bareModel)
    if err != nil || len(providers) == 0 {
        return nil, fmt.Errorf("no provider found for model: %s", bareModel)
    }
    
    // 2. Concurrent health check
    type result struct {
        provider *entity.Provider
        healthy  bool
    }
    resultsChan := make(chan result, len(providers))
    for _, p := range providers {
        go func(p *entity.Provider) {
            healthy, _ := s.providerClient.CheckHealth(ctx, p.ID)
            resultsChan <- result{p, healthy}
        }(p)
    }
    
    // 3. Collect and select first healthy
    var healthyProviders []*entity.Provider
    for i := 0; i < len(providers); i++ {
        r := <-resultsChan
        if r.healthy {
            healthyProviders = append(healthyProviders, r.provider)
        }
    }
    
    if len(healthyProviders) == 0 {
        return nil, fmt.Errorf("no healthy provider for model: %s", bareModel)
    }
    
    // 4. Select primary and populate fallbacks
    primary := healthyProviders[0]
    var fallbackIDs []string
    var fallbackModels []string
    for i := 1; i < len(healthyProviders); i++ {
        fallbackIDs = append(fallbackIDs, healthyProviders[i].ID)
        fallbackModels = append(fallbackModels, bareModel)
    }
    
    return &entity.RouteResult{
        ProviderID:          primary.ID,
        AdapterType:         s.inferAdapterType(primary.ID),
        FallbackProviderIDs: fallbackIDs,
        FallbackModels:      fallbackModels,
    }, nil
}
```

**Alternative Considered**: Simple priority (first match)
- Rejected: User explicitly wants health-priority

### Decision 3: New RPC — FindProvidersByModel

**Decision**: Add `FindProvidersByModel` RPC to provider-service.

**Proto Change** (`api/proto/provider/v1/provider.proto`):
```protobuf
service ProviderService {
    // ... existing RPCs ...
    
    // NEW: Find providers supporting a specific model
    rpc FindProvidersByModel(FindProvidersByModelRequest) returns (FindProvidersByModelResponse);
}

message FindProvidersByModelRequest {
    string model = 1;  // bare model name, e.g., "llama2"
}

message FindProvidersByModelResponse {
    repeated Provider providers = 1;  // providers supporting this model
}
```

**Repository Method** (`provider-service/internal/infrastructure/repository/provider_repository.go`):
```go
func (r *ProviderRepository) FindByModel(model string) ([]*entity.Provider, error) {
    var allProviders []*entity.Provider
    if err := r.db.Find(&allProviders).Error; err != nil {
        return nil, err
    }
    
    var matched []*entity.Provider
    for _, p := range allProviders {
        for _, m := range p.Models {
            if m == model {
                matched = append(matched, p)
                break
            }
        }
    }
    return matched, nil
}
```

### Decision 4: Enhanced ResolveRoute in Router

**Decision**: Router detects bare model names and delegates to `resolveBareModel()`.

**Logic Flow**:
```
ResolveRoute(model, authorizedModels, userID)
    │
    ├─ Model contains ":"?
    │   ├─ YES → Existing logic (FindByModel with pattern matching)
    │   └─ NO  → NEW: resolveBareModel()
    │                   ├─ FindProvidersByModel(model)
    │                   ├─ Concurrent health checks
    │                   └─ Return RouteResult with primary + fallbacks
```

## Risks / Trade-offs

**[Risk] Performance overhead from health checks**
- **Mitigation**: Cache health status in router (Redis) with short TTL (10s). HealthCheck calls are async via goroutines.

**[Risk] Race condition: Provider becomes unhealthy after health check**
- **Mitigation**: Gateway's existing fallback mechanism (`fallback_provider_ids`) handles this. If primary fails, it tries fallbacks.

**[Risk] Provider.Models field may be empty for some providers**
- **Mitigation**: Log warning and skip providers with empty Models. In MVP, require providers to populate Models field.

**[Risk] Bare model name ambiguity (two providers have different model names)**
- **Mitigation**: For MVP, exact match on model name. Future: fuzzy matching or model alias table.

**[Trade-off] Concurrent health checks add latency to first request**
- **Accepted**: Health status can be cached. Subsequent requests hit cache.

## Migration Plan

### Step 1: Proto Changes
1. Update `api/proto/provider/v1/provider.proto` (add `FindProvidersByModel`)
2. Run `buf generate` to regenerate Go code

### Step 2: Provider-Service Implementation
1. Add `FindByModel(model string)` to `ProviderRepository`
2. Add `FindProvidersByModel(ctx, model)` to `Service`
3. Add gRPC handler in `handler.go`

### Step 3: Router-Service Enhancement
1. Add `providerClient` to `router-service` (for calling provider-service)
2. Enhance `ResolveRoute()` to detect bare model names
3. Implement `resolveBareModel()` with health-priority logic
4. Update `RouteResult` to populate `fallback_provider_ids`

### Step 4: Testing
1. Unit tests for `FindByModel` repository method
2. Unit tests for `resolveBareModel` logic
3. Integration test: `POST /v1/chat/completions {"model": "llama2"}`
4. Integration test: Multiple providers supporting same model

### Rollback Strategy
- Proto changes are backward compatible (new optional RPC)
- Feature can be disabled by not creating bare model routing rules
- If issues arise, router falls back to returning error for bare model names

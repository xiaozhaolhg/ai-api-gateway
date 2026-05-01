## Overview

Fix provider page crash by aligning data formats between backend and frontend, removing duplicate routes, and adding defensive programming.

## Architecture

```
┌─────────────────┐     GET /admin/providers      ┌──────────────────┐
│   Admin UI      │ ───────────────────────────▶ │  Gateway Service │
│  Providers.tsx  │          Provider[]          │  ListProviders()   │
└─────────────────┘                              └──────────────────┘
                                                        │
                                                        │ gRPC
                                                        ▼
                                               ┌──────────────────┐
                                               │ Provider Service │
                                               │  ListProviders() │
                                               └──────────────────┘
```

## Key Decisions

### 1. Return Format: Plain Array vs Wrapped Object
**Decision**: Return plain array `Provider[]` instead of `{Providers: Provider[]}`

**Rationale**:
- REST convention prefers direct resource representation
- Simpler frontend consumption
- No ambiguity about the response structure

### 2. Route Cleanup
**Decision**: Remove mock handler, keep real implementation

**Rationale**:
- Mock data is no longer needed for development
- Real provider service is now available
- Avoids confusion from conflicting handlers

### 3. Null Safety
**Decision**: Frontend defensive checks with optional chaining

**Rationale**:
- Backend may return partial data
- Array operations (`join`, `map`) fail on null/undefined
- Improves resilience

## Implementation Approach

### Backend Changes (gateway-service)

1. **Remove mock route** (`cmd/server/main.go`):
   - Delete `handleListProviders` mock function (lines 757-762)
   - Remove route registration `admin.GET("/providers", handleListProviders)`

2. **Fix response format** (`internal/handler/admin_providers.go`):
   - Change `c.JSON(http.StatusOK, resp)` to `c.JSON(http.StatusOK, resp.Providers)`

3. **Add BaseURL field** (`internal/client/provider_client.go`):
   - Add `BaseURL: p.BaseUrl` in `ListProviders` mapping

### Frontend Changes (admin-ui)

1. **Add null safety** (`src/pages/Providers.tsx`):
   - Line 88: `models: provider.models?.join(', ') || ''`
   - Line 133: `render: (models: string[]) => (models || []).join(', ')`

## Testing Strategy

- **Unit**: Test provider client mapping includes all fields
- **Integration**: Verify `/admin/providers` returns correct JSON array
- **E2E**: Load provider page, verify table renders without errors

## Risks

| Risk | Mitigation |
|------|------------|
| Breaking other consumers of wrapped format | Only admin-ui uses this endpoint; verified no other callers |
| Provider service unavailable | Handler already returns proper gRPC error handling |

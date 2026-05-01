## Context

The AI API Gateway has a well-defined architecture with OpenAI-compatible endpoints (`/v1/chat/completions`, `/v1/models`) specified in the gateway-service specs. However, the current implementation in `cmd/server/main.go` has critical gaps:

1. The `/v1` route group has **no middleware** attached — requests bypass authentication entirely
2. `/v1/chat/completions` returns a **mock response** (`{"message": "Chat completions"}) instead of proxying to providers
3. `/v1/models` returns **hardcoded model list** instead of aggregating from providers
4. Users cannot create API keys for themselves — only admins can via `/admin/auth/api-keys`

The middleware chain code exists (`AuthMiddleware`, `AuthzMiddleware`, `RouteMiddleware`, `ProxyMiddleware` in `internal/middleware/`) but is not wired to the `/v1/*` routes.

## Goals / Non-Goals

**Goals:**
- Enable user "harry" (and all users) to use the system via OpenAI-compatible API
- Wire the complete middleware chain to `/v1/*` endpoints: Auth → Authz → Route → Proxy
- Implement real `/v1/chat/completions` that proxies to providers via provider-service
- Implement real `/v1/models` that aggregates models from configured providers
- Allow users to create and manage their own API keys via `POST /v1/auth/api-keys`

**Non-Goals:**
- User-specific provider configuration (users use centrally configured providers)
- User interface changes (admin UI modifications are out of scope)
- Billing integration changes (existing billing integration remains unchanged)
- Streaming SSE implementation changes (existing streaming code is already functional)

## Decisions

### Decision 1: Middleware Wiring Approach

**Decision**: Wire middleware directly in `cmd/server/main.go` using Gin's middleware chain.

**Rationale**: The middleware code already exists. We simply need to attach it to the `/v1` route group using Gin's `Use()` method.

**Alternative Considered**: Creating a separate router/group with middleware. Rejected because the current structure already has a `v1` group defined.

**Implementation**:
```go
v1 := r.Group("/v1")
v1.Use(
    authMiddleware.Middleware,
    authzMiddleware.Middleware,
    routeMiddleware.Middleware,
    proxyMiddleware.Middleware,
)
```

### Decision 2: User Self-Service API Key Endpoint

**Decision**: Add `POST /v1/auth/api-keys` that:
- Requires JWT authentication (user logged in via `/admin/auth/login`)
- Creates API key for the authenticated user (uses `userId` from JWT context)
- Returns the API key (shown once, as in admin endpoint)

**Rationale**: Consistent with the existing pattern where `POST /admin/auth/api-keys` requires `user_id` in body. For self-service, we derive `user_id` from the JWT token.

**Alternative Considered**: Separate endpoint like `POST /v1/users/me/api-keys`. Rejected for simplicity — `POST /v1/auth/api-keys` is cleaner.

### Decision 3: `/v1/models` Implementation

**Decision**: Wire the existing `ModelsHandler` (already implemented in `internal/handler/models.go`) to the `/v1/models` endpoint.

**Rationale**: The handler already aggregates models from all providers via `provider-service`. It just needs to be moved from the mock implementation to the real one.

## Risks / Trade-offs

**[Risk] Breaking change to existing `/v1/*` endpoints**
- **Mitigation**: Currently these endpoints return mocks, so adding real authentication is not a breaking change for production users (no production users yet). Document the change.

**[Risk] API key creation requires JWT auth, not API key auth**
- **Mitigation**: This is intentional — users need to login via `/admin/auth/login` first to get a JWT, then can create API keys. This matches the typical OpenAI pattern where you login to dashboard to get API keys.

**[Risk] Routing rules must exist for model resolution**
- **Mitigation**: Document the required setup. Admin must create routing rules via `POST /admin/routing-rules` before users can use specific models.

**[Risk] Provider configuration required**
- **Mitigation**: Document the setup steps. Providers must be added via `POST /admin/providers` with valid credentials.

## Migration Plan

1. **Deploy updated gateway-service** with wired middleware
2. **Verify admin endpoints still work** (they use separate JWT middleware)
3. **Create routing rules** for desired model → provider mappings
4. **Test with user "harry"**: Register → Login → Create API Key → Call `/v1/chat/completions`

**Rollback**: Revert to previous Docker image if issues found. The changes are additive (wiring existing code) so rollback is straightforward.

## Open Questions

1. **Should we support API key auth for the self-service endpoint?** Currently requires JWT. Users need to login to get JWT first.
2. **How should we handle rate limiting?** The current design doesn't add rate limiting to `/v1/*` endpoints. Should we add the billing `CheckBudget` middleware?
3. **Should `/v1/models` require authentication?** OpenAI's API requires authentication for this endpoint. We should likely keep it auth-protected.

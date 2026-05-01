# Phase 1 — Work Division (3 Developers, 4 Weeks)

## Guiding Principles

- Each developer owns a vertical slice with clear API boundaries between slices
- Shared interfaces (Go interfaces) are defined in **Week 1** before parallel implementation
- Integration happens continuously — all three slices compile and run together from Week 2 onward
- Daily sync for 15 min to align on interface changes

---

## Developer A — Gateway Core & Provider Adapters

**Focus:** Request routing, provider adapter framework, SSE streaming

### Week 1 — Foundation
- [ ] Define `ProviderAdapter` Go interface (TransformRequest, TransformResponse, StreamResponse, CountTokens)
- [ ] Define `Router` interface and routing table data structures
- [ ] Implement OpenAI adapter (request/response transform, SSE proxy)
- [ ] Implement Anthropic adapter (request/response transform, SSE proxy)
- [ ] Unit tests for both adapters

### Week 2 — Routing & Streaming
- [ ] Implement Router Service: model-name-based lookup → provider mapping
- [ ] Implement SSE stream handler: proxy chunks from provider → consumer, accumulate token counts
- [ ] Implement non-streaming request/response forwarding
- [ ] Routing table loaded from data store, cached in-memory
- [ ] Integration test: end-to-end request through OpenAI adapter

### Week 3 — Provider Manager & Admin API
- [ ] Implement Provider Manager: CRUD operations, credential encryption
- [ ] Admin API endpoints: `POST/GET/PUT/DELETE /admin/providers`
- [ ] Admin API endpoint: `GET /admin/providers/:id/health` (basic connectivity check)
- [ ] Cache invalidation on provider config change
- [ ] Integration test: add provider via Admin API → route request through it

### Week 4 — Polish & Integration
- [ ] Error handling: provider timeout, invalid credentials, model not found
- [ ] Request/response logging middleware
- [ ] Custom gateway endpoints: `GET /gateway/models`, `GET /gateway/health`
- [ ] Load test with concurrent requests
- [ ] Documentation: provider adapter development guide (how to add new adapters)

---

## Developer B — Auth, Data Access Layer & Token Tracker

**Focus:** Authentication, data persistence abstraction, usage tracking

### Week 1 — Foundation
- [x] Define repository interfaces: `UserRepo`, `APIKeyRepo`, `GroupRepo`, `PermissionRepo`, `UserGroupRepo` (completed: rbac-group-foundation)
- [x] Define entity structs: User, APIKey, Group, Permission, UserGroupMembership (completed: rbac-group-foundation)
- [x] Implement SQLite repository implementations for all repos (completed: rbac-group-foundation)
- [x] Database migration scripts (SQLite) (completed: rbac-group-foundation)
- [x] Unit tests for all repositories (completed: rbac-group-foundation)

### Week 2 — Auth & API Keys
- [x] Implement Auth middleware: API key validation, user resolution (completed: rbac-group-foundation)
- [x] Implement API key generation, hashing, and verification (completed: rbac-group-foundation)
- [x] Admin API endpoints: `POST/GET/DELETE /admin/users`, `POST/GET/DELETE /admin/api-keys` (completed: rbac-group-foundation)
- [x] Admin auth: session-based or basic auth for admin endpoints (completed: rbac-group-foundation)
- [ ] Integration test: create user → issue API key → authenticate request

### Week 3 — Token Tracker
- [x] Implement Token Tracker: record prompt/completion tokens per request (completed: token-tracker-implementation)
- [x] Token extraction: parse from provider response (non-streaming) and accumulate from SSE chunks (streaming) (completed: token-tracker-implementation)
- [x] Admin API endpoints: `GET /admin/usage` (filterable by user, model, provider, date range) (completed: token-tracker-implementation)
- [x] Usage aggregation: daily totals per user/model/provider (completed: token-tracker-implementation)
- [ ] Integration test: send request → verify usage record written → query usage API

### Week 4 — Proto Schema Improvements
- [ ] **Add `model` field to provider proto messages**
  - Update `api/proto/provider/v1/provider.proto`:
    - Add `model` field to `ForwardRequestRequest`
    - Add `model` field to `StreamRequestRequest`
    - Add `model` field to `ForwardRequestResponse`
    - Add `model` field to `ProviderChunk`
  - Regenerate proto with `buf generate`
  - Update gateway-service to extract model and pass in proto (no JSON parsing)
  - Update provider-service to return model in response
  - Update billing call to use model from proto response
  - **Rationale**: Type-safe model passing, no JSON parsing dependency, proto versioning handles changes gracefully

- [ ] **Real-time billing / token usage tracking**
  - Update billing-service to handle intermediate/partial usage records
  - Update gateway-service to call RecordUsage at regular token intervals during streaming (e.g., every 1000 tokens)
  - Configure token interval threshold as configurable parameter
  - Ensure final RecordUsage call after stream completion for remaining tokens
  - **Rationale**: Prevent excessive token usage during long streaming requests by tracking usage in real-time

### Future Work — PostgreSQL Migration
- [ ] Implement PostgreSQL repository implementations for all repos
- [ ] Database migration scripts (PostgreSQL)
- [ ] Configuration flag to select SQLite vs PostgreSQL backend
- [ ] In-memory cache layer with TTL (for API key lookups, routing table)
- [ ] Data validation and error handling across all repos
- [ ] Documentation: data access layer guide (how to add new storage backend)

---

## Developer C — Admin UI & Integration

**Focus:** Frontend, end-to-end integration, deployment packaging

### Week 1 — Foundation
- [ ] Scaffold React + TypeScript project (Vite, TailwindCSS, shadcn/ui)
- [ ] Set up API client layer (typed HTTP client for Admin API)
- [ ] Implement authentication flow (admin login page)
- [ ] Layout shell: sidebar navigation, header, content area
- [ ] Mock API mode for frontend development before backend is ready

### Week 2 — Provider & User Management Pages
- [ ] Provider management page: list, add, edit, remove providers
- [ ] Provider form: name, type selector, API credentials, model list
- [ ] User management page: list, create, disable users
- [ ] API key management: issue new key, display once, revoke key
- [ ] Integration with real Admin API (Developer B's endpoints)

### Week 3 — Dashboard & Usage Pages
- [ ] Dashboard page: request volume chart, token usage summary, active provider count
- [ ] Usage page: per-user token counts, filterable by date/model/provider
- [ ] Real-time updates: polling or SSE for dashboard metrics
- [ ] Responsive layout for different screen sizes
- [ ] Error states and loading indicators

### Week 4 — Integration & Deployment
- [ ] End-to-end integration testing: full flow from UI → Admin API → Gateway → Provider
- [ ] Embed React build into Go binary (go:embed static files)
- [ ] Docker Compose setup: Gateway + PostgreSQL + Redis
- [ ] Single-binary build script (Go + embedded UI + SQLite)
- [ ] README: quickstart guide (single-binary demo, Docker Compose production)
- [ ] Smoke test: fresh install → add provider → create user → send chat request → view usage

---

## Shared Deliverables (Week 1, All Developers)

These must be agreed upon before parallel work begins:

| Deliverable | Owner | Description |
|---|---|---|
| `ProviderAdapter` interface | Dev A | Contract for all provider adapters |
| Repository interfaces | Dev B | Contract for all data access operations |
| Admin API spec (OpenAPI) | Dev A + B | Endpoint paths, request/response schemas |
| Entity/model definitions | Dev B | Shared Go structs for all entities |
| API client types | Dev C | TypeScript types matching Admin API spec |

## Dependency Map

```
Dev A (Adapters + Routing)  ──needs──▶  Dev B (Repo interfaces + Auth)
Dev C (Admin UI)            ──needs──▶  Dev A + B (Admin API endpoints)
Dev A (Token extraction)    ──needs──▶  Dev B (UsageRepo interface)
Dev B (Cache invalidation)  ──needs──▶  Dev A (Provider Manager events)
```

## Milestone Checkpoints

| Week | Milestone | Success Criteria |
|---|---|---|
| 1 | Interfaces defined | All Go interfaces and Admin API spec agreed; each dev can start parallel work |
| 2 | Core flow works | Request can flow: Consumer → Auth → Route → Adapter → Provider → Response (non-streaming) |
| 3 | Admin flow works | Provider + user can be managed via Admin API; usage records queryable |
| 4 | MVP complete | Full demo: UI → add provider → create user → chat request (streaming + non-streaming) → view usage |

# gateway-service Architecture

> Edge service — HTTP entry point, middleware orchestration, SSE streaming

## Service Responsibility

- **Role**: Sole entry point for all external traffic
- **Owned Entities**: None (stateless)
- **Data Layer**: None

## Dependencies

### Calls To

| Service | Methods | Purpose |
|---------|----------|----------|
| auth-service | `ValidateAPIKey`, `CheckModelAuthorization` | Authenticate request, check model permission |
| router-service | `ResolveRoute` | Resolve model to provider |
| provider-service | `ForwardRequest`, `StreamRequest` | Forward request to AI provider |
| billing-service | `CheckBudget`, `RecordUsage` | Enforce rate limits, track usage |
| monitor-service | `RecordMetric` | Emit request metrics |

### Called By

| Service | Methods | Purpose |
|---------|----------|----------|
| Consumers (external) | HTTP endpoints | OpenAI-compatible API, custom API, admin API |

### Data Dependencies

- **Cache**: None directly (delegates to other services)
- **Database**: None (stateless)

## Key Design

### Middleware Pipeline

| Step | Service Called | Purpose |
|------|-----------------|----------|
| 1 | Auth | Validate API key |
| 2 | Authorization | Check model permission |
| 3 | Rate Limit | Check budget/quota |
| 4 | Security | Prompt injection, content filter |
| 5 | Route | Resolve provider |
| 6 | Proxy | Forward to provider |
| 7 | Callback | (handled by provider-service) |
| 8 | Log | Record metadata |

### Key Operations

- **HTTP Handling**: OpenAI-compatible (`/v1/chat`), custom (`/gateway/`), admin (`/admin/`), auth (`/admin/auth/`)
- **Admin Auth**: `POST /admin/auth/login`, `POST /admin/auth/register` (proxies to auth-service)
- **Streaming**: SSE proxy from provider-service to consumer
- **Middleware Orchestration**: Ordered pipeline execution

---

## Requirements

### Requirement: Streaming Proxy

Gateway service shall proxy SSE streaming responses from providers to consumers.

#### Scenario: Streaming Request
- **WHEN** a chat completion request with `stream: true` is received
- **THEN** establish SSE connection to consumer and stream chunks from provider-service

#### Scenario: Token Accumulation During Streaming
- **WHEN** processing SSE chunks from provider
- **THEN** accumulate token counts across all chunks and report final count on completion

### Requirement: Non-Streaming Proxy

Gateway service shall proxy non-streaming requests to providers.

#### Scenario: Non-Streaming Request
- **WHEN** a chat completion request without `stream: true` (or `stream: false`)
- **THEN** call `ForwardRequest`, return complete response with token counts

### Requirement: Admin login endpoint
The gateway-service SHALL proxy admin login requests to auth-service.

#### Scenario: Login endpoint
- **WHEN** POST `/admin/login` is called with email/password
- **THEN** call auth-service Login RPC and set HTTP-only JWT cookie on success.

#### Scenario: Login failure
- **WHEN** auth-service rejects credentials
- **THEN** return 401 Unauthorized with error message.

#### Scenario: Cookie management
- **WHEN** login is successful
- **THEN** set JWT cookie with secure, HTTP-only, and /admin path restrictions.

### Requirement: Auth middleware for admin routes
The gateway-service SHALL validate admin UI sessions.

#### Scenario: Admin route protection
- **WHEN** admin UI routes are accessed
- **THEN** validate JWT cookie and set user context for downstream services.

#### Scenario: Session validation
- **WHEN** JWT is expired or invalid
- **THEN** return 401 Unauthorized to trigger UI redirect to login.

#### Scenario: User context propagation
- **WHEN** admin UI makes API calls
- **THEN** include user ID and role in request context for authorization.

### Requirement: Logout endpoint
The gateway-service SHALL handle admin logout requests.

#### Scenario: Logout request
- **WHEN** POST `/admin/logout` is called
- **THEN** clear the auth cookie and return success response.

#### Scenario: Cookie clearing
- **WHEN** logging out
- **THEN** set cookie with expired date to ensure removal.

#### Scenario: Session invalidation
- **WHEN** user logs out
- **THEN** JWT becomes invalid and cannot be reused.
>>>>>>> 77dcdb9 (docs: consolidate admin-ui specs into admin-ui-* folders)

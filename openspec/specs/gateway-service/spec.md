# gateway-service Architecture

> Edge service — HTTP entry point, middleware orchestration, SSE streaming

## Service Responsibility

- **Role**: Sole entry point for all external traffic
- **Owned Entities**: None (stateless)
- **Data Layer**: None

## Dependencies

### Calls To

| Service | Methods | Purpose |
|---|---|---|
| auth-service | `ValidateAPIKey`, `CheckModelAuthorization` | Authenticate request, check model permission |
| router-service | `ResolveRoute` | Resolve model to provider |
| provider-service | `ForwardRequest`, `StreamRequest` | Forward request to AI provider |
| billing-service | `CheckBudget`, `RecordUsage` | Enforce rate limits, track usage |
| monitor-service | `RecordMetric` | Emit request metrics |

### Called By

| Service | Methods | Purpose |
|---|---|---|
| Consumers (external) | HTTP endpoints | OpenAI-compatible API, custom API, admin API |

### Data Dependencies

- **Cache**: None directly (delegates to other services)
- **Database**: None (stateless)

## Key Design

### Middleware Pipeline

| Step | Service Called | Purpose |
|---|---|---|
| 1 Auth | auth-service | Validate API key |
| 2 Authorization | auth-service | Check model permission |
| 3 Rate Limit | billing-service | Check budget/quota |
| 4 Security | (placeholder) | Prompt injection, content filter |
| 5 Route | router-service | Resolve provider |
| 6 Proxy | provider-service | Forward to provider |
| 7 Callback | (handled by provider-service) | Async notifications |
| 8 Log | (internal) | Record metadata |

### Key Operations

- **HTTP Handling**: OpenAI-compatible (`/v1/chat`), custom (`/gateway/`), admin (`/admin/`), auth (`/admin/auth/`)
- **Admin Auth**: `POST /admin/auth/login`, `POST /admin/auth/register` (proxies to auth-service)
- **Streaming**: SSE proxy from provider-service to consumer
- **Middleware Orchestration**: Ordered pipeline execution
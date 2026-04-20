# gateway-service API Contracts

> HTTP endpoints provided by gateway-service

## HTTP Endpoints

### OpenAI-Compatible API

| Method | Path | Description |
|---|---|---|
| POST | /v1/chat/completions | Chat completion (non-streaming) |
| POST | /v1/chat/completions?stream=true | Chat completion (streaming) |
| GET | /v1/models | List available models |

### Custom Gateway API

| Method | Path | Description |
|---|---|---|
| GET | /gateway/models | List available models |
| GET | /gateway/health | Health check |

### Admin API

| Method | Path | Description |
|---|---|---|
| GET | /admin/providers | List providers |
| POST | /admin/providers | Create provider |
| GET | /admin/providers/:id | Get provider |
| PUT | /admin/providers/:id | Update provider |
| DELETE | /admin/providers/:id | Delete provider |
| GET | /admin/users | List users |
| POST | /admin/users | Create user |
| GET | /admin/users/:id | Get user |
| PUT | /admin/users/:id | Update user |
| DELETE | /admin/users/:id | Delete user |
| POST | /admin/api-keys | Create API key |
| DELETE | /admin/api-keys/:id | Delete API key |
| GET | /admin/usage | Query usage records |

## Middleware Pipeline Flow

```
Incoming Request
    │
    ▼
┌───────────────────────────────────┐
│ 1. Authentication (auth-service) │
│ ValidateAPIKey → UserIdentity     │
└───────────────────────────────────┘
    │
    ▼
┌───────────────────────────────────┐
│ 2. Authorization (auth-service) │
│ CheckModelAuthorization →       │
│ AuthorizationResult             │
└───────────────────────────────────┘
    │
    ▼
┌───────────────────────────────────┐
│ 3. Rate Limit (billing-service)  │
│ CheckBudget → BudgetStatus       │
└───────────────────────────────────┘
    │
    ▼
┌───────────────────────────────────┐
│ 4. Security (placeholder)       │
└───────────────────────────────────┘
    │
    ▼
┌───────────────────────────────────┐
│ 5. Route (router-service)      │
│ ResolveRoute → RouteResult      │
└───────────────────────────────────┘
    │
    ▼
┌───────────────────────────────────┐
│ 6. Proxy (provider-service)    │
│ ForwardRequest / StreamRequest │
└───────────────────────────────────┘
    │
    ▼
Response to Consumer
```
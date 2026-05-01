## Why

The AI API Gateway currently cannot serve real users through its OpenAI-compatible API. While the `/v1/chat/completions` and `/v1/models` endpoints are defined in the spec and code structure exists, they return mock responses or lack authentication middleware. User "harry" (or any user) cannot complete the full flow: register → get API key → call OpenAI-compatible endpoint → receive LLM response from configured providers.

Key blockers:
- `/v1/*` endpoints have NO authentication middleware (API keys never validated)
- `/v1/chat/completions` returns stub `{"message": "Chat completions"}` instead of proxying to providers
- Users cannot create their own API keys (admin-only via `/admin/auth/api-keys`)
- Middleware chain (Auth → Authz → Route → Proxy) exists in code but is not wired to `/v1/*` routes

## What Changes

- **Wire middleware to `/v1/*` endpoints**: Add `AuthMiddleware` → `AuthzMiddleware` → `RouteMiddleware` → `ProxyMiddleware` to the `/v1` route group in `gateway-service/cmd/server/main.go`
- **Implement `/v1/chat/completions` handler**: Replace mock response with real proxy logic using the existing `ProxyMiddleware`
- **Implement `/v1/models` handler**: Return real model list aggregated from providers via `ModelsHandler`
- **Add user self-service API key endpoint**: New `POST /v1/auth/api-keys` endpoint for authenticated users to create their own API keys (returns key once)
- **Document routing rule setup**: Ensure initial routing rules exist so model names resolve to providers (e.g., `gpt-4` → OpenAI provider)

## Capabilities

### New Capabilities
- `gateway-openai-api-wiring`: Wire authentication, authorization, routing, and proxy middleware to `/v1/*` OpenAI-compatible endpoints
- `gateway-user-api-key-self-service`: Enable users to create and manage their own API keys via `POST /v1/auth/api-keys`

### Modified Capabilities
- `gateway-service`: Middleware wiring changes for `/v1/*` endpoints; `/v1/chat/completions` and `/v1/models` now functional
- `auth-service`: `CreateAPIKey` logic already supports user-scoped key creation; may need `ListAPIKeys` / `DeleteAPIKey` for user's own keys

## Impact

- **Code**: `gateway-service/cmd/server/main.go` (route wiring), `gateway-service/internal/middleware/` (middleware chain)
- **APIs**: 
  - `POST /v1/chat/completions` — now functional (OpenAI-compatible)
  - `GET /v1/models` — now functional (OpenAI-compatible)
  - `POST /v1/auth/api-keys` — NEW (user self-service)
- **Dependencies**: Requires routing rules in `router-service` and providers in `provider-service` to be configured
- **Database**: No schema changes (uses existing `api_keys`, `users`, `routing_rules` tables)

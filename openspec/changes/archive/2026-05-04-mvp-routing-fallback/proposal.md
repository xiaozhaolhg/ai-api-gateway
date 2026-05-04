## Why

The current MVP routing system lacks fallback support: if the primary LLM provider fails, requests fail ungracefully. There is no way to configure routing rules via the admin API, and model names cannot be mapped between primary and fallback providers. This causes downtime during provider outages and blocks MVP production readiness. Basic reactive fallback (primary + one fallback) is needed to ensure request continuity for critical use cases.

## What Changes

- Gateway service: Add `/admin/routing-rules` CRUD endpoints to manage routing configuration
- Router service: Implement `getFallbackProviders()` to return primary + fallback providers with model mapping
- Router service: Add `fallback_model` field to `RoutingRule` entity and gRPC proto
- Gateway service: Implement fallback retry logic in `ProxyMiddleware` with model rewriting
- Provider service: Add `opencode-zen` adapter type inference to support the new fallback provider

## Capabilities

### New Capabilities
<!-- No new capabilities introduced - all changes modify existing service behaviors -->

### Modified Capabilities
- `router-service`: Route resolution requirements now include fallback provider chains and per-provider model mapping
- `gateway-service`: Proxy middleware requirements now include fallback retry logic; new admin API requirements added for routing rule management
- `provider-service`: Adapter type inference requirements now include `opencode-zen` provider type

## Impact

- **Code**: `router-service` (routing logic, proto, entities), `gateway-service` (admin handlers, proxy middleware, router client), `provider-service` (adapter inference)
- **APIs**: New gateway admin endpoints (`/admin/routing-rules`), updated router gRPC proto (`RoutingRule`, `RouteResult` messages)
- **Dependencies**: No new external dependencies; relies on existing Redis caching, gRPC communication between services
- **Systems**: Affects request routing flow, admin UI will need updates to expose routing rule management (out of scope for this change)
- **Database**: `router-service` SQLite/PostgreSQL requires `fallback_model` column addition to `routing_rules` table

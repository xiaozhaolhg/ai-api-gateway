# router-service

## Purpose

Routing domain — resolves model name to provider, manages fallback chains.

> **Note**: See existing implementation in `openspec/specs/router-service-architecture/spec.md`

## Service Responsibility

- **Role**: Model to provider resolution
- **Owned Entities**: RoutingRule
- **Data Layer**: router-db (SQLite/PostgreSQL)

## Dependencies

### Calls To

| Service | Methods | Purpose |
|---|---|---|
| provider-service | `GetProviderByType` | Verify provider exists |
| provider-service | `FindProvidersByModel` | Find providers supporting bare model names |
| provider-service | `HealthCheck` | Check provider health status |

### Called By

| Service | Methods | Purpose |
|---|---|---|
| gateway-service | `ResolveRoute` | Resolve model to provider |
| gateway-service | `RefreshRoutingTable` | Invalidate cache |

### Data Dependencies

- **Database**: router-db (RoutingRule)
- **Cache**: Redis (routing table)

## Key Design

### Route Resolution

1. Receive model name from gateway-service
2. Check if model contains ":" separator (provider:model) or is bare model name
3. If bare model: use FindProvidersByModel + health checks to resolve
4. If provider:model: match against RoutingRule patterns (wildcard support)
5. Return RouteResult with provider_id, adapter_type, fallback_provider_ids

### Key Operations

- **ResolveRoute**: model → provider
- **CreateRoutingRule/UpdateRoutingRule/DeleteRoutingRule**: Rule management
- **RefreshRoutingTable**: Cache invalidation

## Requirements

### Requirement: Redis Caching
Router service SHALL cache resolved routes in Redis with TTL, supporting both system and per-user rules.

#### Scenario: Cache Hit
- **WHEN** a `ResolveRoute` request is received and the route exists in Redis cache
- **THEN** return the cached route immediately without querying the database

#### Scenario: Cache Miss
- **WHEN** a `ResolveRoute` request is received and the route is not in cache
- **THEN** query the database, cache the result in Redis with 5-minute TTL, then return

#### Scenario: Cache Invalidation
- **WHEN** `RefreshRoutingTable` is called after provider configuration changes
- **THEN** clear all routing-related cache keys to force fresh lookups

#### Scenario: Cache Key Includes User Context
- **WHEN** caching a resolved route for a user with `user_id`
- **THEN** the cache key SHALL include the `user_id` to avoid cross-user cache collisions
- **AND** system rules and user rules SHALL have distinct cache keys

#### Scenario: Cache Invalidation for User Rules
- **WHEN** `RefreshRoutingTable` is called after a user's routing rule is created/updated/deleted
- **THEN** clear cache keys for that specific user's rules
- **AND** clear system-wide cache keys if system rules changed

### Requirement: Authorized Models Filtering

Router service SHALL filter routes based on authorized models passed from gateway.

#### Scenario: Authorized Route Resolution
- **WHEN** `ResolveRoute` is called with a model and `authorized_models` list
- **THEN** only return routes for providers serving models in the authorized list

#### Scenario: Unauthorized Model Request
- **WHEN** the requested model is not in `authorized_models`
- **THEN** return NOT_FOUND error without querying providers

### Requirement: Fallback Route Resolution
Router service SHALL return fallback provider IDs and corresponding model mappings when resolving routes, with support for per-user routing rules.

#### Scenario: Route resolution with user rule (user OVERRIDES system)
- **WHEN** `ResolveRoute` is called with a `user_id` and a matching user-specific routing rule exists
- **THEN** return `RouteResult` using the user's rule (provider_id, fallback_provider_ids)
- **AND** ignore any system-wide rule for the same model pattern

#### Scenario: Route resolution without user rule (fallback to system)
- **WHEN** `ResolveRoute` is called with a `user_id` but no user-specific rule matches
- **THEN** fall back to system-wide routing rules (where `user_id` is NULL or empty)
- **AND** return `RouteResult` with system rule's provider and fallback chain

#### Scenario: Route resolution without user_id (system rule only)
- **WHEN** `ResolveRoute` is called without `user_id` (or empty string)
- **THEN** only match system-wide routing rules
- **AND** return `RouteResult` with `fallback_provider_ids` populated if configured

#### Scenario: Route resolution with fallback chain
- **WHEN** `ResolveRoute` is called and the matching rule has `fallback_provider_ids` configured
- **THEN** return `RouteResult` with the ordered `fallback_provider_ids` list
- **AND** the gateway/provider service SHALL try each provider in order on failure

#### Scenario: Route resolution without fallback
- **WHEN** `ResolveRoute` is called and no fallback provider is configured for the matching rule
- **THEN** return `RouteResult` with empty `fallback_provider_ids`

### Requirement: Per-User Route Resolution
Router service ResolveRoute SHALL accept an optional `user_id` parameter to support per-user routing rules.

#### Scenario: ResolveRoute with user_id
- **WHEN** `ResolveRoute` is called with a `user_id` parameter
- **THEN** the router-service SHALL first look for routing rules where `user_id` matches
- **AND** only fall back to system rules (user_id IS NULL) if no user rule matches

#### Scenario: ResolveRoute without user_id
- **WHEN** `ResolveRoute` is called without `user_id` (or empty string)
- **THEN** router-service SHALL only match system-wide rules (user_id IS NULL)
- **AND** return NOT_FOUND if no system rule matches

### Requirement: Bare Model Name Detection

The router-service `ResolveRoute` SHALL detect bare model names (model strings that do not contain the ":" separator) and delegate to bare model resolution logic.

#### Scenario: Bare model name detection
- **WHEN** `ResolveRoute` is called with model="llama2" (no ":" separator)
- **THEN** the router SHALL invoke bare model resolution instead of pattern matching

#### Scenario: Provider-prefixed model name
- **WHEN** `ResolveRoute` is called with model="ollama:llama2" (contains ":")
- **THEN** the router SHALL use existing pattern matching logic (no bare model resolution)

### Requirement: Health Check Integration

The router-service SHALL use the provider-service `HealthCheck` RPC to determine provider health status when resolving bare model names.

#### Scenario: Health check via existing RPC
- **WHEN** resolving a bare model name with multiple supporting providers
- **THEN** the router SHALL call `HealthCheck` RPC for each provider concurrently
- **AND** wait for all health check results before selecting primary provider

#### Scenario: HealthCheck RPC failure
- **WHEN** `HealthCheck` RPC fails for a provider (network error, timeout)
- **THEN** that provider SHALL be treated as unhealthy
- **AND** excluded from primary/fallback selection

### Requirement: Health-Priority Selection

The router-service SHALL select the healthiest provider as primary, with remaining healthy providers as fallbacks.

#### Scenario: All providers healthy
- **WHEN** multiple providers support a model and all are healthy
- **THEN** select the first provider (by sorted order) as primary
- **AND** populate `fallback_provider_ids` with remaining providers in order

#### Scenario: Some providers unhealthy
- **WHEN** some providers are unhealthy
- **THEN** select the first healthy provider as primary
- **AND** populate `fallback_provider_ids` only with other healthy providers
- **AND** exclude unhealthy providers entirely

#### Scenario: No healthy providers
- **WHEN** no providers are healthy for a model
- **THEN** return an error: "no healthy provider found for model: <model>"

## References

- Architecture: `openspec/specs/router-service-architecture/spec.md`
- Deployment: `openspec/specs/router-service-deployment/spec.md`
- Testing: `openspec/specs/router-service-testing/spec.md`
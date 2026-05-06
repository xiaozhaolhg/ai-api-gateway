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
2. Match against RoutingRule patterns (wildcard support)
3. Return RouteResult with provider_id, adapter_type

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
- **THEN** the router-service SHALL only match system-wide rules (user_id IS NULL)
- **AND** return NOT_FOUND if no system rule matches

## References

- Architecture: `openspec/specs/router-service-architecture/spec.md`
- Deployment: `openspec/specs/router-service-deployment/spec.md`
- Testing: `openspec/specs/router-service-testing/spec.md`
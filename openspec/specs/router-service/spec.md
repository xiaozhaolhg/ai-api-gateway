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

Router service SHALL cache resolved routes in Redis with TTL.

#### Scenario: Cache Hit
- **WHEN** a `ResolveRoute` request is received and the route exists in Redis cache
- **THEN** return the cached route immediately without querying the database

#### Scenario: Cache Miss
- **WHEN** a `ResolveRoute` request is received and the route is not in cache
- **THEN** query the database, cache the result in Redis with 5-minute TTL, then return

#### Scenario: Cache Invalidation
- **WHEN** `RefreshRoutingTable` is called after provider configuration changes
- **THEN** clear all routing-related cache keys to force fresh lookups

### Requirement: Authorized Models Filtering

Router service SHALL filter routes based on authorized models passed from gateway.

#### Scenario: Authorized Route Resolution
- **WHEN** `ResolveRoute` is called with a model and `authorized_models` list
- **THEN** only return routes for providers serving models in the authorized list

#### Scenario: Unauthorized Model Request
- **WHEN** the requested model is not in `authorized_models`
- **THEN** return NOT_FOUND error without querying providers

## References

- Architecture: `openspec/specs/router-service-architecture/spec.md`
- Deployment: `openspec/specs/router-service-deployment/spec.md`
- Testing: `openspec/specs/router-service-testing/spec.md`
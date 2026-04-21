# router-service Architecture

> Routing domain — resolves model name to provider, manages fallback chains

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

## References

- Architecture: `openspec/specs/router-service-architecture/spec.md`
- Deployment: `openspec/specs/router-service-deployment/spec.md`
- Testing: `openspec/specs/router-service-testing/spec.md`
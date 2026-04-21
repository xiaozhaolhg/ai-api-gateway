## Why

Design documents under docs/ are currently unmanaged by OpenSpec. Migrating them enables tracking design changes, version control per-service architecture, and linking design to implementation. System spec must capture service relationships for overall view; service spec must show individual service dependencies.

## What Changes

- Create `openspec/specs/system/` with architecture.md (service relationships) and api-contracts.md (overall API definitions)
- Create service-level spec folders: gateway-service, auth-service, provider-service, billing-service, monitor-service
- Each service spec contains architecture.md showing design + dependencies to other services
- Reference router-service-architecture specs instead of duplicating (already exists)

## Capabilities

### New Capabilities

- `system-architecture`: Service calling relationships, dependency diagrams, layered architecture
- `system-api-contracts`: Shared gRPC API contracts definition used across services
- `gateway-service`: Edge layer, middleware pipeline, service orchestration
- `auth-service`: Identity, access control, model authorization
- `provider-service`: Provider management, adapters, callback mechanism
- `billing-service`: Usage tracking, pricing, budgets
- `monitor-service`: Metrics, alerting, health monitoring
- `router-service`: Reference to existing specs (do not duplicate)

### Modified Capabilities

- None — this is a design migration, not a requirement change

## Impact

- Spec artifacts: openspec/specs/{service}/{architecture.md, api-contracts.md}
- No code changes — spec management only
- router-service-architecture: reference existing, do not recreate
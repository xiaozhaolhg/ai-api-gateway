## Why

The current permission system requires administrators to assign individual model/provider permissions to each group, creating a combinatorial explosion of rules as the gateway scales. With dozens of models across multiple providers, configuring access per-group-per-resource is tedious, error-prone, and opaque to end users. The work division explicitly calls for tier-based access control to simplify permission management, reduce admin overhead, and provide clear access levels.

Today, `CheckModelAuthorization` is an MVP stub returning `*` (all models allowed), and granular `Permission` entities exist but are not enforced. A tiered system replaces this gap with a practical, enforceable model: groups are assigned to a tier, and tiers define which providers and models members can access.

## What Changes

- Introduce a **Tier** entity in auth-service with predefined access levels (Basic, Standard, Premium, Enterprise), each specifying allowed provider and model patterns
- Groups reference a tier instead of accumulating individual permission rules; the tier's patterns are resolved at authorization time
- Add admin CRUD for tiers (create, update, delete, list) and a custom tier builder with **bulk selection interface** for administrators to select multiple models from multiple providers in one attempt
- Update `CheckModelAuthorization` to resolve the user's group → tier → allowed patterns, replacing the current MVP stub
- Update admin UI permission assignment to show tier selection with a preview of allowed models/providers per tier, featuring **multi-select capabilities** for models across providers
- Migrate existing granular `Permission` records to tier-based assignments
- **Group-tier binding**: Groups now have a `tier_id` field; group creation accepts `tier_id`, group update can assign or remove tier assignment

## Capabilities

### New Capabilities
- `auth-tier-management`: Tier entity, CRUD operations, tier-based permission resolution logic, migration from granular permissions, and bulk model/provider selection interface for efficient tier configuration

### Modified Capabilities
- `auth-service`: `CheckModelAuthorization` changes from MVP stub to tier-based enforcement; group entity gains tier reference
- `admin-ui-permissions`: Permission assignment UI changes from individual rule editor to tier selection with model/provider preview

## Impact

- **auth-service**: New Tier entity, TierRepository, TierService, updated AuthService, database migration (tiers table, groups.tier_id column)
- **gateway-service**: No API change — consumes existing `CheckModelAuthorization` gRPC; behavior changes internally in auth-service
- **admin-ui**: New tier management page, updated group permission assignment flow
- **Proto**: New Tier-related messages and RPCs in `api/proto/auth/v1/auth.proto`
- **Database**: Migration adds `tiers` table and `tier_id` foreign key on `groups`

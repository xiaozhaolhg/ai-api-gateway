## Why

The current RBAC implementation is incomplete: all Group and Permission gRPC handlers return `nil` stubs, `CheckModelAuthorization` allows every active user to access all models (`"*"`), and `UserIdentity.group_ids` is always empty. Meanwhile, gateway admin endpoints return hardcoded mock data instead of calling auth/billing gRPC backends, making the admin UI non-functional for real operations. Completing the group/permission domain layer and wiring the gateway to actual service backends are prerequisites for any real authorization enforcement.

## What Changes

- Add Group, Permission, and UserGroupMembership domain entities to auth-service with GORM persistence and database migration
- Add GroupRepository, PermissionRepository, and UserGroupRepository interfaces and SQLite implementations
- Wire gateway admin handlers (`/admin/users`, `/admin/api-keys`, `/admin/usage`) to call auth-service and billing-service gRPC backends instead of returning mock data
- Fix `admin_auth.go` mock login to use auth-service `Login` RPC
- Add gateway HTTP routes for group and permission management that proxy to auth-service
- Populate `UserIdentity.group_ids` in `ValidateAPIKey` from UserGroupMembership table

## Capabilities

### New Capabilities
- `auth-group-management`: Group CRUD, user-group membership management, and group-scoped model patterns / token limits / rate limits in auth-service
- `auth-permission-management`: Permission CRUD with resource types (model, provider, admin_feature), actions, and allow/deny effects in auth-service
- `gateway-admin-wiring`: Gateway admin endpoints wired to auth-service and billing-service gRPC backends

### Modified Capabilities
- `auth-service`: Add group/permission entities, repositories, and populate group_ids in ValidateAPIKey response
- `gateway-service`: Replace mock admin handlers with real gRPC calls; add group/permission admin routes

## Impact

- **auth-service**: New entity files, repository interfaces/implementations, migration additions, handler changes (replacing nil stubs)
- **gateway-service**: Handler rewrites for admin endpoints, new route registrations, removal of mock data
- **api/proto/auth/v1/auth.proto**: No changes this sprint (existing proto messages suffice for Group/Permission)
- **billing-service**: No changes this sprint (group-scoped limits deferred to Sprint 3)
- **admin-ui**: No changes this sprint (existing API client already has group/permission endpoints)

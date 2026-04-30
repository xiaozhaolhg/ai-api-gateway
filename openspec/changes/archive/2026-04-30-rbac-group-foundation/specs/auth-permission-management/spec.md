## ADDED Requirements

### Requirement: Permission CRUD operations
The auth-service SHALL provide GrantPermission, RevokePermission, ListPermissions, and CheckPermission gRPC handlers that persist Permission entities with group_id, resource_type, resource_id, action, effect, and status fields.

#### Scenario: Grant a permission to a group
- **WHEN** GrantPermission is called with group_id, resource_type="model", resource_id="gpt-4", action="access", effect="allow"
- **THEN** a Permission entity is persisted with a generated UUID, the provided fields, status="active", and created_at/updated_at timestamps

#### Scenario: Grant permission with deny effect
- **WHEN** GrantPermission is called with resource_type="model", resource_id="gpt-4", action="access", effect="deny"
- **THEN** the Permission entity is persisted with effect="deny"

#### Scenario: Revoke a permission
- **WHEN** RevokePermission is called with an existing permission ID
- **THEN** the Permission entity is deleted

#### Scenario: List permissions by group
- **WHEN** ListPermissions is called with group_id and pagination
- **THEN** all permissions for that group are returned with total count

#### Scenario: Check permission for a user
- **WHEN** CheckPermission is called with user_id, resource_type, resource_id, and action
- **THEN** the system resolves the user's groups, collects all active permissions matching the resource_type/resource_id/action, and returns allowed=true if any permission has effect="allow" and no matching permission has effect="deny"

### Requirement: Permission resource types
The Permission entity SHALL support resource_type values: "model" (model access), "provider" (provider access), and "admin_feature" (admin UI feature access).

#### Scenario: Model resource type permission
- **WHEN** a Permission is created with resource_type="model" and resource_id="gpt-*"
- **THEN** the permission applies to model access authorization for models matching the glob pattern

#### Scenario: Admin feature resource type permission
- **WHEN** a Permission is created with resource_type="admin_feature" and resource_id="user_management"
- **THEN** the permission applies to admin UI feature access for the user_management feature

### Requirement: Permission effect semantics
When evaluating permissions, deny effects SHALL override allow effects for the same resource_type, resource_id, and action.

#### Scenario: Deny overrides allow
- **WHEN** a user's groups have both an "allow" and a "deny" permission for resource_type="model", resource_id="gpt-4", action="access"
- **THEN** CheckPermission SHALL return allowed=false

# auth-service

## Purpose

Identity, access control, and model authorization domain for the AI API Gateway.

## Service Responsibility

- **Role**: Identity validation, user management, model authorization
- **Owned Entities**: User, Group, APIKey, Permission
- **Data Layer**: auth-db (SQLite/PostgreSQL)

## Dependencies

### Calls To

| Service | Methods | Purpose |
|---------|----------|----------|
| (none) | — | Does not call other internal services |

### Called By

| Service | Methods | Purpose |
|---------|----------|----------|
| gateway-service | `ValidateAPIKey`, `CheckModelAuthorization` | Authenticate requests, check model permissions |
| gateway-service | `CreateUser`, `UpdateUser`, `DeleteUser` | User CRUD |
| gateway-service | `CreateAPIKey`, `DeleteAPIKey` | API key management |
| gateway-service | `Register`, `Login` | User registration and login |

### Data Dependencies

- **Database**: auth-db (User, Group, APIKey, Permission)
- **Cache**: Redis (API key → user lookup, group membership)

## Key Design

### Authentication Flow

1. Receive API key from gateway-service
2. Hash key and lookup in database
3. Return UserIdentity with user_id, role, group_ids

### Model Authorization Flow

1. Receive user_id, group_ids, model from gateway-service
2. Check Permission entities for group → model mapping
3. Return AuthorizationResult with allowed and authorized_models list

### Key Operations

- **ValidateAPIKey**: API key → UserIdentity
- **CheckModelAuthorization**: user/group + model → allowed/reason
- **Register**: username/email + password → JWT token (new user)
- **Login**: username/email + password → JWT token
- **CreateUser/UpdateUser/DeleteUser**: User CRUD
- **CreateAPIKey/DeleteAPIKey**: API key lifecycle (key returned once)
- **CreateGroup/AddUserToGroup**: Group management (Phase 2+)
- **GrantPermission/RevokePermission**: Model access control (Phase 2+)

---

## Requirements

### Requirement: Login RPC support
The auth-service SHALL support email/password authentication for admin UI users.

#### Scenario: Login request
- **WHEN** gateway-service calls Login RPC with email/password
- **THEN** validate credentials and return JWT with user identity.

#### Scenario: Password validation
- **WHEN** validating login credentials
- **THEN** compare with stored password hash using secure algorithm.

#### Scenario: JWT generation
- **WHEN** login is successful
- **THEN** generate JWT with user ID, role, and expiration.

### Requirement: User password management
The auth-service SHALL store and manage user passwords securely.

#### Scenario: User creation with password
- **WHEN** creating new user with password
- **THEN** hash password before storage (never store plaintext).

#### Scenario: Password update
- **WHEN** updating user password
- **THEN** validate current password and hash new password.

#### Scenario: Password reset
- **WHEN** resetting user password
- **THEN** generate temporary password and force change on next login.

### Requirement: Extended role support
The auth-service SHALL support three user roles: admin, user, viewer.

#### Scenario: Role validation
- **WHEN** creating or updating user
- **THEN** validate role is one of: admin, user, viewer.

#### Scenario: Default role assignment
- **WHEN** creating user without explicit role
- **THEN** assign 'user' as default role.

#### Scenario: Role-based permissions
- **WHEN** checking user permissions
- **THEN** use role to determine access level in authorization logic.

### Requirement: ValidateAPIKey returns group membership
The auth-service ValidateAPIKey RPC SHALL query UserGroupMembership records for the resolved user and populate the group_ids field in the UserIdentity response.

#### Scenario: User with group memberships
- **WHEN** ValidateAPIKey resolves a user who belongs to groups "group-1" and "group-2"
- **THEN** the UserIdentity response SHALL include group_ids=["group-1","group-2"]

#### Scenario: User with no group memberships
- **WHEN** ValidateAPIKey resolves a user who does not belong to any group
- **THEN** the UserIdentity response SHALL include group_ids=[] (empty, not nil)

### Requirement: Group and Permission entity persistence
The auth-service SHALL persist Group, Permission, and UserGroupMembership entities in the auth-db database and auto-migrate them on startup.

#### Scenario: Database migration on startup
- **WHEN** auth-service starts
- **THEN** Group, Permission, and UserGroupMembership tables are created/updated via GORM AutoMigrate

### Requirement: Group management RPCs implemented
The auth-service SHALL implement CreateGroup, UpdateGroup, DeleteGroup, ListGroups, AddUserToGroup, and RemoveUserFromGroup gRPC handlers with full persistence logic.

#### Scenario: CreateGroup handler
- **WHEN** CreateGroup RPC is called with valid parameters
- **THEN** a Group entity is persisted and returned with generated ID and timestamps

#### Scenario: ListGroups handler
- **WHEN** ListGroups RPC is called with pagination
- **THEN** groups are returned from the database with total count

### Requirement: Permission management RPCs implemented
The auth-service SHALL implement GrantPermission, RevokePermission, ListPermissions, and CheckPermission gRPC handlers with full persistence logic.

#### Scenario: GrantPermission handler
- **WHEN** GrantPermission RPC is called with valid parameters
- **THEN** a Permission entity is persisted and returned with generated ID and timestamps

#### Scenario: CheckPermission handler
- **WHEN** CheckPermission RPC is called with user_id, resource_type, resource_id, action
- **THEN** the system resolves user groups, collects matching permissions, and returns allowed=true only if an allow permission exists with no matching deny

### Requirement: Group CRUD operations
The auth-service SHALL provide CreateGroup, UpdateGroup, DeleteGroup, and ListGroup gRPC handlers that persist Group entities with name, description, parent_group_id, model_patterns, token_limit, and rate_limit fields.

#### Scenario: Create a new group
- **WHEN** CreateGroup is called with name "developers" and description "Developer team"
- **THEN** a Group entity is persisted with a generated UUID, the provided fields, and empty model_patterns/token_limit/rate_limit defaults

#### Scenario: Create group with model patterns and limits
- **WHEN** CreateGroup is called with name "power-users", model_patterns=["gpt-4","claude-*"], token_limit={prompt_tokens:100000,completion_tokens:100000,period:"daily"}, rate_limit={requests_per_minute:60,requests_per_day:10000}
- **THEN** the Group entity is persisted with all provided configuration

#### Scenario: Update a group
- **WHEN** UpdateGroup is called with an existing group ID and new name "senior-devs"
- **THEN** the Group entity is updated with the new name and updated_at timestamp

#### Scenario: Delete a group
- **WHEN** DeleteGroup is called with an existing group ID
- **THEN** the Group entity and all associated UserGroupMembership records are removed

#### Scenario: List groups with pagination
- **WHEN** ListGroups is called with page=1, page_size=10
- **THEN** up to 10 groups are returned with total count

### Requirement: User-Group membership management
The auth-service SHALL provide AddUserToGroup and RemoveUserToGroup gRPC handlers that manage UserGroupMembership records linking users to groups.

#### Scenario: Add user to group
- **WHEN** AddUserToGroup is called with user_id and group_id
- **THEN** a UserGroupMembership record is created with a generated UUID and added_at timestamp

#### Scenario: Add user to group that they are already in
- **WHEN** AddUserToGroup is called with a user_id and group_id that already has a membership
- **THEN** the operation SHALL return an error indicating duplicate membership

#### Scenario: Remove user from group
- **WHEN** RemoveUserFromGroup is called with user_id and group_id
- **THEN** the UserGroupMembership record is deleted

#### Scenario: Remove user from group they are not in
- **WHEN** RemoveUserFromGroup is called with a user_id and group_id that has no membership
- **THEN** the operation SHALL return success (idempotent)

### Requirement: Group-scoped model patterns
A Group entity SHALL carry a model_patterns field (list of glob patterns) that defines which models members of the group are authorized to access. This field is stored but not enforced in this sprint.

#### Scenario: Group with model patterns
- **WHEN** a Group is created with model_patterns=["gpt-4","claude-*"]
- **THEN** the patterns are persisted and retrievable via GetByID/List

### Requirement: Group-scoped token limits
A Group entity SHALL carry an optional token_limit field with prompt_tokens, completion_tokens, and period. This field is stored but not enforced in this sprint.

#### Scenario: Group with token limit
- **WHEN** a Group is created with token_limit={prompt_tokens:50000,completion_tokens:50000,period:"daily"}
- **THEN** the limit is persisted and retrievable via GetByID/List

### Requirement: Group-scoped rate limits
A Group entity SHALL carry an optional rate_limit field with requests_per_minute and requests_per_day. This field is stored but not enforced in this sprint.

#### Scenario: Group with rate limit
- **WHEN** a Group is created with rate_limit={requests_per_minute:30,requests_per_day:5000}
- **THEN** the limit is persisted and retrievable via GetByID/List

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

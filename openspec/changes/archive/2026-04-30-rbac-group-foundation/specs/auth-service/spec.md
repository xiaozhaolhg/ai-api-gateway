## MODIFIED Requirements

### Requirement: ValidateAPIKey returns group membership
The auth-service ValidateAPIKey RPC SHALL query UserGroupMembership records for the resolved user and populate the group_ids field in the UserIdentity response.

#### Scenario: User with group memberships
- **WHEN** ValidateAPIKey resolves a user who belongs to groups "group-1" and "group-2"
- **THEN** the UserIdentity response SHALL include group_ids=["group-1","group-2"]

#### Scenario: User with no group memberships
- **WHEN** ValidateAPIKey resolves a user who does not belong to any group
- **THEN** the UserIdentity response SHALL include group_ids=[] (empty, not nil)

## ADDED Requirements

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

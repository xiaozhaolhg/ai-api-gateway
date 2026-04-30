## 1. Auth-Service Domain Layer

- [x] 1.1 Create Group entity (`auth-service/internal/domain/entity/group.go`) with ID, Name, Description, ParentGroupID, ModelPatterns, TokenLimit, RateLimit, CreatedAt, UpdatedAt
- [x] 1.2 Create Permission entity (`auth-service/internal/domain/entity/permission.go`) with ID, GroupID, ResourceType, ResourceID, Action, Effect, Status, CreatedAt, UpdatedAt
- [x] 1.3 Create UserGroupMembership entity (`auth-service/internal/domain/entity/user_group.go`) with ID, UserID, GroupID, AddedAt
- [x] 1.4 Add GroupRepository, PermissionRepository, UserGroupRepository interfaces to `auth-service/internal/domain/port/repository.go`
- [x] 1.5 Implement GroupRepository GORM implementation (`auth-service/internal/infrastructure/repository/group_repository.go`)
- [x] 1.6 Implement PermissionRepository GORM implementation (`auth-service/internal/infrastructure/repository/permission_repository.go`)
- [x] 1.7 Implement UserGroupRepository GORM implementation (`auth-service/internal/infrastructure/repository/user_group_repository.go`)
- [x] 1.8 Update migration (`auth-service/internal/infrastructure/migration/migration.go`) to add Group, Permission, UserGroupMembership to AutoMigrate

## 2. Auth-Service Application Layer

- [x] 2.1 Add GroupService with CreateGroup, UpdateGroup, DeleteGroup, ListGroups methods
- [x] 2.2 Add PermissionService with GrantPermission, RevokePermission, ListPermissions, CheckPermission methods
- [x] 2.3 Add UserGroupService with AddUserToGroup, RemoveUserFromGroup, GetUserGroups methods
- [x] 2.4 Update AuthService to accept GroupRepository, PermissionRepository, UserGroupRepository dependencies
- [x] 2.5 Update ValidateAPIKey to populate group_ids from UserGroupMembership table
- [x] 2.6 Wire new repositories in `auth-service/cmd/server/main.go` dependency injection

## 3. Auth-Service gRPC Handlers

- [x] 3.1 Implement CreateGroup handler (replace nil stub)
- [x] 3.2 Implement UpdateGroup handler (replace nil stub)
- [x] 3.3 Implement DeleteGroup handler (replace nil stub)
- [x] 3.4 Implement ListGroups handler (replace nil stub)
- [x] 3.5 Implement AddUserToGroup handler (replace nil stub)
- [x] 3.6 Implement RemoveUserFromGroup handler (replace nil stub)
- [x] 3.7 Implement GrantPermission handler (replace nil stub)
- [x] 3.8 Implement RevokePermission handler (replace nil stub)
- [x] 3.9 Implement ListPermissions handler (replace nil stub)
- [x] 3.10 Implement CheckPermission handler (replace nil stub)

## 4. Gateway Admin Handler Wiring

- [x] 4.1 Replace handleListUsers mock with auth-service ListUsers gRPC call
- [x] 4.2 Replace handleCreateUser mock with auth-service CreateUser gRPC call
- [x] 4.3 Replace handleUpdateUser mock with auth-service UpdateUser gRPC call
- [x] 4.4 Replace handleDeleteUser mock with auth-service DeleteUser gRPC call
- [x] 4.5 Add handleListAPIKeys calling auth-service ListAPIKeys gRPC
- [x] 4.6 Add handleCreateAPIKey calling auth-service CreateAPIKey gRPC
- [x] 4.7 Add handleDeleteAPIKey calling auth-service DeleteAPIKey gRPC
- [x] 4.8 Replace handleGetUsage mock with billing-service GetUsage gRPC call
- [x] 4.9 Initialize billingClient in gateway main.go alongside authClient

## 5. Gateway Group/Permission Admin Routes

- [x] 5.1 Add GET /admin/auth/groups route → ListGroups gRPC
- [x] 5.2 Add POST /admin/auth/groups route → CreateGroup gRPC
- [x] 5.3 Add PUT /admin/auth/groups/:id route → UpdateGroup gRPC
- [x] 5.4 Add DELETE /admin/auth/groups/:id route → DeleteGroup gRPC
- [x] 5.5 Add POST /admin/auth/groups/:id/members route → AddUserToGroup gRPC
- [x] 5.6 Add DELETE /admin/auth/groups/:id/members/:user_id route → RemoveUserFromGroup gRPC
- [x] 5.7 Add GET /admin/auth/permissions route → ListPermissions gRPC
- [x] 5.8 Add POST /admin/auth/permissions route → GrantPermission gRPC
- [x] 5.9 Add DELETE /admin/auth/permissions/:id route → RevokePermission gRPC

## 6. Gateway Auth Fix

- [x] 6.1 Remove mock login from admin_auth.go (already using auth-service in main.go handleLogin)
- [x] 6.2 Ensure /admin/auth/login always calls auth-service Login gRPC (verify no fallback to mock)

## 7. Testing

- [x] 7.1 Add unit tests for GroupRepository (CRUD + pagination)
- [x] 7.2 Add unit tests for PermissionRepository (CRUD + group filtering)
- [x] 7.3 Add unit tests for UserGroupRepository (add/remove/query)
- [x] 7.4 Add unit tests for GroupService (create, update, delete, list)
- [x] 7.5 Add unit tests for PermissionService (grant, revoke, check with deny-override)
- [x] 7.6 Add unit tests for UserGroupService (add, remove, duplicate detection)
- [x] 7.7 Verify auth-service compiles and starts with new entities/migrations
- [x] 7.8 Verify gateway compiles and admin endpoints return real data

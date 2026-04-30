## ADDED Requirements

### Requirement: Gateway admin user handlers call auth-service
The gateway admin handlers for user CRUD operations SHALL call auth-service gRPC methods instead of returning hardcoded mock data.

#### Scenario: List users via gateway
- **WHEN** GET /admin/auth/users is called with valid JWT
- **THEN** the gateway calls auth-service ListUsers gRPC and returns the real user list

#### Scenario: Create user via gateway
- **WHEN** POST /admin/auth/users is called with name, email, role, password
- **THEN** the gateway calls auth-service CreateUser gRPC and returns the created user

#### Scenario: Update user via gateway
- **WHEN** PUT /admin/auth/users/:id is called with updated fields
- **THEN** the gateway calls auth-service UpdateUser gRPC and returns the updated user

#### Scenario: Delete user via gateway
- **WHEN** DELETE /admin/auth/users/:id is called
- **THEN** the gateway calls auth-service DeleteUser gRPC and returns success

### Requirement: Gateway admin API key handlers call auth-service
The gateway admin handlers for API key operations SHALL call auth-service gRPC methods.

#### Scenario: Create API key via gateway
- **WHEN** POST /admin/auth/api-keys is called with user_id and name
- **THEN** the gateway calls auth-service CreateAPIKey gRPC and returns the API key (shown once)

#### Scenario: List API keys via gateway
- **WHEN** GET /admin/auth/api-keys/:user_id is called
- **THEN** the gateway calls auth-service ListAPIKeys gRPC and returns the key list

#### Scenario: Delete API key via gateway
- **WHEN** DELETE /admin/auth/api-keys/:id is called
- **THEN** the gateway calls auth-service DeleteAPIKey gRPC and returns success

### Requirement: Gateway admin usage handler calls billing-service
The gateway admin handler for usage queries SHALL call billing-service gRPC instead of returning mock data.

#### Scenario: Get usage via gateway
- **WHEN** GET /admin/auth/usage is called
- **THEN** the gateway calls billing-service GetUsage gRPC and returns real usage records

### Requirement: Gateway group admin endpoints
The gateway SHALL expose HTTP endpoints for group management that proxy to auth-service gRPC.

#### Scenario: List groups via gateway
- **WHEN** GET /admin/auth/groups is called with valid JWT
- **THEN** the gateway calls auth-service ListGroups gRPC and returns the group list

#### Scenario: Create group via gateway
- **WHEN** POST /admin/auth/groups is called with name and description
- **THEN** the gateway calls auth-service CreateGroup gRPC and returns the created group

#### Scenario: Update group via gateway
- **WHEN** PUT /admin/auth/groups/:id is called with updated fields
- **THEN** the gateway calls auth-service UpdateGroup gRPC and returns the updated group

#### Scenario: Delete group via gateway
- **WHEN** DELETE /admin/auth/groups/:id is called
- **THEN** the gateway calls auth-service DeleteGroup gRPC and returns success

#### Scenario: Add user to group via gateway
- **WHEN** POST /admin/auth/groups/:id/members is called with user_id
- **THEN** the gateway calls auth-service AddUserToGroup gRPC and returns success

#### Scenario: Remove user from group via gateway
- **WHEN** DELETE /admin/auth/groups/:id/members/:user_id is called
- **THEN** the gateway calls auth-service RemoveUserFromGroup gRPC and returns success

### Requirement: Gateway permission admin endpoints
The gateway SHALL expose HTTP endpoints for permission management that proxy to auth-service gRPC.

#### Scenario: List permissions via gateway
- **WHEN** GET /admin/auth/permissions is called with valid JWT
- **THEN** the gateway calls auth-service ListPermissions gRPC and returns the permission list

#### Scenario: Grant permission via gateway
- **WHEN** POST /admin/auth/permissions is called with group_id, resource_type, resource_id, action, effect
- **THEN** the gateway calls auth-service GrantPermission gRPC and returns the created permission

#### Scenario: Revoke permission via gateway
- **WHEN** DELETE /admin/auth/permissions/:id is called
- **THEN** the gateway calls auth-service RevokePermission gRPC and returns success

### Requirement: Admin login uses auth-service
The gateway admin login handler SHALL validate credentials via auth-service Login gRPC instead of accepting any input.

#### Scenario: Valid login
- **WHEN** POST /admin/auth/login is called with valid email and password
- **THEN** the gateway calls auth-service Login gRPC, receives a token, sets auth cookie, and returns token + user

#### Scenario: Invalid login
- **WHEN** POST /admin/auth/login is called with invalid credentials
- **THEN** the gateway returns 401 Unauthorized

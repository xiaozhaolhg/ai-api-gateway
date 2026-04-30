## MODIFIED Requirements

### Requirement: Admin user management endpoints
The gateway admin user management endpoints SHALL call auth-service gRPC backends instead of returning hardcoded mock data.

#### Scenario: List users
- **WHEN** GET /admin/auth/users is called
- **THEN** the gateway calls auth-service ListUsers gRPC and returns the response as JSON

#### Scenario: Create user
- **WHEN** POST /admin/auth/users is called with user data
- **THEN** the gateway calls auth-service CreateUser gRPC and returns the created user as JSON

#### Scenario: Update user
- **WHEN** PUT /admin/auth/users/:id is called with updated fields
- **THEN** the gateway calls auth-service UpdateUser gRPC and returns the updated user as JSON

#### Scenario: Delete user
- **WHEN** DELETE /admin/auth/users/:id is called
- **THEN** the gateway calls auth-service DeleteUser gRPC and returns success

## ADDED Requirements

### Requirement: Admin API key management endpoints
The gateway SHALL expose API key management endpoints that call auth-service gRPC.

#### Scenario: Create API key
- **WHEN** POST /admin/auth/api-keys is called with user_id and name
- **THEN** the gateway calls auth-service CreateAPIKey gRPC and returns the key (shown once)

#### Scenario: List API keys
- **WHEN** GET /admin/auth/api-keys/:user_id is called
- **THEN** the gateway calls auth-service ListAPIKeys gRPC and returns the key list

#### Scenario: Delete API key
- **WHEN** DELETE /admin/auth/api-keys/:id is called
- **THEN** the gateway calls auth-service DeleteAPIKey gRPC and returns success

### Requirement: Admin usage endpoint calls billing-service
The gateway admin usage endpoint SHALL call billing-service gRPC instead of returning mock data.

#### Scenario: Get usage
- **WHEN** GET /admin/auth/usage is called
- **THEN** the gateway calls billing-service GetUsage gRPC and returns real usage records

### Requirement: Admin group management endpoints
The gateway SHALL expose group management endpoints that proxy to auth-service gRPC.

#### Scenario: Group CRUD
- **WHEN** any of GET/POST/PUT/DELETE /admin/auth/groups[/:id] is called
- **THEN** the gateway calls the corresponding auth-service Group gRPC method and returns the result

#### Scenario: Group membership management
- **WHEN** POST /admin/auth/groups/:id/members or DELETE /admin/auth/groups/:id/members/:user_id is called
- **THEN** the gateway calls AddUserToGroup or RemoveUserFromGroup gRPC respectively

### Requirement: Admin permission management endpoints
The gateway SHALL expose permission management endpoints that proxy to auth-service gRPC.

#### Scenario: Permission CRUD
- **WHEN** any of GET/POST/DELETE /admin/auth/permissions[/:id] is called
- **THEN** the gateway calls the corresponding auth-service Permission gRPC method and returns the result

### Requirement: Admin login uses auth-service
The gateway admin login handler SHALL validate credentials via auth-service Login gRPC.

#### Scenario: Login with valid credentials
- **WHEN** POST /admin/auth/login is called with valid email and password
- **THEN** the gateway calls auth-service Login gRPC, sets auth cookie, and returns token + user

#### Scenario: Login with invalid credentials
- **WHEN** POST /admin/auth/login is called with invalid credentials
- **THEN** the gateway returns 401 Unauthorized

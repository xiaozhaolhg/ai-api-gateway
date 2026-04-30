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

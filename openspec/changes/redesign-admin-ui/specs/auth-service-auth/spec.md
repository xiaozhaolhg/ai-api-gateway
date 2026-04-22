## MODIFIED Requirements

### Requirement: Login RPC support
The auth-service SHALL support email/password authentication for admin UI users.

#### Scenario: Login request
- **WHEN** gateway-service calls Login RPC with email/password
- **THEN** validate credentials and return JWT with user identity

#### Scenario: Password validation
- **WHEN** validating login credentials
- **THEN** compare with stored password hash using secure algorithm

#### Scenario: JWT generation
- **WHEN** login is successful
- **THEN** generate JWT with user ID, role, and expiration

### Requirement: User password management
The auth-service SHALL store and manage user passwords securely.

#### Scenario: User creation with password
- **WHEN** creating new user with password
- **THEN** hash password before storage (never store plaintext)

#### Scenario: Password update
- **WHEN** updating user password
- **THEN** validate current password and hash new password

#### Scenario: Password reset
- **WHEN** resetting user password
- **THEN** generate temporary password and force change on next login

### Requirement: Extended role support
The auth-service SHALL support three user roles: admin, user, viewer.

#### Scenario: Role validation
- **WHEN** creating or updating user
- **THEN** validate role is one of: admin, user, viewer

#### Scenario: Default role assignment
- **WHEN** creating user without explicit role
- **THEN** assign 'user' as default role

#### Scenario: Role-based permissions
- **WHEN** checking user permissions
- **THEN** use role to determine access level in authorization logic

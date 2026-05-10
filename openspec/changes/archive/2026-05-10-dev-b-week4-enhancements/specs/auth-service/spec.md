## MODIFIED Requirements

### Requirement: Login RPC support
The auth-service SHALL support email/password and username/password authentication for admin UI users.

#### Scenario: Login request with email
- **WHEN** gateway-service calls Login RPC with email/password
- **THEN** validate credentials and return JWT with user identity.

#### Scenario: Login request with username
- **WHEN** gateway-service calls Login RPC with username/password
- **THEN** validate credentials and return JWT with user identity.

#### Scenario: Password validation
- **WHEN** validating login credentials
- **THEN** compare with stored password hash using secure algorithm.

#### Scenario: JWT generation
- **WHEN** login is successful
- **THEN** generate JWT with user ID, role, and expiration.

### Requirement: User password management
The auth-service SHALL store and manage user passwords securely with username support.

#### Scenario: User creation with password and username
- **WHEN** creating new user with password and username
- **THEN** hash password before storage (never store plaintext)
- **AND** validate username uniqueness.

#### Scenario: User creation with password and email only
- **WHEN** creating new user with password but no username
- **THEN** hash password before storage (never store plaintext)
- **AND** allow creation with empty username field.

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

### Requirement: User CRUD operations with mandatory immutable username
The auth-service SHALL provide CreateUser, UpdateUser, DeleteUser, and ListUsers gRPC handlers that persist User entities with name, email, mandatory immutable username, role, and status fields.

#### Scenario: Create a new user with mandatory username
- **WHEN** CreateUser is called with name, email, mandatory username, role, and password
- **THEN** a User entity is persisted with a generated UUID, provided fields, and created_at timestamp
- **AND** username uniqueness is enforced.
- **AND** username is required and cannot be empty.

#### Scenario: Update user without changing username
- **WHEN** UpdateUser is called with an existing user ID and new name/email/role but no username
- **THEN** User entity is updated with new fields and updated_at timestamp
- **AND** username remains unchanged (immutable).

#### Scenario: Attempt to update username
- **WHEN** UpdateUser is called with an existing user ID and new username
- **THEN** operation SHALL return error indicating username cannot be changed
- **AND** username remains unchanged.

#### Scenario: List users with usernames
- **WHEN** ListUsers is called with pagination
- **THEN** users are returned with username field populated
- **AND** response includes total count.

### Requirement: Username uniqueness enforcement
The auth-service SHALL enforce unique usernames across all users.

#### Scenario: Database constraint
- **WHEN** attempting to create user with duplicate username
- **THEN** database rejects insertion with unique constraint violation
- **AND** system returns appropriate error message.

#### Scenario: Case sensitivity
- **WHEN** creating usernames "JohnDoe" and "johndoe"
- **THEN** system treats them as different usernames
- **AND** allows both to exist simultaneously.

### Requirement: Group CRUD operations with description support
The auth-service SHALL provide CreateGroup, UpdateGroup, DeleteGroup, and ListGroup gRPC handlers that persist Group entities with name, description, parent_group_id, model_patterns, token_limit, and rate_limit fields.

#### Scenario: Create a new group with description
- **WHEN** CreateGroup is called with name "developers" and description "Developer team"
- **THEN** a Group entity is persisted with a generated UUID, provided fields including description, and empty model_patterns/token_limit/rate_limit defaults
- **AND** response includes description field.

#### Scenario: Update group description
- **WHEN** UpdateGroup is called with an existing group ID and new description
- **THEN** Group entity is updated with new description and updated_at timestamp
- **AND** response includes updated description field.

#### Scenario: Create group without description
- **WHEN** CreateGroup is called with name but no description
- **THEN** a Group entity is persisted with empty description field
- **AND** response includes empty description field.

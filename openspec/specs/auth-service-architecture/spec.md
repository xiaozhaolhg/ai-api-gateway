## Purpose

Identity, access control, and model authorization service with gRPC interface and SQLite persistence.
## Requirements
### Requirement: DDD four-layer architecture
The auth-service SHALL implement four-layer Clean Architecture: Domain, Application, Infrastructure, and Handler with dependency direction from outer to inner layers.

#### Scenario: Domain layer has no external dependencies
- **WHEN** the domain layer is imported
- **THEN** it SHALL NOT import any code from application, infrastructure, or handler layers

#### Scenario: Infrastructure implements domain interfaces
- **WHEN** a new repository is needed
- **THEN** it SHALL be implemented by creating a struct that implements the Repository interface defined in the domain layer

### Requirement: User entity and repository
The auth-service SHALL own the User entity with fields: id, name, **username**, email, role, status, created_at. It SHALL provide a UserRepository interface for CRUD operations.

#### Scenario: Create user
- **WHEN** a CreateUser request is received via gRPC
- **THEN** the service SHALL persist a new User entity to SQLite
- **AND** return the created User with generated id and timestamp

#### Scenario: List users
- **WHEN** a ListUsers request is received
- **THEN** the service SHALL return all users from SQLite with pagination support

### Requirement: APIKey entity and repository
The auth-service SHALL own the APIKey entity with fields: id, user_id, key_hash, name, scopes, created_at, expires_at. It SHALL provide an APIKeyRepository interface.

#### Scenario: Create API key
- **WHEN** a CreateAPIKey request is received
- **THEN** the service SHALL generate a random API key, store its SHA-256 hash, and return the plain key once
- **AND** the plain key SHALL NOT be stored or retrievable after creation

#### Scenario: Validate API key
- **WHEN** a ValidateAPIKey request is received with an API key string
- **THEN** the service SHALL hash the key, look up the hash in SQLite, and return UserIdentity with user_id, role, and group_ids

### Requirement: Model authorization
The auth-service SHALL implement CheckModelAuthorization that determines whether a user/group can access a specific model.

#### Scenario: MVP authorization — all active users allowed
- **WHEN** CheckModelAuthorization is called for any active user in MVP
- **THEN** the service SHALL return allowed=true with authorized_models containing all known models

#### Scenario: Disabled user denied
- **WHEN** CheckModelAuthorization is called for a user with status "disabled"
- **THEN** the service SHALL return allowed=false with reason "user disabled"

### Requirement: gRPC server implementation
The auth-service SHALL implement the AuthService gRPC server as defined in the api/ module proto definitions.

#### Scenario: All proto RPCs are implemented
- **WHEN** the auth-service starts
- **THEN** it SHALL register all RPCs defined in auth.proto: ValidateAPIKey, CheckModelAuthorization, GetUser, CreateUser, UpdateUser, DeleteUser, ListUsers, CreateAPIKey, DeleteAPIKey, ListAPIKeys, Register, Login

### Requirement: SQLite persistence
The auth-service SHALL use SQLite as its database, with tables for users and api_keys, managed via migrations.

#### Scenario: Database initialized on startup
- **WHEN** the auth-service starts
- **THEN** it SHALL create the SQLite database file and run migrations if the database is new
- **AND** existing data SHALL be preserved on restart


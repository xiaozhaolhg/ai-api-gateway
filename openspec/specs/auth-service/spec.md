# auth-service Architecture

> Identity, access control, and model authorization domain

## Service Responsibility

- **Role**: Identity validation, user management, model authorization
- **Owned Entities**: User, Group, APIKey, Permission
- **Data Layer**: auth-db (SQLite/PostgreSQL)

## Dependencies

### Calls To

| Service | Methods | Purpose |
|---|---|---|
| (none) | — | Does not call other internal services |

### Called By

| Service | Methods | Purpose |
|---|---|---|
| gateway-service | `ValidateAPIKey`, `CheckModelAuthorization` | Authenticate requests, check model permissions |
| gateway-service | `CreateUser`, `UpdateUser`, `DeleteUser` | User CRUD |
| gateway-service | `CreateAPIKey`, `DeleteAPIKey` | API key management |

### Data Dependencies

- **Database**: auth-db (User, Group, APIKey, Permission)
- **Cache**: Redis (API key → user lookup, group memberships)

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
- **CreateUser/UpdateUser/DeleteUser**: User CRUD
- **CreateAPIKey/DeleteAPIKey**: API key lifecycle (key returned once)
- **CreateGroup/AddUserToGroup**: Group management (Phase 2+)
- **GrantPermission/RevokePermission**: Model access control (Phase 2+)
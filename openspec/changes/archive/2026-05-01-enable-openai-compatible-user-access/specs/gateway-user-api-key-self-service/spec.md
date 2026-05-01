## ADDED Requirements

### Requirement: User self-service API key creation
The gateway-service SHALL provide an endpoint for authenticated users to create API keys for their own account.

#### Scenario: Authenticated user creates API key
- **WHEN** an authenticated user (with valid JWT) makes a POST request to `/v1/auth/api-keys` with a request body containing `name` (string)
- **THEN** the gateway SHALL call `auth-service.CreateAPIKey` with the user's ID from the JWT context
- **AND** return HTTP 201 with `api_key_id`, `api_key` (shown only once), and `name`

#### Scenario: Missing authentication
- **WHEN** an unauthenticated request is made to `/v1/auth/api-keys`
- **THEN** the gateway SHALL return HTTP 401 with error message "authorization required"

#### Scenario: Missing key name
- **WHEN** a request is made to `/v1/auth/api-keys` without a `name` field
- **THEN** the gateway SHALL return HTTP 400 with error message "invalid request"

### Requirement: User self-service API key listing
The gateway-service SHALL provide an endpoint for authenticated users to list their own API keys.

#### Scenario: List own API keys
- **WHEN** an authenticated user makes a GET request to `/v1/auth/api-keys`
- **THEN** the gateway SHALL call `auth-service.ListAPIKeys` with the user's ID from the JWT context
- **AND** return HTTP 200 with the list of API keys (with `id`, `name`, `created_at`, but NOT the full key)

#### Scenario: API key list masks sensitive data
- **WHEN** the API key list is returned
- **THEN** the response SHALL NOT include the `api_key` field (only `api_key_id` for identification)

### Requirement: User self-service API key deletion
The gateway-service SHALL provide an endpoint for authenticated users to delete their own API keys.

#### Scenario: Delete own API key
- **WHEN** an authenticated user makes a DELETE request to `/v1/auth/api-keys/:id` with a key ID that belongs to the user
- **THEN** the gateway SHALL call `auth-service.DeleteAPIKey` with the key ID
- **AND** return HTTP 200 with success message

#### Scenario: Delete another user's API key
- **WHEN** an authenticated user makes a DELETE request to `/v1/auth/api-keys/:id` with a key ID that does NOT belong to them
- **THEN** the gateway SHALL return HTTP 403 with error message "forbidden"

### Requirement: User self-service API key endpoint authentication
The `/v1/auth/api-keys` endpoints SHALL use JWT authentication (same as `/admin/*` endpoints), NOT API key authentication.

#### Scenario: JWT required
- **WHEN** a request is made to `/v1/auth/api-keys` with an API key in the `Authorization` header
- **THEN** the gateway SHALL return HTTP 401 (these endpoints require JWT from login, not API key)

#### Scenario: Valid JWT accepted
- **WHEN** a request is made to `/v1/auth/api-keys` with a valid JWT in the `Authorization: Bearer <token>` header or `auth_token` cookie
- **THEN** the request SHALL proceed and the user ID from the JWT SHALL be used for the operation

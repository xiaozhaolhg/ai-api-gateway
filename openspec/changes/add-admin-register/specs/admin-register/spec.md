## ADDED Requirements

### Requirement: User registration in auth-service
The auth-service SHALL implement user registration and expose a `Register` RPC that accepts username or email.

#### Scenario: Successful registration with email
- **WHEN** client calls Register with valid email, name, and password
- **THEN** system creates user with hashed password, returns user info

#### Scenario: Successful registration with username
- **WHEN** client calls Register with valid username, name, and password
- **THEN** system creates user with hashed password, returns user info

#### Scenario: Duplicate email or username
- **WHEN** client calls Register with email or username that already exists
- **THEN** system returns 409 conflict error

#### Scenario: Invalid password
- **WHEN** client calls Register with weak password (< 8 chars)
- **THEN** system returns 400 validation error

### Requirement: Gateway proxies register to auth-service
The gateway-service SHALL proxy `POST /admin/register` requests to auth-service Register RPC.

#### Scenario: Successful registration
- **WHEN** client POSTs to /admin/register with valid credentials
- **THEN** gateway proxies to auth-service, returns JWT token

### Requirement: Registration form in admin-ui
The admin-ui SHALL include a registration page accessible from login page.

#### Scenario: User registers successfully
- **WHEN** user fills form and submits
- **THEN** user is created, auto-logged in, and redirected to dashboard
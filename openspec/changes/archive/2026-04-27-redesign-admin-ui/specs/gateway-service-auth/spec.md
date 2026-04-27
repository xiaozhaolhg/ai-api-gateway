## MODIFIED Requirements

### Requirement: Admin login endpoint
The gateway-service SHALL proxy admin login requests to auth-service.

#### Scenario: Login endpoint
- **WHEN** POST `/admin/login` is called with email/password
- **THEN** call auth-service Login RPC and set HTTP-only JWT cookie on success

#### Scenario: Login failure
- **WHEN** auth-service rejects credentials
- **THEN** return 401 Unauthorized with error message

#### Scenario: Cookie management
- **WHEN** login is successful
- **THEN** set JWT cookie with secure, HTTP-only, and /admin path restrictions

### Requirement: Auth middleware for admin routes
The gateway-service SHALL validate admin UI sessions.

#### Scenario: Admin route protection
- **WHEN** admin UI routes are accessed
- **THEN** validate JWT cookie and set user context for downstream services

#### Scenario: Session validation
- **WHEN** JWT is expired or invalid
- **THEN** return 401 Unauthorized to trigger UI redirect to login

#### Scenario: User context propagation
- **WHEN** admin UI makes API calls
- **THEN** include user ID and role in request context for authorization

### Requirement: Logout endpoint
The gateway-service SHALL handle admin logout requests.

#### Scenario: Logout request
- **WHEN** POST `/admin/logout` is called
- **THEN** clear the auth cookie and return success response

#### Scenario: Cookie clearing
- **WHEN** logging out
- **THEN** set cookie with expired date to ensure removal

#### Scenario: Session invalidation
- **WHEN** user logs out
- **THEN** JWT becomes invalid and cannot be reused

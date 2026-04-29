## ADDED Requirements

### Requirement: Admin UI authentication flow
The admin UI SHALL authenticate users via email/password login and maintain session state.

#### Scenario: User login
- **WHEN** user navigates to `/admin/login` and submits email/password
- **THEN** gateway-service validates credentials via auth-service and sets HTTP-only JWT cookie

#### Scenario: Session validation
- **WHEN** user accesses protected admin routes
- **THEN** UI checks auth context and redirects to login if not authenticated

#### Scenario: User logout
- **WHEN** user clicks logout
- **THEN** UI clears auth context and gateway-service clears the auth cookie

### Requirement: Auth context and route guards
The admin UI SHALL protect routes based on authentication status.

#### Scenario: Protected route access
- **WHEN** unauthenticated user accesses `/admin/*` (except `/admin/login`)
- **THEN** redirect to `/admin/login`

#### Scenario: Authenticated root access
- **WHEN** authenticated user accesses `/admin`
- **THEN** redirect to `/admin/dashboard`

#### Scenario: Auth context persistence
- **WHEN** page reloads or browser closes/reopens
- **THEN** auth context is restored from cookie if valid

### Requirement: Auth service login endpoint
The auth-service SHALL support email/password authentication for admin UI users.

#### Scenario: Login request
- **WHEN** gateway-service calls auth-service Login RPC with email/password
- **THEN** auth-service validates credentials and returns JWT with user identity

#### Scenario: Invalid credentials
- **WHEN** login credentials are invalid
- **THEN** auth-service returns error and gateway-service forwards to UI

#### Scenario: Password storage
- **WHEN** user is created or password updated
- **THEN** auth-service stores password hash (not plaintext)

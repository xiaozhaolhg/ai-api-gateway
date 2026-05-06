## MODIFIED Requirements

### Requirement: Session management
The admin-ui SHALL manage authentication session state using React Context, with proactive expired session detection and automatic redirect to login.

#### Scenario: Session persistence across page reload
- **WHEN** the user reloads the page while having a valid token in localStorage
- **THEN** the admin-ui SHALL restore the auth context from localStorage
- **AND** allow access to protected routes without re-login.

#### Scenario: 401 response triggers logout
- **WHEN** an API call returns 401 Unauthorized
- **THEN** the `UnifiedAPIClient` SHALL invoke the `onUnauthorized` callback
- **AND** the callback SHALL call `logout()` which clears the token and localStorage
- **AND** `ProtectedRoute` SHALL detect `isAuthenticated=false` and redirect to `/login`

#### Scenario: Token expiry check triggers logout
- **WHEN** the periodic token expiry check (every 60s) detects the JWT `exp` claim is past (with 30s early expiry for clock skew)
- **THEN** the system SHALL display a warning message "Session expired. Redirecting to login..."
- **AND** call `logout()` which clears the token and localStorage
- **AND** `ProtectedRoute` SHALL detect `isAuthenticated=false` and redirect to `/login`

#### Scenario: Valid token continues
- **WHEN** the JWT `exp` claim is still valid (with 30s early expiry buffer)
- **THEN** the system SHALL continue normal operation without redirect

#### Scenario: Invalid token format
- **WHEN** the token format is invalid (cannot decode `exp` claim)
- **THEN** the system SHALL skip the client-side expiry check (not force logout)
- **AND** rely on the 401 interceptor to reject truly invalid tokens server-side

#### Scenario: Logout
- **WHEN** the user clicks the logout button
- **THEN** the admin-ui SHALL call `POST /admin/auth/logout`
- **AND** clear the auth context and localStorage
- **AND** redirect to the login page.

### Requirement: Typed API client
The admin-ui SHALL include a typed HTTP client for all gateway-service admin API endpoints with auth header injection, centralized error handling, and 401 response interceptor.

#### Scenario: API client matches admin endpoints
- **WHEN** the API client module is inspected
- **THEN** it SHALL have typed methods for all /admin/* endpoints defined in gateway-service API contracts
- **AND** it SHALL include an Authorization header with JWT token on every request
- **AND** it SHALL display error feedback via antd message API on request failure

#### Scenario: 401 response interceptor
- **WHEN** an API call returns HTTP 401
- **THEN** the `request()` method SHALL invoke `this.onUnauthorized?.()` callback
- **AND** display an error message "Session expired. Please login again."
- **AND** throw an error with the message "Unauthorized"

#### Scenario: Callback not set
- **WHEN** `onUnauthorized` is `undefined`
- **THEN** the system SHALL NOT throw an error (graceful handling)

#### Scenario: Normal error handling preserved
- **WHEN** an API call returns HTTP 500 (or any non-401 error)
- **THEN** the system SHALL display the error message via `message.error()`
- **AND** throw the error (without invoking `onUnauthorized`)

### Requirement: Auth context and route guards
The admin UI SHALL protect routes based on authentication status, with proactive session expiry detection.

#### Scenario: Protected route access
- **WHEN** unauthenticated user accesses `/admin/*` (except `/admin/login`)
- **THEN** redirect to `/admin/login`.

#### Scenario: Authenticated root access
- **WHEN** authenticated user accesses `/admin`
- **THEN** redirect to `/admin/dashboard`.

#### Scenario: Auth context persistence
- **WHEN** page reloads or browser closes/reopens
- **THEN** auth context is restored from localStorage if valid.

#### Scenario: Periodic expiry check
- **WHEN** `AuthProvider` mounts with a valid token
- **THEN** the system SHALL set up a `setInterval` that calls `checkTokenExpiry` every 60000 milliseconds (60 seconds)
- **AND** clear the interval on unmount or when token changes to null

#### Scenario: JWT payload decoding
- **WHEN** checking token expiry
- **THEN** the system SHALL split the token by `.` and take the second segment (payload)
- **AND** decode it using `atob()` (base64 decode)
- **AND** parse the result as JSON to get the payload object

#### Scenario: Early expiry buffer
- **WHEN** comparing the current time to the token's `exp` claim
- **THEN** the system SHALL treat the token as expired 30000 milliseconds (30 seconds) before the actual `exp` time
- **AND** this buffer accounts for clock skew between browser and server

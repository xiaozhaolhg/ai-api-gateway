## Purpose

Authentication and session management for the admin UI, including login, logout, registration, and protected routes.

## Requirements

### Requirement: Login page
The admin-ui SHALL provide a login page at `/login` with username and password fields.

#### Scenario: Successful login
- **WHEN** the user submits valid credentials on the login page
- **THEN** the admin-ui SHALL call `POST /admin/auth/login` via gateway-service
- **AND** store the returned JWT token and user info in auth context
- **AND** redirect to the dashboard page (`/`)

#### Scenario: Failed login
- **WHEN** the user submits invalid credentials
- **THEN** the admin-ui SHALL display an error message using antd `message.error`
- **AND** remain on the login page with fields preserved

### Requirement: Registration page
The admin-ui SHALL provide a registration page at `/register` accessible from login page.

#### Scenario: Successful registration
- **WHEN** the user submits valid registration form (name, email/username, password)
- **THEN** the admin-ui SHALL call `POST /admin/auth/register` via gateway-service
- **AND** store the returned JWT token and user info in auth context
- **AND** redirect to the dashboard page (`/`)

#### Scenario: Registration failure
- **WHEN** the user submits with duplicate email/username or weak password
- **THEN** the admin-ui SHALL display an error message
- **AND** remain on the registration page with fields preserved

### Requirement: Session management
The admin-ui SHALL manage authentication session state using React Context.

#### Scenario: Session persistence across page reload
- **WHEN** the user reloads the page while having a valid token in localStorage
- **THEN** the admin-ui SHALL restore the auth context from localStorage
- **AND** allow access to protected routes without re-login

#### Scenario: Session expiry
- **WHEN** the stored token is expired or invalid
- **THEN** the admin-ui SHALL clear the auth context
- **AND** redirect to the login page

#### Scenario: Logout
- **WHEN** the user clicks the logout button
- **THEN** the admin-ui SHALL call `POST /admin/auth/logout`
- **AND** clear the auth context and localStorage
- **AND** redirect to the login page

### Requirement: Protected routes
The admin-ui SHALL protect all admin routes behind authentication.

#### Scenario: Unauthenticated access to protected route
- **WHEN** an unauthenticated user navigates to any route other than `/login`
- **THEN** the admin-ui SHALL redirect to `/login`
- **AND** preserve the intended destination URL for post-login redirect

#### Scenario: Authenticated access to protected route
- **WHEN** an authenticated user navigates to a protected route
- **THEN** the admin-ui SHALL render the requested page

### Requirement: Auth header injection
The admin-ui API client SHALL include the JWT token in all API requests.

#### Scenario: API request with auth header
- **WHEN** the API client makes any request to gateway-service
- **THEN** it SHALL include an `Authorization: Bearer <token>` header
- **AND** if the token is missing, redirect to login

### Requirement: Current user info
The admin-ui SHALL display the current user's identity in the header.

#### Scenario: User menu display
- **WHEN** the user is authenticated
- **THEN** the admin-ui SHALL display the user's name and role in the header
- **AND** provide a logout action in the user menu

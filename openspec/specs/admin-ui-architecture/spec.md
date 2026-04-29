## Purpose

React single-page application for provider management, user management, API key management, and usage dashboards.

## Requirements

### Requirement: React SPA with Vite
The admin-ui SHALL be a React single-page application built with Vite, TypeScript, Ant Design 6.3.6, and TailwindCSS for custom overrides.

#### Scenario: Development server runs
- **WHEN** `npm run dev` is executed
- **THEN** the Vite dev server SHALL start with a proxy to gateway-service for API calls

#### Scenario: Production build
- **WHEN** `npm run build` is executed
- **THEN** static assets SHALL be generated in the dist/ directory

### Requirement: Provider management page
The admin-ui SHALL provide a page for managing LLM providers (list, add, edit, remove) using Ant Design components.

#### Scenario: List providers
- **WHEN** the provider management page is loaded
- **THEN** it SHALL display all providers in an antd Table with name, type, base URL, models, status, and actions columns
- **AND** show an antd Spin while loading and an antd Empty when no providers exist

#### Scenario: Add provider
- **WHEN** the user fills in the provider form and submits
- **THEN** the admin-ui SHALL call POST /admin/providers via gateway-service
- **AND** the new provider SHALL appear in the table
- **AND** an antd success message SHALL be displayed

#### Scenario: Edit provider
- **WHEN** the user clicks edit on a provider
- **THEN** the admin-ui SHALL open an antd Modal with the provider's current values
- **AND** on submit, call PUT /admin/providers/:id
- **AND** update the provider in the table

#### Scenario: Delete provider
- **WHEN** the user confirms deletion of a provider via antd Popconfirm
- **THEN** the admin-ui SHALL call DELETE /admin/providers/:id
- **AND** the provider SHALL be removed from the table

### Requirement: User and API key management page
The admin-ui SHALL provide pages for managing users and API keys using Ant Design components.

#### Scenario: Create user
- **WHEN** the user fills in name, email, and role and submits
- **THEN** the admin-ui SHALL call POST /admin/users
- **AND** the new user SHALL appear in the user list

#### Scenario: Edit user
- **WHEN** the user clicks edit on a user
- **THEN** the admin-ui SHALL open an antd Modal with the user's current values
- **AND** on submit, call PUT /admin/users/:id

#### Scenario: Issue API key
- **WHEN** the user clicks "Issue API Key" for a user
- **THEN** the admin-ui SHALL call POST /admin/api-keys
- **AND** the generated key SHALL be displayed in an antd Alert with a copy button and a warning that it cannot be retrieved again

#### Scenario: Revoke API key
- **WHEN** the user confirms revocation via antd Popconfirm
- **THEN** the admin-ui SHALL call DELETE /admin/api-keys/:id
- **AND** the key SHALL be removed from the list

#### Scenario: API key user selector populated from users API
- **WHEN** the API keys page is loaded
- **THEN** the user selector SHALL be populated from GET /admin/users
- **AND** display user names as options (NOT hardcoded values)

### Requirement: Usage dashboard page
The admin-ui SHALL provide a dashboard showing token usage statistics using Ant Design components.

#### Scenario: View usage summary
- **WHEN** the usage dashboard is loaded
- **THEN** it SHALL display total prompt/completion tokens, cost, and request count using antd Statistic cards
- **AND** it SHALL support filtering by user, model, provider, and date range using antd Form and DatePicker

### Requirement: Health dashboard page
The admin-ui SHALL provide a page showing provider health status using real API data.

#### Scenario: View provider health
- **WHEN** the health dashboard is loaded
- **THEN** it SHALL call GET /admin/health via gateway-service
- **AND** display each provider's status (healthy/degraded/down) using antd Badge, latency, and error rate

### Requirement: Typed API client
The admin-ui SHALL include a typed HTTP client for all gateway-service admin API endpoints with auth header injection and centralized error handling.

#### Scenario: API client matches admin endpoints
- **WHEN** the API client module is inspected
- **THEN** it SHALL have typed methods for all /admin/* endpoints defined in gateway-service API contracts
- **AND** it SHALL include an Authorization header with JWT token on every request
- **AND** it SHALL display error feedback via antd message API on request failure

### Requirement: Ant Design component system
The admin-ui SHALL use Ant Design 6.3.6 as the primary component library for all UI elements.

#### Scenario: Component usage
- **WHEN** any page is rendered
- **THEN** it SHALL use antd components for tables (Table), forms (Form, Input, Select), modals (Modal), confirmations (Popconfirm), status badges (Badge, Tag), loading states (Spin), empty states (Empty), messages (message), and notifications (notification)

### Requirement: Shared AppShell layout
The admin-ui SHALL use an Ant Design Layout with Sider for navigation across all pages.

#### Scenario: AppShell renders sidebar and content
- **WHEN** any authenticated page is rendered
- **THEN** the AppShell SHALL display an antd Layout with a collapsible Sider containing a grouped Menu
- **AND** a Content area with Breadcrumb navigation
- **AND** a Header with user menu and language switcher

### Requirement: Error and empty states
The admin-ui SHALL display proper feedback for loading, error, and empty states on all pages.

#### Scenario: Loading state
- **WHEN** data is being fetched
- **THEN** the page SHALL display an antd Spin component

#### Scenario: Error state
- **WHEN** an API request fails
- **THEN** the page SHALL display an antd Alert or message.error with the error description

#### Scenario: Empty state
- **WHEN** a list has no items
- **THEN** the page SHALL display an antd Empty component with a description and action link

---

## Authentication & Session Management

### Requirement: Admin UI authentication flow
The admin UI SHALL authenticate users via email/password login and maintain session state.

#### Scenario: User login
- **WHEN** user navigates to `/admin/login` and submits email/password
- **THEN** gateway-service validates credentials via auth-service and sets HTTP-only JWT cookie.

#### Scenario: Session validation
- **WHEN** user accesses protected admin routes
- **THEN** UI checks auth context and redirects to login if not authenticated.

#### Scenario: User logout
- **WHEN** user clicks logout
- **THEN** UI clears auth context and gateway-service clears the auth cookie.

### Requirement: Auth context and route guards
The admin UI SHALL protect routes based on authentication status.

#### Scenario: Protected route access
- **WHEN** unauthenticated user accesses `/admin/*` (except `/admin/login`)
- **THEN** redirect to `/admin/login`.

#### Scenario: Authenticated root access
- **WHEN** authenticated user accesses `/admin`
- **THEN** redirect to `/admin/dashboard`.

#### Scenario: Auth context persistence
- **WHEN** page reloads or browser closes/reopens
- **THEN** auth context is restored from cookie if valid.

### Requirement: Auth service login endpoint
The auth-service SHALL support email/password authentication for admin UI users.

#### Scenario: Login request
- **WHEN** gateway-service calls auth-service Login RPC with email/password
- **THEN** auth-service validates credentials and returns JWT with user identity.

#### Scenario: Invalid credentials
- **WHEN** login credentials are invalid
- **THEN** auth-service returns error and gateway-service forwards to UI.

#### Scenario: Password storage
- **WHEN** user is created or password updated
- **THEN** auth-service stores password hash (not plaintext).

### Requirement: Login page
The admin-ui SHALL provide a login page at `/login` with username and password fields.

#### Scenario: Successful login
- **WHEN** the user submits valid credentials on the login page
- **THEN** the admin-ui SHALL call `POST /admin/auth/login` via gateway-service
- **AND** store the returned JWT token and user info in auth context
- **AND** redirect to the dashboard page (`/`).

#### Scenario: Failed login
- **WHEN** the user submits invalid credentials
- **THEN** the admin-ui SHALL display an error message using antd `message.error`
- **AND** remain on the login page with fields preserved.

### Requirement: Registration page
The admin-ui SHALL provide a registration page at `/register` accessible from login page.

#### Scenario: Successful registration
- **WHEN** the user submits valid registration form (name, email/username, password)
- **THEN** the admin-ui SHALL call `POST /admin/auth/register` via gateway-service
- **AND** store the returned JWT token and user info in auth context
- **AND** redirect to the dashboard page (`/`).

#### Scenario: Registration failure
- **WHEN** the user submits with duplicate email/username or weak password
- **THEN** the admin-ui SHALL display an error message
- **AND** remain on the registration page with fields preserved.

### Requirement: Session management
The admin-ui SHALL manage authentication session state using React Context.

#### Scenario: Session persistence across page reload
- **WHEN** the user reloads the page while having a valid token in localStorage
- **THEN** the admin-ui SHALL restore the auth context from localStorage
- **AND** allow access to protected routes without re-login.

#### Scenario: Session expiry
- **WHEN** the stored token is expired or invalid
- **THEN** the admin-ui SHALL clear the auth context
- **AND** redirect to the login page.

#### Scenario: Logout
- **WHEN** the user clicks the logout button
- **THEN** the admin-ui SHALL call `POST /admin/auth/logout`
- **AND** clear the auth context and localStorage
- **AND** redirect to the login page.

### Requirement: Protected routes
The admin-ui SHALL protect all admin routes behind authentication.

#### Scenario: Unauthenticated access to protected route
- **WHEN** an unauthenticated user navigates to any route other than `/login`
- **THEN** the admin-ui SHALL redirect to `/login`
- **AND** preserve the intended destination URL for post-login redirect.

#### Scenario: Authenticated access to protected route
- **WHEN** an authenticated user navigates to a protected route
- **THEN** the admin-ui SHALL render the requested page.

### Requirement: Auth header injection
The admin-ui API client SHALL include the JWT token in all API requests.

#### Scenario: API request with auth header
- **WHEN** the API client makes any request to gateway-service
- **THEN** it SHALL include an `Authorization: Bearer <token>` header
- **AND** if the token is missing, redirect to login.

### Requirement: Current user info
The admin-ui SHALL display the current user's identity in the header.

#### Scenario: User menu display
- **WHEN** the user is authenticated
- **THEN** the admin-ui SHALL display the user's name and role in the header
- **AND** provide a logout action in the user menu.


---

## Layout & Navigation

### Requirement: Modern component library integration
The admin UI SHALL use antd 6 components with antd icons.

#### Scenario: Component usage
- **WHEN** building UI elements
- **THEN** use antd Button, Input, Table, Card, Modal, Form components.

#### Scenario: Icon integration
- **WHEN** displaying icons
- **THEN** use antd icons consistent with antd design system.

#### Scenario: Styling consistency
- **WHEN** applying styles
- **THEN** follow antd design tokens and Tailwind CSS classes.

### Requirement: Collapsible sidebar navigation
The admin UI SHALL have a collapsible sidebar with icon-only mode.

#### Scenario: Sidebar toggle
- **WHEN** user clicks collapse button
- **THEN** sidebar collapses to icon-only view, expanding main content area.

#### Scenario: Active state indication
- **WHEN** navigation tab is active
- **THEN** highlight with accent color and background.

#### Scenario: Responsive behavior
- **WHEN** screen width is limited
- **THEN** sidebar automatically collapses to icon-only mode.

### Requirement: Form handling and validation
The admin UI SHALL use antd Form components with react-hook-form for validation.

#### Scenario: Form submission
- **WHEN** user submits create/edit forms
- **THEN** validate with react-hook-form and display errors inline.

#### Scenario: Form reset
- **WHEN** form is cancelled or successfully submitted
- **THEN** reset form state and close modal/drawer.

#### Scenario: Field validation
- **WHEN** user enters invalid data
- **THEN** show real-time validation feedback with proper error messages.

### Requirement: Data fetching and caching
The admin UI SHALL use TanStack Query for API data management.

#### Scenario: Data loading
- **WHEN** page loads
- **THEN** TanStack Query fetches data with loading states.

#### Scenario: Data caching
- **WHEN** navigating between pages
- **THEN** cached data is used when fresh, with background refetch.

#### Scenario: Error handling
- **WHEN** API requests fail
- **THEN** TanStack Query provides error states with retry options.


---

## Dashboard

### Requirement: Admin dashboard overview
The admin UI SHALL display a dashboard with key metrics and system overview.

#### Scenario: Dashboard load
- **WHEN** authenticated user accesses `/admin/dashboard`
- **THEN** display summary cards for providers, users, API keys, and recent usage.

#### Scenario: Metrics display
- **WHEN** dashboard renders
- **THEN** show total counts, recent activity, and cost summary with visual indicators.

#### Scenario: Quick actions
- **WHEN** dashboard loads
- **THEN** provide quick access buttons to add provider, create user, generate API key.

### Requirement: Dashboard data aggregation
The dashboard SHALL aggregate data from multiple services.

#### Scenario: Provider metrics
- **WHEN** dashboard loads
- **THEN** fetch provider count and health status from gateway-service.

#### Scenario: Usage statistics
- **WHEN** dashboard loads
- **THEN** fetch token usage and cost data from billing-service.

#### Scenario: User activity
- **WHEN** dashboard loads
- **THEN** fetch user count and recent API key creation from auth-service.

### Requirement: Dashboard overview page
The admin-ui SHALL provide a dashboard page at `/` as the default landing route.

#### Scenario: Dashboard loads with summary cards
- **WHEN** the dashboard page is loaded
- **THEN** it SHALL display summary cards for: total users, active providers, current month spend, and active alerts
- **AND** each card SHALL use antd `Statistic` component with an icon.

#### Scenario: Dashboard quick navigation
- **WHEN** the dashboard page is loaded
- **THEN** it SHALL display quick-action links to: add provider, create user, issue API key, view alerts
- **AND** each link SHALL navigate to the corresponding page.

#### Scenario: Dashboard data fetching
- **WHEN** the dashboard page is loaded
- **THEN** it SHALL fetch summary data from `GET /admin/users` (count), `GET /admin/providers` (count), `GET /admin/usage` (aggregate), `GET /admin/alerts` (active count)
- **AND** display loading spinners while data is being fetched.



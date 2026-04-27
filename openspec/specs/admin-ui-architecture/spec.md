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


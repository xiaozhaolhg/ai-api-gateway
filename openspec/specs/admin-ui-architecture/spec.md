## Purpose

React single-page application for provider management, user management, API key management, and usage dashboards.
## Requirements
### Requirement: React SPA with Vite
The admin-ui SHALL be a React single-page application built with Vite, TypeScript, and TailwindCSS.

#### Scenario: Development server runs
- **WHEN** `npm run dev` is executed
- **THEN** the Vite dev server SHALL start with a proxy to gateway-service for API calls

#### Scenario: Production build
- **WHEN** `npm run build` is executed
- **THEN** static assets SHALL be generated in the dist/ directory

### Requirement: Provider management page
The admin-ui SHALL provide a page for managing LLM providers (list, add, edit, remove).

#### Scenario: List providers
- **WHEN** the provider management page is loaded
- **THEN** it SHALL display all providers with name, type, status, and available models

#### Scenario: Add provider
- **WHEN** the user fills in the provider form and submits
- **THEN** the admin-ui SHALL call POST /admin/providers via gateway-service
- **AND** the new provider SHALL appear in the list

#### Scenario: Delete provider
- **WHEN** the user confirms deletion of a provider
- **THEN** the admin-ui SHALL call DELETE /admin/providers/:id
- **AND** the provider SHALL be removed from the list

### Requirement: User and API key management page
The admin-ui SHALL provide pages for managing users and API keys.

#### Scenario: Create user
- **WHEN** the user fills in name, email, and role and submits
- **THEN** the admin-ui SHALL call POST /admin/users
- **AND** the new user SHALL appear in the user list

#### Scenario: Issue API key
- **WHEN** the user clicks "Issue API Key" for a user
- **THEN** the admin-ui SHALL call POST /admin/api-keys
- **AND** the generated key SHALL be displayed once with a warning that it cannot be retrieved again

#### Scenario: Revoke API key
- **WHEN** the user clicks "Revoke" on an API key
- **THEN** the admin-ui SHALL call DELETE /admin/api-keys/:id
- **AND** the key SHALL be removed from the list

### Requirement: Usage dashboard page
The admin-ui SHALL provide a dashboard showing token usage statistics.

#### Scenario: View usage summary
- **WHEN** the usage dashboard is loaded
- **THEN** it SHALL display total prompt/completion tokens, cost, and request count
- **AND** it SHALL support filtering by user, model, provider, and date range

### Requirement: Health dashboard page
The admin-ui SHALL provide a page showing provider health status.

#### Scenario: View provider health
- **WHEN** the health dashboard is loaded
- **THEN** it SHALL display each provider's status (healthy/degraded/down), latency, and error rate

### Requirement: Typed API client
The admin-ui SHALL include a typed HTTP client for all gateway-service admin API endpoints.

#### Scenario: API client matches admin endpoints
- **WHEN** the API client module is inspected
- **THEN** it SHALL have typed methods for all /admin/* endpoints defined in gateway-service API contracts


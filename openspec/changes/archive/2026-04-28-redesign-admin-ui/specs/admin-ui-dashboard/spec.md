## ADDED Requirements

### Requirement: Admin dashboard overview
The admin UI SHALL display a dashboard with key metrics and system overview.

#### Scenario: Dashboard load
- **WHEN** authenticated user accesses `/admin/dashboard`
- **THEN** display summary cards for providers, users, API keys, and recent usage

#### Scenario: Metrics display
- **WHEN** dashboard renders
- **THEN** show total counts, recent activity, and cost summary with visual indicators

#### Scenario: Quick actions
- **WHEN** dashboard loads
- **THEN** provide quick access buttons to add provider, create user, generate API key

### Requirement: Dashboard data aggregation
The dashboard SHALL aggregate data from multiple services.

#### Scenario: Provider metrics
- **WHEN** dashboard loads
- **THEN** fetch provider count and health status from gateway-service

#### Scenario: Usage statistics
- **WHEN** dashboard loads
- **THEN** fetch token usage and cost data from billing-service

#### Scenario: User activity
- **WHEN** dashboard loads
- **THEN** fetch user count and recent API key creation from auth-service

## Purpose

Dashboard overview page for the admin UI, providing summary statistics and quick navigation.

## Requirements

### Requirement: Dashboard overview page
The admin-ui SHALL provide a dashboard page at `/` as the default landing route.

#### Scenario: Dashboard loads with summary cards
- **WHEN** the dashboard page is loaded
- **THEN** it SHALL display summary cards for: total users, active providers, current month spend, and active alerts
- **AND** each card SHALL use antd `Statistic` component with an icon

#### Scenario: Dashboard quick navigation
- **WHEN** the dashboard page is loaded
- **THEN** it SHALL display quick-action links to: add provider, create user, issue API key, view alerts
- **AND** each link SHALL navigate to the corresponding page

#### Scenario: Dashboard data fetching
- **WHEN** the dashboard page is loaded
- **THEN** it SHALL fetch summary data from `GET /admin/users` (count), `GET /admin/providers` (count), `GET /admin/usage` (aggregate), `GET /admin/alerts` (active count)
- **AND** display loading spinners while data is being fetched

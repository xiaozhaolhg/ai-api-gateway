## ADDED Requirements

### Requirement: Role-based navigation filtering
The admin UI SHALL filter navigation options based on user role.

#### Scenario: Admin role navigation
- **WHEN** user with 'admin' role loads
- **THEN** show all navigation tabs: Dashboard, Providers, Users, API Keys, Usage, Health, Settings

#### Scenario: User role navigation
- **WHEN** user with 'user' role loads
- **THEN** show limited tabs: Dashboard, API Keys (own), Usage (own), Health, Settings

#### Scenario: Viewer role navigation
- **WHEN** user with 'viewer' role loads
- **THEN** show read-only tabs: Dashboard, Usage (own), Health

### Requirement: Role-based access control
The admin UI SHALL enforce access restrictions at the page level.

#### Scenario: Admin full access
- **WHEN** admin role user accesses any admin page
- **THEN** allow full CRUD operations

#### Scenario: User limited access
- **WHEN** user role accesses Providers or Users pages
- **THEN** redirect to dashboard with access denied message

#### Scenario: Viewer read-only access
- **WHEN** viewer role attempts create/edit/delete operations
- **THEN** disable buttons and show read-only state

### Requirement: Role-based data filtering
The admin UI SHALL filter data based on user role and ownership.

#### Scenario: User API key access
- **WHEN** user role views API Keys page
- **THEN** only show API keys belonging to that user

#### Scenario: User usage data
- **WHEN** user role views Usage page
- **THEN** only show usage data for that user

#### Scenario: Viewer usage data
- **WHEN** viewer role views Usage page
- **THEN** only show usage data for that user with read-only controls

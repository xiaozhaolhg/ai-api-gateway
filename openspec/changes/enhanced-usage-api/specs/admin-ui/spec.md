## ADDED Requirements

### Requirement: Enhanced usage page with user and group specific views
The admin-ui SHALL provide enhanced usage analytics with user-specific, group-specific, and export functionality beyond the current basic usage table.

#### Scenario: User-specific usage view
- **WHEN** admin navigates to `/usage/users/:id` or clicks "View Usage" from user management
- **THEN** admin-ui displays usage data filtered for that specific user with charts and detailed breakdowns

#### Scenario: Group-specific usage view  
- **WHEN** admin navigates to `/usage/groups/:id` or clicks "View Usage" from group management
- **THEN** admin-ui displays usage data filtered for that specific group with member breakdowns

#### Scenario: Usage dashboard with charts
- **WHEN** admin views the usage page
- **THEN** admin-ui displays token consumption charts, cost trends, and usage patterns over time

#### Scenario: Advanced filtering options
- **WHEN** admin uses usage page filters
- **THEN** admin-ui provides filters for user, group, provider, model, date range with real-time updates

### Requirement: Usage data export functionality
The admin-ui SHALL provide export functionality for usage data in multiple formats with configurable filters.

#### Scenario: CSV export with filters
- **WHEN** admin clicks "Export CSV" button with applied filters
- **THEN** admin-ui downloads CSV file with filtered usage data including headers and metadata

#### Scenario: JSON export with filters
- **WHEN** admin clicks "Export JSON" button with applied filters  
- **THEN** admin-ui downloads JSON file with filtered usage data and export metadata

#### Scenario: Export progress indicator
- **WHEN** admin initiates export for large datasets
- **THEN** admin-ui shows progress indicator and estimated completion time

#### Scenario: Export format selection
- **WHEN** admin clicks export dropdown
- **THEN** admin-ui shows options for CSV, JSON, and Excel formats with size estimates

### Requirement: Usage analytics and visualization
The admin-ui SHALL provide comprehensive usage analytics with interactive charts and visual representations.

#### Scenario: Token consumption chart
- **WHEN** admin views usage dashboard
- **THEN** admin-ui displays line chart showing token consumption over selected time period

#### Scenario: Cost breakdown chart
- **WHEN** admin views usage dashboard
- **THEN** admin-ui displays pie chart showing cost breakdown by provider and model

#### Scenario: User/group comparison charts
- **WHEN** admin selects multiple users or groups
- **THEN** admin-ui displays comparison charts for usage patterns and costs

#### Scenario: Interactive date range selector
- **WHEN** admin adjusts date range on charts
- **THEN** charts update dynamically with loading states and smooth transitions

### Requirement: Enhanced API client for usage endpoints
The admin-ui API client SHALL support the new usage API endpoints with proper typing and error handling.

#### Scenario: Get user-specific usage
- **WHEN** frontend calls `apiClient.getUsageByUser(userId, filters)`
- **THEN** API client calls `/admin/usage/users/:id` endpoint with proper error handling

#### Scenario: Get group-specific usage
- **WHEN** frontend calls `apiClient.getUsageByGroup(groupId, filters)`
- **THEN** API client calls `/admin/usage/groups/:id` endpoint with proper error handling

#### Scenario: Export usage data
- **WHEN** frontend calls `apiClient.exportUsage(filters, format)`
- **THEN** API client calls `/admin/usage/export` endpoint and handles file download

#### Scenario: Type-safe usage data
- **WHEN** usage data is received from API
- **THEN** TypeScript types ensure proper data structure and type safety

### Requirement: Usage page navigation integration
The admin-ui SHALL integrate usage views with existing user and group management pages.

#### Scenario: User management usage links
- **WHEN** admin views user management page
- **THEN** each user row has "View Usage" link that navigates to user-specific usage view

#### Scenario: Group management usage links
- **WHEN** admin views group management page
- **THEN** each group row has "View Usage" link that navigates to group-specific usage view

#### Scenario: Breadcrumb navigation
- **WHEN** admin navigates to user/group specific usage views
- **THEN** breadcrumbs show navigation path: Usage > Users > [User Name]

#### Scenario: Quick access from dashboard
- **WHEN** admin views main dashboard
- **THEN** dashboard provides quick access buttons to top users/groups by usage

### Requirement: Usage data performance optimization
The admin-ui SHALL handle large usage datasets efficiently with pagination and virtualization.

#### Scenario: Virtualized usage table
- **WHEN** usage data contains thousands of records
- **THEN** admin-ui uses virtual scrolling to maintain performance without loading all data

#### Scenario: Lazy loading for charts
- **WHEN** admin views usage analytics
- **THEN** chart data loads progressively with skeleton states during data fetching

#### Scenario: Cached usage data
- **WHEN** admin revisits usage page with same filters
- **THEN** admin-ui uses cached data when available and shows refresh timestamp

#### Scenario: Background data refresh
- **WHEN** admin has usage page open for extended periods
- **THEN** admin-ui periodically refreshes data in background with user notifications

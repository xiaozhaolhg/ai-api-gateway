## ADDED Requirements

### Requirement: Usage API HTTP endpoints
The auth-service SHALL provide HTTP endpoints for accessing usage data through the gateway-service, with support for user-specific, group-specific, and export functionality.

#### Scenario: Get usage by user ID
- **WHEN** gateway-service calls `GET /admin/usage/users/:id` with user_id and optional query parameters (start_date, end_date, page, page_size)
- **THEN** auth-service calls billing-service GetUsageByUser gRPC and returns paginated usage records for that user

#### Scenario: Get usage by group ID
- **WHEN** gateway-service calls `GET /admin/usage/groups/:id` with group_id and optional query parameters (start_date, end_date, page, page_size)
- **THEN** auth-service calls billing-service GetUsageByGroup gRPC and returns paginated usage records for that group

#### Scenario: Export usage data as CSV
- **WHEN** gateway-service calls `GET /admin/usage/export?format=csv` with optional filters (user_id, group_id, start_date, end_date)
- **THEN** auth-service calls billing-service for filtered usage data and returns CSV-formatted response with appropriate headers

#### Scenario: Export usage data as JSON
- **WHEN** gateway-service calls `GET /admin/usage/export?format=json` with optional filters (user_id, group_id, start_date, end_date)
- **THEN** auth-service calls billing-service for filtered usage data and returns JSON-formatted response

#### Scenario: Usage API with date range filtering
- **WHEN** usage endpoints are called with start_date and end_date parameters
- **THEN** auth-service validates date format and passes filters to billing-service gRPC calls

#### Scenario: Usage API with pagination
- **WHEN** usage endpoints are called with page and page_size parameters
- **THEN** auth-service validates pagination parameters and returns paginated results with total count

### Requirement: Usage API authentication and authorization
The auth-service SHALL enforce authentication and authorization for all usage API endpoints using existing admin middleware.

#### Scenario: Admin authentication required
- **WHEN** any usage endpoint is called without valid admin authentication
- **THEN** auth-service returns 401 Unauthorized error

#### Scenario: Admin authorization required
- **WHEN** a non-admin user attempts to access usage endpoints
- **THEN** auth-service returns 403 Forbidden error

#### Scenario: User can access own usage data
- **WHEN** a regular user calls `GET /admin/usage/users/:id` with their own user_id
- **THEN** auth-service allows access to their own usage data (future enhancement)

### Requirement: Usage API error handling
The auth-service SHALL provide comprehensive error handling for usage API endpoints with appropriate HTTP status codes and error messages.

#### Scenario: Invalid user ID
- **WHEN** `GET /admin/usage/users/:id` is called with non-existent user_id
- **THEN** auth-service returns 404 Not Found with error message "User not found"

#### Scenario: Invalid group ID
- **WHEN** `GET /admin/usage/groups/:id` is called with non-existent group_id
- **THEN** auth-service returns 404 Not Found with error message "Group not found"

#### Scenario: Invalid date format
- **WHEN** usage endpoints are called with invalid start_date or end_date format
- **THEN** auth-service returns 400 Bad Request with error message "Invalid date format"

#### Scenario: Billing service unavailable
- **WHEN** billing-service gRPC call fails or times out
- **THEN** auth-service returns 503 Service Unavailable with error message "Usage data temporarily unavailable"

#### Scenario: Export format not supported
- **WHEN** `GET /admin/usage/export` is called with unsupported format parameter
- **THEN** auth-service returns 400 Bad Request with error message "Unsupported export format"

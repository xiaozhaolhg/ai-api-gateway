## MODIFIED Requirements

### Requirement: Admin API
The gateway-service SHALL provide admin HTTP endpoints for all management operations.

#### Scenario: Admin endpoints cover all service entities
- **WHEN** the gateway-service admin API is inspected
- **THEN** it SHALL include endpoints for: providers, users, API keys, usage, authentication, routing rules, groups, permissions, budgets, pricing rules, alert rules, alerts, and health

## ADDED Requirements

### Requirement: Admin authentication endpoints
The gateway-service SHALL provide admin authentication endpoints.

#### Scenario: Login
- **WHEN** a POST request is made to `/admin/auth/login` with username and password
- **THEN** the gateway-service SHALL validate credentials via auth-service
- **AND** return a JWT token and user identity

#### Scenario: Logout
- **WHEN** a POST request is made to `/admin/auth/logout`
- **THEN** the gateway-service SHALL invalidate the session

#### Scenario: Current user
- **WHEN** a GET request is made to `/admin/auth/me` with a valid JWT
- **THEN** the gateway-service SHALL return the current user's identity

### Requirement: Routing rule admin endpoints
The gateway-service SHALL provide admin endpoints for routing rule management.

#### Scenario: List routing rules
- **WHEN** a GET request is made to `/admin/routing-rules`
- **THEN** the gateway-service SHALL return all routing rules from router-service

#### Scenario: Create routing rule
- **WHEN** a POST request is made to `/admin/routing-rules` with rule data
- **THEN** the gateway-service SHALL create the rule via router-service

#### Scenario: Get routing rule
- **WHEN** a GET request is made to `/admin/routing-rules/:id`
- **THEN** the gateway-service SHALL return the specified routing rule

#### Scenario: Update routing rule
- **WHEN** a PUT request is made to `/admin/routing-rules/:id` with updated data
- **THEN** the gateway-service SHALL update the rule via router-service

#### Scenario: Delete routing rule
- **WHEN** a DELETE request is made to `/admin/routing-rules/:id`
- **THEN** the gateway-service SHALL delete the rule via router-service

### Requirement: Group admin endpoints
The gateway-service SHALL provide admin endpoints for group management.

#### Scenario: List groups
- **WHEN** a GET request is made to `/admin/groups`
- **THEN** the gateway-service SHALL return all groups from auth-service

#### Scenario: Create group
- **WHEN** a POST request is made to `/admin/groups` with group data
- **THEN** the gateway-service SHALL create the group via auth-service

#### Scenario: Get group
- **WHEN** a GET request is made to `/admin/groups/:id`
- **THEN** the gateway-service SHALL return the specified group

#### Scenario: Update group
- **WHEN** a PUT request is made to `/admin/groups/:id` with updated data
- **THEN** the gateway-service SHALL update the group via auth-service

#### Scenario: Delete group
- **WHEN** a DELETE request is made to `/admin/groups/:id`
- **THEN** the gateway-service SHALL delete the group via auth-service

#### Scenario: Add group member
- **WHEN** a POST request is made to `/admin/groups/:id/members` with a user ID
- **THEN** the gateway-service SHALL add the user to the group via auth-service

#### Scenario: Remove group member
- **WHEN** a DELETE request is made to `/admin/groups/:id/members/:userId`
- **THEN** the gateway-service SHALL remove the user from the group via auth-service

### Requirement: Permission admin endpoints
The gateway-service SHALL provide admin endpoints for permission management.

#### Scenario: List permissions
- **WHEN** a GET request is made to `/admin/permissions`
- **THEN** the gateway-service SHALL return all permissions from auth-service

#### Scenario: Create permission
- **WHEN** a POST request is made to `/admin/permissions` with permission data
- **THEN** the gateway-service SHALL create the permission via auth-service

#### Scenario: Get permission
- **WHEN** a GET request is made to `/admin/permissions/:id`
- **THEN** the gateway-service SHALL return the specified permission

#### Scenario: Update permission
- **WHEN** a PUT request is made to `/admin/permissions/:id` with updated data
- **THEN** the gateway-service SHALL update the permission via auth-service

#### Scenario: Delete permission
- **WHEN** a DELETE request is made to `/admin/permissions/:id`
- **THEN** the gateway-service SHALL delete the permission via auth-service

### Requirement: Budget admin endpoints
The gateway-service SHALL provide admin endpoints for budget management.

#### Scenario: List budgets
- **WHEN** a GET request is made to `/admin/budgets`
- **THEN** the gateway-service SHALL return all budgets from billing-service

#### Scenario: Create budget
- **WHEN** a POST request is made to `/admin/budgets` with budget data
- **THEN** the gateway-service SHALL create the budget via billing-service

#### Scenario: Get budget
- **WHEN** a GET request is made to `/admin/budgets/:id`
- **THEN** the gateway-service SHALL return the specified budget

#### Scenario: Update budget
- **WHEN** a PUT request is made to `/admin/budgets/:id` with updated data
- **THEN** the gateway-service SHALL update the budget via billing-service

#### Scenario: Delete budget
- **WHEN** a DELETE request is made to `/admin/budgets/:id`
- **THEN** the gateway-service SHALL delete the budget via billing-service

### Requirement: Pricing rule admin endpoints
The gateway-service SHALL provide admin endpoints for pricing rule management.

#### Scenario: List pricing rules
- **WHEN** a GET request is made to `/admin/pricing-rules`
- **THEN** the gateway-service SHALL return all pricing rules from billing-service

#### Scenario: Create pricing rule
- **WHEN** a POST request is made to `/admin/pricing-rules` with rule data
- **THEN** the gateway-service SHALL create the rule via billing-service

#### Scenario: Get pricing rule
- **WHEN** a GET request is made to `/admin/pricing-rules/:id`
- **THEN** the gateway-service SHALL return the specified pricing rule

#### Scenario: Update pricing rule
- **WHEN** a PUT request is made to `/admin/pricing-rules/:id` with updated data
- **THEN** the gateway-service SHALL update the rule via billing-service

#### Scenario: Delete pricing rule
- **WHEN** a DELETE request is made to `/admin/pricing-rules/:id`
- **THEN** the gateway-service SHALL delete the rule via billing-service

### Requirement: Alert rule admin endpoints
The gateway-service SHALL provide admin endpoints for alert rule management.

#### Scenario: List alert rules
- **WHEN** a GET request is made to `/admin/alert-rules`
- **THEN** the gateway-service SHALL return all alert rules from monitor-service

#### Scenario: Create alert rule
- **WHEN** a POST request is made to `/admin/alert-rules` with rule data
- **THEN** the gateway-service SHALL create the rule via monitor-service

#### Scenario: Get alert rule
- **WHEN** a GET request is made to `/admin/alert-rules/:id`
- **THEN** the gateway-service SHALL return the specified alert rule

#### Scenario: Update alert rule
- **WHEN** a PUT request is made to `/admin/alert-rules/:id` with updated data
- **THEN** the gateway-service SHALL update the rule via monitor-service

#### Scenario: Delete alert rule
- **WHEN** a DELETE request is made to `/admin/alert-rules/:id`
- **THEN** the gateway-service SHALL delete the rule via monitor-service

### Requirement: Alert admin endpoints
The gateway-service SHALL provide admin endpoints for alert lifecycle management.

#### Scenario: List active alerts
- **WHEN** a GET request is made to `/admin/alerts`
- **THEN** the gateway-service SHALL return all alerts from monitor-service

#### Scenario: Acknowledge alert
- **WHEN** a PUT request is made to `/admin/alerts/:id/acknowledge`
- **THEN** the gateway-service SHALL acknowledge the alert via monitor-service

### Requirement: Health admin endpoint
The gateway-service SHALL provide an admin endpoint for provider health status.

#### Scenario: Get provider health
- **WHEN** a GET request is made to `/admin/health`
- **THEN** the gateway-service SHALL return provider health status from monitor-service

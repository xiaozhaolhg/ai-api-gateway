# gateway-service API Contracts

> HTTP endpoints provided by gateway-service

## HTTP Endpoints

### OpenAI-Compatible API

| Method | Path | Description |
|---|---|---|
| POST | /v1/chat/completions | Chat completion (non-streaming) |
| POST | /v1/chat/completions?stream=true | Chat completion (streaming) |
| GET | /v1/models | List available models |

### Custom Gateway API

| Method | Path | Description |
|---|---|---|
| GET | /gateway/models | List available models |
| GET | /gateway/health | Health check |

### Admin API

#### Authentication

| Method | Path | Description |
|---|---|---|
| POST | /admin/auth/login | Login with username/password, returns JWT token |
| POST | /admin/auth/logout | Logout (invalidate session) |
| GET | /admin/auth/me | Get current user info |

#### Providers

| Method | Path | Description |
|---|---|---|
| GET | /admin/providers | List providers |
| POST | /admin/providers | Create provider |
| GET | /admin/providers/:id | Get provider |
| PUT | /admin/providers/:id | Update provider |
| DELETE | /admin/providers/:id | Delete provider |

#### Users

| Method | Path | Description |
|---|---|---|
| GET | /admin/users | List users |
| POST | /admin/users | Create user |
| GET | /admin/users/:id | Get user |
| PUT | /admin/users/:id | Update user |
| DELETE | /admin/users/:id | Delete user |

#### API Keys

| Method | Path | Description |
|---|---|---|
| GET | /admin/api-keys/:userId | List API keys for a user |
| POST | /admin/api-keys | Create API key (returns key that is only shown once) |
| DELETE | /admin/api-keys/:id | Delete/revoke API key |

#### Routing Rules

| Method | Path | Description |
|---|---|---|
| GET | /admin/routing-rules | List routing rules |
| POST | /admin/routing-rules | Create routing rule |
| GET | /admin/routing-rules/:id | Get routing rule |
| PUT | /admin/routing-rules/:id | Update routing rule |
| DELETE | /admin/routing-rules/:id | Delete routing rule |

#### Groups

| Method | Path | Description |
|---|---|---|
| GET | /admin/groups | List groups |
| POST | /admin/groups | Create group |
| GET | /admin/groups/:id | Get group |
| PUT | /admin/groups/:id | Update group |
| DELETE | /admin/groups/:id | Delete group |
| POST | /admin/groups/:id/members | Add member to group |
| DELETE | /admin/groups/:id/members/:userId | Remove member from group |

#### Permissions

| Method | Path | Description |
|---|---|---|
| GET | /admin/permissions | List permissions |
| POST | /admin/permissions | Create permission |
| GET | /admin/permissions/:id | Get permission |
| PUT | /admin/permissions/:id | Update permission |
| DELETE | /admin/permissions/:id | Delete permission |

#### Budgets

| Method | Path | Description |
|---|---|---|
| GET | /admin/budgets | List budgets |
| POST | /admin/budgets | Create budget |
| GET | /admin/budgets/:id | Get budget |
| PUT | /admin/budgets/:id | Update budget |
| DELETE | /admin/budgets/:id | Delete budget |

#### Pricing Rules

| Method | Path | Description |
|---|---|---|
| GET | /admin/pricing-rules | List pricing rules |
| POST | /admin/pricing-rules | Create pricing rule |
| GET | /admin/pricing-rules/:id | Get pricing rule |
| PUT | /admin/pricing-rules/:id | Update pricing rule |
| DELETE | /admin/pricing-rules/:id | Delete pricing rule |

#### Alert Rules

| Method | Path | Description |
|---|---|---|
| GET | /admin/alert-rules | List alert rules |
| POST | /admin/alert-rules | Create alert rule |
| GET | /admin/alert-rules/:id | Get alert rule |
| PUT | /admin/alert-rules/:id | Update alert rule |
| DELETE | /admin/alert-rules/:id | Delete alert rule |

#### Alerts

| Method | Path | Description |
|---|---|---|
| GET | /admin/alerts | List active alerts |
| PUT | /admin/alerts/:id/acknowledge | Acknowledge alert |

#### Usage

| Method | Path | Description |
|---|---|---|
| GET | /admin/usage | Query usage records (supports userId, startDate, endDate filters) |

#### Health

| Method | Path | Description |
|---|---|---|
| GET | /admin/health | Get provider health status (latency, error rate, status) |

## Middleware Pipeline Flow

```
Incoming Request
    │
    ▼
┌───────────────────────────────────┐
│ 1. Authentication (auth-service) │
│ ValidateAPIKey → UserIdentity     │
└───────────────────────────────────┘
    │
    ▼
┌───────────────────────────────────┐
│ 2. Authorization (auth-service) │
│ CheckModelAuthorization →       │
│ AuthorizationResult             │
└───────────────────────────────────┘
    │
    ▼
┌───────────────────────────────────┐
│ 3. Rate Limit (billing-service)  │
│ CheckBudget → BudgetStatus       │
└───────────────────────────────────┘
    │
    ▼
┌───────────────────────────────────┐
│ 4. Security (placeholder)       │
└───────────────────────────────────┘
    │
    ▼
┌───────────────────────────────────┐
│ 5. Route (router-service)      │
│ ResolveRoute → RouteResult      │
└───────────────────────────────────┘
    │
    ▼
┌───────────────────────────────────┐
│ 6. Proxy (provider-service)    │
│ ForwardRequest / StreamRequest │
└───────────────────────────────────┘
    │
    ▼
Response to Consumer
```

## Requirements

### Requirement: Admin API
The gateway-service SHALL provide admin HTTP endpoints for all management operations.

#### Scenario: Admin endpoints cover all service entities
- **WHEN** the gateway-service admin API is inspected
- **THEN** it SHALL include endpoints for: providers, users, API keys, usage, authentication, routing rules, groups, permissions, budgets, pricing rules, alert rules, alerts, and health

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
- **THEN** the gateway-service SHALL update the pricing rule via billing-service

#### Scenario: Delete pricing rule
- **WHEN** a DELETE request is made to `/admin/pricing-rules/:id`
- **THEN** the gateway-service SHALL delete the pricing rule via billing-service

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
# admin-user-routing-rules

## Purpose

Admin API for managing per-user routing rules on behalf of users.

## Service Responsibility

- **Role**: Provide admin-facing API for managing user routing rules
- **Owned Entities**: None (gateway-service owns the API, router-service owns the data)
- **Data Layer**: None (proxies to router-service)

## Dependencies

### Calls To

| Service | Methods | Purpose |
|---------|----------|----------|
| router-service | `CreateRoutingRule`, `UpdateRoutingRule`, `DeleteRoutingRule`, `ListRoutingRules` | Rule management |

### Called By

| Service | Methods | Purpose |
|---------|----------|----------|
| Admin users | HTTP `/admin/users/:userId/routing-rules` | Manage user rules |

## Requirements

### Requirement: Admin per-user routing rules management
The gateway-service SHALL provide admin API endpoints for configuring per-user routing rules on behalf of users.

#### Scenario: Admin lists user's routing rules
- **WHEN** an admin makes a GET request to `/admin/users/:userId/routing-rules`
- **THEN** the gateway-service SHALL verify the admin role
- **AND** query router-service for routing rules where `user_id` matches the specified `userId`
- **AND** return HTTP 200 with the list of routing rules for that user

#### Scenario: Admin creates routing rule for user
- **WHEN** an admin makes a POST request to `/admin/users/:userId/routing-rules` with body containing `model_pattern`, `provider_id`, `priority`, and optional `fallback_provider_ids`
- **THEN** the gateway-service SHALL verify the admin role
- **AND** call `router-service.CreateRoutingRule` with the specified `userId` as `user_id`
- **AND** return HTTP 201 with the created routing rule

#### Scenario: Admin updates user's routing rule
- **WHEN** an admin makes a PUT request to `/admin/users/:userId/routing-rules/:id` with updated fields
- **THEN** the gateway-service SHALL verify the admin role
- **AND** call `router-service.UpdateRoutingRule` with the updated data
- **AND** return HTTP 200 with the updated routing rule

#### Scenario: Admin deletes user's routing rule
- **WHEN** an admin makes a DELETE request to `/admin/users/:userId/routing-rules/:id`
- **THEN** the gateway-service SHALL verify the admin role
- **AND** call `router-service.DeleteRoutingRule` with the rule ID
- **AND** return HTTP 200 with success message

#### Scenario: Non-admin access attempt
- **WHEN** a non-admin user makes a request to `/admin/users/:userId/routing-rules`
- **THEN** the gateway-service SHALL return HTTP 403 with error message "admin access required"
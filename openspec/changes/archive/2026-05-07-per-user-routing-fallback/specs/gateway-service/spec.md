## ADDED Requirements

### Requirement: User self-service routing rules API
The gateway-service SHALL provide user self-service API endpoints for managing per-user routing rules (JWT authentication required).

#### Scenario: List own routing rules
- **WHEN** an authenticated user (with valid JWT) makes a GET request to `/v1/routing-rules`
- **THEN** the gateway-service SHALL extract the user ID from the JWT context
- **AND** call `router-service.ListRoutingRules` with the user's ID
- **AND** return HTTP 200 with the list of routing rules for that user

#### Scenario: Create own routing rule
- **WHEN** an authenticated user makes a POST request to `/v1/routing-rules` with body `{model_pattern, provider_id, priority, fallback_provider_ids}`
- **THEN** the gateway-service SHALL extract the user ID from the JWT context
- **AND** call `router-service.CreateRoutingRule` with `user_id` set to the authenticated user's ID
- **AND** return HTTP 201 with the created routing rule

#### Scenario: Update own routing rule
- **WHEN** an authenticated user makes a PUT request to `/v1/routing-rules/:id` with updated fields
- **THEN** the gateway-service SHALL verify the routing rule belongs to the user (by checking `user_id` matches JWT user ID)
- **AND** call `router-service.UpdateRoutingRule` with the updated data
- **AND** return HTTP 200 with the updated routing rule

#### Scenario: Delete own routing rule
- **WHEN** an authenticated user makes a DELETE request to `/v1/routing-rules/:id`
- **THEN** the gateway-service SHALL verify the routing rule belongs to the user
- **AND** call `router-service.DeleteRoutingRule` with the rule ID
- **AND** return HTTP 200 with success message

#### Scenario: Access another user's rule
- **WHEN** an authenticated user makes a request to `/v1/routing-rules/:id` where the rule belongs to a different user
- **THEN** the gateway-service SHALL return HTTP 403 with error message "forbidden"

### Requirement: Admin per-user routing rules API
The gateway-service SHALL provide admin API endpoints for configuring per-user routing rules on behalf of users.

#### Scenario: Admin lists user's routing rules
- **WHEN** an admin makes a GET request to `/admin/users/:userId/routing-rules`
- **THEN** the gateway-service SHALL verify the admin role
- **AND** call `router-service.ListRoutingRules` with the specified `userId`
- **AND** return HTTP 200 with the list of routing rules for that user

#### Scenario: Admin creates routing rule for user
- **WHEN** an admin makes a POST request to `/admin/users/:userId/routing-rules` with body `{model_pattern, provider_id, priority, fallback_provider_ids}`
- **THEN** the gateway-service SHALL verify the admin role
- **AND** call `router-service.CreateRoutingRule` with `user_id` set to the specified `userId`
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

### Requirement: Pass user_id to router-service
The gateway-service SHALL pass the authenticated user's ID to router-service when resolving routes.

#### Scenario: Resolve route with user context
- **WHEN** an authenticated user makes a request to `/v1/chat/completions`
- **THEN** the gateway-service SHALL extract the `user_id` from the API key or JWT context
- **AND** pass `user_id` to `router-service.ResolveRoute`
- **AND** the router-service SHALL use per-user routing rules if available

#### Scenario: Resolve route without user context (admin/system request)
- **WHEN** a request is made without user authentication (e.g., system-internal request)
- **THEN** the gateway-service SHALL call `router-service.ResolveRoute` without `user_id`
- **AND** only system-wide routing rules SHALL be used

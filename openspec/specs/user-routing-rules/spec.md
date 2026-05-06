# user-routing-rules

## Purpose

User self-service API for managing per-user routing rules with fallback chain support.

## Service Responsibility

- **Role**: Provide user-facing API for managing their own routing rules
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
| Consumers (authenticated users) | HTTP `/v1/routing-rules` | Manage own rules |

## Requirements

### Requirement: User routing rules management
The gateway-service SHALL provide user self-service API endpoints for managing per-user routing rules with fallback chain support.

#### Scenario: List own routing rules
- **WHEN** an authenticated user (with valid JWT) makes a GET request to `/v1/routing-rules`
- **THEN** the gateway-service SHALL query router-service for routing rules where `user_id` matches the authenticated user's ID
- **AND** return HTTP 200 with the list of routing rules (including `model_pattern`, `provider_id`, `priority`, `fallback_provider_ids`)

#### Scenario: Create own routing rule
- **WHEN** an authenticated user makes a POST request to `/v1/routing-rules` with body containing `model_pattern`, `provider_id`, `priority`, and optional `fallback_provider_ids`
- **THEN** the gateway-service SHALL call `router-service.CreateRoutingRule` with the user's ID as `user_id`
- **AND** return HTTP 201 with the created routing rule

#### Scenario: Update own routing rule
- **WHEN** an authenticated user makes a PUT request to `/v1/routing-rules/:id` with updated fields
- **THEN** the gateway-service SHALL verify the rule belongs to the user (by checking `user_id`)
- **AND** call `router-service.UpdateRoutingRule` with the updated data
- **AND** return HTTP 200 with the updated routing rule

#### Scenario: Delete own routing rule
- **WHEN** an authenticated user makes a DELETE request to `/v1/routing-rules/:id`
- **THEN** the gateway-service SHALL verify the rule belongs to the user (by checking `user_id`)
- **AND** call `router-service.DeleteRoutingRule` with the rule ID
- **AND** return HTTP 200 with success message

#### Scenario: Access another user's rule
- **WHEN** an authenticated user makes a request to `/v1/routing-rules/:id` where the rule belongs to a different user
- **THEN** the gateway-service SHALL return HTTP 403 with error message "forbidden"

### Requirement: User routing rule fallback chain
The user routing rules SHALL support ordered fallback provider chains for automatic failover.

#### Scenario: Create rule with fallback chain
- **WHEN** a user creates a routing rule with `fallback_provider_ids: ["ollama", "opencode_zen"]`
- **THEN** the rule SHALL be stored with the ordered fallback chain
- **AND** on primary provider failure, the system SHALL try fallback providers in order

#### Scenario: Fallback triggers on all errors
- **WHEN** the primary provider returns a 5xx error, timeout, or specific error code
- **THEN** the system SHALL attempt the next provider in the `fallback_provider_ids` chain
- **AND** continue trying until success or all providers exhausted

### Requirement: User rules override system rules
When a user has a routing rule that matches a model request, that rule SHALL OVERRIDE any system-wide routing rule for the same model pattern.

#### Scenario: User rule takes precedence
- **WHEN** a user (Harry) has a rule `gpt-4* → ollama` with fallback `opencode_zen`
- **AND** the system has a rule `gpt-4* → opencode_zen`
- **AND** Harry makes a request for model `gpt-4`
- **THEN** the system SHALL use `ollama` as the primary provider for Harry's request
- **AND** fall back to `opencode_zen` if `ollama` fails
- **AND** the system rule SHALL be ignored for Harry's request
## ADDED Requirements

### Requirement: Routing Rules Admin API
Gateway service SHALL expose CRUD endpoints for routing rules under `/admin/routing-rules`.

#### Scenario: Create routing rule
- **WHEN** POST `/admin/routing-rules` is called with valid routing rule data (model pattern, primary provider, fallback provider, fallback model)
- **THEN** call router-service `CreateRoutingRule` gRPC and return the created rule with HTTP 201

#### Scenario: List routing rules
- **WHEN** GET `/admin/routing-rules` is called
- **THEN** call router-service `ListRoutingRules` gRPC and return the rule list with HTTP 200

#### Scenario: Update routing rule
- **WHEN** PUT `/admin/routing-rules/:id` is called with updated rule data
- **THEN** call router-service `UpdateRoutingRule` gRPC and return the updated rule with HTTP 200

#### Scenario: Delete routing rule
- **WHEN** DELETE `/admin/routing-rules/:id` is called
- **THEN** call router-service `DeleteRoutingRule` gRPC and return success with HTTP 204

### Requirement: Fallback Retry Logic
Gateway service SHALL implement fallback retry in `ProxyMiddleware` when the primary provider fails with retryable errors.

#### Scenario: Primary provider success
- **WHEN** the primary provider responds successfully to a non-streaming or streaming request
- **THEN** return the response to the consumer without attempting fallback

#### Scenario: Primary provider failure with fallback available
- **WHEN** the primary provider fails with a retryable error (network timeout, 5xx, 429)
- **THEN** rewrite the request model to the fallback provider's configured model
- **AND** retry the request with the fallback provider

#### Scenario: All providers fail
- **WHEN** all primary and fallback providers fail with retryable errors
- **THEN** return an error response with code `all_providers_failed` and HTTP 502

#### Scenario: Non-retryable error
- **WHEN** the primary provider returns a non-retryable error (4xx client error)
- **THEN** return the error immediately without attempting fallback

## MODIFIED Requirements

<!-- No existing requirements modified for this change -->

## REMOVED Requirements

<!-- No requirements removed for this change -->

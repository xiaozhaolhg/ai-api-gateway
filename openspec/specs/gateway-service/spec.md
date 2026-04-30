# gateway-service

## Purpose

Edge service — HTTP entry point, middleware orchestration, SSE streaming.

## Service Responsibility

- **Role**: Sole entry point for all external traffic
- **Owned Entities**: None (stateless)
- **Data Layer**: None

## Dependencies

### Calls To

| Service | Methods | Purpose |
|---------|----------|----------|
| auth-service | `ValidateAPIKey`, `CheckModelAuthorization` | Authenticate request, check model permission |
| router-service | `ResolveRoute` | Resolve model to provider |
| provider-service | `ForwardRequest`, `StreamRequest` | Forward request to AI provider |
| billing-service | `CheckBudget`, `RecordUsage` | Enforce rate limits, track usage |
| monitor-service | `RecordMetric` | Emit request metrics |

### Called By

| Service | Methods | Purpose |
|---------|----------|----------|
| Consumers (external) | HTTP endpoints | OpenAI-compatible API, custom API, admin API |

### Data Dependencies

- **Cache**: None directly (delegates to other services)
- **Database**: None (stateless)

## Key Design

### Middleware Pipeline

| Step | Service Called | Purpose |
|------|-----------------|----------|
| 1 | Auth | Validate API key |
| 2 | Authorization | Check model permission |
| 3 | Rate Limit | Check budget/quota |
| 4 | Security | Prompt injection, content filter |
| 5 | Route | Resolve provider |
| 6 | Proxy | Forward to provider |
| 7 | Callback | (handled by provider-service) |
| 8 | Log | Record metadata |

### Key Operations

- **HTTP Handling**: OpenAI-compatible (`/v1/chat`), custom (`/gateway/`), admin (`/admin/`), auth (`/admin/auth/`)
- **Admin Auth**: `POST /admin/auth/login`, `POST /admin/auth/register` (proxies to auth-service)
- **Streaming**: SSE proxy from provider-service to consumer
- **Middleware Orchestration**: Ordered pipeline execution

---

## Requirements

### Requirement: Streaming Proxy

Gateway service SHALL proxy SSE streaming responses from providers to consumers.

#### Scenario: Streaming Request
- **WHEN** a chat completion request with `stream: true` is received
- **THEN** establish SSE connection to consumer and stream chunks from provider-service

#### Scenario: Token Accumulation During Streaming
- **WHEN** processing SSE chunks from provider
- **THEN** accumulate token counts across all chunks and report final count on completion

### Requirement: Token recording after non-streaming requests
The gateway-service SHALL call billing-service `RecordUsage` RPC after completing a non-streaming LLM request, passing user_id (from JWT), group_id (from ValidateAPIKey response), provider_id, model, and token counts from the provider response.

#### Scenario: Non-streaming request completes
- **WHEN** a non-streaming request to a provider completes
- **THEN** the gateway extracts prompt_tokens and completion_tokens from the ForwardRequestResponse.TokenCounts
- **AND** calls `billingClient.RecordUsage()` with user_id, group_id, provider_id, model, and token counts

#### Scenario: Multiple providers called (non-streaming)
- **WHEN** a request fans out to multiple providers/models
- **THEN** the gateway SHALL call `RecordUsage` separately for each provider/model combination
- **AND** each call includes the correct provider_id and model for that specific call

### Requirement: Token recording after streaming requests
The gateway-service SHALL accumulate token counts across all SSE chunks and call billing-service `RecordUsage` after the stream completes.

#### Scenario: Streaming request completes
- **WHEN** a streaming LLM request completes (all SSE chunks received)
- **THEN** the gateway sums totalPromptTokens and totalCompletionTokens across all chunks
- **AND** calls `billingClient.RecordUsage()` with the accumulated totals

#### Scenario: Multiple streaming providers
- **WHEN** a streaming request calls multiple providers
- **THEN** the gateway SHALL call `RecordUsage` separately for each provider after their stream completes
- **AND** each call uses the correct provider_id and accumulated tokens for that provider

### Requirement: Token extraction from provider responses
The gateway-service SHALL extract token counts from provider responses for both non-streaming (JSON response) and streaming (SSE chunks) flows.

#### Scenario: Non-streaming token extraction
- **WHEN** a non-streaming response is received from provider-service
- **THEN** the gateway reads `resp.TokenCounts.PromptTokens` and `resp.TokenCounts.CompletionTokens`
- **AND** makes these available for `RecordUsage` call

#### Scenario: Streaming token accumulation
- **WHEN** SSE chunks are received from provider-service
- **THEN** the gateway updates running totals: `totalPromptTokens += chunk.AccumulatedTokens.PromptTokens`
- **AND** `totalCompletionTokens += chunk.AccumulatedTokens.CompletionTokens`
- **AND** after stream completion, uses these totals for `RecordUsage`

### Requirement: Group ID propagation for token recording
The gateway-service SHALL use the user's group_id (from `ValidateAPIKey` response) when calling `RecordUsage`.

#### Scenario: Single group user
- **WHEN** a user belongs to one group
- **THEN** the gateway uses that group's ID as `group_id` in `RecordUsage`

#### Scenario: Multi-group user (MVP)
- **WHEN** a user belongs to multiple groups
- **THEN** the gateway uses the FIRST group ID from `UserIdentity.group_ids` as `group_id` in `RecordUsage`
- **Note**: Future enhancement will record separate UsageRecord for each group

### Requirement: Non-Streaming Proxy

Gateway service SHALL proxy non-streaming requests to providers.

#### Scenario: Non-Streaming Request
- **WHEN** a chat completion request without `stream: true` (or `stream: false`)
- **THEN** call `ForwardRequest`, return complete response with token counts

### Requirement: Admin login endpoint
The gateway-service SHALL proxy admin login requests to auth-service.

#### Scenario: Login endpoint
- **WHEN** POST `/admin/login` is called with email/password
- **THEN** call auth-service Login RPC and set HTTP-only JWT cookie on success.

#### Scenario: Login failure
- **WHEN** auth-service rejects credentials
- **THEN** return 401 Unauthorized with error message.

#### Scenario: Cookie management
- **WHEN** login is successful
- **THEN** set JWT cookie with secure, HTTP-only, and /admin path restrictions.

### Requirement: Admin login uses auth-service
The gateway admin login handler SHALL validate credentials via auth-service Login gRPC.

#### Scenario: Valid login
- **WHEN** POST /admin/auth/login is called with valid email and password
- **THEN** the gateway calls auth-service Login gRPC, receives a token, sets auth cookie, and returns token + user

#### Scenario: Invalid login
- **WHEN** POST /admin/auth/login is called with invalid credentials
- **THEN** the gateway returns 401 Unauthorized

### Requirement: Admin user management endpoints
The gateway admin user management endpoints SHALL call auth-service gRPC backends instead of returning hardcoded mock data.

#### Scenario: List users
- **WHEN** GET /admin/auth/users is called
- **THEN** the gateway calls auth-service ListUsers gRPC and returns the response as JSON

#### Scenario: Create user
- **WHEN** POST /admin/auth/users is called with user data
- **THEN** the gateway calls auth-service CreateUser gRPC and returns the created user as JSON

#### Scenario: Update user
- **WHEN** PUT /admin/auth/users/:id is called with updated fields
- **THEN** the gateway calls auth-service UpdateUser gRPC and returns the updated user as JSON

#### Scenario: Delete user
- **WHEN** DELETE /admin/auth/users/:id is called
- **THEN** the gateway calls auth-service DeleteUser gRPC and returns success

### Requirement: Admin API key management endpoints
The gateway SHALL expose API key management endpoints that call auth-service gRPC.

#### Scenario: Create API key
- **WHEN** POST /admin/auth/api-keys is called with user_id and name
- **THEN** the gateway calls auth-service CreateAPIKey gRPC and returns the key (shown once)

#### Scenario: List API keys
- **WHEN** GET /admin/auth/api-keys/:user_id is called
- **THEN** the gateway calls auth-service ListAPIKeys gRPC and returns the key list

#### Scenario: Delete API key
- **WHEN** DELETE /admin/auth/api-keys/:id is called
- **THEN** the gateway calls auth-service DeleteAPIKey gRPC and returns success

### Requirement: Admin usage endpoint calls billing-service
The gateway admin usage endpoint SHALL call billing-service gRPC instead of returning mock data.

#### Scenario: Get usage
- **WHEN** GET /admin/auth/usage is called
- **THEN** the gateway calls billing-service GetUsage gRPC and returns real usage records

### Requirement: Admin group management endpoints
The gateway SHALL expose group management endpoints that proxy to auth-service gRPC.

#### Scenario: Group CRUD
- **WHEN** any of GET/POST/PUT/DELETE /admin/auth/groups[/:id] is called
- **THEN** the gateway calls the corresponding auth-service Group gRPC method and returns the result

#### Scenario: Group membership management
- **WHEN** POST /admin/auth/groups/:id/members or DELETE /admin/auth/groups/:id/members/:user_id is called
- **THEN** the gateway calls AddUserToGroup or RemoveUserFromGroup gRPC respectively

### Requirement: Admin permission management endpoints
The gateway SHALL expose permission management endpoints that proxy to auth-service gRPC.

#### Scenario: Permission CRUD
- **WHEN** any of GET/POST/DELETE /admin/auth/permissions[/:id] is called
- **THEN** the gateway calls the corresponding auth-service Permission gRPC method and returns the result

### Requirement: Auth middleware for admin routes
The gateway-service SHALL validate admin UI sessions.

#### Scenario: Admin route protection
- **WHEN** admin UI routes are accessed
- **THEN** validate JWT cookie and set user context for downstream services.

#### Scenario: Session validation
- **WHEN** JWT is expired or invalid
- **THEN** return 401 Unauthorized to trigger UI redirect to login.

#### Scenario: User context propagation
- **WHEN** admin UI makes API calls
- **THEN** include user ID and role in request context for authorization.

### Requirement: Logout endpoint
The gateway-service SHALL handle admin logout requests.

#### Scenario: Logout request
- **WHEN** POST `/admin/logout` is called
- **THEN** clear the auth cookie and return success response.

#### Scenario: Cookie clearing
- **WHEN** logging out
- **THEN** set cookie with expired date to ensure removal.

#### Scenario: Session invalidation
- **WHEN** user logs out
- **THEN** JWT becomes invalid and cannot be reused.

### Requirement: Error Handling

Gateway service SHALL handle errors with appropriate HTTP status codes and structured error responses.

#### Scenario: Provider timeout
- **WHEN** a provider request exceeds 30 seconds
- **THEN** return HTTP 504 Gateway Timeout
- **AND** return error code "gateway_timeout" with descriptive message

#### Scenario: Invalid API key
- **WHEN** an API key fails validation
- **THEN** return HTTP 401 Unauthorized
- **AND** return error code "invalid_api_key"

#### Scenario: Model not found
- **WHEN** no routing rule exists for the requested model
- **THEN** return HTTP 404 Not Found
- **AND** return error code "model_not_found"

#### Scenario: Provider unavailable
- **WHEN** the provider service is unreachable
- **THEN** return HTTP 502 Bad Gateway
- **AND** return error code "provider_error"

#### Scenario: Rate limit exceeded
- **WHEN** the user exceeds their rate limit quota
- **THEN** return HTTP 429 Too Many Requests
- **AND** return error code "rate_limit_exceeded"

#### Scenario: Authorization denied
- **WHEN** the user is not authorized for the requested model
- **THEN** return HTTP 403 Forbidden
- **AND** return error code "insufficient_permissions"

### Requirement: Structured Logging

Gateway service SHALL log requests and responses in structured JSON format with correlation IDs and sensitive data masking.

#### Scenario: Request logging
- **WHEN** an HTTP request is received
- **THEN** log in JSON format with request_id, method, path, user_id
- **AND** include duration, status code, and timestamp

#### Scenario: Sensitive data masking
- **WHEN** logging request or response bodies
- **THEN** mask sensitive fields (api_key, credentials, password, token)
- **AND** replace values with "***"

#### Scenario: Correlation ID propagation
- **WHEN** a request includes X-Request-ID header
- **THEN** use that ID for all related log entries
- **AND** propagate the ID to downstream gRPC calls

### Requirement: Models Endpoint

Gateway service SHALL provide an endpoint to aggregate models from all configured providers.

#### Scenario: List all models
- **WHEN** GET /gateway/models is called
- **THEN** query all providers concurrently for their models
- **AND** return aggregated list in OpenAI-compatible format

#### Scenario: Models caching
- **WHEN** models are listed successfully
- **THEN** cache the result for 5 minutes
- **AND** return cached result on subsequent requests

#### Scenario: Provider unavailable during listing
- **WHEN** one provider is unavailable during models listing
- **THEN** return models from available providers
- **AND** log warning about unavailable provider

### Requirement: Health Endpoint

Gateway service SHALL provide deep health checks that verify dependent services.

#### Scenario: Deep health check
- **WHEN** GET /gateway/health is called
- **THEN** check health of auth, router, provider, and billing services
- **AND** return status and latency for each service

#### Scenario: Healthy status
- **WHEN** all dependent services are responding normally
- **THEN** return overall status "healthy"
- **AND** return HTTP 200

#### Scenario: Degraded status
- **WHEN** one service is responding slowly (>500ms)
- **THEN** return overall status "degraded"
- **AND** return HTTP 200

#### Scenario: Unhealthy status
- **WHEN** one or more services are down
- **THEN** return overall status "unhealthy"
- **AND** return HTTP 503

### Requirement: Graceful Shutdown

Gateway service SHALL handle shutdown signals gracefully, completing active requests before exiting.

#### Scenario: SIGTERM received
- **WHEN** SIGTERM or SIGINT is received
- **THEN** stop accepting new connections
- **AND** wait up to 30 seconds for active requests

#### Scenario: SSE stream during shutdown
- **WHEN** shutdown occurs while SSE streams are active
- **THEN** send final [DONE] marker to all streams
- **AND** close connections gracefully

#### Scenario: Shutdown timeout
- **WHEN** active requests exceed 30 second shutdown window
- **THEN** force close remaining connections
- **AND** exit process

### Requirement: Request Timeouts

All outbound gRPC calls SHALL have configurable timeouts.

#### Scenario: Auth service timeout
- **WHEN** calling auth-service
- **THEN** apply 5 second timeout
- **AND** return 504 on timeout

#### Scenario: Router service timeout
- **WHEN** calling router-service
- **THEN** apply 5 second timeout
- **AND** return 504 on timeout

#### Scenario: Provider service timeout
- **WHEN** calling provider-service for request forwarding
- **THEN** apply 30 second timeout
- **AND** return 504 on timeout

#### Scenario: Billing service timeout
- **WHEN** calling billing-service
- **THEN** apply 10 second timeout
- **AND** return empty result on timeout (fail open)

### Requirement: Billing Integration

Gateway service SHALL integrate with billing-service for usage queries.

#### Scenario: Query usage records
- **WHEN** GET /admin/usage is called with admin auth
- **THEN** forward query to billing-service GetUsage RPC
- **AND** return formatted usage records

#### Scenario: Billing service unavailable
- **WHEN** billing-service is unavailable during usage query
- **THEN** return empty list with warning
- **AND** log error for monitoring

### Requirement: SSE Heartbeat

Gateway service SHALL send periodic heartbeat messages during SSE streaming.

#### Scenario: SSE keepalive
- **WHEN** SSE stream is active for more than 30 seconds
- **THEN** send heartbeat comment line ": ping"
- **AND** continue until stream completes

#### Scenario: Heartbeat format
- **WHEN** sending heartbeat
- **THEN** use SSE comment format (starts with colon)
- **AND** flush immediately to keep connection alive

### Requirement: Lazy Connection

Gateway service SHALL defer gRPC connections until first request.

#### Scenario: Startup without dependencies
- **WHEN** gateway service starts
- **THEN** do not immediately connect to dependent services
- **AND** accept HTTP requests

#### Scenario: First request connection
- **WHEN** first request requires a downstream service
- **THEN** establish gRPC connection on demand
- **AND** cache connection for reuse

#### Scenario: Connection failure handling
- **WHEN** lazy connection fails
- **THEN** return HTTP 503 Service Unavailable
- **AND** include service name in error details

### Requirement: Load Testing

Gateway service SHALL include automated load testing with k6.

#### Scenario: Concurrent request load
- **WHEN** running k6 load test with 100 concurrent VU
- **THEN** sustain load for 5 minutes
- **AND** achieve p95 latency under 2 seconds

#### Scenario: Error rate threshold
- **WHEN** running load test
- **THEN** maintain error rate below 1%
- **AND** report any failures

#### Scenario: Streaming load
- **WHEN** testing with 50 concurrent SSE streams
- **THEN** sustain streams for 5 minutes
- **AND** verify no connection drops

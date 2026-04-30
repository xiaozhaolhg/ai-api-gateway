## ADDED Requirements

### Gateway Service Error Handling

Gateway service shall handle errors with appropriate HTTP status codes and structured error responses.

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

### Gateway Service Structured Logging

Gateway service shall log requests and responses in structured JSON format with correlation IDs and sensitive data masking.

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

### Gateway Service Models Endpoint

Gateway service shall provide an endpoint to aggregate models from all configured providers.

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

### Gateway Service Health Endpoint

Gateway service shall provide deep health checks that verify dependent services.

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

### Gateway Service Graceful Shutdown

Gateway service shall handle shutdown signals gracefully, completing active requests before exiting.

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

### Gateway Service Request Timeouts

All outbound gRPC calls shall have configurable timeouts.

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

### Gateway Service Billing Integration

Gateway service shall integrate with billing-service for usage queries.

#### Scenario: Query usage records
- **WHEN** GET /admin/usage is called with admin auth
- **THEN** forward query to billing-service GetUsage RPC
- **AND** return formatted usage records

#### Scenario: Billing service unavailable
- **WHEN** billing-service is unavailable during usage query
- **THEN** return empty list with warning
- **AND** log error for monitoring

### Gateway Service SSE Heartbeat

Gateway service shall send periodic heartbeat messages during SSE streaming.

#### Scenario: SSE keepalive
- **WHEN** SSE stream is active for more than 30 seconds
- **THEN** send heartbeat comment line ": ping"
- **AND** continue until stream completes

#### Scenario: Heartbeat format
- **WHEN** sending heartbeat
- **THEN** use SSE comment format (starts with colon)
- **AND** flush immediately to keep connection alive

### Gateway Service Lazy Connection

Gateway service shall defer gRPC connections until first request.

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

### Gateway Service Load Testing

Gateway service shall include automated load testing with k6.

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

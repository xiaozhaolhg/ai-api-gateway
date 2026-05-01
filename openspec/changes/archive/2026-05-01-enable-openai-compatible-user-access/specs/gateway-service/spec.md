## MODIFIED Requirements

### Requirement: Non-Streaming Proxy

Gateway service SHALL proxy non-streaming requests to providers.

#### Scenario: Non-Streaming Request
- **WHEN** a chat completion request without `stream: true` (or `stream: false`)
- **THEN** call `ForwardRequest`, return complete response with token counts

#### Scenario: Non-Streaming Request via /v1/chat/completions
- **WHEN** a POST request is made to `/v1/chat/completions` with `stream: false` or omitted
- **THEN** the gateway SHALL validate API key, check authorization, resolve route, and proxy to provider via `ProxyMiddleware`
- **AND** return the complete OpenAI-compatible response

### Requirement: Streaming Proxy

Gateway service SHALL proxy SSE streaming responses from providers to consumers.

#### Scenario: Streaming Request via /v1/chat/completions
- **WHEN** a POST request is made to `/v1/chat/completions` with `stream: true`
- **THEN** the gateway SHALL validate API key, check authorization, resolve route, and proxy SSE stream to consumer
- **AND** accumulate token counts and record usage after stream completes

### Requirement: Models Endpoint

Gateway service SHALL provide an endpoint to aggregate models from all configured providers.

#### Scenario: List models via /v1/models
- **WHEN** a GET request is made to `/v1/models`
- **THEN** the gateway SHALL validate API key and return models from all configured providers in OpenAI-compatible format
- **AND** use the existing `ModelsHandler` to aggregate models via provider-service

#### Scenario: Models caching
- **WHEN** models are listed successfully
- **THEN** cache the result for 5 minutes
- **AND** return cached result on subsequent requests

#### Scenario: Provider unavailable during listing
- **WHEN** one provider is unavailable during models listing
- **THEN** return models from available providers
- **AND** log warning about unavailable provider

## ADDED Requirements

### Requirement: OpenAI-compatible API middleware wiring
The gateway-service SHALL wire the complete middleware chain (Auth → Authz → Route → Proxy) to the `/v1` route group.

#### Scenario: Middleware chain execution
- **WHEN** a request is made to any `/v1/*` endpoint
- **THEN** the following middleware SHALL execute in order:
  1. `AuthMiddleware` - validate API key via auth-service
  2. `AuthzMiddleware` - check model authorization via auth-service
  3. `RouteMiddleware` - resolve model to provider via router-service
  4. `ProxyMiddleware` - forward request to provider via provider-service

#### Scenario: Middleware chain implementation
- **WHEN** the gateway-service starts
- **THEN** the `/v1` route group in `cmd/server/main.go` SHALL use the `Use()` method to attach all middleware
- **AND** the `ProxyMiddleware` SHALL be the final handler that processes the request

### Requirement: /v1/* endpoint authentication
The gateway-service SHALL require valid API key authentication for all `/v1/*` OpenAI-compatible endpoints.

#### Scenario: Valid API key on /v1/chat/completions
- **WHEN** a POST request is made to `/v1/chat/completions` with valid API key in `Authorization: Bearer <key>` header
- **THEN** the request SHALL proceed through the middleware chain

#### Scenario: Missing API key on /v1/models
- **WHEN** a GET request is made to `/v1/models` without an API key
- **THEN** the gateway SHALL return HTTP 401 with error code `invalid_api_key`

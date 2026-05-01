# gateway-openai-api-wiring

## Purpose

Middleware wiring for OpenAI-compatible endpoints (`/v1/*`) in gateway-service.

## Service Responsibility

- Wire authentication, authorization, routing, and proxy middleware to `/v1/*` endpoints
- Ensure proper middleware execution order: Auth → Authz → Route → Proxy

## Dependencies

### Calls To

| Service | Methods | Purpose |
|---------|----------|----------|
| auth-service | `ValidateAPIKey`, `CheckModelAuthorization` | Authenticate and authorize requests |
| router-service | `ResolveRoute` | Resolve model to provider |
| provider-service | `ForwardRequest`, `StreamRequest` | Forward request to AI provider |

### Called By

| Service | Methods | Purpose |
|---------|----------|----------|
| Consumers (external) | HTTP endpoints | OpenAI-compatible API |

## Requirements

### Requirement: OpenAI-compatible endpoint authentication
The gateway-service SHALL validate API keys for all `/v1/*` OpenAI-compatible endpoints using the existing `AuthMiddleware`.

#### Scenario: Valid API key
- **WHEN** a request is made to `/v1/chat/completions` with a valid API key in the `Authorization: Bearer <key>` header
- **THEN** the request SHALL proceed to the authorization middleware

#### Scenario: Missing API key
- **WHEN** a request is made to `/v1/chat/completions` without an API key
- **THEN** the gateway SHALL return HTTP 401 with error code `invalid_api_key`

#### Scenario: Invalid API key
- **WHEN** a request is made to `/v1/chat/completions` with an invalid or expired API key
- **THEN** the gateway SHALL return HTTP 401 with error code `invalid_api_key`

### Requirement: OpenAI-compatible endpoint authorization
The gateway-service SHALL check model authorization for all `/v1/*` OpenAI-compatible endpoints using the existing `AuthzMiddleware`.

#### Scenario: Authorized model access
- **WHEN** an authenticated request is made to `/v1/chat/completions` with a model the user is authorized to access
- **THEN** the request SHALL proceed to the route middleware

#### Scenario: Unauthorized model access
- **WHEN** an authenticated request is made to `/v1/chat/completions` with a model the user is NOT authorized to access
- **THEN** the gateway SHALL return HTTP 403 with error code `insufficient_permissions`

### Requirement: Model to provider routing for OpenAI-compatible endpoints
The gateway-service SHALL resolve the requested model to a provider using the `RouteMiddleware` for all `/v1/*` OpenAI-compatible endpoints.

#### Scenario: Model resolution success
- **WHEN** a request is made to `/v1/chat/completions` with model `"gpt-4"`
- **THEN** the gateway SHALL resolve the model to a provider via `router-service.ResolveRoute`
- **AND** proceed to proxy the request to the resolved provider

#### Scenario: Model not found
- **WHEN** a request is made to `/v1/chat/completions` with a model that has no routing rule
- **THEN** the gateway SHALL return HTTP 404 with error code `model_not_found`

### Requirement: Request proxy for OpenAI-compatible endpoints
The gateway-service SHALL proxy requests to the resolved provider using the `ProxyMiddleware` for all `/v1/*` OpenAI-compatible endpoints.

#### Scenario: Non-streaming request
- **WHEN** a request is made to `/v1/chat/completions` with `stream: false` or omitted
- **THEN** the gateway SHALL call `provider-service.ForwardRequest`
- **AND** return the complete response in OpenAI-compatible format

#### Scenario: Streaming request
- **WHEN** a request is made to `/v1/chat/completions` with `stream: true`
- **THEN** the gateway SHALL call `provider-service.StreamRequest`
- **AND** return SSE chunks in OpenAI-compatible format

### Requirement: Middleware chain wiring
The gateway-service SHALL wire the complete middleware chain (Auth → Authz → Route → Proxy) to the `/v1` route group in `cmd/server/main.go`.

#### Scenario: Middleware order
- **WHEN** the gateway-service starts
- **THEN** the `/v1` route group SHALL have middleware executed in order: `AuthMiddleware` → `AuthzMiddleware` → `RouteMiddleware` → `ProxyMiddleware`

#### Scenario: Existing admin endpoints unaffected
- **WHEN** the middleware is wired to `/v1` routes
- **THEN** the `/admin/*` routes SHALL continue to use the existing `jwtAuthMiddleware` without changes

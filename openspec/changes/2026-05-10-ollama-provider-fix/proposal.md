# Proposal: Fix Ollama Provider Integration and Admin-UI Deployment

## Problem Statement

The system has multiple issues preventing successful LLM provider integration and admin-UI deployment:

1. **Gateway-Provider Context Mismatch**: The `route.go` middleware sets provider ID in request context but `proxy.go` reads from gin context, causing "Provider not resolved" errors.

2. **Admin-UI API Proxy Configuration**: The nginx proxy incorrectly routes `/admin` instead of `/admin/auth`, causing 404 errors on auth endpoints.

3. **Provider-Service HTTP Request Bug**: The `makeHTTPRequest` and `makeStreamingHTTPRequest` functions pass `nil` as the request body instead of the actual body data, causing all provider requests to fail.

4. **Ollama Adapter Model Name Handling**: The Ollama adapter doesn't properly handle model names with provider prefixes (e.g., "ollama:qwen3.5") or map them to valid Ollama model identifiers.

5. **Docker Network Access**: On Linux, containers cannot resolve `host.docker.internal` which prevents provider-service from connecting to Ollama on the host machine.

**Impact**:
- Chat completions API returns 403/404 errors
- Admin UI cannot authenticate users
- Provider requests fail with empty body
- Ollama models cannot be accessed from containers

## Proposed Solution

### 1. Gateway Context Fix
Add `c.Set("providerId", result.ProviderID)` in `route.go` middleware to ensure proxy middleware can read provider ID from gin context.

### 2. Admin-UI Nginx Configuration Fix
Change nginx location from `/admin` to `/admin/auth` to only proxy auth-related endpoints, allowing other admin routes to be handled by the React app.

### 3. Admin-UI BaseURL Fix
Update `config.ts` to default `baseURL` to empty string for relative paths in production, and add Dockerfile ARG support for build-time configuration.

### 4. Provider-Service HTTP Body Fix
Fix both `makeHTTPRequest` and `makeStreamingHTTPRequest` to use `bytes.NewReader(body)` instead of `nil` for the request body.

### 5. Ollama Adapter Model Handling Fix
- Strip provider prefix from model names (e.g., "ollama:qwen3.5" → "qwen3.5")
- Add version suffix for model names without version (e.g., "qwen3.5" → "qwen3.5:0.8b")
- Use `/api/chat` format with `messages` array instead of deprecated `prompt` field

### 6. Docker Network Fix
Add `extra_hosts` configuration to provider-service to enable container-to-host networking on Linux.

## Scope

### In Scope:
- Fix gateway middleware context propagation
- Fix admin-UI nginx and build configuration
- Fix provider-service HTTP request body handling
- Fix Ollama adapter request transformation
- Fix Docker network configuration for host access

### Out of Scope:
- Changes to authentication service logic
- Database schema changes
- Adding new provider types
- Changes to billing/monitoring services

## Success Criteria

1. `POST /v1/chat/completions` with Ollama provider returns valid responses
2. Admin-UI auth endpoints (/admin/auth/*) return 200 instead of 404
3. All provider requests include proper request body data
4. Container can connect to Ollama on host machine via Docker network

## Dependencies

- Ollama service running on host machine port 11434
- Go 1.23+ for building services
- Docker with docker-compose support

## Risks

- Ollama model name mapping is a temporary hardcoded solution
- Network configuration may differ across Docker versions
- Admin-UI relative path may break in some deployment scenarios
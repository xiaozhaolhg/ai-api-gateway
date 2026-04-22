## Why

The current codebase has a working prototype (`router-service`) that routes chat requests to LLM providers, but it is a single-process monolith with no authentication, no usage tracking, no admin interface, and no service boundaries. The business requires a governed, observable, multi-tenant AI gateway — a single binary cannot deliver that. We must decompose the system into independently deployable microservices so that auth, routing, provider management, billing, and monitoring can evolve at their own pace and be operated by separate teams.

## What Changes

- **BREAKING**: Freeze existing `router-service/` as `router-service-legacy/` — it will no longer be the deployed service
- Create shared `api/` Go module with protobuf definitions and buf-generated gRPC stubs for all 5 internal services
- Create `auth-service/` — identity, API key validation, user/group management, model authorization (gRPC server, SQLite)
- Create `router-service/` (fresh) — model-to-provider route resolution (gRPC server, SQLite)
- Create `provider-service/` — provider CRUD, provider adapters (OpenAI-compatible, Anthropic, Gemini, Ollama, OpenCode Zen), request forwarding, SSE streaming, async callback dispatch to billing and monitor (gRPC server, SQLite)
- Create `gateway-service/` — HTTP entry point, middleware pipeline (auth → authz → ratelimit → security → route → proxy → callback → log), OpenAI-compatible API, custom gateway API, admin API (gRPC client to all 5 services)
- Create `billing-service/` — usage recording via provider callback, cost estimation, budget enforcement, pricing rules (gRPC server, SQLite)
- Create `monitor-service/` — metrics collection via provider callback, provider health tracking, alerting (gRPC server, SQLite)
- Create `admin-ui/` — React SPA for provider management, user/API key management, usage dashboards (separate nginx container)
- Create `docker-compose.yaml` for local development (7 services + Redis)
- Create `go.work` for multi-module Go workspace

## Capabilities

### New Capabilities

- `api-proto`: Shared protobuf definitions and buf-generated gRPC stubs for all inter-service communication
- `auth-service-architecture`: Identity, API key validation, user/group management, model authorization
- `auth-service-deployment`: Docker build, config, SQLite, health check
- `auth-service-testing`: Unit, integration, and gRPC contract tests
- `router-service-architecture`: Model-to-provider route resolution, routing rule management
- `router-service-deployment`: Docker build, config, SQLite, health check
- `router-service-testing`: Unit, integration, and gRPC contract tests
- `provider-service-architecture`: Provider CRUD, adapter framework, request forwarding, SSE streaming, callback dispatch
- `provider-service-deployment`: Docker build, config, SQLite, health check
- `provider-service-testing`: Unit, integration, adapter, and gRPC contract tests
- `gateway-service-architecture`: HTTP entry point, middleware pipeline, gRPC client orchestration
- `gateway-service-deployment`: Docker build, config, health check
- `gateway-service-testing`: Unit, integration, middleware, and e2e tests
- `billing-service-architecture`: Usage recording, cost estimation, budget enforcement, pricing rules
- `billing-service-deployment`: Docker build, config, SQLite, health check
- `billing-service-testing`: Unit, integration, and gRPC contract tests
- `monitor-service-architecture`: Metrics collection, provider health, alerting
- `monitor-service-deployment`: Docker build, config, SQLite, health check
- `monitor-service-testing`: Unit, integration, and gRPC contract tests
- `system-deployment`: Docker Compose orchestration, go.work, top-level Makefile
- `admin-ui-architecture`: React SPA, API client, pages for providers/users/usage/dashboard
- `admin-ui-deployment`: Docker build (nginx), config, health check

### Modified Capabilities

- `router-service-architecture`: Requirements change from single-process HTTP server to gRPC server with SQLite persistence
- `router-service-deployment`: Deployment changes from standalone Gin server to gRPC service in Docker Compose

## Impact

- **Code**: Existing `router-service/` frozen and renamed; 6 new service directories created; shared `api/` module added
- **APIs**: New gRPC APIs for all 5 internal services; gateway-service exposes OpenAI-compatible HTTP endpoints and admin API
- **Dependencies**: buf (proto generation), grpc-go, protobuf-go, SQLite drivers, Docker Compose, Redis (optional for MVP)
- **Build**: Multi-module Go workspace (go.work); each service has independent go.mod, Dockerfile, Makefile
- **Operations**: 7 containers in Docker Compose (6 Go services + admin-ui); SQLite per service; Redis shared

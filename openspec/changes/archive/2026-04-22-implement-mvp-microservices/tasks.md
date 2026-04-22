## 1. Repository Setup & Legacy Freeze

- [x] 1.1 Rename `router-service/` to `router-service-legacy/` and commit
- [x] 1.2 Create `go.work` file at repo root
- [x] 1.3 Create top-level `Makefile` with targets: build, test, up, down, proto
- [x] 1.4 Create `docker-compose.yaml` skeleton with service definitions and network config

## 2. Shared API Module (api/)

- [x] 2.1 Create `api/` directory with `go.mod` (module: github.com/ai-api-gateway/api)
- [x] 2.2 Create `api/buf.yaml` and `api/buf.gen.yaml` with buf configuration
- [x] 2.3 Create `api/proto/common/v1/common.proto` with shared messages (Empty, TokenCounts, ProviderResponseCallback, BudgetStatus)
- [x] 2.4 Create `api/proto/auth/v1/auth.proto` matching auth-service API contracts
- [x] 2.5 Create `api/proto/router/v1/router.proto` matching router-service API contracts
- [x] 2.6 Create `api/proto/provider/v1/provider.proto` matching provider-service API contracts
- [x] 2.7 Create `api/proto/billing/v1/billing.proto` matching billing-service API contracts
- [x] 2.8 Create `api/proto/monitor/v1/monitor.proto` matching monitor-service API contracts
- [x] 2.9 Run `buf lint` and fix any linting errors
- [x] 2.10 Run `buf generate` and verify generated Go stubs in `api/gen/`
- [x] 2.11 Add `api/` entry to `go.work`

## 3. Auth Service (auth-service/)

- [x] 3.1 Create `auth-service/` directory with `go.mod`, `cmd/server/main.go`
- [x] 3.2 Create `auth-service/internal/domain/entity/` with User and APIKey entities
- [x] 3.3 Create `auth-service/internal/domain/port/` with UserRepository and APIKeyRepository interfaces
- [x] 3.4 Create `auth-service/internal/application/` with auth logic (key hashing, validation, authorization)
- [x] 3.5 Create `auth-service/internal/infrastructure/config/` with config loading (YAML + env vars)
- [x] 3.6 Create `auth-service/internal/infrastructure/repository/` with SQLite UserRepository implementation
- [x] 3.7 Create `auth-service/internal/infrastructure/repository/` with SQLite APIKeyRepository implementation
- [x] 3.8 Create `auth-service/internal/infrastructure/migration/` with GORM AutoMigrate for SQLite schema
- [x] 3.9 Create `auth-service/internal/handler/` with gRPC server implementing AuthService proto (placeholder - needs proto generation)
- [x] 3.10 Create `auth-service/configs/config.yaml` with default settings
- [x] 3.11 Create `auth-service/Dockerfile` (multi-stage, alpine, non-root)
- [x] 3.12 Create `auth-service/Makefile` with build, test, docker targets
- [x] 3.13 Add `auth-service/` entry to `go.work` (already in go.work from task 1.2)
- [ ] 3.14 Write unit tests for domain entities and key hashing
- [ ] 3.15 Write integration tests for SQLite repositories
- [ ] 3.16 Write gRPC server tests

## 4. Provider Service (provider-service/)

- [x] 4.1 Create `provider-service/` directory with `go.mod`, `cmd/server/main.go`
- [x] 4.2 Create `provider-service/internal/domain/entity/` with Provider entity
- [x] 4.3 Create `provider-service/internal/domain/port/` with ProviderRepository and ProviderAdapter interfaces
- [x] 4.4 Create `provider-service/internal/application/` with forwarding logic and callback dispatch
- [x] 4.5 Create `provider-service/internal/infrastructure/adapter/` with OpenAI-compatible adapter
- [x] 4.6 Create `provider-service/internal/infrastructure/adapter/` with Anthropic adapter
- [x] 4.7 Create `provider-service/internal/infrastructure/adapter/` with Ollama adapter
- [x] 4.8 Create `provider-service/internal/infrastructure/adapter/` with OpenCode Zen adapter
- [x] 4.9 Create `provider-service/internal/infrastructure/adapter/` with Gemini adapter
- [x] 4.10 Create `provider-service/internal/infrastructure/adapter/factory.go` with adapter factory/registry
- [x] 4.11 Create `provider-service/internal/infrastructure/repository/` with SQLite ProviderRepository
- [x] 4.12 Create `provider-service/internal/infrastructure/crypto/` with AES-256-GCM credential encryption
- [x] 4.13 Create `provider-service/internal/infrastructure/config/` with config loading
- [x] 4.14 Create `provider-service/internal/infrastructure/migration/` with GORM AutoMigrate for SQLite schema
- [x] 4.15 Create `provider-service/internal/handler/` with gRPC server implementing ProviderService proto (placeholder - needs proto generation)
- [x] 4.16 Create `provider-service/configs/config.yaml` with default settings
- [x] 4.17 Create `provider-service/Dockerfile` (multi-stage, alpine, non-root)
- [x] 4.18 Create `provider-service/Makefile` with build, test, docker targets
- [x] 4.19 Add `provider-service/` entry to `go.work` (already in go.work from task 1.2)
- [ ] 4.20 Write unit tests for domain entities and adapter interface
- [ ] 4.21 Write integration tests for SQLite ProviderRepository
- [ ] 4.22 Write adapter integration tests (Ollama, OpenAI-compatible)
- [ ] 4.23 Write gRPC server tests (ForwardRequest, CreateProvider, StreamRequest)
- [ ] 4.24 Write callback dispatch tests (async fire-and-forget, failure isolation)

## 5. Router Service (router-service/ — fresh build)

- [x] 5.1 Create `router-service/` directory with `go.mod`, `cmd/server/main.go`
- [x] 5.2 Create `router-service/internal/domain/entity/` with RoutingRule entity
- [x] 5.3 Create `router-service/internal/domain/port/` with RoutingRuleRepository and Router interfaces
- [x] 5.4 Create `router-service/internal/application/` with route resolution logic (pattern matching, priority)
- [x] 5.5 Create `router-service/internal/infrastructure/repository/` with SQLite RoutingRuleRepository
- [x] 5.6 Create `router-service/internal/infrastructure/config/` with config loading
- [x] 5.7 Create `router-service/internal/infrastructure/migration/` with GORM AutoMigrate for SQLite schema
- [x] 5.8 Create `router-service/internal/handler/` with gRPC server implementing RouterService proto (placeholder - needs proto generation)
- [x] 5.9 Create `router-service/configs/config.yaml` with default settings
- [x] 5.10 Create `router-service/Dockerfile` (multi-stage, alpine, non-root)
- [x] 5.11 Create `router-service/Makefile` with build, test, docker targets
- [x] 5.12 Add `router-service/` entry to `go.work` (already in go.work from task 1.2)
- [ ] 5.13 Write unit tests for routing logic (pattern matching, priority ordering)
- [ ] 5.14 Write integration tests for SQLite RoutingRuleRepository
- [ ] 5.15 Write gRPC server tests (ResolveRoute, CreateRoutingRule, RefreshRoutingTable)

## 6. Gateway Service (gateway-service/)

- [x] 6.1 Create `gateway-service/` directory with `go.mod`, `cmd/server/main.go`
- [x] 6.2 Create `gateway-service/internal/client/` with gRPC client wrappers for all 5 services (with retry)
- [x] 6.3 Create `gateway-service/internal/middleware/auth.go` — API key validation via auth-service
- [x] 6.4 Create `gateway-service/internal/middleware/authz.go` — model authorization via auth-service
- [x] 6.5 Create `gateway-service/internal/middleware/ratelimit.go` — placeholder pass-through
- [x] 6.6 Create `gateway-service/internal/middleware/security.go` — placeholder pass-through
- [x] 6.7 Create `gateway-service/internal/middleware/route.go` — route resolution via router-service
- [x] 6.8 Create `gateway-service/internal/middleware/proxy.go` — request forwarding via provider-service
- [x] 6.9 Create `gateway-service/internal/middleware/log.go` — request metadata logging
- [x] 6.10 Create `gateway-service/internal/handler/chat.go` — /v1/chat/completions (streaming + non-streaming)
- [x] 6.11 Create `gateway-service/internal/handler/models.go` — /v1/models
- [x] 6.12 Create `gateway-service/internal/handler/admin_providers.go` — /admin/providers CRUD
- [x] 6.13 Create `gateway-service/internal/handler/admin_users.go` — /admin/users and /admin/api-keys
- [x] 6.14 Create `gateway-service/internal/handler/admin_usage.go` — /admin/usage
- [x] 6.15 Create `gateway-service/internal/handler/health.go` — /health and /gateway/health
- [x] 6.16 Create `gateway-service/internal/infrastructure/config/` with config loading
- [x] 6.17 Create `gateway-service/configs/config.yaml` with default settings and service addresses
- [x] 6.18 Create `gateway-service/Dockerfile` (multi-stage, alpine, non-root)
- [x] 6.19 Create `gateway-service/Makefile` with build, test, docker targets
- [x] 6.20 Add `gateway-service/` entry to `go.work` (already in go.work from task 1.2)
- [ ] 6.21 Write unit tests for middleware pipeline (auth, authz, route, proxy) with mock gRPC clients
- [ ] 6.22 Write integration tests for HTTP endpoints with in-process gRPC service stubs
- [ ] 6.23 Write SSE streaming tests

## 7. Billing Service (billing-service/)

- [x] 7.1 Create `billing-service/` directory with `go.mod`, `cmd/server/main.go`
- [x] 7.2 Create `billing-service/internal/domain/entity/` with UsageRecord, PricingRule, BillingAccount, Budget entities
- [x] 7.3 Create `billing-service/internal/domain/port/` with repository interfaces for all entities
- [x] 7.4 Create `billing-service/internal/application/` with usage recording, cost calculation, budget checking
- [x] 7.5 Create `billing-service/internal/infrastructure/repository/` with SQLite implementations for all repos
- [x] 7.6 Create `billing-service/internal/infrastructure/config/` with config loading
- [x] 7.7 Create `billing-service/internal/infrastructure/migration/` with GORM AutoMigrate for SQLite schema
- [x] 7.8 Create `billing-service/internal/handler/` with gRPC server implementing BillingService proto
- [x] 7.9 Create `billing-service/configs/config.yaml` with default settings
- [x] 7.10 Create `billing-service/Dockerfile` (multi-stage, alpine, non-root)
- [x] 7.11 Create `billing-service/Makefile` with build, test, docker targets
- [x] 7.12 Add `billing-service/` entry to `go.work` (already in go.work from task 1.2)
- [ ] 7.13 Write unit tests for domain entities, cost calculation, budget checking
- [ ] 7.14 Write integration tests for SQLite repositories
- [ ] 7.15 Write gRPC server tests (OnProviderResponse, RecordUsage, GetUsage, CheckBudget)

## 8. Monitor Service (monitor-service/)

- [x] 8.1 Create `monitor-service/` directory with `go.mod`, `cmd/server/main.go`
- [x] 8.2 Create `monitor-service/internal/domain/entity/` with Metric, AlertRule, Alert, ProviderHealthStatus entities
- [x] 8.3 Create `monitor-service/internal/domain/port/` with repository interfaces for all entities
- [x] 8.4 Create `monitor-service/internal/application/` with metrics collection, health tracking, alert evaluation
- [x] 8.5 Create `monitor-service/internal/infrastructure/repository/` with SQLite implementations for all repos
- [x] 8.6 Create `monitor-service/internal/infrastructure/config/` with config loading
- [x] 8.7 Create `monitor-service/internal/infrastructure/migration/` with GORM AutoMigrate for SQLite schema
- [x] 8.8 Create `monitor-service/internal/handler/` with gRPC server implementing MonitorService proto
- [x] 8.9 Create `monitor-service/configs/config.yaml` with default settings
- [x] 8.10 Create `monitor-service/Dockerfile` (multi-stage, alpine, non-root)
- [x] 8.11 Create `monitor-service/Makefile` with build, test, docker targets
- [x] 8.12 Add `monitor-service/` entry to `go.work` (already in go.work from task 1.2)
- [ ] 8.13 Write unit tests for domain entities and alert rule evaluation
- [ ] 8.14 Write integration tests for SQLite repositories
- [ ] 8.15 Write gRPC server tests (OnProviderResponse, RecordMetric, GetProviderHealth)

## 9. Admin UI (admin-ui/)

- [x] 9.1 Scaffold React + Vite + TypeScript + TailwindCSS project in `admin-ui/`
- [x] 9.2 Create `admin-ui/src/api/` with typed HTTP client for all admin endpoints
- [x] 9.3 Create layout shell: sidebar navigation, header, content area
- [x] 9.4 Create provider management page (list, add, edit, remove)
- [x] 9.5 Create user management page (list, create, disable)
- [x] 9.6 Create API key management page (issue, display once, revoke)
- [x] 9.7 Create usage dashboard page (token counts, cost, filters)
- [x] 9.8 Create provider health dashboard page (status, latency, error rate)
- [x] 9.9 Configure Vite dev proxy to gateway-service
- [x] 9.10 Create `admin-ui/Dockerfile` (multi-stage: node build + nginx serve)
- [x] 9.11 Create `admin-ui/nginx.conf` with SPA fallback routing and API proxy
- [x] 9.12 Create `admin-ui/Makefile` with dev, build, docker targets

## 10. Docker Compose & Integration

- [x] 10.1 Complete `docker-compose.yaml` with all 7 services + Redis, health checks, depends_on (already created in task 1.4)
- [x] 10.2 Configure gateway-service to wait for auth, router, provider services
- [x] 10.3 Configure provider-service callback registration to billing and monitor services on startup
- [ ] 10.4 Verify `docker-compose up` starts all services successfully
- [ ] 10.5 End-to-end smoke test: add provider via admin API → create user → issue API key → send chat request → verify usage record
- [ ] 10.6 End-to-end smoke test: streaming chat request → verify SSE format → verify callback to billing/monitor
- [ ] 10.7 End-to-end smoke test: admin UI → add provider → create user → view usage

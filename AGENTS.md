# PROJECT KNOWLEDGE BASE

**Generated:** 2026-04-21
**Branch:** multi-agent

## OVERVIEW

Enterprise AI Gateway — microservices architecture routing requests to LLM providers. Each service is independently deployable with its own database.

## MICROSERVICES

| Service | Responsibility | Port |
|---------|-------------|------|
| gateway-service | HTTP entry, middleware orchestration | :8080 |
| auth-service | Identity, API keys, model authorization | :50051 |
| router-service | Route model → provider | :50052 |
| provider-service | Provider CRUD, request forwarding | :50053 |
| billing-service | Usage tracking, budgets | :50054 |
| monitor-service | Metrics, alerting | :50055 |

## SERVICE COMMUNICATION

- **External** → gateway-service: HTTPS/REST
- **Internal**: gRPC between services
- **Callbacks**: provider-service → billing/monitor (async observer pattern)

## ARCHITECTURE

```
Consumers → gateway-service → auth-service → router-service → provider-service → LLM Providers
                ↑              ↓              ↓
            Admin UI        billing      monitor
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| OpenSpec specs | `openspec/specs/{service}/` | Design + API contracts |
| Service impl | `{service-name}/` | Service root folders |
| Router impl | `router-service/` | Existing Go implementation |
| Config | `configs/config.yaml` | Runtime config |

## DESIGN LOCATIONS

- **System architecture**: `openspec/specs/system/spec.md`
- **Service specs**: `openspec/specs/{service}/spec.md`
- **API contracts**: `openspec/specs/{service}/api-contracts.md`

## ANTI-PATTERNS

- NO direct database access across service boundaries
- NO hardcoded credentials (use config/env)
- NO URL-based provider detection

## CONVENTIONS

- Model naming: `{provider}:{model}` (e.g., `ollama:llama2`)
- Each service owns its database exclusively
- Cross-service data flows through gRPC APIs
- **Unit Test Completion**: Mark unit test tasks as complete only after running tests and confirming they pass. Do not mark tests as complete before implementation or execution.
- **Integration Test Acceptance**: Integration test acceptance criteria must include running `make down && make clean-images && make up` and verifying all services run as expected before marking tests as complete.
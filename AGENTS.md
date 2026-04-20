# PROJECT KNOWLEDGE BASE

**Generated:** 2026-04-20
**Commit:** 16e19e6
**Branch:** main

## OVERVIEW
AI API Gateway router service. Go-based HTTP API with Gin framework, routing requests to LLM providers (OpenAI, Ollama). Kubernetes-deployable via Helm.

## STRUCTURE
```
./
├── router-service/       # Main Go service
│   ├── cmd/server/    # Entry point
│   ├── internal/     # Domain logic (DDD)
│   │   ├── application/
│   │   ├── domain/
│   │   ├── handler/
│   │   └── infrastructure/
│   ├── configs/      # Config files
│   └── charts/      # Helm charts
├── openspec/          # OpenSpec workflows
└── docs/            # Documentation
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Add route | `internal/domain/port/router.go` | Define interface |
| Add provider | `internal/infrastructure/provider/` | Implement provider |
| Change config | `configs/config.yaml` | Runtime config |
| Deploy | `router-service/Makefile` | Docker/K8s build |

## CONVENTIONS
- DDD: domain -> application -> infrastructure layers
- Providers implement `port.Provider` interface
- Config via `configs/config.yaml` (mounted at `/app/config`)

## ANTI-PATTERNS (THIS PROJECT)
- NO direct provider calls in handler layer
- NO hardcoded API keys (use config/env)

## COMMANDS
```bash
cd router-service && make build          # Docker build
make deploy-docker                      # Run locally
make deploy-kind                      # Deploy to KinD
```

## NOTES
- 12 Go source files, ~1160 LOC total
- Uses Gin HTTP framework
- No test files yet
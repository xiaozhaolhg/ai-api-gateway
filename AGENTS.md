# PROJECT KNOWLEDGE BASE

**Generated:** 2026-04-21
**Commit:** 7922dbe
**Branch:** main

## OVERVIEW
AI API Gateway router service. Go-based HTTP API with Gin framework, routing requests to LLM providers (Ollama, OpenCode Zen). Uses factory pattern for extensible provider registration. Kubernetes-deployable via Docker.

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
│   │       ├── config/
│   │       └── provider/
│   ├── configs/      # Config files
│   └── Makefile      # Build/deploy scripts
├── openspec/          # OpenSpec workflows
└── docs/            # Documentation
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Add route | `internal/domain/port/router.go` | Define interface |
| Add provider | `internal/infrastructure/provider/` | Implement provider + factory |
| Add factory | `internal/infrastructure/provider/factory.go` | Implement ProviderFactory interface |
| Change config | `router-service/configs/config.yaml` | Runtime config (map-based) |
| Deploy | `router-service/Makefile` | Docker/K8s build |

## CONVENTIONS
- DDD: domain -> application -> infrastructure layers
- Providers implement `port.Provider` interface
- Factories implement `ProviderFactory` interface
- Provider registration via factory pattern in handler
- Config via `configs/config.yaml` (map-based provider structure)
- Model naming: `{provider_name}:{model_name}` (e.g., `ollama:llama2`)

## FACTORY PATTERN
The router uses a factory pattern for provider registration:
1. **ProviderFactory Interface**: Type(), Create(), Validate(), Defaults(), Description()
2. **ProviderRegistry**: Thread-safe registry for factory management
3. **Built-in Factories**: OllamaFactory, OpenCodeZenFactory
4. **Handler Integration**: Registry initialized in Setup(), providers created from config

### Adding a New Provider
1. Implement `port.Provider` interface in `internal/infrastructure/provider/<name>.go`
2. Implement `ProviderFactory` interface with factory struct
3. Register factory in handler Setup() function
4. Add provider config to `configs/config.yaml` under `provider.providers`
5. Model prefix will be `{factory_type}:{model_name}`

## ANTI-PATTERNS (THIS PROJECT)
- NO direct provider calls in handler layer
- NO hardcoded API keys (use config/env)
- NO URL-based provider detection (use factory pattern)
- NO hardcoded provider instantiation in handler

## COMMANDS
```bash
cd router-service && make build          # Docker build
make deploy-docker                      # Run locally
make deploy-kind                        # Deploy to KinD
```

## API ENDPOINTS
- `POST /v1/chat/completions` - Chat completion (streaming/non-streaming)
- `GET /v1/models` - List available models
- `GET /v1/providers` - List registered providers and status
- `GET /health` - Health check

## CONFIG STRUCTURE
```yaml
server:
  port: "8080"
provider:
  providers:
    ollama:
      enabled: true
      endpoint: "http://localhost:11434"
      api_key: ""
    opencode_zen:
      enabled: true
      endpoint: "https://opencode.ai/zen"
      api_key: "${API_KEY}"
```

## NOTES
- Factory pattern enables extensible provider registration
- Thread-safe provider registry with sync.RWMutex
- Map-based config for dynamic provider support
- Provider discovery via /v1/providers endpoint
- Breaking change from v1.x: config structure and model prefixes
- See docs/MIGRATION.md for migration guide

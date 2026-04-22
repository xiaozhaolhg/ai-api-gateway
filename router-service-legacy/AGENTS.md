# ROUTER SERVICE

Go microservice with Gin HTTP framework. DDD architecture.

## LAYERS
```
internal/
├── application/     # Orchestration (service)
├── domain/         # Business logic (port interfaces)
├── handler/        # HTTP handlers
└── infrastructure/ # External integrations (provider impls)
```

## PACKAGES
| Package | Role |
|---------|------|
| `application/service` | Orchestration, business logic |
| `domain/port` | Interfaces (Router, Provider) |
| `domain/entity` | Domain models |
| `handler` | HTTP handlers (Gin) |
| `infrastructure/provider` | LLM provider implementations |

## INTERFACES
- `port.Provider`: Chat completion, streaming, list models
- `port.Router`: Provider selection by model name

## ADD PROVIDER
1. Define in `port/provider.go` (new methods)
2. Implement in `infrastructure/provider/<name>.go`
3. Register in `infrastructure/config/config.go`

## ADD ROUTE
1. Define interface in `domain/port/router.go`
2. Implement in `application/service/router.go`
3. Wire in `handler/handler.go`
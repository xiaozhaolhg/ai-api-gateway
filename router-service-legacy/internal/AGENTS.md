# INTERNAL - DOMAIN LOGIC

Core DDD layers. All business logic lives here.

## SUBDIRS
| Dir | Files | Purpose |
|-----|-------|---------|
| `application/` | 4 | Service orchestration |
| `domain/` | 2 | Interfaces, entities |
| `handler/` | 1 | HTTP layer (thin) |
| `infrastructure/` | 3 | External integrations |

## DOMAIN RULES
- Providers implement `port.Provider`
- All business logic in application layer
- Handler = thin HTTP adapter only
- Config loaded in infrastructure/config

## KEY FILES
| File | Purpose |
|------|---------|
| `port/router.go` | Router interface |
| `port/provider.go` | Provider interface + types |
| `infrastructure/config/config.go` | Config loading |
## Why

The AI API gateway needs a core router-service to unify access to multiple LLM providers through a single, OpenAI-compatible API. Currently, clients must manage separate credentials and API endpoints for each provider (OpenAI, Anthropic, local Ollama, etc.). A centralized router service enables intelligent request routing, provider failover, and a consistent interface—all while running locally with Ollama and OpenCode Zen as the initial providers.

## What Changes

- Create a new Go 1.26 service (`router-service`) following Domain-Driven Design and Clean Architecture
- Implement OpenAI-compatible endpoints (`/v1/chat/completions`, `/v1/models`)
- Support dual providers: local Ollama and OpenCode Zen
- Add streaming support (SSE) for real-time token delivery
- Containerize with multi-stage Docker build and deploy to local KinD cluster via Helm chart
- Mount configuration via YAML file (no persistence required)

## Capabilities

### New Capabilities
- `router-service-architecture`: Internal structure, layers (domain/application/infrastructure/interface), provider abstractions, routing logic
- `router-service-deployment`: Dockerfile, Helm chart, Kubernetes manifests, KinD deployment workflow
- `router-service-testing`: Test strategy, unit tests for domain logic, integration tests for provider adapters, e2e tests for API compatibility

### Modified Capabilities
- None (greenfield service, no existing specs to modify)

## Impact

- **New service**: `router-service/` directory with Go code following Clean Architecture
- **Configuration**: YAML-based config mounted from host for provider endpoints and routing rules
- **Dependencies**: Go 1.26, Gin HTTP framework, httplib for upstream calls
- **Deployment**: Docker image (~15MB), KinD cluster with NGINX Ingress via Helm
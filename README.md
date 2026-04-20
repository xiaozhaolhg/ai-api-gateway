# AI API Gateway

A Go-based HTTP API gateway for routing LLM requests to multiple providers (OpenAI, Ollama, OpenCode Zen). Built with the Gin framework and deployed via Kubernetes.

## Features

- **Multi-provider support**: Route requests to Ollama, OpenCode Zen, and other providers
- **Factory pattern**: Extensible provider registration system
- **Model routing**: Automatic routing based on model name prefixes
- **Streaming support**: Server-sent events for streaming responses
- **Provider discovery**: `/v1/providers` endpoint for provider status and configuration

## Quick Start

### Configuration

The router uses a YAML configuration file. Example:

```yaml
server:
  port: "8080"
provider:
  providers:
    ollama:
      enabled: true
      endpoint: "http://host.docker.internal:11434"
      api_key: ""
    opencode_zen:
      enabled: true
      endpoint: "https://opencode.ai/zen"
      api_key: "${OPCODE_API_KEY}"
```

### Model Naming

Models are prefixed with the provider name:
- Ollama: `ollama:llama2`, `ollama:mistral`
- OpenCode Zen: `opencode_zen:gpt-4`, `opencode_zen:claude-3`

### API Endpoints

- `POST /v1/chat/completions` - Chat completion (streaming and non-streaming)
- `GET /v1/models` - List available models
- `GET /v1/providers` - List registered providers and their status
- `GET /health` - Health check

### Example Usage

```bash
# Chat completion with Ollama
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ollama:llama2",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# List models
curl http://localhost:8080/v1/models

# List providers
curl http://localhost:8080/v1/providers
```

## Development

### Build

```bash
cd router-service
make build
```

### Run locally

```bash
cd router-service
make deploy-docker
```

### Deploy to Kubernetes

```bash
cd router-service
make deploy-kind
```

## Breaking Changes

### v2.0.0 - Factory Pattern Migration

- **Config structure**: Changed from struct-based to map-based provider configuration
- **Model prefix**: OpenCode provider renamed from `opencode:` to `opencode_zen:`
- **Provider detection**: Changed from URL-based to factory-based registration

See [MIGRATION.md](docs/MIGRATION.md) for migration guide.

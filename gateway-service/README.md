# AI API Gateway

A high-performance gateway service that routes AI model requests to multiple LLM providers with unified authentication, rate limiting, and billing.

## Features

- **Multi-Provider Support**: Route requests to Ollama, OpenAI, and custom providers
- **Unified API**: OpenAI-compatible `/v1/chat/completions` endpoint
- **Streaming**: Server-Sent Events (SSE) for real-time responses
- **Authentication**: JWT-based auth with API key support
- **Rate Limiting**: Per-user and per-model rate limits
- **Billing**: Usage tracking and cost estimation
- **Real-Time Streaming Usage**: Record usage at configurable intervals during streaming to prevent cost overruns
- **Health Monitoring**: Deep health checks with dependency status
- **Graceful Shutdown**: Zero-downtime deployments

## Quick Start

### Prerequisites

- Go 1.21+
- Access to backend services (auth, router, provider, billing)

### Run Locally

```bash
# Install dependencies
go mod download

# Run the service
go run cmd/server/main.go

# Or build and run
go build -o gateway cmd/server/main.go
./gateway
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `HOST` | `0.0.0.0` | HTTP server host |
| `AUTH_SERVICE_ADDRESS` | `localhost:50051` | Auth service gRPC address |
| `ROUTER_SERVICE_ADDRESS` | `localhost:50052` | Router service gRPC address |
| `PROVIDER_SERVICE_ADDRESS` | `localhost:50053` | Provider service gRPC address |
| `BILLING_SERVICE_ADDRESS` | `localhost:50054` | Billing service gRPC address |
| `MONITOR_SERVICE_ADDRESS` | `localhost:50055` | Monitor service gRPC address |

### Configuration

Configuration is loaded from `configs/config.yaml`:

```yaml
server:
  port: "8080"
  host: "0.0.0.0"

auth_service:
  address: "auth-service:50051"

streaming_token_interval: 1000  # Record usage every 1000 completion tokens (0 = disable)
```

| Variable | Default | Description |
|----------|---------|-------------|
| `STREAMING_TOKEN_INTERVAL` | `1000` | Override streaming token interval (0 = disable) |
| `PORT` | `8080` | HTTP server port |
| `HOST` | `0.0.0.0` | HTTP server host |
| `AUTH_SERVICE_ADDRESS` | `localhost:50051` | Auth service gRPC address |
| `ROUTER_SERVICE_ADDRESS` | `localhost:50052` | Router service gRPC address |
| `PROVIDER_SERVICE_ADDRESS` | `localhost:50053` | Provider service gRPC address |
| `BILLING_SERVICE_ADDRESS` | `localhost:50054` | Billing service gRPC address |
| `MONITOR_SERVICE_ADDRESS` | `localhost:50055` | Monitor service gRPC address |

```yaml
# Full example config.yaml:
server:
  port: "8080"
  host: "0.0.0.0"

auth_service:
  address: "auth-service:50051"
router_service:
  address: "router-service:50052"
provider_service:
  address: "provider-service:50053"
billing_service:
  address: "billing-service:50054"
monitor_service:
  address: "monitor-service:50055"

streaming_token_interval: 1000  # Record usage every 1000 completion tokens (0 disables intermediate recording)
```

## API Endpoints

### Health & Status

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Simple liveness check |
| `/gateway/health` | GET | Deep health check with dependency status |

**Example: Deep Health Check**
```bash
curl http://localhost:8080/gateway/health
```

```json
{
  "status": "healthy",
  "gateway": "healthy",
  "services": {
    "auth": {"status": "healthy", "latency": "15ms"},
    "router": {"status": "healthy", "latency": "8ms"},
    "provider": {"status": "healthy", "latency": "23ms"},
    "billing": {"status": "healthy", "latency": "12ms", "optional": true}
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Models

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/gateway/models` | GET | List all available models from all providers |
| `/v1/models` | GET | OpenAI-compatible models list |

**Example: List Models**
```bash
curl http://localhost:8080/gateway/models
```

```json
{
  "object": "list",
  "data": [
    {
      "id": "ollama:llama2",
      "object": "model",
      "created": 1705312200,
      "owned_by": "Ollama"
    },
    {
      "id": "openai:gpt-4",
      "object": "model",
      "created": 1705312200,
      "owned_by": "OpenAI"
    }
  ]
}
```

### Chat Completions

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/chat/completions` | POST | Create chat completion (OpenAI-compatible) |

**Request:**
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "ollama:llama2",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ],
    "max_tokens": 150,
    "temperature": 0.7
  }'
```

**Bare Model Name (New):**
You can also use bare model names without the provider prefix:
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "llama2",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ],
    "max_tokens": 150,
    "temperature": 0.7
  }'
```
The gateway will automatically resolve the bare model name to an available provider.

**Response:**
```json
{
  "id": "chatcmpl-abc123",
  "object": "chat.completion",
  "created": 1705312200,
  "model": "ollama:llama2",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "I'm doing well, thank you! How can I help you today?"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 12,
    "completion_tokens": 15,
    "total_tokens": 27
  }
}
```

### Streaming

Add `"stream": true` to request for SSE streaming:

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "ollama:llama2",
    "messages": [{"role": "user", "content": "Tell me a story"}],
    "stream": true
  }'
```

**SSE Response:**
```
data: {"choices":[{"delta":{"content":"Once"}}]}

data: {"choices":[{"delta":{"content":" upon"}}]}

data: {"choices":[{"delta":{"content":" a time"}}]}

...

data: {"prompt_tokens":10,"completion_tokens":150,"total_tokens":160,"done":true}
```

### Streaming Usage Tracking

For streaming requests, the gateway records usage at configurable intervals (default every 1000 completion tokens) to prevent cost overruns. This ensures real-time visibility into ongoing usage costs.

| Setting | Default | Description |
|---------|---------|-------------|
| `streaming_token_interval` (YAML) | `1000` | Tokens between intermediate usage recordings (0 = disable) |
| `STREAMING_TOKEN_INTERVAL` (env) | `1000` | Override interval via environment variable |

**Behavior:**
- Intermediate usage is recorded as tokens accumulate to the configured interval
- A final usage record is always sent after stream completion
- If the stream errors mid-way, all accumulated tokens are still recorded
- All records are aggregated by the billing service for accurate cost tracking

### Admin Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/admin/auth/login` | POST | User login |
| `/admin/auth/register` | POST | User registration |
| `/admin/auth/me` | GET | Get current user |
| `/admin/auth/users` | GET | List all users (admin only) |
| `/admin/providers` | GET | List providers |
| `/admin/providers` | POST | Create provider |
| `/admin/providers/:id` | PUT | Update provider |
| `/admin/providers/:id` | DELETE | Delete provider |
| `/admin/auth/usage` | GET | Get usage statistics |

## Error Handling

The gateway uses structured error responses with HTTP status codes:

| Status | Error Code | Description |
|--------|------------|-------------|
| 400 | `bad_request` | Invalid request parameters |
| 401 | `invalid_api_key` | Missing or invalid API key |
| 403 | `insufficient_permissions` | Insufficient permissions |
| 404 | `model_not_found` | Model not found |
| 429 | `rate_limit_exceeded` | Rate limit exceeded |
| 502 | `provider_error` | Provider service error |
| 503 | `service_unavailable` | Service unavailable |
| 504 | `gateway_timeout` | Request timeout |
| 500 | `internal_error` | Internal server error |

**Error Response Format:**
```json
{
  "error": {
    "code": "model_not_found",
    "message": "Model 'gpt-5' not found",
    "details": "Available models: ollama:llama2, openai:gpt-4"
  }
}
```

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Client    │────▶│   Gateway    │────▶│    Auth     │
│             │◄────│   Service    │◄────│   Service   │
└─────────────┘     └──────────────┘     └─────────────┘
                            │
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
   ┌─────────┐       ┌──────────┐       ┌──────────┐
   │ Router  │────▶ │ Provider │────▶  │  LLM     │
   │ Service │       │ Service  │       │ Provider│
   └─────────┘       └──────────┘       └──────────┘
        │                  │
        ▼                  ▼
   ┌─────────┐       ┌──────────┐
   │ Billing │       │ Monitor  │
   │ Service │       │ Service  │
   └─────────┘       └──────────┘
```

## Logging

The gateway uses structured JSON logging with:
- **Request ID**: Unique ID for each request (propagated via `X-Request-ID` header)
- **Correlation ID**: For tracing across services
- **Sensitive Data Masking**: Passwords, tokens, API keys are redacted
- **Log Level**: Configurable (debug, info, warn, error)

**Example Log Entry:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/v1/chat/completions",
  "status_code": 200,
  "duration": "1.23s",
  "user_id": "user-123",
  "provider_id": "ollama",
  "model": "llama2",
  "client_ip": "192.168.1.1"
}
```

## Testing

### Unit Tests

```bash
go test ./...
```

### Load Tests (k6)

```bash
cd tests/load
k6 run gateway_load_test.js
```

See [tests/load/README.md](tests/load/README.md) for details.

## Provider Adapters

To add a new LLM provider, implement a provider adapter:

```go
type ProviderAdapter interface {
    ForwardRequest(ctx context.Context, req *ForwardRequest) (*ProviderResponse, error)
    StreamRequest(ctx context.Context, req *ForwardRequest) (<-chan StreamChunk, error)
    HealthCheck(ctx context.Context) error
}
```

See [Provider Adapter Guide](../docs/provider-adapter-guide.md) for complete documentation.

## Deployment

### Docker

```bash
docker build -t gateway-service .
docker run -p 8080:8080 gateway-service
```

### Kubernetes

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

### Graceful Shutdown

The gateway handles SIGINT and SIGTERM signals:
1. Stops accepting new requests
2. Waits for active requests to complete (max 30s)
3. Closes gRPC connections
4. Exits cleanly

## Monitoring

### Health Endpoints

- `/health` - Liveness probe
- `/gateway/health` - Readiness probe with dependency checks

### Metrics

The gateway exposes metrics via the monitor service:
- Request count and latency
- Error rates
- Token usage
- Provider availability

## Development

### Project Structure

```
gateway-service/
├── cmd/server/          # Application entry point
├── internal/
│   ├── client/          # gRPC clients for backend services
│   ├── errors/          # Error types and handling
│   ├── handler/         # HTTP handlers
│   ├── infrastructure/  # Config, logging, etc.
│   ├── middleware/      # Gin middleware (auth, logging, proxy)
│   └── util/            # Utility functions
├── configs/             # Configuration files
├── tests/               # Test suites
└── docs/                # Documentation
```

### Adding a New Endpoint

1. Create handler in `internal/handler/`
2. Add route in `cmd/server/main.go`
3. Add tests
4. Update documentation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a PR

## License

MIT License - see LICENSE file for details.

## Support

- Issues: [GitHub Issues](https://github.com/ai-api-gateway/issues)
- Docs: [Wiki](https://github.com/ai-api-gateway/wiki)
- Discussions: [GitHub Discussions](https://github.com/ai-api-gateway/discussions)

# AI API Gateway

Enterprise microservices gateway for routing LLM requests to multiple providers. Includes gRPC backend services and React admin UI.

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Admin UI  │────▶│  Gateway   │────▶│  Router    │
│  (React)   │     │ Service    │     │ Service    │
└─────────────┘     └──────┬──────┘     └─────────────┘
                         │
            ┌────────────┼────────────┐
            ▼           ▼           ▼
      ┌──────────┐ ┌──────────┐ ┌──────────┐
      │   Auth   │ │Provider │ │ Billing │
      │ Service  │ │ Service │ │ Service │
      └──────────┘ └──────────┘ └──────────┘
```

## Services

| Service | Port | Protocol | Responsibility |
|---------|------|----------|----------------|
| gateway-service | 8080 | HTTP/REST | HTTP entry, middleware orchestration |
| auth-service | 50051 | gRPC | Identity, API keys, model authorization |
| router-service | 50052 | gRPC | Route model → provider |
| provider-service | 50053 | gRPC | Provider CRUD, request forwarding |
| billing-service | 50054 | gRPC | Usage tracking, budgets |
| monitor-service | 50055 | gRPC | Metrics, alerting |
| admin-ui | 3000 | HTTP | Admin dashboard |
| redis | 6379 | - | Session cache |

## Quick Start

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- Node.js 18+ (for admin-ui)

### Start All Services

```bash
# Build and start all containers
make up

# Or start without rebuilding
docker compose up -d --no-build
```

### Verify Services

```bash
# Check running containers
docker compose ps

# Test health endpoints
curl http://localhost:8080/health
curl http://localhost:3000/health
```

## Development

### Build

```bash
# Build all services
make build

# Build single service
make -C auth-service build
```

### Test

```bash
# Test all services
make test

# Test single service
cd auth-service && go test ./...
```

### Validate Changes

After making code changes, validate with Docker:

```bash
# Stop all containers
make down

# Clean and rebuild images
make clean-images
make up

# Verify services are healthy
curl http://localhost:8080/health
curl http://localhost:3000/health
```

### Run Locally (without Docker)

```bash
# Start Redis
docker run -d -p 6379:6379 redis:7-alpine

# Start auth-service
cd auth-service && go run ./cmd/server

# Start gateway-service
cd gateway-service && go run ./cmd/server

# Start admin-ui
cd admin-ui && npm run dev
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| AUTH_SERVICE_ADDRESS | Auth service gRPC endpoint | localhost:50051 |
| ROUTER_SERVICE_ADDRESS | Router service gRPC endpoint | localhost:50052 |
| PROVIDER_SERVICE_ADDRESS | Provider service gRPC endpoint | localhost:50053 |

### Provider Config

Edit `configs/config.yaml`:

```yaml
server:
  port: "8080"
provider:
  providers:
    ollama:
      enabled: true
      endpoint: "http://host.docker.internal:11434"
    opencode_zen:
      enabled: true
      endpoint: "https://opencode.ai/zen"
```

## API Endpoints

### Gateway Service (port 8080)

| Endpoint | Method | Description |
|----------|--------|-------------|
| /health | GET | Health check |
| /v1/chat/completions | POST | Chat completion |
| /v1/models | GET | List models |
| /v1/providers | GET | List providers |

### Admin UI (port 3000)

| Route | Description |
|-------|-------------|
| /admin/login | Login page |
| /admin/dashboard | Main dashboard |
| /admin/providers | Provider management |
| /admin/users | User management |
| /admin/api-keys | API key management |
| /admin/usage | Usage analytics |
| /admin/health | Service health |

## Auth Flow

### Login

```bash
curl -X POST http://localhost:8080/admin/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"securepass"}'
```

Returns JWT cookie (HTTP-only, Secure, /admin path)

## Model Naming

Provider prefix format: `{provider}:{model}` (e.g., `ollama:llama2`)

- Ollama: `ollama:llama2`, `ollama:mistral`
- OpenCode Zen: `opencode_zen:gpt-4`, `opencode_zen:claude-3`

### Bare Model Names (New)

You can now use bare model names without the provider prefix:
- `llama2` → automatically resolves to an available provider
- `gpt-4` → automatically resolves to an available provider

When multiple providers support the same model, the healthiest provider is selected automatically, with the remaining healthy providers used as fallbacks.

**Example:**
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-xxx..." \
  -d '{
    "model": "llama2",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

## Testing

```bash
# Run all tests
make test

# Run auth-service tests with coverage
cd auth-service && go test -cover ./...
```

## Docker

### Build Images

```bash
docker build -t ai-api-gateway/auth-service:latest ./auth-service
docker build -t ai-api-gateway/gateway-service:latest ./gateway-service
docker build -t ai-api-gateway/admin-ui:latest ./admin-ui
```

### Logs

```bash
docker compose logs -f auth-service
docker compose logs -f gateway-service
```

## Project Structure

```
.
├── api/                    # Shared protobuf definitions
├── auth-service/           # Identity & authentication
├── gateway-service        # HTTP gateway
├── router-service         # Request routing
├── provider-service      # Provider management
├── billing-service      # Usage tracking
├── monitor-service       # Metrics & alerting
├── admin-ui/             # React admin dashboard
├── configs/              # Configuration files
├── openspec/             # Change specifications
└── go.mod               # Unified Go module
```

## License

MIT
## User Flow

### Register & Login

1. **Register** (if new user):
   ```bash
   curl -X POST http://localhost:8080/admin/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username":"myuser","email":"user@example.com","name":"My User","password":"mypass"}'
   ```

2. **Login** to get JWT token:
   ```bash
   curl -X POST http://localhost:8080/admin/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"mypass"}'
   ```
   Response includes `token` field.

### Create API Key

Use the JWT token to create an API key:

```bash
curl -X POST http://localhost:8080/v1/auth/api-keys \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"my-key"}'
```

Response:
```json
{
  "api_key_id": "abc123...",
  "api_key": "sk-xxx...",
  "name": "my-key"
}
```

**Important**: The `api_key` is shown only once! Save it securely.

### Use `/v1/chat/completions`

Use the API key to make requests:

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-xxx..." \
  -d '{
    "model": "ollama:llama2",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

### List Available Models

```bash
curl http://localhost:8080/v1/models \
  -H "Authorization: Bearer sk-xxx..."
```

### Manage API Keys

List your API keys:
```bash
curl http://localhost:8080/v1/auth/api-keys \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Delete an API key:
```bash
curl -X DELETE http://localhost:8080/v1/auth/api-keys/{key_id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

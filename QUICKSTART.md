# AI API Gateway — Quick Start Guide

## Overview

AI API Gateway is a microservices-based gateway that routes requests to multiple LLM providers with authentication, usage tracking, and monitoring.

## Quick Start Options

### Option 1: Single Binary Demo (Fastest)

Build and run the gateway with embedded admin UI:

```bash
# Build single binary with embedded UI
make build-single

# Run the gateway
./gateway-service/bin/gateway

# Access admin UI
# Open http://localhost:8080
# Default admin credentials: admin@example.com / password123
```

### Option 2: Docker Compose (Production)

Run all services with Docker Compose:

```bash
# Start all services
make up

# Access admin UI
# Open http://localhost:8080
# Default admin credentials: admin@example.com / password123

# View logs
docker compose logs -f gateway-service
```

## Environment Variables

| Variable | Default | Description |
|-----------|----------|-------------|
| `SERVER_PORT` | 8080 | Gateway HTTP port |
| `AUTH_SERVICE_ADDRESS` | auth-service:50051 | Auth service gRPC address |
| `ROUTER_SERVICE_ADDRESS` | router-service:50052 | Router service gRPC address |
| `PROVIDER_SERVICE_ADDRESS` | provider-service:50053 | Provider service gRPC address |
| `BILLING_SERVICE_ADDRESS` | billing-service:50054 | Billing service gRPC address |
| `MONITOR_SERVICE_ADDRESS` | monitor-service:50055 | Monitor service gRPC address |

## Admin UI Features

- **Dashboard**: System overview with metrics
- **Providers**: Manage LLM providers (OpenAI, Anthropic, Ollama, etc.)
- **Users**: User management and API keys
- **Groups**: User groups and permissions
- **Usage**: Token usage and cost tracking
- **Budgets**: Per-user spending limits
- **Pricing Rules**: Model pricing configuration
- **Routing Rules**: Model-to-provider routing
- **Alerts**: System alerts and notifications
- **Health**: Provider health monitoring

## API Usage

### 1. Get API Key

```bash
# Login to admin UI and create API key
# Or use default key: sk-test-key-12345
```

### 2. Make Chat Request

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer sk-test-key-12345" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ollama:llama2",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

### 3. List Available Models

```bash
curl http://localhost:8080/gateway/models \
  -H "Authorization: Bearer sk-test-key-12345"
```

## Development

### Build from Source

```bash
# Build all services
make build

# Run tests
make test

# Run admin UI tests
make test-ui

# Generate protobuf
make proto
```

### Development Mode

Run gateway with admin UI in development mode:

```bash
# Terminal 1: Start backend services
make up

# Terminal 2: Start admin UI dev server
cd admin-ui
npm run dev

# Access admin UI at http://localhost:3000
# API calls proxy to backend at http://localhost:8080
```

## Configuration

Configuration is loaded from `gateway-service/configs/config.yaml`. Key sections:

```yaml
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
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Service Connection Issues

```bash
# Check service logs
docker compose logs auth-service
docker compose logs router-service
docker compose logs provider-service
docker compose logs billing-service
docker compose logs monitor-service
```

### Build Issues

```bash
# Clean and rebuild
make clean
make build-single
```

## Next Steps

1. **Add Providers**: Configure your LLM providers in the admin UI
2. **Create Users**: Add users and generate API keys
3. **Set Budgets**: Configure spending limits per user
4. **Configure Pricing**: Set up pricing rules for models
5. **Monitor Usage**: Track token usage and costs through the dashboard

## Support

- **Documentation**: Full API docs at `/docs/`
- **Issues**: Report bugs on GitHub
- **Architecture**: See `docs/developer_completion_analysis.md` for implementation status

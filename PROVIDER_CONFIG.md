# Provider Configuration

## Adding Providers via Admin API

### Prerequisites

- Gateway service running on `http://localhost:8080`
- Valid JWT token (obtain via `/admin/auth/login`)

### Login to Get JWT Token

```bash
curl -X POST http://localhost:8080/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"securepass"}'
```

Response includes `token` field - use this as `AUTH_TOKEN`.

### Create Ollama Provider

```bash
curl -X POST http://localhost:8080/admin/providers \
  -H "Content-Type: application/json" \
  -H "Cookie: auth_token=AUTH_TOKEN" \
  -d '{
    "name": "Ollama Local",
    "type": "ollama",
    "base_url": "http://host.docker.internal:11434",
    "credentials": "{}"
  }'
```

### Create OpenAI-Compatible Provider (OpenCode Zen)

```bash
curl -X POST http://localhost:8080/admin/providers \
  -H "Content-Type: application/json" \
  -H "Cookie: auth_token=AUTH_TOKEN" \
  -d '{
    "name": "OpenCode Zen",
    "type": "opencode_zen",
    "base_url": "https://opencode.ai/zen",
    "credentials": "{\"api_key\": \"YOUR_API_KEY\"}"
  }'
```

### List Configured Providers

```bash
curl http://localhost:8080/admin/providers \
  -H "Cookie: auth_token=AUTH_TOKEN"
```

## Provider Types

| Type | Name | Base URL | Credentials |
|------|------|---------|-------------|
| `ollama` | Ollama Local | `http://host.docker.internal:11434` | `{}` (none needed) |
| `opencode_zen` | OpenCode Zen | `https://opencode.ai/zen` | `{"api_key": "..."}` |

## Verifying Provider Health

```bash
curl http://localhost:8080/admin/providers/{provider_id}/health \
  -H "Cookie: auth_token=AUTH_TOKEN"
```

## Model Naming

After adding a provider, models are automatically available using the format:

`{provider-type}:{model-name}`

Examples:
- Ollama: `ollama:llama2`, `ollama:mistral`
- OpenCode Zen: `opencode_zen:gpt-4`, `opencode_zen:claude-3`

## Troubleshooting

### Provider Not Found

If you get `failed to get provider: record not found`:

1. Check provider is created: `GET /admin/providers`
2. Verify `type` field matches the provider prefix in model name
3. Ensure provider is enabled in database

### Models Not Showing

If `GET /v1/models` returns empty:

1. Check provider has models configured
2. Verify provider health: `GET /admin/providers/{id}/health`
3. Check provider-service logs: `docker logs ai-api-gateway-provider-service-1`

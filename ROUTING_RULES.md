# Routing Rule Setup

## Model Naming Convention

Models follow the format: `{provider}:{model-name}`

Examples:
- `ollama:llama2` - Llama 2 model on Ollama
- `opencode_zen:gpt-4` - GPT-4 on OpenCode Zen

## How Routing Works

When a request is made to `/v1/chat/completions`, the `RouteMiddleware` extracts the provider from the model name by splitting on `:`.

For example:
- Model: `ollama:llama2` → Provider: `ollama`
- Model: `opencode_zen:gpt-4` → Provider: `opencode_zen`

## Setting Up Routing Rules

### Option 1: Direct Database Insert

Connect to the router-service database and insert routing rules:

```sql
INSERT INTO routing_rules (model_pattern, provider_id, priority) VALUES
('gpt-4', 'openai', 100),
('gpt-3.5-*', 'openai', 90),
('claude-*', 'anthropic', 100),
('llama*', 'ollama', 80);
```

### Option 2: Via Admin API (if available)

```bash
curl -X POST http://localhost:8080/admin/routing-rules \
  -H "Content-Type: application/json" \
  -H "Cookie: auth_token=YOUR_JWT" \
  -d '{
    "model_pattern": "gpt-4",
    "provider_id": "openai",
    "priority": 100
  }'
```

## Example Routing Rules

| Model Pattern | Provider | Priority | Description |
|---------------|----------|---------|-------------|
| `gpt-4` | `openai` | 100 | Exact match for GPT-4 |
| `gpt-3.5-*` | `openai` | 90 | Wildcard for GPT-3.5 models |
| `claude-*` | `anthropic` | 100 | All Claude models |
| `llama*` | `ollama` | 80 | All Llama models on Ollama |

## Verifying Routing

After setting up routing rules, test with:

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{
    "model": "ollama:llama2",
    "messages": [{"role": "user", "content": "Hello"}]
  }'
```

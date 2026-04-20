# Migration Guide: Factory Pattern (v2.0.0)

This guide helps you migrate from the old URL-based provider detection to the new factory pattern architecture.

## Breaking Changes

### 1. Config Structure Change

**Old Config (v1.x):**
```yaml
server:
  port: "8080"
provider:
  ollama:
    enabled: true
    endpoint: "http://host.docker.internal:11434"
    api_key: ""
  opencode_zen:
    enabled: false
    endpoint: "https://api.opencode.com"
    api_key: ""
```

**New Config (v2.0):**
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
      api_key: ""
```

**Key Changes:**
- Provider configuration is now nested under `provider.providers`
- Structure is map-based (dynamic) instead of struct-based (static)
- OpenCode Zen endpoint updated to `https://opencode.ai/zen`

### 2. Model Prefix Change

**Old Model Names:**
- OpenCode models: `opencode:gpt-4`, `opencode:claude-3`

**New Model Names:**
- OpenCode models: `opencode_zen:gpt-4`, `opencode_zen:claude-3`

**Key Changes:**
- OpenCode provider name changed from `opencode` to `opencode_zen`
- Model prefix changed from `opencode:` to `opencode_zen:`
- Ollama models unchanged: `ollama:llama2`, `ollama:mistral`

### 3. Provider Detection Change

**Old Behavior (v1.x):**
- Providers detected based on URL patterns (localhost:11434, host.docker.internal:11434, etc.)
- Hardcoded URL matching in handler

**New Behavior (v2.0):**
- Providers instantiated via factory pattern based on config keys
- Registry-based provider creation
- Extensible without code changes

## Migration Steps

### Step 1: Backup Current Config

```bash
cp router-service/configs/config.yaml router-service/configs/config.yaml.backup
```

### Step 2: Update Config Structure

Convert your config from the old structure to the new map-based structure:

```yaml
# Old
provider:
  ollama:
    enabled: true
    endpoint: "http://localhost:11434"
    api_key: ""
  opencode_zen:
    enabled: false
    endpoint: "https://api.opencode.com"
    api_key: ""

# New
provider:
  providers:
    ollama:
      enabled: true
      endpoint: "http://localhost:11434"
      api_key: ""
    opencode_zen:
      enabled: true
      endpoint: "https://opencode.ai/zen"
      api_key: ""
```

### Step 3: Update Client Model Names

Find and replace model prefixes in your client code:

```bash
# Old
model: "opencode:gpt-4"

# New
model: "opencode_zen:gpt-4"
```

### Step 4: Deploy New Version

```bash
cd router-service
make build
make deploy-kind
```

### Step 5: Verify Deployment

Check that providers are initialized correctly:

```bash
# Check provider status
curl http://localhost:8080/v1/providers

# Check models
curl http://localhost:8080/v1/models

# Test chat completion
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "ollama:llama2",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

## Rollback Procedure

If you need to rollback to the previous version:

```bash
# Restore old config
cp router-service/configs/config.yaml.backup router-service/configs/config.yaml

# Deploy previous version (if available)
# or rebuild from previous commit
git checkout <previous-commit>
cd router-service
make build
make deploy-kind
```

## Troubleshooting

### Issue: "Unknown provider type" error

**Cause:** Typo in provider config key

**Solution:** Ensure provider keys match registered factory types:
- `ollama` (correct)
- `opencode_zen` (correct)
- `opencode` (incorrect - old name)

### Issue: Models not appearing

**Cause:** Provider not enabled in config

**Solution:** Set `enabled: true` for the provider in config

### Issue: "opencode:" models not working

**Cause:** Using old model prefix

**Solution:** Update to `opencode_zen:` prefix

## New Features

### /v1/providers Endpoint

New endpoint for provider discovery and status:

```bash
curl http://localhost:8080/v1/providers
```

Response:
```json
{
  "providers": [
    {
      "type": "ollama",
      "description": "Ollama - Run LLMs locally",
      "configured": true,
      "enabled": true,
      "endpoint": "http://localhost:11434",
      "has_api_key": false,
      "defaults": {
        "endpoint": "http://localhost:11434",
        "enabled": false
      }
    },
    {
      "type": "opencode_zen",
      "description": "OpenCode Zen - Curated AI models for coding agents",
      "configured": true,
      "enabled": true,
      "endpoint": "https://opencode.ai/zen",
      "has_api_key": true,
      "defaults": {
        "endpoint": "https://opencode.ai/zen",
        "enabled": false
      }
    }
  ]
}
```

### Extensibility

The new factory pattern makes it easy to add new providers without modifying core handler code. See AGENTS.md for details on adding custom providers.
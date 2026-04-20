## Why

Current provider registration uses fragile URL-based matching to determine which provider to instantiate. This breaks when Ollama runs on different ports, prevents multiple provider instances, and requires code changes for each new provider. The architecture needs to be extensible for future provider integrations without modifying core registration logic.

## What Changes

- **BREAKING**: Replace URL-based provider detection with factory pattern registration
- Implement ProviderFactory interface with Type(), Create(), Validate(), Defaults(), Description() methods
- Add ProviderRegistry for factory management and provider instantiation
- Convert ProviderConfig from struct to map-based structure for dynamic provider support
- Create OllamaFactory and OpenCodeZenFactory implementations
- Update handler to use registry instead of hardcoded switch statement
- Add /v1/providers endpoint for provider discovery and status
- Update config.yaml to use new map-based provider structure
- Change OpenCode provider name from "opencode" to "opencode_zen" for consistency
- Change OpenCode model prefix from "opencode:" to "opencode_zen:"

## Capabilities

### New Capabilities
- `router-service-provider-factory`: Extensible provider registration system using factory pattern

### Modified Capabilities
- None (implementation-only change, no spec-level requirement changes)

## Impact

**Affected Code:**
- `internal/handler/handler.go` - Replace URL matching with registry-based instantiation
- `internal/infrastructure/config/config.go` - Convert ProviderConfig to map-based structure
- `internal/infrastructure/provider/` - Add factory.go, update ollama.go and opencode.go with factories
- `router-service/configs/config.yaml` - Update to new map-based structure

**API Changes:**
- Add GET /v1/providers endpoint for provider discovery
- Model prefix changes: opencode:* → opencode_zen:*

**Breaking Changes:**
- Config structure change (requires config migration)
- Model prefix change for OpenCode models
- No backwards compatibility (clean break to new architecture)
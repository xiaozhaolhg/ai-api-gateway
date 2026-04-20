## 1. Factory Interface and Registry

- [x] 1.1 Create ProviderFactory interface in internal/infrastructure/provider/factory.go with Type(), Create(), Validate(), Defaults(), Description() methods
- [x] 1.2 Implement ProviderRegistry in internal/infrastructure/provider/factory.go with Register(), Create(), ListTypes(), GetFactory(), GetDefaults(), ValidateSettings() methods
- [x] 1.3 Add thread-safe operations to ProviderRegistry using sync.RWMutex
- [x] 1.4 Add unit tests for ProviderFactory interface and ProviderRegistry in internal/infrastructure/provider/factory_test.go

## 2. Ollama Factory Implementation

- [x] 2.1 Create OllamaFactory struct in internal/infrastructure/provider/ollama.go
- [x] 2.2 Implement Type() method to return "ollama"
- [x] 2.3 Implement Description() method to return "Ollama - Run LLMs locally"
- [x] 2.4 Implement Validate() method to require non-empty endpoint
- [x] 2.5 Implement Defaults() method to return endpoint "http://localhost:11434", enabled false, api_key empty
- [x] 2.6 Implement Create() method to return NewOllamaProvider(settings)
- [x] 2.7 Add unit tests for OllamaFactory in internal/infrastructure/provider/ollama_test.go

## 3. OpenCode Zen Factory Implementation

- [x] 3.1 Create OpenCodeZenFactory struct in internal/infrastructure/provider/opencode.go
- [x] 3.2 Implement Type() method to return "opencode_zen"
- [x] 3.3 Implement Description() method to return "OpenCode Zen - Curated AI models for coding agents"
- [x] 3.4 Implement Validate() method to require non-empty endpoint
- [x] 3.5 Implement Defaults() method to return endpoint "https://opencode.ai/zen", enabled false, api_key empty
- [x] 3.6 Implement Create() method to return NewOpenCodeZenProvider(settings)
- [x] 3.7 Add unit tests for OpenCodeZenFactory in internal/infrastructure/provider/opencode_test.go

## 4. Config Structure Update

- [x] 4.1 Update ProviderConfig struct in internal/infrastructure/config/config.go to use map[string]ProviderSettings for Providers field
- [x] 4.2 Update GetEnabledProviders() method to return map[string]ProviderSettings instead of slice
- [x] 4.3 Update resolveEnvVars() function to handle map-based provider structure
- [x] 4.4 Add unit tests for updated config loading in internal/infrastructure/config/config_test.go
- [x] 4.5 Update router-service/configs/config.yaml to new map-based structure with providers.ollama and providers.opencode_zen keys
- [x] 4.6 Update opencode_zen endpoint to "https://opencode.ai/zen" and enable it

## 5. OpenCode Provider Naming Update

- [x] 5.1 Update OpenCodeZenProvider Name() method in internal/infrastructure/provider/opencode.go to return "opencode_zen"
- [x] 5.2 Update extractOpenCodeModelName() function to strip "opencode_zen:" prefix (13 characters)
- [x] 5.3 Update ListModels() method to use "opencode_zen:" prefix instead of "opencode:"
- [x] 5.4 Update unit tests in internal/infrastructure/provider/opencode_test.go to use new prefix

## 6. Handler Integration

- [x] 6.1 Update internal/handler/handler.go Setup() function to initialize ProviderRegistry
- [x] 6.2 Register OllamaFactory and OpenCodeZenFactory with registry in Setup()
- [x] 6.3 Replace URL-based provider matching with registry.Create() calls
- [x] 6.4 Add logging for successful provider initialization
- [x] 6.5 Add error handling and logging for failed provider creation
- [x] 6.6 Add warning log when no providers are enabled
- [x] 7.1 Add providersHandler() method to internal/handler/handler.go
- [x] 7.2 Register GET /v1/providers endpoint in Setup() function
- [x] 7.3 Implement providersHandler to return registered types, descriptions, configuration status, and defaults
- [x] 7.4 Add unit tests for providersHandler in internal/handler/handler_test.go

## 8. Integration Testing

- [x] 8.1 Create integration test for factory registration and provider instantiation
- [x] 8.2 Create integration test for config loading with new map-based structure
- [x] 8.3 Create integration test for /v1/providers endpoint
- [x] 8.4 Create integration test for model routing with new provider names
- [x] 8.5 Test with both providers enabled
- [x] 8.6 Test with single provider enabled
- [x] 8.7 Test with no providers enabled

## 9. Documentation and Migration

- [x] 9.1 Update README.md with new config structure example
- [x] 9.2 Add migration guide for old config to new config structure
- [x] 9.3 Document breaking changes (model prefix change, config structure change)
- [x] 9.4 Document /v1/providers endpoint usage
- [x] 9.5 Update AGENTS.md with factory pattern information

## 10. Bug Fixes Discovered During Verification

- [x] 10.1 Fix OpenCode Zen streaming implementation to use bufio.Scanner for SSE parsing instead of json.Decoder
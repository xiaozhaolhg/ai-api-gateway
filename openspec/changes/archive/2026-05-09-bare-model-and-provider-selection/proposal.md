## Why

Users currently cannot use bare model names (e.g., "llama2") when calling the AI API Gateway. They must know the internal `provider:model` format (e.g., "ollama:llama2"), which exposes internal routing concepts. Additionally, when multiple providers support the same model, there's no intelligent selection mechanism — the system returns the first matching rule without considering provider health or availability.

This change enables **bare model name support** (resolving "llama2" → provider automatically) and **health-priority provider selection** (choosing the healthiest provider when multiple support the same model).

## What Changes

- **NEW RPC**: `FindProvidersByModel` in provider-service — reverse lookup to find all providers supporting a given model
- **Router enhancement**: Detect bare model names (no ":" separator) and automatically resolve via `FindProvidersByModel`
- **Health-priority selection**: When multiple providers support a model, check health via existing `HealthCheck` RPC and select the healthiest
- **Fallback integration**: Non-primary providers automatically populate `fallback_provider_ids` and `fallback_models` fields
- **No "provider:model" format required**: Users only need to know model names like "llama2", "gpt-4"

## Capabilities

### New Capabilities

- `bare-model-resolution`: Resolve bare model names (without provider prefix) to the appropriate provider, enabling user-friendly model specification
- `provider-health-selection`: Select providers based on health status when multiple providers support the same model

### Modified Capabilities

- `router-service`: Routing logic updated to handle bare model names and perform health-priority selection
- `provider-service`: New `FindProvidersByModel` RPC and enhanced `Provider` entity query capabilities

## Impact

- **APIs**: New `FindProvidersByModel` RPC in `provider.proto`; modified `ResolveRoute` behavior in `router.proto`
- **Code**: 
  - `provider-service`: New repository method, service method, gRPC handler
  - `router-service`: Enhanced `ResolveRoute` logic with health checking
  - `gateway-service`: No changes required (behavior is transparent)
- **Dependencies**: HealthCheck RPC (already defined in provider.proto) will be actively used
- **Proto generation**: `buf generate` required after `provider.proto` changes

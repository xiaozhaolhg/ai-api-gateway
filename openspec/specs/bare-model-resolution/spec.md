## Purpose

Enables users to specify models by bare name (e.g., "llama2") without knowing the internal `provider:model` format. The router-service automatically resolves bare model names to the appropriate provider.

## Scope

- **In Scope**: Bare model name resolution in router-service, new `FindProvidersByModel` RPC in provider-service
- **Out of Scope**: Model authorization enforcement (Phase 2+), Gateway-service changes

## Requirements

### Requirement: Bare Model Name Detection

The router-service `ResolveRoute` SHALL detect bare model names (model strings that do not contain the ":" separator) and delegate to bare model resolution logic.

#### Scenario: Bare model name detection
- **WHEN** `ResolveRoute` is called with model="llama2" (no ":" separator)
- **THEN** the router SHALL invoke bare model resolution instead of pattern matching

#### Scenario: Provider-prefixed model name
- **WHEN** `ResolveRoute` is called with model="ollama:llama2" (contains ":")
- **THEN** the router SHALL use existing pattern matching logic (no bare model resolution)

### Requirement: FindProvidersByModel RPC

The provider-service SHALL expose a `FindProvidersByModel` RPC that returns all providers supporting a given bare model name.

#### Scenario: Model supported by multiple providers
- **WHEN** `FindProvidersByModel` is called with model="llama2"
- **THEN** return all providers where `Models` field contains "llama2"
- **AND** return providers sorted by ID (deterministic order)

#### Scenario: Model not supported by any provider
- **WHEN** `FindProvidersByModel` is called with model="nonexistent-model"
- **THEN** return an empty providers list (not an error)

#### Scenario: Provider with empty Models field
- **WHEN** a provider has empty `Models` field
- **THEN** that provider SHALL NOT be included in results (even if name matches)

### Requirement: Bare Model Resolution Result

The bare model resolution SHALL return a `RouteResult` with the primary provider and fallback providers populated.

#### Scenario: Single provider supports model
- **WHEN** only one provider supports "llama2"
- **THEN** return `RouteResult` with that provider as `provider_id`
- **AND** `fallback_provider_ids` SHALL be empty

#### Scenario: Multiple providers support model
- **WHEN** multiple providers support "llama2"
- **THEN** return `RouteResult` with the healthiest provider as `provider_id`
- **AND** populate `fallback_provider_ids` with remaining healthy providers
- **AND** populate `fallback_models` with the bare model name for each fallback

### Requirement: Transparent Resolution

The bare model resolution SHALL be transparent to gateway-service — the `ResolveRoute` response format is unchanged.

#### Scenario: Gateway receives response for bare model
- **WHEN** gateway calls `ResolveRoute` with model="llama2"
- **THEN** gateway receives `RouteResult` with valid `provider_id`, `adapter_type`, `fallback_provider_ids`
- **AND** gateway proceeds with normal proxy flow (no changes needed in gateway)

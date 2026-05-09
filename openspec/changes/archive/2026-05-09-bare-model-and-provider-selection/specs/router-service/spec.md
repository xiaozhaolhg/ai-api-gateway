## MODIFIED Requirements

### Requirement: Fallback Route Resolution

Router service ResolveRoute SHALL accept an optional `user_id` parameter to support per-user routing rules. When the model name does NOT contain a ":" separator (bare model name), the router SHALL use bare model resolution via `FindProvidersByModel` RPC.

#### Scenario: Route resolution with bare model name
- **WHEN** `ResolveRoute` is called with model="llama2" (no ":" separator)
- **THEN** the router SHALL invoke `bare model resolution` logic
- **AND** call provider-service `FindProvidersByModel("llama2")`
- **AND** return `RouteResult` with healthiest provider as primary

#### Scenario: Route resolution with user rule (user OVERRIDES system)
- **WHEN** `ResolveRoute` is called with a `user_id` and a matching user-specific routing rule exists
- **THEN** return `RouteResult` using the user's rule (provider_id, fallback_provider_ids)
- **AND** ignore any system-wide rule for the same model pattern

#### Scenario: Route resolution without user rule (fallback to system)
- **WHEN** `ResolveRoute` is called with a `user_id` but no user-specific rule matches
- **THEN** fall back to system-wide routing rules (where `user_id` is NULL or empty)
- **AND** return `RouteResult` with system rule's provider and fallback chain

#### Scenario: Route resolution without user_id (system rule only)
- **WHEN** `ResolveRoute` is called without `user_id` (or empty string)
- **THEN** only match system-wide rules (user_id IS NULL)
- **AND** return `RouteResult` with `fallback_provider_ids` populated if configured

#### Scenario: Route resolution with fallback chain
- **WHEN** `ResolveRoute` is called and the matching rule has `fallback_provider_ids` configured
- **THEN** return `RouteResult` with the ordered `fallback_provider_ids` list
- **AND** the gateway/provider service SHALL try each provider in order on failure

#### Scenario: Route resolution without fallback
- **WHEN** `ResolveRoute` is called and no fallback provider is configured for the matching rule
- **THEN** return `RouteResult` with empty `fallback_provider_ids`

#### Scenario: Bare model resolution with single provider
- **WHEN** `FindProvidersByModel("llama2")` returns only one provider
- **THEN** return `RouteResult` with that provider as primary
- **AND** `fallback_provider_ids` SHALL be empty

#### Scenario: Bare model resolution with multiple providers
- **WHEN** `FindProvidersByModel("llama2")` returns multiple providers
- **THEN** perform concurrent health checks via `HealthCheck` RPC
- **AND** select the first healthy provider as primary
- **AND** populate `fallback_provider_ids` with remaining healthy providers

#### Scenario: No healthy providers for bare model
- **WHEN** no providers are healthy for a bare model name
- **THEN** return an error: "no healthy provider found for model: <model>"

### Requirement: Per-User Route Resolution

Router service ResolveRoute SHALL accept an optional `user_id` parameter to support per-user routing rules.

#### Scenario: ResolveRoute with user_id
- **WHEN** `ResolveRoute` is called with a `user_id` parameter
- **THEN** the router-service SHALL first look for routing rules where `user_id` matches
- **AND** only fall back to system rules (user_id IS NULL) if no user rule matches

#### Scenario: ResolveRoute without user_id
- **WHEN** `ResolveRoute` is called without `user_id` (or empty string)
- **THEN** the router-service SHALL only match system-wide rules (user_id IS NULL)
- **AND** return NOT_FOUND if no system rule matches

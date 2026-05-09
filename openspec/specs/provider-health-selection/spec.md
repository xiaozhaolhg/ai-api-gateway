## Purpose

Selects the healthiest provider when multiple providers support the same model, using the existing `HealthCheck` RPC to determine provider availability.

## Scope

- **In Scope**: Health-priority selection logic in router-service, concurrent health checks, integration with `HealthCheck` RPC
- **Out of Scope**: Load balancing algorithms (round-robin, least-connections), cost-based selection, latency-based selection

## Requirements

### Requirement: Health Check Integration

The router-service SHALL use the provider-service `HealthCheck` RPC to determine provider health status when resolving bare model names.

#### Scenario: Health check via existing RPC
- **WHEN** resolving a bare model name with multiple supporting providers
- **THEN** the router SHALL call `HealthCheck` RPC for each provider concurrently
- **AND** wait for all health check results before selecting primary provider

#### Scenario: HealthCheck RPC failure
- **WHEN** `HealthCheck` RPC fails for a provider (network error, timeout)
- **THEN** that provider SHALL be treated as unhealthy
- **AND** excluded from primary/fallback selection

### Requirement: Health-Priority Selection

The router-service SHALL select the healthiest provider as primary, with remaining healthy providers as fallbacks.

#### Scenario: All providers healthy
- **WHEN** multiple providers support a model and all are healthy
- **THEN** select the first provider (by sorted order) as primary
- **AND** populate `fallback_provider_ids` with remaining providers in order

#### Scenario: Some providers unhealthy
- **WHEN** some providers are unhealthy
- **THEN** select the first healthy provider as primary
- **AND** populate `fallback_provider_ids` only with other healthy providers
- **AND** exclude unhealthy providers entirely

#### Scenario: No healthy providers
- **WHEN** no providers are healthy for a model
- **THEN** return an error: "no healthy provider found for model: <model>"

### Requirement: Concurrent Health Checks

Health checks for multiple providers SHALL be performed concurrently to minimize latency.

#### Scenario: Concurrent execution
- **WHEN** 3 providers support "llama2"
- **THEN** launch 3 concurrent goroutines to call `HealthCheck`
- **AND** collect results via channel
- **AND** total elapsed time ≈ slowest single health check (not sum)

#### Scenario: Health check timeout
- **WHEN** a `HealthCheck` call exceeds 5 seconds
- **THEN** treat that provider as unhealthy
- **AND** continue with remaining providers

### Requirement: Health Status Caching (Phase 2+)

**Note**: For MVP, health checks are performed on each resolution. Caching is planned for Phase 2+.

#### Scenario: Future caching behavior
- **WHEN** health status caching is implemented
- **THEN** cache health status in Redis with 10-second TTL
- **AND** use cached status if available and not expired

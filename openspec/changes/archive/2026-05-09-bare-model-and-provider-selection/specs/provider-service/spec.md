## ADDED Requirements

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

#### Scenario: Exact model name match
- **WHEN** `FindProvidersByModel` is called
- **THEN** perform exact match on model name (not wildcard)
- **AND** return only providers with exact model name in `Models` list

## MODIFIED Requirements

### Requirement: Provider management

The provider-service SHALL provide complete provider lifecycle management including CRUD operations, credential encryption, health checks, and gRPC API. **NEW: Including `FindProvidersByModel` RPC for reverse model-to-provider lookup.**

#### Scenario: Provider CRUD fully implemented
- **WHEN** the provider-service is running
- **THEN** all gRPC handlers for ProviderService RPCs SHALL be fully implemented (not stubs)
- **AND** CreateProvider/UpdateProvider/DeleteProvider SHALL persist changes to the database
- **AND** ListProviders/GetProvider SHALL return provider data with credentials masked
- **AND** `FindProvidersByModel` SHALL return providers supporting a given model

#### Scenario: Integration test with mock server
- **WHEN** an integration test is run
- **THEN** it SHALL use a mock HTTP server to simulate an external provider
- **AND** verify the full flow: add provider via Admin API → route request through it
- **AND** verify `FindProvidersByModel` returns correct providers for a model

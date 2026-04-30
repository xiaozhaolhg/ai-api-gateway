# provider-service

## Purpose

Provider domain — provider CRUD, adapters, response callback dispatch.

## Service Responsibility

- **Role**: Provider management, request forwarding, callback dispatch
- **Owned Entities**: Provider, EncryptedCredential
- **Data Layer**: provider-db (SQLite/PostgreSQL)

## Dependencies

### Calls To

| Service | Methods | Purpose |
|---|---|---|
| (external) | HTTPS | Call external AI providers (OpenAI, Anthropic, Gemini) |

### Called By

| Service | Methods | Purpose |
|---|---|---|
| gateway-service | `ForwardRequest`, `StreamRequest` | Forward consumer requests |
| gateway-service | `CreateProvider`, `UpdateProvider`, `DeleteProvider` | Provider CRUD |
| router-service | `GetProviderByType` | Verify provider exists |

### Data Dependencies

- **Database**: provider-db (Provider)
- **Cache**: Redis (provider config)

## Key Design

### Provider Adapters

- **OpenAI Adapter**: Transform to OpenAI format, parse response
- **Anthropic Adapter**: Transform to Anthropic format, parse response  
- **Gemini Adapter**: Transform to Gemini format, parse response

### Callback Dispatch (Observer Pattern)

After each provider response:
1. Extract token counts, latency, status
2. Dispatch ProviderResponseCallback to billing-service (async)
3. Dispatch ProviderResponseCallback to monitor-service (async)
4. Non-blocking — fire and forget

### Key Operations

- **ForwardRequest**: Non-streaming request to provider
- **StreamRequest**: Streaming request (SSE proxy)
- **CreateProvider/UpdateProvider/DeleteProvider**: Provider lifecycle
- **RegisterSubscriber/UnregisterSubscriber**: Callback registration

### Data Encryption

- Credentials encrypted at rest using AES-256-GCM
- Encryption key managed via config/env

## Requirements

### Requirement: Provider management
The provider-service SHALL provide complete provider lifecycle management including CRUD operations, credential encryption, health checks, and gRPC API.

#### Scenario: Provider CRUD fully implemented
- **WHEN** the provider-service is running
- **THEN** all gRPC handlers for ProviderService RPCs SHALL be fully implemented (not stubs)
- **AND** CreateProvider/UpdateProvider/DeleteProvider SHALL persist changes to the database
- **AND** ListProviders/GetProvider SHALL return provider data with credentials masked

#### Scenario: Integration test with mock server
- **WHEN** an integration test is run
- **THEN** it SHALL use a mock HTTP server to simulate an external provider
- **AND** verify the full flow: add provider via Admin API → route request through it
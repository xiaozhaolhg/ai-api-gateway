## MODIFIED Requirements

### Requirement: DDD four-layer architecture
The provider-service SHALL implement four-layer Clean Architecture: Domain, Application, Infrastructure, and Handler with dependency direction from outer to inner layers.

#### Scenario: Domain layer has no external dependencies
- **WHEN** the domain layer is imported
- **THEN** it SHALL NOT import any code from application, infrastructure, or handler layers

#### Scenario: Infrastructure implements domain interfaces
- **WHEN** a new LLM provider adapter is needed
- **THEN** it SHALL be implemented by creating a struct that implements the ProviderAdapter interface defined in the domain layer

### Requirement: Provider entity and repository
The provider-service SHALL own the Provider entity with fields: id, name, type, base_url, credentials (encrypted), models, status, created_at, updated_at. It SHALL provide a ProviderRepository interface for CRUD operations.

#### Scenario: Create provider
- **WHEN** a CreateProvider request is received via gRPC
- **THEN** the service SHALL encrypt credentials with AES-256-GCM, persist the Provider to SQLite, and return the Provider (with credentials redacted)

#### Scenario: List providers
- **WHEN** a ListProviders request is received
- **THEN** the service SHALL return all providers with credentials redacted

### Requirement: Provider adapter interface
The provider-service SHALL define a ProviderAdapter interface in the domain layer with methods: TransformRequest, TransformResponse, StreamResponse, CountTokens.

#### Scenario: OpenAI-compatible adapter
- **WHEN** a request is forwarded to an OpenAI-type provider
- **THEN** the OpenAI adapter SHALL pass the request through (already in OpenAI format) and parse the response

#### Scenario: Anthropic adapter
- **WHEN** a request is forwarded to an Anthropic-type provider
- **THEN** the Anthropic adapter SHALL transform the request to Anthropic format and parse the Anthropic response back to OpenAI format

#### Scenario: Ollama adapter
- **WHEN** a request is forwarded to an Ollama-type provider
- **THEN** the Ollama adapter SHALL transform the request to Ollama /api/chat format and parse the response back to OpenAI format

### Requirement: Request forwarding (non-streaming)
The provider-service SHALL implement ForwardRequest that sends a request to the configured provider and returns the response.

#### Scenario: Non-streaming forward
- **WHEN** a ForwardRequest gRPC call is received with provider_id and request_body
- **THEN** the service SHALL look up the provider, decrypt credentials, select the adapter, transform and send the request, and return ForwardRequestResponse with response_body, token_counts, and status_code

### Requirement: Request forwarding (streaming)
The provider-service SHALL implement StreamRequest that sends a streaming request to the provider and returns SSE chunks.

#### Scenario: Streaming forward
- **WHEN** a StreamRequest gRPC call is received
- **THEN** the service SHALL proxy the provider's SSE stream as ProviderChunk messages with accumulated token counts
- **AND** the final chunk SHALL have done=true

### Requirement: Async callback dispatch (observer pattern)
The provider-service SHALL dispatch ProviderResponseCallback to all registered subscribers after each provider response.

#### Scenario: Callback to billing-service
- **WHEN** a provider response completes (success or error)
- **THEN** the service SHALL dispatch OnProviderResponse to billing-service asynchronously (fire-and-forget)
- **AND** callback failure SHALL NOT block the response to the caller

#### Scenario: Callback to monitor-service
- **WHEN** a provider response completes
- **THEN** the service SHALL dispatch OnProviderResponse to monitor-service asynchronously

#### Scenario: Subscriber registration
- **WHEN** a RegisterSubscriber request is received with service_name and callback_endpoint
- **THEN** the service SHALL add the subscriber to its callback list

### Requirement: Credential encryption
The provider-service SHALL encrypt provider credentials at rest using AES-256-GCM with an encryption key from config/env.

#### Scenario: Credentials encrypted on storage
- **WHEN** a provider is created or updated with credentials
- **THEN** the service SHALL encrypt the credentials before storing in SQLite

#### Scenario: Credentials decrypted on use
- **WHEN** a request is forwarded to a provider
- **THEN** the service SHALL decrypt the credentials to obtain the API key for the provider request

### Requirement: gRPC server implementation
The provider-service SHALL implement the ProviderService gRPC server as defined in the api/ module proto definitions.

#### Scenario: All proto RPCs are implemented
- **WHEN** the provider-service starts
- **THEN** it SHALL register all RPCs defined in provider.proto: ForwardRequest, StreamRequest, GetProvider, CreateProvider, UpdateProvider, DeleteProvider, ListProviders, ListModels, GetProviderByType, RegisterSubscriber, UnregisterSubscriber

### Requirement: SQLite persistence
The provider-service SHALL use SQLite as its database with a providers table managed via migrations.

#### Scenario: Database initialized on startup
- **WHEN** the provider-service starts
- **THEN** it SHALL create the SQLite database file and run migrations if the database is new

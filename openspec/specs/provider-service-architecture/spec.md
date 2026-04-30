## Purpose

Provider CRUD, adapter framework, request forwarding, SSE streaming, and callback dispatch service with gRPC interface and SQLite persistence.
## Requirements
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

#### Scenario: Provider entity structure
- **WHEN** the Provider entity is created
- **THEN** it SHALL contain all required fields: id, name, type, base_url, credentials, models, status, created_at, updated_at

### Requirement: StreamingResult and TokenCounts entities
The provider-service SHALL define `StreamingResult` and `TokenCounts` entities in the domain layer for streaming response handling and token accumulation.

#### Scenario: TokenCounts entity structure
- **WHEN** the TokenCounts entity is used
- **THEN** it SHALL contain fields: PromptTokens, CompletionTokens, AccumulatedTokens
- **AND** it SHALL provide a Total() method that returns the accumulated or computed total

#### Scenario: StreamingResult entity structure  
- **WHEN** the StreamingResult entity is used
- **THEN** it SHALL contain fields: TransformedData []byte, TokenCounts TokenCounts, IsFinal bool

#### Scenario: Create provider
- **WHEN** a CreateProvider request is received via gRPC
- **THEN** the service SHALL encrypt credentials with AES-256-GCM, persist the Provider to SQLite, and return the Provider (with credentials redacted)

#### Scenario: Update provider with timestamp update
- **WHEN** an UpdateProvider request is received via gRPC
- **THEN** the service SHALL update the UpdatedAt field to current time
- **AND** encrypt the credentials field if provided (skip if empty)
- **AND** return the Provider with credentials masked as `***`

#### Scenario: List providers
- **WHEN** a ListProviders request is received
- **THEN** the service SHALL return all providers with credentials redacted

#### Scenario: Get provider by ID with masked credentials
- **WHEN** a GetProvider request is received with a valid provider ID
- **THEN** the service SHALL return the provider with credentials field set to `***`

#### Scenario: Delete provider
- **WHEN** a DeleteProvider request is received with a valid provider ID
- **THEN** the service SHALL delete the provider from the database
- **AND** return success

#### Scenario: Duplicate name detection
- **WHEN** a CreateProvider request is received with a name that already exists
- **THEN** the service SHALL return an error (provider already exists)

#### Scenario: UUID v4 auto-generation
- **WHEN** a CreateProvider request is received with an empty ID
- **THEN** the service SHALL generate a UUID v4 for the provider ID

#### Scenario: Timestamp tracking on create
- **WHEN** a CreateProvider request is received
- **THEN** the service SHALL set CreatedAt and UpdatedAt to current time

#### Scenario: Timestamp tracking on update
- **WHEN** an UpdateProvider request is received
- **THEN** the service SHALL update UpdatedAt to current time (preserve CreatedAt)

### Requirement: Provider adapter interface
The provider-service SHALL define a ProviderAdapter interface in the domain layer with methods: TransformRequest, TransformResponse, CountTokens, TestConnection. The TransformResponse method SHALL support both streaming and non-streaming via a unified interface with `isStreaming` flag.

#### Scenario: Unified TransformResponse interface
- **WHEN** TransformResponse is called with `isStreaming=false`
- **THEN** the adapter SHALL process the complete response body and return (transformedData, tokenCounts, isFinal=true, error)
- **WHEN** TransformResponse is called with `isStreaming=true`
- **THEN** the adapter SHALL process a single SSE chunk and return (transformedChunk, updatedTokenCounts, isFinal, error)

#### Scenario: Token accumulation during streaming
- **WHEN** streaming chunks are processed
- **THEN** the adapter SHALL accumulate token counts progressively via the `accumulatedTokens` parameter
- **AND** the final chunk SHALL contain the complete token counts

#### Scenario: OpenAI-compatible adapter
- **WHEN** a request is forwarded to an OpenAI-type provider
- **THEN** the OpenAI adapter SHALL pass the request through (already in OpenAI format) and parse the response
- **AND** for streaming, the adapter SHALL detect the `[DONE]` marker for final chunk identification

#### Scenario: Anthropic adapter
- **WHEN** a request is forwarded to an Anthropic-type provider
- **THEN** the Anthropic adapter SHALL transform the request to Anthropic format and parse the Anthropic response back to OpenAI format
- **AND** for streaming, the adapter SHALL detect the `message_stop` event for final chunk identification

#### Scenario: Ollama adapter
- **WHEN** a request is forwarded to an Ollama-type provider
- **THEN** the Ollama adapter SHALL transform the request to Ollama /api/chat format and parse the response back to OpenAI format
- **AND** TestConnection SHALL call the `/api/tags` endpoint to verify Ollama is running

#### Scenario: TestConnection method
- **WHEN** TestConnection(credentials string) is called on an adapter
- **THEN** the adapter SHALL make a lightweight request to verify connectivity to the external provider
- **AND** return nil if successful, error if failed (authentication error, unreachable URL, etc.)

#### Scenario: OpenAI adapter TestConnection
- **WHEN** TestConnection is called on an OpenAI adapter with valid credentials
- **THEN** the adapter SHALL make a lightweight request (e.g., list models) to verify connectivity
- **AND** return nil if successful, error if failed

#### Scenario: Anthropic adapter TestConnection
- **WHEN** TestConnection is called on an Anthropic adapter with valid credentials
- **THEN** the adapter SHALL make a test request to verify connectivity
- **AND** return nil if successful, error if failed

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


## MODIFIED Requirements

### Requirement: Provider adapter interface
The provider-service SHALL define a ProviderAdapter interface in the domain layer with methods: TransformRequest, TransformResponse, CountTokens, TestConnection.

#### Scenario: Unified TransformResponse interface
- **WHEN** TransformResponse is called with `isStreaming=false`
- **THEN** the adapter SHALL process the complete response body and return (transformedData, tokenCounts, isFinal=true, error)
- **WHEN** TransformResponse is called with `isStreaming=true`
- **THEN** the adapter SHALL process a single SSE chunk and return (transformedChunk, updatedTokenCounts, isFinal, error)

#### Scenario: Token accumulation during streaming
- **WHEN** streaming chunks are processed
- **THEN** the adapter SHALL accumulate token counts progressively via the `accumulatedTokens` parameter
- **AND** the final chunk SHALL contain the complete token counts

#### Scenario: TestConnection method
- **WHEN** TestConnection(credentials string) is called on an adapter
- **THEN** the adapter SHALL make a lightweight request to verify connectivity to the external provider
- **AND** return nil if successful, error if failed (authentication error, unreachable URL, etc.)

#### Scenario: OpenAI-compatible adapter
- **WHEN** a request is forwarded to an OpenAI-type provider
- **THEN** the OpenAI adapter SHALL pass the request through (already in OpenAI format) and parse the response
- **AND** for streaming, the adapter SHALL detect the `[DONE]` marker for final chunk identification
- **AND** TestConnection SHALL call the list models endpoint to verify connectivity

#### Scenario: Anthropic adapter
- **WHEN** a request is forwarded to an Anthropic-type provider
- **THEN** the Anthropic adapter SHALL transform the request to Anthropic format and parse the Anthropic response back to OpenAI format
- **AND** for streaming, the adapter SHALL detect the `message_stop` event for final chunk identification
- **AND** TestConnection SHALL make a lightweight test request to verify connectivity

#### Scenario: Ollama adapter
- **WHEN** a request is forwarded to an Ollama-type provider
- **THEN** the Ollama adapter SHALL transform the request to Ollama /api/chat format and parse the response back to OpenAI format
- **AND** TestConnection SHALL call the `/api/tags` endpoint to verify Ollama is running

### Requirement: gRPC server implementation
The provider-service SHALL implement the ProviderService gRPC server as defined in the api/ module proto definitions.

#### Scenario: All proto RPCs are implemented
- **WHEN** the provider-service starts
- **THEN** it SHALL register all RPCs defined in provider.proto: ForwardRequest, StreamRequest, GetProvider, CreateProvider, UpdateProvider, DeleteProvider, ListProviders, ListModels, GetProviderByType, RegisterSubscriber, UnregisterSubscriber
- **AND** all handlers SHALL be fully implemented (not stubs)

#### Scenario: CreateProvider encrypts credentials
- **WHEN** CreateProvider gRPC is called
- **THEN** the handler SHALL call the application service to encrypt credentials and persist the provider
- **AND** return the Provider with credentials masked as `***`

#### Scenario: UpdateProvider encrypts credentials
- **WHEN** UpdateProvider gRPC is called with new credentials
- **THEN** the handler SHALL call the application service to encrypt credentials and update the provider
- **AND** return the Provider with credentials masked as `***`

### Requirement: Provider entity and repository
The provider-service SHALL own the Provider entity with fields: id, name, type, base_url, credentials (encrypted), models, status, created_at, updated_at. It SHALL provide a ProviderRepository interface for CRUD operations.

#### Scenario: UUID v4 auto-generation
- **WHEN** a CreateProvider request is received with an empty ID
- **THEN** the service SHALL generate a UUID v4 for the provider ID

#### Scenario: Timestamp tracking on create
- **WHEN** a CreateProvider request is received
- **THEN** the service SHALL set CreatedAt and UpdatedAt to current time

#### Scenario: Timestamp tracking on update
- **WHEN** an UpdateProvider request is received
- **THEN** the service SHALL update UpdatedAt to current time (preserve CreatedAt)

#### Scenario: Duplicate name detection
- **WHEN** the ProviderRepository needs to check for duplicate names
- **THEN** it SHALL provide a GetByName(name) method that returns the provider if found, error if not found

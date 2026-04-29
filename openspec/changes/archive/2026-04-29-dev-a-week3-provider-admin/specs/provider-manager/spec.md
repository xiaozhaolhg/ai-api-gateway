## ADDED Requirements

### Requirement: Provider CRUD with credential encryption
The provider-service SHALL provide CRUD operations for Provider entities with automatic UUID v4 generation, timestamp tracking, and AES-256-GCM credential encryption.

#### Scenario: Create provider with auto-generated ID
- **WHEN** a CreateProvider request is received via gRPC
- **THEN** the service SHALL generate a UUID v4 for the provider ID if not provided
- **AND** set CreatedAt and UpdatedAt to current time
- **AND** encrypt the credentials field using AES-256-GCM before persisting to database
- **AND** return the Provider with credentials masked as `***`

#### Scenario: Update provider with timestamp update
- **WHEN** an UpdateProvider request is received via gRPC
- **THEN** the service SHALL update the UpdatedAt field to current time
- **AND** encrypt the credentials field if provided (skip if empty)
- **AND** return the Provider with credentials masked as `***`

#### Scenario: List providers with masked credentials
- **WHEN** a ListProviders request is received
- **THEN** the service SHALL return all providers with credentials field set to `***`

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

### Requirement: ProviderAdapter TestConnection method
The provider-service SHALL extend the ProviderAdapter interface with a TestConnection(credentials string) error method that verifies connectivity to the external provider.

#### Scenario: OpenAI adapter TestConnection
- **WHEN** TestConnection is called on an OpenAI adapter with valid credentials
- **THEN** the adapter SHALL make a lightweight request (e.g., list models) to verify connectivity
- **AND** return nil if successful, error if failed

#### Scenario: Anthropic adapter TestConnection
- **WHEN** TestConnection is called on an Anthropic adapter with valid credentials
- **THEN** the adapter SHALL make a test request to verify connectivity
- **AND** return nil if successful, error if failed

#### Scenario: Ollama adapter TestConnection
- **WHEN** TestConnection is called on an Ollama adapter
- **THEN** the adapter SHALL make a request to the Ollama base URL (no credentials needed)
- **AND** return nil if successful, error if failed

#### Scenario: Invalid credentials
- **WHEN** TestConnection is called with invalid credentials
- **THEN** the adapter SHALL return an error indicating authentication failure

### Requirement: Provider health check
The provider-service SHALL provide a health check mechanism that uses the ProviderAdapter's TestConnection method to verify provider connectivity.

#### Scenario: Health check success
- **WHEN** a health check is requested for an active provider
- **THEN** the service SHALL call the provider's adapter TestConnection method
- **AND** return status "healthy" if successful

#### Scenario: Health check failure
- **WHEN** a health check is requested for a provider with invalid credentials or unreachable URL
- **THEN** the service SHALL return status "unhealthy" with error details

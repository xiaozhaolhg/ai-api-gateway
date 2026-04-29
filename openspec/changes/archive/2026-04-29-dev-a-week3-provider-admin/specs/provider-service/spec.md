## MODIFIED Requirements

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

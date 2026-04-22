## MODIFIED Requirements

### Requirement: Shared protobuf definitions for all services
The `api/` module SHALL contain protobuf definitions for all 5 internal gRPC services (auth, router, provider, billing, monitor) plus shared common messages.

#### Scenario: Proto files organized by service
- **WHEN** the api/ module is inspected
- **THEN** proto files SHALL exist at `proto/{service}/v1/{service}.proto` for each service
- **AND** a `proto/common/v1/common.proto` SHALL define shared messages (Empty, TokenCounts, ProviderResponseCallback, etc.)

### Requirement: buf configuration and code generation
The `api/` module SHALL use buf for proto linting and Go gRPC stub generation.

#### Scenario: buf lint passes
- **WHEN** `buf lint` is run
- **THEN** all proto files SHALL pass linting with no errors

#### Scenario: buf generate produces Go stubs
- **WHEN** `buf generate` is run
- **THEN** generated Go files SHALL exist at `gen/{service}/v1/` for each service
- **AND** each generated directory SHALL contain `{service}_grpc.pb.go` and `{service}.pb.go`

### Requirement: api module is importable by all services
The `api/` module SHALL be a standalone Go module (`github.com/ai-api-gateway/api`) importable by all service modules.

#### Scenario: Service imports generated stubs
- **WHEN** a service's go.mod includes `github.com/ai-api-gateway/api` as a dependency
- **THEN** the service SHALL be able to import generated gRPC client and server stubs

### Requirement: Proto definitions match existing API contracts
All protobuf service definitions and messages SHALL match the API contracts defined in `openspec/specs/{service}/api-contracts.md`.

#### Scenario: auth-service proto matches contract
- **WHEN** the auth.proto service definition is compared to `openspec/specs/auth-service/api-contracts.md`
- **THEN** all rpc methods and message fields SHALL match the contract

#### Scenario: provider-service proto matches contract
- **WHEN** the provider.proto service definition is compared to `openspec/specs/provider-service/api-contracts.md`
- **THEN** all rpc methods and message fields SHALL match the contract

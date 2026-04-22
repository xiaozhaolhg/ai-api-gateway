## MODIFIED Requirements

### Requirement: Clean Architecture Layer Structure
The router-service SHALL implement four-layer Clean Architecture: Domain, Application, Infrastructure, and Handler layers with clear dependency direction from outer to inner layers. The Handler layer SHALL be a gRPC server (not an HTTP handler).

#### Scenario: Domain layer has no external dependencies
- **WHEN** the domain layer is imported
- **THEN** it SHALL NOT import any code from application, infrastructure, or handler layers

#### Scenario: Infrastructure implements domain interfaces
- **WHEN** a new routing strategy is needed
- **THEN** it SHALL be implemented by creating a struct that implements the Router interface defined in the domain layer

#### Scenario: Handler is gRPC server
- **WHEN** the router-service starts
- **THEN** it SHALL listen on a gRPC port (not HTTP) and serve RouterService RPCs

### Requirement: Provider Abstraction
The router-service SHALL resolve model names to provider identifiers via gRPC calls to provider-service, not via direct provider instances.

#### Scenario: Request routing via gRPC
- **WHEN** a ResolveRoute request is received with model "ollama:llama3"
- **THEN** the router SHALL look up the RoutingRule in its SQLite database
- **AND** return RouteResult with provider_id and adapter_type

#### Scenario: Route not found
- **WHEN** a ResolveRoute request is received with a model that has no matching routing rule
- **THEN** the router SHALL return a NOT_FOUND gRPC error

### Requirement: RoutingRule entity and repository
The router-service SHALL own the RoutingRule entity with fields: id, model_pattern, provider_id, priority, fallback_provider_id. It SHALL provide a RoutingRuleRepository interface for CRUD operations.

#### Scenario: Create routing rule
- **WHEN** a CreateRoutingRule request is received via gRPC
- **THEN** the service SHALL persist a new RoutingRule to SQLite
- **AND** return the created rule with generated id

#### Scenario: Model pattern matching
- **WHEN** a model name is matched against routing rules
- **THEN** the service SHALL support wildcard patterns (e.g., "gpt-4*") and select the highest-priority match

### Requirement: gRPC server implementation
The router-service SHALL implement the RouterService gRPC server as defined in the api/ module proto definitions.

#### Scenario: All proto RPCs are implemented
- **WHEN** the router-service starts
- **THEN** it SHALL register all RPCs defined in router.proto: ResolveRoute, GetRoutingRules, CreateRoutingRule, UpdateRoutingRule, DeleteRoutingRule, RefreshRoutingTable

### Requirement: SQLite persistence
The router-service SHALL use SQLite as its database with a routing_rules table managed via migrations.

#### Scenario: Database initialized on startup
- **WHEN** the router-service starts
- **THEN** it SHALL create the SQLite database file and run migrations if the database is new
- **AND** existing data SHALL be preserved on restart

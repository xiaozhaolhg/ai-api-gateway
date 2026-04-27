## Purpose

Model-to-provider route resolution via pattern-based routing rules with gRPC interface.
## Requirements
### Requirement: Clean Architecture Layer Structure
The router-service SHALL implement four-layer Clean Architecture: Domain, Application, Infrastructure, and Handler layers with clear dependency direction from outer to inner layers. The Handler layer SHALL be a gRPC server (not an HTTP handler).

#### Scenario: Domain layer has no external dependencies
- **WHEN** the domain layer is imported
- **THEN** it SHALL NOT import any code from application, infrastructure, or handler layers

#### Scenario: Infrastructure implements domain interfaces
- **WHEN** a new routing strategy is needed
- **THEN** it SHALL be implemented by creating a struct that implements the Router interface defined in the domain layer

#### Scenario: Router interface definition
- **WHEN** the Router interface is used
- **THEN** it SHALL define methods: ResolveRoute, CreateRoutingRule, UpdateRoutingRule, DeleteRoutingRule, ListRoutingRules, RefreshRoutingTable

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

### Requirement: RoutingRule and RouteResult entities
The router-service SHALL own the RoutingRule entity with fields: id, model_pattern, provider_id, priority, fallback_provider_id. It SHALL provide a RoutingRuleRepository interface for CRUD operations. It SHALL also define a RouteResult entity with fields: provider_id, adapter_type for route resolution responses.

#### Scenario: Create routing rule
- **WHEN** a CreateRoutingRule request is received via gRPC
- **THEN** the service SHALL persist a new RoutingRule to SQLite
- **AND** return the created rule with generated id

#### Scenario: Update routing rule
- **WHEN** an UpdateRoutingRule request is received with rule id
- **THEN** the service SHALL update the existing RoutingRule in SQLite
- **AND** return the updated rule

#### Scenario: Delete routing rule
- **WHEN** a DeleteRoutingRule request is received with rule id
- **THEN** the service SHALL delete the RoutingRule from SQLite
- **AND** invalidate any cached routes matching this pattern

#### Scenario: List routing rules
- **WHEN** a ListRoutingRules request is received with page and pageSize
- **THEN** the service SHALL return a paginated list of RoutingRules ordered by priority
- **AND** return the total count of all rules

#### Scenario: Refresh routing table
- **WHEN** a RefreshRoutingTable request is received
- **THEN** the service SHALL invalidate the routing table cache
- **AND** ensure subsequent ResolveRoute calls use fresh data

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


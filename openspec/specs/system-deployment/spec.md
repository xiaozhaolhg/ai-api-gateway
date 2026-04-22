## Purpose

Docker Compose orchestration, go.work multi-module structure, and top-level build configuration for the microservices system.
## Requirements
### Requirement: Docker Compose orchestration
The system SHALL provide a docker-compose.yaml that runs all 7 services (6 Go services + admin-ui) plus Redis with proper health checks and startup ordering.

#### Scenario: All services start with docker-compose up
- **WHEN** `docker-compose up` is run
- **THEN** all 7 services SHALL start with correct port mappings
- **AND** health checks SHALL pass for each service
- **AND** gateway-service SHALL wait for auth, router, and provider services to be healthy

#### Scenario: Service networking
- **WHEN** services are running in Docker Compose
- **THEN** each service SHALL be reachable by its service name on the Docker network
- **AND** gateway-service SHALL connect to auth-service:50051, router-service:50052, provider-service:50053, billing-service:50054, monitor-service:50055

### Requirement: go.work multi-module workspace
The repository SHALL use go.work to link all Go modules for local development.

#### Scenario: go.work includes all modules
- **WHEN** go.work is inspected
- **THEN** it SHALL include entries for api/, auth-service/, router-service/, provider-service/, gateway-service/, billing-service/, and monitor-service/

#### Scenario: Cross-module development works
- **WHEN** a developer makes changes to the api/ module
- **THEN** dependent services SHALL see the changes immediately without publishing the module

### Requirement: Top-level Makefile
The repository SHALL provide a top-level Makefile with targets for building, testing, and running all services.

#### Scenario: Build all services
- **WHEN** `make build` is run at the repo root
- **THEN** all Go services SHALL be built and their Docker images created

#### Scenario: Test all services
- **WHEN** `make test` is run at the repo root
- **THEN** all Go service tests SHALL execute

#### Scenario: Start all services
- **WHEN** `make up` is run at the repo root
- **THEN** `docker-compose up` SHALL be executed

#### Scenario: Stop all services
- **WHEN** `make down` is run at the repo root
- **THEN** `docker-compose down` SHALL be executed

### Requirement: Legacy service preserved
The existing router-service SHALL be renamed to router-service-legacy/ and preserved for reference.

#### Scenario: Legacy directory exists
- **WHEN** the repository is inspected
- **THEN** router-service-legacy/ SHALL contain the original code
- **AND** it SHALL NOT be included in Docker Compose or go.work


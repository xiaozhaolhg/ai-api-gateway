## Purpose

Deployment specifications for the router-service microservice with gRPC interface and Docker Compose orchestration.
## Requirements
### Requirement: Docker Multi-Stage Build
The router-service SHALL be built using a multi-stage Dockerfile that produces a minimal production image based on alpine, running a gRPC server (not an HTTP server).

#### Scenario: Build produces minimal image
- **WHEN** docker build is executed on the Dockerfile
- **THEN** the resulting image SHALL be based on alpine and contain only the binary plus CA certificates

#### Scenario: Container runs as non-root
- **WHEN** the container is deployed
- **THEN** it SHALL run with runAsNonRoot and runAsUser set to a non-root UID

### Requirement: Configuration via YAML and environment variables
The router-service SHALL load configuration from a YAML file with environment variable resolution.

#### Scenario: Config loads gRPC port
- **WHEN** the service starts with a config file
- **THEN** it SHALL load gRPC port, SQLite path, and provider-service address from that file

### Requirement: gRPC health check
The router-service SHALL implement gRPC health checking protocol.

#### Scenario: Health check returns SERVING
- **WHEN** a gRPC health check request is made
- **THEN** the service SHALL return SERVING status when ready to accept requests

### Requirement: Docker Compose integration
The router-service SHALL be deployable via Docker Compose with health check and dependency configuration.

#### Scenario: Service starts in Docker Compose
- **WHEN** `docker-compose up` is run
- **THEN** the router-service SHALL start on gRPC port 50052
- **AND** its health check SHALL pass before dependent services (gateway-service) start


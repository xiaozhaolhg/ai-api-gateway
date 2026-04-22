## MODIFIED Requirements

### Requirement: Docker multi-stage build
The billing-service SHALL be built using a multi-stage Dockerfile that produces a minimal production image based on alpine.

#### Scenario: Build produces minimal image
- **WHEN** docker build is executed on the Dockerfile
- **THEN** the resulting image SHALL be based on alpine and contain only the binary plus CA certificates

### Requirement: Container runs as non-root
The billing-service container SHALL run as a non-root user for security.

#### Scenario: Container security context
- **WHEN** the container is deployed
- **THEN** it SHALL run with runAsNonRoot and runAsUser set to a non-root UID

### Requirement: Configuration via YAML and environment variables
The billing-service SHALL load configuration from a YAML file with environment variable resolution.

#### Scenario: Config loads from mounted file
- **WHEN** the service starts with CONFIG_PATH env var pointing to a config file
- **THEN** it SHALL load gRPC port, SQLite path, and other settings from that file

### Requirement: gRPC health check
The billing-service SHALL implement gRPC health checking protocol.

#### Scenario: Health check returns SERVING
- **WHEN** a gRPC health check request is made
- **THEN** the service SHALL return SERVING status when ready to accept requests

### Requirement: Docker Compose integration
The billing-service SHALL be deployable via Docker Compose with health check configuration.

#### Scenario: Service starts in Docker Compose
- **WHEN** `docker-compose up` is run
- **THEN** the billing-service SHALL start on gRPC port 50054

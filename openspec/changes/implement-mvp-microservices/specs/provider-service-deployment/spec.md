## ADDED Requirements

### Requirement: Docker multi-stage build
The provider-service SHALL be built using a multi-stage Dockerfile that produces a minimal production image based on alpine.

#### Scenario: Build produces minimal image
- **WHEN** docker build is executed on the Dockerfile
- **THEN** the resulting image SHALL be based on alpine and contain only the binary plus CA certificates

### Requirement: Container runs as non-root
The provider-service container SHALL run as a non-root user for security.

#### Scenario: Container security context
- **WHEN** the container is deployed
- **THEN** it SHALL run with runAsNonRoot and runAsUser set to a non-root UID

### Requirement: Configuration via YAML and environment variables
The provider-service SHALL load configuration from a YAML file with environment variable resolution, including the encryption key for credentials.

#### Scenario: Encryption key from environment
- **WHEN** the config field for encryption_key contains `${ENCRYPTION_KEY}`
- **THEN** the service SHALL resolve the value from the environment variable

### Requirement: gRPC health check
The provider-service SHALL implement gRPC health checking protocol.

#### Scenario: Health check returns SERVING
- **WHEN** a gRPC health check request is made
- **THEN** the service SHALL return SERVING status when ready to accept requests

### Requirement: Docker Compose integration
The provider-service SHALL be deployable via Docker Compose with health check and dependency configuration.

#### Scenario: Service starts in Docker Compose
- **WHEN** `docker-compose up` is run
- **THEN** the provider-service SHALL start on gRPC port 50053
- **AND** its health check SHALL pass before dependent services start

## Purpose

Deployment specifications for the gateway-service microservice with HTTP interface and Docker Compose orchestration.
## Requirements
### Requirement: Docker multi-stage build
The gateway-service SHALL be built using a multi-stage Dockerfile that produces a minimal production image based on alpine.

#### Scenario: Build produces minimal image
- **WHEN** docker build is executed on the Dockerfile
- **THEN** the resulting image SHALL be based on alpine and contain only the binary plus CA certificates

### Requirement: Container runs as non-root
The gateway-service container SHALL run as a non-root user for security.

#### Scenario: Container security context
- **WHEN** the container is deployed
- **THEN** it SHALL run with runAsNonRoot and runAsUser set to a non-root UID

### Requirement: Configuration via YAML and environment variables
The gateway-service SHALL load configuration from a YAML file with environment variable resolution, including gRPC addresses for all 5 internal services.

#### Scenario: Service addresses from config
- **WHEN** the config file specifies auth-service address as "auth-service:50051"
- **THEN** the gateway SHALL connect to that address via gRPC

### Requirement: HTTP health check
The gateway-service SHALL provide an HTTP health check endpoint.

#### Scenario: Health endpoint returns ok
- **WHEN** GET /health is called
- **THEN** the service SHALL return 200 OK with `{"status": "ok"}`

### Requirement: Docker Compose integration
The gateway-service SHALL be deployable via Docker Compose with health check and dependency configuration.

#### Scenario: Service starts after dependencies
- **WHEN** `docker-compose up` is run
- **THEN** the gateway-service SHALL start on HTTP port 8080
- **AND** it SHALL wait for auth-service, router-service, and provider-service health checks to pass before accepting traffic


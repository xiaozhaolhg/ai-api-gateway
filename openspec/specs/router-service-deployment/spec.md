## Purpose

(TBD) Deployment specifications for the router-service component.

## Requirements

### Requirement: Docker Multi-Stage Build
The router-service SHALL be built using a multi-stage Dockerfile that produces a minimal production image.

#### Scenario: Build produces minimal image
- **WHEN** docker build is executed on the Dockerfile
- **THEN** the resulting image SHALL be based on alpine and contain only the binary plus CA certificates

#### Scenario: Image size is minimal
- **WHEN** docker images is checked for the router-service image
- **THEN** the image size SHALL be under 25MB

### Requirement: Container Runs as Non-Root
The router-service container SHALL run as a non-root user for security.

#### Scenario: Container security context
- **WHEN** the container is deployed in Kubernetes
- **THEN** it SHALL run with securityContext.runAsNonRoot: true and securityContext.runAsUser: 65532

### Requirement: Helm Chart Deployment
The router-service SHALL be deployable via a Helm chart to a KinD cluster.

#### Scenario: Helm install succeeds
- **WHEN** helm install router-service is executed with valid values
- **THEN** all Kubernetes resources (Deployment, Service, Ingress) SHALL be created successfully

#### Scenario: Config mounted from host
- **WHEN** the Helm chart is installed with config.hostPath set
- **THEN** the config file from the host path SHALL be mounted into the container at /app/config/config.yaml

### Requirement: Service Exposed via Ingress
The router-service SHALL be accessible through NGINX Ingress at /v1/* paths.

#### Scenario: Ingress routes traffic to service
- **WHEN** an HTTP request is made to http://localhost/v1/chat/completions
- **THEN** the request SHALL be routed to the router-service pod

#### Scenario: Health endpoint accessible
- **WHEN** an HTTP request is made to http://localhost/health
- **THEN** the router-service SHALL return a 200 OK response

### Requirement: Environment Variables
The router-service SHALL support configuration via environment variables.

#### Scenario: Port configuration via env
- **WHEN** the PORT environment variable is set to 8081
- **THEN** the service SHALL listen on port 8081 instead of default 8080

### Requirement: Resource Limits
The router-service deployment SHALL have defined resource requests and limits.

#### Scenario: Resource limits applied
- **WHEN** the deployment is created in Kubernetes
- **THEN** it SHALL have cpu.limit of 500m and memory.limit of 256Mi

### Requirement: KinD Image Loading
The router-service Docker image SHALL be loadable directly into KinD without an external registry.

#### Scenario: Kind load succeeds
- **WHEN** kind load docker-image router-service:latest is executed
- **THEN** the image SHALL be available in the KinD cluster without pushing to a registry

### Requirement: Build Docker Image
The Makefile SHALL provide a target to build the router-service Docker image.

#### Scenario: Build target exists
- **WHEN** `make build` is executed in router-service/
- **THEN** the Docker image `router-service:latest` SHALL be built

### Requirement: Deploy to Docker
The Makefile SHALL provide a target to build and run the router-service in Docker.

#### Scenario: deploy-docker target
- **WHEN** `make deploy-docker` is executed
- **THEN** the image SHALL be built and a container SHALL be started on port 8080

### Requirement: Deploy to KinD
The Makefile SHALL provide a target to build and deploy the router-service to KinD via Helm.

#### Scenario: deploy-kind target
- **WHEN** `make deploy-kind KIND_CLUSTER=<name>` is executed
- **THEN** the image SHALL be built, loaded into KinD, and deployed via Helm to namespace ai-gateway

### Requirement: Configurable Kind Cluster
The Makefile SHALL allow specifying the KinD cluster name at runtime.

#### Scenario: Custom cluster name
- **WHEN** `make deploy-kind KIND_CLUSTER=my-cluster` is executed
- **THEN** the image SHALL be loaded into the specified cluster
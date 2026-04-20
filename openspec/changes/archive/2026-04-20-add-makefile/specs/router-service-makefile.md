## Purpose

Define Makefile targets for building and deploying the router-service.

## Changes from Parent Spec

This is a delta spec adding Makefile tooling to the router-service-deployment spec.

## Requirements

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

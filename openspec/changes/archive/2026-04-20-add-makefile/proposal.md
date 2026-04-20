## Why

Currently, deploying the router-service requires manual docker build and kubectl/helm commands. Adding a Makefile with standardized targets will simplify the development workflow and provide consistent deployment options for both local Docker testing and KinD cluster deployment.

## What Changes

- Create Makefile in router-service/ with two deployment targets
- Target `deploy-docker`: builds the latest image and runs it in Docker with local config
- Target `deploy-kind`: builds the latest image, loads it into KinD cluster, and deploys via Helm
- Remove the deprecated k8s.yaml manifest file

## Capabilities

### New Capabilities

- **router-service-deployment**: Define Makefile targets for building and deploying the router-service

### Modified Capabilities

- None

## Impact

- File: `router-service/Makefile` (new)
- File: `router-service/deploy/k8s.yaml` (deleted)
- File: `router-service/charts/templates/configmap.yaml` (new - ConfigMap for config mount)
- File: `router-service/charts/templates/deployment.yaml` (modified - added volume/volumeMount)
- File: `router-service/charts/values.yaml` (modified - added config.data section)
- Directory: `router-service/charts/router-service/` → `router-service/charts/` (flattened)
- Commands available: `make deploy-docker`, `make deploy-kind`, `make check-prereqs`

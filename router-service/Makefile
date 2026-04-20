# Makefile for router-service
#
# Targets:
#   build         - Build Docker image
#   deploy-docker - Build and run in Docker
#   deploy-kind   - Build, load to KinD, and deploy via Helm
#   check-prereqs - Verify required tools are installed
#
# Prerequisites:
#   docker - Required for all targets
#   kind   - Required for deploy-kind
#   helm   - Required for deploy-kind
#
# Variables:
#   KIND_CLUSTER - KinD cluster name (default: desktop, e.g., make deploy-kind KIND_CLUSTER=my-cluster)
#
# Ports:
#   8080 - Used by deploy-docker for the router-service HTTP endpoint

IMAGE_NAME := router-service:latest
KIND_CLUSTER ?= desktop
CHART_PATH := ./charts
NAMESPACE := ai-gateway
CONFIG_PATH := $(PWD)/../configs/config.yaml

.PHONY: build deploy-docker deploy-kind check-prereqs

# Build Docker image
build:
	docker build -t $(IMAGE_NAME) .

# Build and run in Docker
deploy-docker: build
	docker rm -f router-service 2>/dev/null || true
	docker run -d --name router-service -p 8080:8080 \
		-v $(CONFIG_PATH):/app/config/config.yaml:ro \
		-e CONFIG_PATH=/app/config/config.yaml \
		$(IMAGE_NAME)

# Check that required tools are installed
check-prereqs:
	@which kind > /dev/null 2>&1 || (echo "Error: kind is not installed. Install from https://kind.sigs.k8s.io" && exit 1)
	@which helm > /dev/null 2>&1 || (echo "Error: helm is not installed. Install from https://helm.sh" && exit 1)
	@kind get clusters 2>/dev/null | grep -q $(KIND_CLUSTER) || (echo "Error: KinD cluster '$(KIND_CLUSTER)' not found. Create it with: kind create cluster --name $(KIND_CLUSTER)" && exit 1)

# Build, load to KinD, and deploy via Helm
deploy-kind: build check-prereqs
	kind load docker-image $(IMAGE_NAME) --name $(KIND_CLUSTER)
	helm upgrade --install router-service $(CHART_PATH) \
		--namespace $(NAMESPACE) \
		--create-namespace \
		--set image.tag=latest

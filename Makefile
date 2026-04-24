.PHONY: build test up down proto clean clean-images help
.SILENT: clean help

SERVICES = auth-service router-service provider-service gateway-service billing-service monitor-service

# Show help
help:
	@echo "AI API Gateway - Makefile Commands"
	@echo ""
	@echo "Usage: make <target> [SERVICE=<service>]"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build all services or SERVICE=<name>"
	@echo "  test         Test all services or SERVICE=<name>"
	@echo "  up           Start all services with Docker Compose"
	@echo "  down         Stop all services"
	@echo "  proto        Generate protobuf stubs"
	@echo "  clean        Clean build artifacts ( SERVICE=<name> )"
	@echo "  clean-images Clean Docker images ( SERVICE=<name> )"
	@echo "  help         Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make test SERVICE=gateway-service"
	@echo "  make build SERVICE=auth-service"
	@echo "  make clean-images SERVICE=gateway-service"

# Build all services or a specific one
build:
ifneq ($(SERVICE),)
	@echo "Building $(SERVICE)..."
	$(MAKE) -C $(SERVICE) build
else
	@echo "Building all services..."
	@for dir in $(SERVICES); do \
		if [ -d "$$dir" ]; then \
			echo "Building $$dir..."; \
			$(MAKE) -C $$dir build; \
		fi; \
	done
endif

# Test all services or a specific one
test:
ifneq ($(SERVICE),)
	@echo "Testing $(SERVICE)..."
	cd $(SERVICE) && go test -v ./...
else
	@echo "Testing all services..."
	@for dir in $(SERVICES); do \
		if [ -d "$$dir" ]; then \
			echo "Testing $$dir..."; \
			cd $$dir && go test -v ./... && cd ..; \
		fi; \
	done
endif

# Start all services with Docker Compose
up:
	docker compose up -d

# Stop all services with Docker Compose
down:
	docker compose down

# Generate protobuf stubs
proto:
	@echo "Generating protobuf stubs..."
	$(MAKE) -C api proto

# Clean build artifacts for all services or a specific one
clean:
ifneq ($(SERVICE),)
	@echo "Cleaning $(SERVICE)..."
	$(MAKE) -C $(SERVICE) clean || true
else
	@echo "Cleaning all build artifacts..."
	@for dir in $(SERVICES); do \
		if [ -d "$$dir" ]; then \
			echo "Cleaning $$dir..."; \
			$(MAKE) -C $$dir clean || true; \
		fi; \
	done
endif

# Clean all images or a specific service image
clean-images:
ifneq ($(SERVICE),)
	@echo "Cleaning image for $(SERVICE)..."
	docker rmi ai-api-gateway-$(SERVICE):latest >/dev/null 2>&1 || true
	@docker image prune -f --filter "label=service=$(SERVICE)" >/dev/null 2>&1 || true
else
	@echo "Cleaning all images..."
	@for image in $(shell docker images | grep ai-api-gateway | awk '{print $$1}'); do \
		echo "Cleaning $$image..."; \
		docker rmi $$image:latest >/dev/null 2>&1 || true; \
	done
endif
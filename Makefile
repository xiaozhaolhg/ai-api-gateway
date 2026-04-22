.PHONY: build test up down proto clean

# Build all services
build:
	@echo "Building all services..."
	@for dir in auth-service router-service provider-service gateway-service billing-service monitor-service; do \
		if [ -d "$$dir" ]; then \
			echo "Building $$dir..."; \
			$(MAKE) -C $$dir build; \
		fi; \
	done

# Test all services
test:
	@echo "Testing all services..."
	@for dir in auth-service router-service provider-service gateway-service billing-service monitor-service; do \
		if [ -d "$$dir" ]; then \
			echo "Testing $$dir..."; \
			cd $$dir && go test -v ./... && cd ..; \
		fi; \
	done

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

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@for dir in auth-service router-service provider-service gateway-service billing-service monitor-service; do \
		if [ -d "$$dir" ]; then \
			echo "Cleaning $$dir..."; \
			$(MAKE) -C $$dir clean || true; \
		fi; \
	done

clean-images:
	@echo "Cleaning images..."
	@for image in $(shell docker images | grep ai-api-gateway | awk '{print $$1}'); do \
		echo "Cleaning $$image..."; \
		docker rmi $$image:latest || true; \
	done
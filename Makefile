.PHONY: build test up down proto clean

# Build all services
build:
	@echo "Building all services..."
	@for dir in api auth-service router-service provider-service gateway-service billing-service monitor-service; do \
		if [ -d "$$dir" ]; then \
			echo "Building $$dir..."; \
			$(MAKE) -C $$dir build; \
		fi; \
	done

# Test all services
test:
	@echo "Testing all services..."
	@for dir in api auth-service router-service provider-service gateway-service billing-service monitor-service; do \
		if [ -d "$$dir" ]; then \
			echo "Testing $$dir..."; \
			$(MAKE) -C $$dir test; \
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

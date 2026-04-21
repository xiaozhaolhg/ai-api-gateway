.PHONY: build test clean docker proto

# Build the service
build:
	go build -o bin/auth-service ./cmd/server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker:
	docker build -t ai-api-gateway/auth-service:latest .

# Run Docker container
docker-run:
	docker run -p 50051:50051 ai-api-gateway/auth-service:latest

# Generate protobuf stubs (requires buf)
proto:
	buf generate

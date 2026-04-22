.PHONY: build test clean docker proto

# Build the service
build:
	go build -o bin/gateway-service ./cmd/server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker:
	docker build -t ai-api-gateway/gateway-service:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 ai-api-gateway/gateway-service:latest

# Generate protobuf stubs (requires buf)
proto:
	buf generate

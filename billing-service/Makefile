.PHONY: build test clean docker proto

# Build the service
build:
	go build -o bin/billing-service ./cmd/server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker:
	docker build -t ai-api-gateway/billing-service:latest .

# Run Docker container
docker-run:
	docker run -p 50054:50054 ai-api-gateway/billing-service:latest

# Generate protobuf stubs (requires buf)
proto:
	buf generate

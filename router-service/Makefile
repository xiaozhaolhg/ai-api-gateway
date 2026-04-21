.PHONY: build test clean docker proto

# Build the service
build:
	go build -o bin/router-service ./cmd/server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker:
	docker build -t ai-api-gateway/router-service:latest .

# Run Docker container
docker-run:
	docker run -p 50052:50052 ai-api-gateway/router-service:latest

# Generate protobuf stubs (requires buf)
proto:
	buf generate

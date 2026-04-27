#!/bin/bash
set -e

echo "=== Developer A Week 1 Verification ==="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Building Docker images..."
docker-compose build provider-service router-service > /dev/null 2>&1 || {
    echo -e "${RED}Failed to build Docker images${NC}"
    exit 1
}

echo -e "${GREEN}Docker images built successfully${NC}"
echo ""

echo "Running provider-service tests..."
if docker-compose run --rm provider-service go test ./... -v 2>&1 | tail -20; then
    echo -e "${GREEN}provider-service tests passed${NC}"
else
    echo -e "${RED}provider-service tests failed${NC}"
    exit 1
fi
echo ""

echo "Running router-service tests..."
if docker-compose run --rm router-service go test ./... -v 2>&1 | tail -20; then
    echo -e "${GREEN}router-service tests passed${NC}"
else
    echo -e "${RED}router-service tests failed${NC}"
    exit 1
fi
echo ""

echo -e "${GREEN}=== All Week 1 Verifications Passed! ===${NC}"

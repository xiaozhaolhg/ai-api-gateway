---
name: gateway-sync-upstream
description: Fetches and rebases the current branch with the latest upstream/main changes for the AI API Gateway project, then verifies the integration by running comprehensive tests and service health checks. Use before starting new work to ensure branch is up-to-date. Do NOT use when branch is already synced, when working on main branch, or when rebasing would disrupt ongoing work.
---

# Gateway Sync Upstream

Fetches and rebases the current branch with the latest upstream/main changes for the AI API Gateway project.

## Purpose

Ensure the current branch is up-to-date with upstream/main before starting new work, then verify that all services can run successfully after the sync.

## Process

### 1. Fetch Upstream
```bash
git fetch upstream main
```

If user environment is in WSL, execute:
```
export GIT_DISCOVERY_ACROSS_FILESYSTEM=1
```

### 2. Check Current Branch and Working Directory
```bash
git branch --show-current
git status --short
```

Note the branch name for conflict resolution context.

**If working directory is NOT clean (has uncommitted changes):**

Ask user: "Working directory has uncommitted changes. What would you like to do?"

Options:
- **stash**: `git stash push -m "WIP before upstream sync"` - Save changes temporarily
- **commit**: `git add <files> && git commit -m "WIP: describe changes"` - Commit changes first
- **drop**: `git checkout -- . && git clean -fd` - Discard all changes (WARNING: data loss)
- **abort**: Stop sync and let user handle manually

Wait for user choice before proceeding.

### 3. Rebase onto Upstream
```bash
git rebase upstream/main
```

### 4. Handle Conflicts
If conflicts occur:
- Call `git-rebase-workflow` to resolve
- Apply learned patterns for gateway project:
  - Preserve upstream structural changes (imports, client initialization like billingClient, authClient)
  - Merge local logic into new structure
  - Verify dependencies (clients initialized before use)
  - Keep upstream's file organization

### 5. Verify Result
```bash
git log --oneline -5
git diff upstream/main
```

Ensure rebase completed successfully.

**If changes were stashed in Step 2:**
```bash
git stash pop
```
Restore stashed changes after successful rebase.

### 6. Comprehensive Integration Verification

#### 6.1. Build Verification
```bash
# Build all services to ensure no compilation errors
go build ./...
go test -build-only ./...
```

#### 6.2. Service Health Verification
```bash
# Clean up any existing containers and images
make down
make clean-images

# Start all services
make up

# Wait for services to be ready (30-60 seconds)
sleep 45

# Check service health
curl -f http://localhost:8080/health || echo "Gateway service health check failed"
curl -f http://localhost:50051/grpc.health.v1.Health/Check || echo "Auth service health check failed"
curl -f http://localhost:50052/grpc.health.v1.Health/Check || echo "Router service health check failed"
curl -f http://localhost:50053/grpc.health.v1.Health/Check || echo "Provider service health check failed"
curl -f http://localhost:50054/grpc.health.v1.Health/Check || echo "Billing service health check failed"
curl -f http://localhost:50055/grpc.health.v1.Health/Check || echo "Monitor service health check failed"

# Check Docker container status
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

#### 6.3. Integration Tests
```bash
# Run unit tests for all services
go test -v ./...

# Run integration tests (if available)
make test-integration || echo "Integration tests not available or failed"

# Test API endpoints
curl -X GET http://localhost:8080/gateway/models || echo "Models endpoint test failed"
curl -X GET http://localhost:8080/gateway/health || echo "Gateway health endpoint test failed"
```

#### 6.4. Service Communication Verification
```bash
# Test service-to-service communication
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model": "test:model", "messages": [{"role": "user", "content": "test"}]}' \
  || echo "Chat completions endpoint test failed"

# Check service logs for errors
docker logs gateway-service --tail=20 | grep -i error || echo "No errors in gateway-service logs"
docker logs auth-service --tail=20 | grep -i error || echo "No errors in auth-service logs"
```

### 7. Cleanup and Final Verification
```bash
# Optional: Clean up test environment
make down

# Final git status check
git status

# Ensure working directory is clean
git diff --exit-code
```

## Gateway-Specific Conflict Patterns

### Client Initialization Conflicts
When upstream adds new gRPC clients (e.g., billingClient, monitorClient):
- Accept upstream's client initialization pattern
- Merge local handler usage into new structure
- Verify client is initialized before handler routes are registered

### Import Conflicts
When upstream adds new proto imports:
- Accept all upstream imports
- Add any local-only imports that don't duplicate
- Remove unused imports after resolution

### Handler Structure Conflicts
When upstream reorganizes handler functions:
- Follow upstream's direct gRPC call pattern
- Adapt local logic to upstream structure
- Preserve local business rules

## Success Indicators

### Git Status
- `git status` shows clean working directory
- `git log` shows linear history from upstream/main
- No conflict markers in any file

### Build and Compilation
- `go build ./...` passes for all services
- `go test -build-only ./...` succeeds
- No compilation errors or missing imports

### Service Health
- All Docker containers start successfully (`make up`)
- All services respond to health checks:
  - Gateway: `http://localhost:8080/health`
  - Auth: `localhost:50051/grpc.health.v1.Health/Check`
  - Router: `localhost:50052/grpc.health.v1.Health/Check`
  - Provider: `localhost:50053/grpc.health.v1.Health/Check`
  - Billing: `localhost:50054/grpc.health.v1.Health/Check`
  - Monitor: `localhost:50055/grpc.health.v1.Health/Check`
- `docker ps` shows all services running with healthy status

### Integration Tests
- Unit tests pass (`go test -v ./...`)
- Integration tests pass (`make test-integration` if available)
- API endpoints respond correctly:
  - `/gateway/models` returns model list
  - `/gateway/health` returns health status
  - `/v1/chat/completions` accepts requests (may return provider errors)

### Service Communication
- No error messages in service logs
- Service-to-service communication works
- gRPC clients can connect to respective services

## Failure Recovery

### Rebase Failures
If rebase fails repeatedly:
- Abort: `git rebase --abort`
- Create fresh branch from upstream/main
- Cherry-pick specific commits if needed
- Ask user for manual intervention

### Service Health Failures
If services fail to start or health checks fail:
- Check Docker logs: `docker logs <service-name>`
- Verify port availability: `netstat -tlnp | grep :8080`
- Clean up containers: `make down && make clean-images`
- Check for resource constraints: `docker system df`
- Verify configuration files exist: `ls -la configs/`

### Build Failures
If compilation fails:
- Check Go modules: `go mod tidy && go mod download`
- Verify proto generation: `make proto`
- Check for missing dependencies: `go mod why <package>`
- Clean build cache: `go clean -modcache && go clean -cache`

### Integration Test Failures
If tests fail:
- Check service dependencies: ensure all services are running
- Verify test configuration: check test config files
- Check network connectivity between services
- Review test logs for specific error messages
- Run individual service tests: `go test ./<service>/...`

### Rollback Strategy
If verification fails completely:
1. Stop all services: `make down`
2. Reset to previous working state: `git reset --hard HEAD~1`
3. Restore stashed changes if needed: `git stash pop`
4. Document the failure for future reference
5. Consider creating a new branch from a known good state

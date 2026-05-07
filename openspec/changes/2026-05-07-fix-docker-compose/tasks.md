# Tasks: Fix Docker Compose Configuration

**Owner**: Developer C (cynkiller)  
**Collaborators**: Developer A (testing), Developer B (service config)  
**Status**: Planning

## Phase 1: Fix monitor-service Port Mapping (High Priority)

- [ ] **Task 1.1**: Add ports to monitor-service
  - [ ] Edit `docker-compose.yml`
  - [ ] Add `ports:` section to monitor-service
  - [ ] Map `"50055:50055"` (host:container)
  - **Acceptance**: `grep -A5 "monitor-service:" docker-compose.yml` shows ports section

- [ ] **Task 1.2**: Verify port is accessible
  - [ ] Start monitor-service only: `docker compose up -d monitor-service`
  - [ ] Wait for startup (30s)
  - [ ] Test: `curl http://localhost:50055/health` or `grpcurl -plaintext localhost:50055 list`
  - [ ] Verify port is accessible from host
  - **Acceptance**: Port 50055 returns response (not connection refused)

## Phase 2: Add Health Checks (High Priority)

- [ ] **Task 2.1**: Add healthcheck to gateway-service
  - [ ] Edit `docker-compose.yml`
  - [ ] Add `healthcheck:` section to gateway-service
  - [ ] Use HTTP health endpoint: `http://localhost:8080/health`
  - **Acceptance**: `docker compose ps` shows "Up (healthy)" for gateway

- [ ] **Task 2.2**: Add healthcheck to gRPC services (auth, router, provider, billing, monitor)
  - [ ] Edit `docker-compose.yml`
  - [ ] Add `healthcheck:` to each gRPC service
  - [ ] Use `grpc_health_probe` tool (or equivalent)
  - **Acceptance**: All gRPC services show "Up (healthy)"

- [ ] **Task 2.3**: Add healthcheck to redis
  - [ ] Edit `docker-compose.yml`
  - [ ] Use `redis-cli ping` for health check
  - **Acceptance**: Redis shows "Up (healthy)"

## Phase 3: Verify Service Dependencies (High Priority)

- [ ] **Task 3.1**: Verify depends_on relationships
  - [ ] Check gateway-service depends_on: auth, router, provider
  - [ ] Check provider-service depends_on: billing, monitor
  - [ ] Check router-service depends_on: redis
  - [ ] Verify condition is `service_healthy` where possible
  - **Acceptance**: `docker compose config` validates without errors

- [ ] **Task 3.2**: Test startup order
  - [ ] Run `docker compose up -d --build`
  - [ ] Monitor logs: `docker compose logs -f`
  - [ ] Verify services start in correct order
  - [ ] Check no "connection refused" errors in logs
  - **Acceptance**: All services start without dependency errors

## Phase 4: Full Deployment Test (High Priority)

- [ ] **Task 4.1**: Clean environment
  - [ ] Stop all containers: `docker compose down -v`
  - [ ] Remove old images (optional): `docker compose rm -f`
  - [ ] Verify clean state: `docker compose ps`
  - **Acceptance**: No containers running, clean state

- [ ] **Task 4.2**: Build and start all services
  - [ ] Run `docker compose up -d --build`
  - [ ] Wait for all services to be healthy (max 3 minutes)
  - [ ] Check status: `docker compose ps`
  - **Acceptance**: All 7 services show "Up (healthy)"

- [ ] **Task 4.3**: Verify all endpoints accessible
  - [ ] Test gateway: `curl http://localhost:8080/health`
  - [ ] Test gateway deep health: `curl http://localhost:8080/gateway/health`
  - [ ] Test auth: `grpcurl -plaintext localhost:50051 list` (or equivalent)
  - [ ] Test monitor: `grpcurl -plaintext localhost:50055 list`
  - [ ] Verify all services respond correctly
  - **Acceptance**: All health endpoints return healthy status

## Phase 5: Documentation & Cleanup (Medium Priority)

- [ ] **Task 5.1**: Update QUICKSTART.md
  - [ ] Verify `QUICKSTART.md` reflects fixed Docker Compose
  - [ ] Add note about monitor-service port now exposed
  - [ ] Document health check endpoints
  - **Acceptance**: Documentation matches implementation

- [ ] **Task 5.2**: Code cleanup
  - [ ] Search for TODO comments related to Docker Compose
  - [ ] Verify no hardcoded localhost references (use service names)
  - [ ] Check all environment variables are documented
  - **Acceptance**: No TODOs related to Docker Compose

## Summary

| Phase | Tasks | Priority | Dependencies |
|-------|-------|----------|--------------|
| Phase 1: Fix Ports | 2 | **High** | None |
| Phase 2: Health Checks | 3 | **High** | Phase 1 |
| Phase 3: Dependencies | 2 | **High** | Phase 2 |
| Phase 4: Full Test | 3 | **High** | Phase 3 |
| Phase 5: Documentation | 2 | **Medium** | Phase 4 |
| **Total** | **12** | | |

## Timeline Estimate

| Phase | Estimated Time | Owner |
|-------|----------------|--------|
| Phase 1 | 0.5 hour | Developer C |
| Phase 2 | 1 hour | Developer C |
| Phase 3 | 1 hour | Developer C + Dev A |
| Phase 4 | 1.5 hours | Developer C + Dev A |
| Phase 5 | 0.5 hour | Developer C |
| **Total** | **4.5 hours** | |

## Critical Path

```
Phase 1 (Fix Ports) → Phase 2 (Health Checks) → Phase 3 (Dependencies)
                                                        ↓
                                              Phase 4 (Full Test) → Phase 5 (Document)
```

**Blocker**: Phase 4 (full test) cannot start until Phases 1-3 are complete.

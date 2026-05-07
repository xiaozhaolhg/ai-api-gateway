# Proposal: Fix Docker Compose Configuration for Production

## Problem Statement

The `docker-compose.yml` has incomplete service configurations that prevent full system deployment:

1. **Monitor-service port mapping missing**: The `monitor-service` section lacks `ports:` configuration, making it inaccessible from host for health checks and debugging.
2. **Service dependency incomplete**: The `provider-service` depends on `billing-service` and `monitor-service`, but health check verification is needed to ensure proper startup order.
3. **Health check endpoints**: Need verification that all services' health endpoints are properly configured and accessible.

**Impact**:
- Cannot access monitor-service from host (port 50055 not exposed)
- Unclear if services start in correct order
- Docker Compose deployment may fail silently if dependencies aren't ready

## Proposed Solution

Fix `docker-compose.yml` to:
1. Add port mapping for monitor-service (`:50055:50055`)
2. Add health check configurations for all services
3. Verify service dependencies and startup order
4. Test full `docker compose up -d` deployment

## Scope

- **In Scope**:
  - Add `ports:` to monitor-service in `docker-compose.yml`
  - Add `healthcheck:` configurations for all services
  - Verify `depends_on:` relationships are correct
  - Test `docker compose up -d` and verify all services healthy
  - Update `QUICKSTART.md` if needed

- **Out of Scope**:
  - Changes to individual service Dockerfiles
  - Adding new services
  - Changes to Kubernetes manifests (if exist)
  - Database migration scripts

## Success Criteria

- `monitor-service` port 50055 is accessible from host
- All services have health check endpoints configured
- `docker compose up -d` brings up all services successfully
- `curl http://localhost:50055/health` returns healthy status
- All services pass health checks within 60 seconds of startup
- `docker compose ps` shows all services as "Up (healthy)"

## Dependencies

- **Docker Compose**: Must be installed and configured
- **All service images**: Must be built (or use `build:` contexts)
- **Network**: `ai-gateway` network must be properly configured
- **Volumes**: Data persistence volumes must be available

## Owner

- **Primary**: Developer C (cynkiller)
- **Collaborators**: Developer A (testing), Developer B (service config verification)

## References

- Original task: `docs/phase1_work_division.md` → Developer C → Week 4 → "Docker Compose setup: Gateway + PostgreSQL + Redis"
- Current config: `docker-compose.yml`
- Related: `QUICKSTART.md` (deployment guide)

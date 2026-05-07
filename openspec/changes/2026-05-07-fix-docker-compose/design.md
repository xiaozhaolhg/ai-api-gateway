# Design: Fix Docker Compose Configuration

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Docker Compose Network                      │
│                    (ai-gateway)                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────┐     ┌──────────────┐     ┌─────────────┐│
│  │  Gateway   │────▶│  Router     │────▶│  Provider   ││
│  │  Service   │     │  Service    │     │  Service    ││
│  │  :8080     │     │  :50052     │     │  :50053     ││
│  └──────┬──────┘     └──────┬───────┘     └──────┬──────┘│
│         │                 │                 │              │
│         ▼                 ▼                 ▼              │
│  ┌─────────────┐     ┌──────────────┐     ┌─────────────┐│
│  │   Auth     │     │  Billing    │     │  Monitor    ││
│  │  Service   │     │  Service    │     │  Service    ││
│  │  :50051     │     │  :50054     │     │  :50055 ←───││
│  └─────────────┘     └──────────────┘     └─────────────┘│
│                                                              │
│  ┌─────────────┐     ┌──────────────┐                    │
│  │   Redis    │     │  PostgreSQL │                    │
│  │   :6379    │     │  :5432       │                    │
│  └─────────────┘     └──────────────┘                    │
└─────────────────────────────────────────────────────────────┘
```

**Current Issues**:
- Monitor-service `:50055` port NOT exposed (missing `ports:` config)
- Health checks needed for all services
- Startup order must be verified

## 1. Docker Compose Structure

### 1.1 Current monitor-service Configuration

**File**: `docker-compose.yml` (monitor-service section)

**Current (Broken)**:
```yaml
monitor-service:
  build:
    context: .
    dockerfile: ./monitor-service/Dockerfile
  ports:
    # ← MISSING PORT MAPPING
  networks:
    - ai-gateway
```

**After (Fixed)**:
```yaml
monitor-service:
  build:
    context: .
    dockerfile: ./monitor-service/Dockerfile
  ports:
    - "50055:50055"           # NEW: Expose gRPC port
  networks:
    - ai-gateway
  healthcheck:                   # NEW: Health check
    test: ["CMD", "grpc_health_probe", "-addr=localhost:50055"]
    interval: 30s
    timeout: 10s
    retries: 3
    start_period: 40s
```

### 1.2 Complete Service Configuration (Reference)

```yaml
services:
  gateway-service:
    build:
      context: .
      dockerfile: ./gateway-service/Dockerfile
    ports:
      - "8080:8080"
    depends_on:
      auth-service:
        condition: service_healthy
      router-service:
        condition: service_healthy
      provider-service:
        condition: service_healthy
    networks:
      - ai-gateway
    environment:
      - AUTH_SERVICE_ADDRESS=auth-service:50051
      - ROUTER_SERVICE_ADDRESS=router-service:50052
      - PROVIDER_SERVICE_ADDRESS=provider-service:50053
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  auth-service:
    build:
      context: .
      dockerfile: ./auth-service/Dockerfile
    ports:
      - "50051:50051"
    volumes:
      - auth-service-data:/data
    networks:
      - ai-gateway
    healthcheck:
      test: ["CMD", "grpc_health_probe", "-addr=localhost:50051"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  router-service:
    build:
      context: .
      dockerfile: ./router-service/Dockerfile
    ports:
      - "50052:50052"
    volumes:
      - router-service-data:/data
    depends_on:
      redis:
        condition: service_started
    networks:
      - ai-gateway
    environment:
      - REDIS_ADDRESS=redis:6379
    healthcheck:
      test: ["CMD", "grpc_health_probe", "-addr=localhost:50052"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  provider-service:
    build:
      context: .
      dockerfile: ./provider-service/Dockerfile
    ports:
      - "50053:50053"
    volumes:
      - provider-service-data:/data
    depends_on:
      billing-service:
        condition: service_healthy
      monitor-service:
        condition: service_healthy
    networks:
      - ai-gateway
    environment:
      - BILLING_SERVICE_ADDRESS=billing-service:50054
      - MONITOR_SERVICE_ADDRESS=monitor-service:50055
    healthcheck:
      test: ["CMD", "grpc_health_probe", "-addr=localhost:50053"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  billing-service:
    build:
      context: .
      dockerfile: ./billing-service/Dockerfile
    ports:
      - "50054:50054"
    volumes:
      - billing-service-data:/data
    networks:
      - ai-gateway
    healthcheck:
      test: ["CMD", "grpc_health_probe", "-addr=localhost:50054"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  monitor-service:                     # FIXED
    build:
      context: .
      dockerfile: ./monitor-service/Dockerfile
    ports:
      - "50055:50055"           # NEW
    networks:
      - ai-gateway
    healthcheck:                   # NEW
      test: ["CMD", "grpc_health_probe", "-addr=localhost:50055"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    networks:
      - ai-gateway
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 30s
      timeout: 5s
      retries: 3

networks:
  ai-gateway:
    driver: bridge

volumes:
  auth-service-data:
  router-service-data:
  provider-service-data:
  billing-service-data:
```

## 2. Service Dependencies

### 2.1 Dependency Graph

```
gateway-service
  ├── depends_on: auth-service (healthy)
  ├── depends_on: router-service (healthy)
  └── depends_on: provider-service (healthy)

provider-service
  ├── depends_on: billing-service (healthy)
  └── depends_on: monitor-service (healthy)

router-service
  └── depends_on: redis (started)

billing-service
  └── (no dependencies)

monitor-service
  └── (no dependencies)

auth-service
  └── (no dependencies)
```

### 2.2 Startup Order

1. **First**: redis, auth-service, billing-service, monitor-service (parallel)
2. **Second**: router-service (needs redis), provider-service (needs billing + monitor)
3. **Third**: gateway-service (needs auth + router + provider)

**Total startup time**: ~2-3 minutes (with health checks)

## 3. Health Check Strategy

### 3.1 Health Check Methods

| Service | Method | Tool | Endpoint |
|---------|--------|------|----------|
| gateway-service | HTTP | curl | `http://localhost:8080/health` |
| auth-service | gRPC | grpc_health_probe | `localhost:50051` |
| router-service | gRPC | grpc_health_probe | `localhost:50052` |
| provider-service | gRPC | grpc_health_probe | `localhost:50053` |
| billing-service | gRPC | grpc_health_probe | `localhost:50054` |
| monitor-service | gRPC | grpc_health_probe | `localhost:50055` |
| redis | Native | redis-cli ping | `localhost:6379` |

### 3.2 Health Check Configuration

All services use:
- **interval**: 30s (check every 30 seconds)
- **timeout**: 10s (fail if no response in 10s)
- **retries**: 3 (mark unhealthy after 3 failures)
- **start_period**: 40s (wait 40s before first check)

## 4. Verification Steps

### 4.1 Pre-Deployment Check

```bash
# Verify Docker Compose syntax
docker compose config

# Check all required files exist
ls -la gateway-service/Dockerfile
ls -la auth-service/Dockerfile
ls -la router-service/Dockerfile
ls -la provider-service/Dockerfile
ls -la billing-service/Dockerfile
ls -la monitor-service/Dockerfile
```

### 4.2 Deployment Test

```bash
# Clean previous state
docker compose down -v
docker compose rm -f

# Build and start all services
docker compose up -d --build

# Wait for startup (monitor progress)
docker compose logs -f

# In another terminal: check service status
docker compose ps

# Expected output: all services "Up (healthy)"
```

### 4.3 Post-Deployment Verification

```bash
# Test gateway (HTTP)
curl http://localhost:8080/health
curl http://localhost:8080/gateway/health

# Test auth (gRPC via grpcurl or similar)
grpcurl -plaintext localhost:50051 list

# Test monitor (gRPC)
grpcurl -plaintext localhost:50055 list

# Check all containers are running
docker compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
```

## 5. File Structure

```
.
├── docker-compose.yml          # UPDATE: Add monitor-service ports + healthchecks
├── gateway-service/Dockerfile    # (no changes needed)
├── auth-service/Dockerfile      # (no changes needed)
├── router-service/Dockerfile     # (no changes needed)
├── provider-service/Dockerfile   # (no changes needed)
├── billing-service/Dockerfile    # (no changes needed)
├── monitor-service/Dockerfile    # (no changes needed)
└── QUICKSTART.md              # UPDATE: Reflect fixed Docker Compose
```

## 6. Rationale

| Approach | Pros | Cons |
|----------|------|------|
| **Fix docker-compose.yml (chosen)** | Complete solution, health checks, proper startup order | Requires testing all services |
| Minimal fix (only ports) | Quick, minimal changes | No health checks, startup order issues |
| Kubernetes only | Better for production | More complex, not Phase 1 goal |

**Decision**: Fix `docker-compose.yml` completely with health checks and proper dependencies to ensure reliable deployment.

## 7. Backward Compatibility

- **Existing deployments**: Will need to `docker compose down` and `up` again (breaking change for running systems)
- **Port mappings**: Adding monitor-service port doesn't break existing functionality
- **Health checks**: New feature, doesn't affect services without health endpoint

**Migration**: Stop all containers → Apply changes → Start all containers.

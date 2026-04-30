# Design: Gateway Week 4 Completion

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     gateway-service                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  HTTP Entry                                                      │
│  ├── /health (liveness)                                        │
│  ├── /gateway/health (readiness with deps)                      │
│  ├── /gateway/models (aggregated)                               │
│  ├── /v1/chat/completions (with full middleware)                │
│  └── /admin/* (with auth)                                       │
│                                                                  │
├─────────────────────────────────────────────────────────────────┤
│  Middleware Chain                                                │
│  1. LogMiddleware (structured JSON)                            │
│  2. AuthMiddleware (API key → JWT)                               │
│  3. RateLimitMiddleware (placeholder)                           │
│  4. RouteMiddleware (model → provider)                          │
│  5. ProxyMiddleware (with timeout & error translation)           │
├─────────────────────────────────────────────────────────────────┤
│  Error Handling                                                  │
│  ├── Domain errors (internal/errors/)                          │
│  ├── Error translation middleware                              │
│  └── HTTP response format                                      │
├─────────────────────────────────────────────────────────────────┤
│  gRPC Clients                                                    │
│  ├── AuthClient (lazy connect)                                  │
│  ├── RouterClient (lazy connect)                              │
│  ├── ProviderClient (lazy connect)                            │
│  └── BillingClient (NEW - real implementation)                │
└─────────────────────────────────────────────────────────────────┘
```

## 1. Error Handling Design

### Error Types (internal/errors/)

```go
type ErrorCode string

const (
    ErrProviderTimeout       ErrorCode = "gateway_timeout"
    ErrInvalidCredentials    ErrorCode = "invalid_api_key"
    ErrModelNotFound         ErrorCode = "model_not_found"
    ErrProviderUnavailable   ErrorCode = "provider_error"
    ErrRateLimitExceeded     ErrorCode = "rate_limit_exceeded"
    ErrAuthDenied            ErrorCode = "insufficient_permissions"
    ErrInternal              ErrorCode = "internal_error"
)

type GatewayError struct {
    Code    ErrorCode
    Message string
    Details string
    Cause   error
}
```

### gRPC to HTTP Mapping

| gRPC Status | HTTP Status | Error Code | When |
|-------------|-------------|------------|------|
| DEADLINE_EXCEEDED | 504 | gateway_timeout | Provider >30s |
| UNAUTHENTICATED | 401 | invalid_api_key | Invalid API key |
| NOT_FOUND | 404 | model_not_found | No route for model |
| UNAVAILABLE | 502 | provider_error | Provider service down |
| RESOURCE_EXHAUSTED | 429 | rate_limit_exceeded | Quota exceeded |
| PERMISSION_DENIED | 403 | insufficient_permissions | Auth denied |
| INTERNAL | 500 | internal_error | Unexpected error |

### Response Format

```json
{
  "error": {
    "code": "gateway_timeout",
    "message": "Provider request timed out after 30s",
    "details": "ollama:llama2 connection timeout"
  }
}
```

## 2. Logging Middleware Design

### Log Levels

- **Access**: Every HTTP request (method, path, status, duration, user_id)
- **Request**: Request body (configurable, with sensitive field masking)
- **Response**: Response body (only on error)
- **Audit**: Admin operations (create/delete provider)

### Structured JSON Format

```json
{
  "timestamp": "2026-01-15T10:30:00Z",
  "level": "info",
  "request_id": "req_abc123",
  "method": "POST",
  "path": "/v1/chat/completions",
  "status": 200,
  "duration_ms": 1250,
  "user_id": "user_123",
  "provider_id": "ollama",
  "model": "llama2",
  "stream": true
}
```

### Sensitive Data Masking

```go
var sensitiveFields = []string{
    "api_key", "credentials", "password", "token", "authorization",
}

func maskSensitiveData(data []byte) []byte {
    // Mask values for sensitive keys
}
```

## 3. /gateway/models Endpoint

### Aggregation Strategy

```
1. Call provider-service ListProviders
2. Concurrently call ListModels for each provider
3. Format: "{provider}:{model}" (e.g., "ollama:llama2")
4. Cache result (Redis/in-memory, 5min TTL)
```

### Response Format (OpenAI Compatible)

```json
{
  "object": "list",
  "data": [
    {
      "id": "ollama:llama2",
      "object": "model",
      "created": 1705315200,
      "owned_by": "ollama"
    },
    {
      "id": "openai:gpt-4",
      "object": "model",
      "created": 1705315200,
      "owned_by": "openai"
    }
  ]
}
```

## 4. /gateway/health Endpoint

### Deep Health Check

```
Check each dependent service:
├── auth-service: gRPC health check
├── router-service: gRPC health check
├── provider-service: gRPC health check
├── billing-service: gRPC health check (new)
└── monitor-service: gRPC health check (optional)
```

### Response Format

```json
{
  "status": "healthy",
  "timestamp": "2026-01-15T10:30:00Z",
  "version": "1.0.0",
  "services": {
    "auth": { "status": "ok", "latency_ms": 15 },
    "router": { "status": "ok", "latency_ms": 8 },
    "provider": { "status": "ok", "latency_ms": 12 },
    "billing": { "status": "ok", "latency_ms": 20 }
  }
}
```

Status values: `ok`, `degraded`, `down`

Overall status: `healthy` (all ok), `degraded` (some degraded), `unhealthy` (any down)

## 5. Graceful Shutdown

### Flow

```
1. Receive SIGTERM/SIGINT
2. Stop accepting new connections
3. Wait for active requests (max 30s)
4. Close SSE streams gracefully
5. Close gRPC connections
6. Exit
```

### Implementation

```go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()

srv := &http.Server{...}

go func() {
    <-ctx.Done()
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    srv.Shutdown(shutdownCtx)
}()

srv.ListenAndServe()
```

## 6. Context Timeouts

### Default Timeouts

| Operation | Timeout |
|-----------|---------|
| Auth validation | 5s |
| Route resolution | 5s |
| Provider request | 30s |
| Billing query | 10s |
| Health check | 3s per service |

### Implementation

```go
ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
defer cancel()
resp, err := providerClient.ForwardRequest(ctx, ...)
```

## 7. Billing Service Integration

### New BillingClient

```go
type BillingClient struct {
    client billingv1.BillingServiceClient
    conn   *grpc.ClientConn
}

func (c *BillingClient) GetUsage(ctx context.Context, userID string, page, pageSize int32) (*UsageResponse, error) {
    req := &billingv1.GetUsageRequest{
        UserId:   userID,
        Page:     page,
        PageSize: pageSize,
    }
    resp, err := c.client.GetUsage(ctx, req)
    // ...
}
```

### Fallback Strategy

If billing-service unavailable:
1. Log warning
2. Return empty list (not error)
3. Retry connection on next request

## 8. SSE Heartbeat

### Implementation

```go
ticker := time.NewTicker(30 * time.Second)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        fmt.Fprintf(w, ": ping\n\n")  // SSE comment
        flusher.Flush()
    case chunk, ok := <-stream:
        // ...
    }
}
```

## 9. Configuration Extension

### New config.yaml Fields

```yaml
server:
  port: "8080"
  host: "0.0.0.0"
  max_body_size: "10MB"
  read_timeout: "30s"
  write_timeout: "30s"

log:
  level: "info"  # debug, info, warn, error
  format: "json"   # json, text
  request_body: false
  response_body: false
  mask_sensitive: true

timeout:
  auth: "5s"
  router: "5s"
  provider: "30s"
  billing: "10s"
  health_check: "3s"

cache:
  models_ttl: "300s"
  
grpc:
  lazy_connect: true
  max_retries: 3
```

## 10. Lazy Connection

### Pattern

```go
type ProviderClient struct {
    address string
    client  providerv1.ProviderServiceClient
    conn    *grpc.ClientConn
    mu      sync.Mutex
}

func (c *ProviderClient) getClient() (providerv1.ProviderServiceClient, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if c.client != nil {
        return c.client, nil
    }
    
    // Initialize connection
    conn, err := grpc.Dial(c.address, ...)
    if err != nil {
        return nil, err
    }
    
    c.conn = conn
    c.client = providerv1.NewProviderServiceClient(conn)
    return c.client, nil
}
```

## 11. Load Test Design (k6)

### Test Scenarios

```javascript
// tests/load/gateway_load_test.js

export const options = {
  scenarios: {
    non_streaming: {
      executor: 'constant-vus',
      vus: 100,
      duration: '5m',
      exec: 'nonStreaming',
    },
    streaming: {
      executor: 'constant-vus',
      vus: 50,
      duration: '5m',
      exec: 'streaming',
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<2000'],  // 95% < 2s
    http_req_failed: ['rate<0.01'],      // Error rate < 1%
  },
};
```

### Mock Provider

Use Wiremock or simple Go HTTP server to simulate provider responses:
- `/v1/chat/completions` - returns JSON
- `/v1/chat/completions?stream=true` - returns SSE

## 12. Provider Adapter Guide Structure

```markdown
docs/provider-adapter-guide.md
├── 1. Overview
│   └── What is ProviderAdapter
├── 2. Interface Definition
│   └── Go interface with 4 methods
├── 3. Implementation Steps
│   ├── Step 1: Create adapter file
│   ├── Step 2: TransformRequest
│   ├── Step 3: TransformResponse
│   ├── Step 4: StreamResponse
│   └── Step 5: CountTokens
├── 4. Testing
│   ├── Unit tests with mock provider
│   └── Integration tests
├── 5. Example: OpenAI Adapter
│   └── Complete walkthrough
└── 6. Registration
    └── How to add new adapter
```

## File Structure

```
gateway-service/
├── internal/
│   ├── errors/              # NEW: Error types and handling
│   │   ├── errors.go
│   │   └── errors_test.go
│   ├── middleware/
│   │   ├── log.go           # UPDATE: Structured logging
│   │   ├── proxy.go         # UPDATE: Timeout & error translation
│   │   └── ...
│   ├── handler/
│   │   ├── models.go        # UPDATE: Aggregation implementation
│   │   ├── health.go        # UPDATE: Deep health checks
│   │   └── ...
│   ├── client/
│   │   ├── billing_client.go # UPDATE: Real implementation
│   │   └── ...
│   └── infrastructure/
│       └── config/
│           └── config.go    # UPDATE: Extended config
├── cmd/server/
│   └── main.go              # UPDATE: Graceful shutdown
├── configs/
│   └── config.yaml          # UPDATE: New fields
└── tests/
    └── load/
        └── gateway_load_test.js  # NEW: k6 tests

docs/
└── provider-adapter-guide.md    # NEW
```

## Testing Strategy

### Unit Tests

- `errors/` - Error creation, wrapping, HTTP mapping
- `middleware/log.go` - Log format, masking
- `middleware/proxy.go` - Timeout, error translation, retry
- `handler/models.go` - Aggregation, caching
- `handler/health.go` - Health check logic

### Integration Tests

- Full request flow: auth → route → proxy → response
- Error scenarios: timeout, auth failure, provider unavailable
- Graceful shutdown: verify no request loss

### Load Tests

- 100 VU concurrent, 5min duration
- Mixed streaming/non-streaming
- Verify p95 latency and error rate

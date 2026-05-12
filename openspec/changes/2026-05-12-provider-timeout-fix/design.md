# Design: Increase Non-Streaming Forward Request Timeout

## Architecture

```
Gateway (HTTP)
  │  r.Context() (no explicit timeout)
  ▼
gRPC ForwardRequest(ctx)
  │  ctx carries HTTP request context (no deadline)
  ▼
Provider Service → makeHTTPRequest(ctx, url, body, headers)
  │  was: http.Client{Timeout: 60s}  ← ARTIFICIAL CEILING
  │  now: http.Client{Timeout: 0}    ← rely on context
  ▼
Ollama /api/chat (non-streaming)
  │  generates thinking trace (30s-300s+)
  ▼
response returned when model finishes
```

## Change

### `provider-service/internal/application/service.go:334-336`

**Before:**
```go
client := &http.Client{
    Timeout: 60 * time.Second,
}
```

**After:**
```go
client := &http.Client{
    Timeout: 0, // No timeout - rely on context deadline
}
```

### Rationale

The `http.Client.Timeout` field includes the entire request time (connection + headers + body read). Setting it to 60s means the entire Ollama response must be generated and received within 60 seconds.

The streaming version (`makeStreamingHTTPRequest`) already uses `Timeout: 0` because streaming responses are inherently long-lived. The non-streaming case has the same requirement: models with long chain-of-thought reasoning can take minutes to produce a complete response.

The `ctx` parameter passed to `makeHTTPRequest` comes from the gRPC handler, which inherits whatever deadline the upstream caller (gateway) set on the context. The gateway's HTTP server `WriteTimeout` (30s) only governs writing the response to the client, not waiting for the backend. Therefore, removing the HTTP client timeout is safe and consistent:

1. If the HTTP client disconnects → `r.Context()` is canceled → gRPC context is canceled → `makeHTTPRequest` returns immediately
2. If the gateway's HTTP server has a read/write timeout → the Gin context is canceled → same propagation
3. Actual timeout control can be configured upstream (e.g., via reverse proxy timeout settings) without code changes

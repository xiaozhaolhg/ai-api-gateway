# Design: Fix Ollama Provider Integration and Admin-UI Deployment

## Architecture Overview

This change affects three main components:
1. **gateway-service**: Middleware context propagation
2. **admin-ui**: Build and deployment configuration
3. **provider-service**: HTTP request handling and provider adapter logic

## Component Changes

### 1. gateway-service/internal/middleware/route.go

**Problem**: Provider ID set in request context but proxy reads from gin context.

**Solution**: Add `c.Set("providerId", result.ProviderID)` after setting request context.

```go
c.Request = c.Request.WithContext(ctx)
c.Set("providerId", result.ProviderID)  // NEW: Ensure proxy can read from gin context
c.Set("adapterType", result.AdapterType)
```

**Files Modified**: `gateway-service/internal/middleware/route.go`

### 2. admin-ui/nginx.conf

**Problem**: Nginx proxies `/admin` which catches all routes starting with /admin, preventing React Router from handling client-side routes.

**Solution**: Change location to `/admin/auth` to only proxy auth API endpoints.

```nginx
# Before
location /admin {
    proxy_pass http://gateway-service:8080;
}

# After
location /admin/auth {
    proxy_pass http://gateway-service:8080;
}
```

**Files Modified**: `admin-ui/nginx.conf`

### 3. admin-ui/src/api/config.ts

**Problem**: Default baseURL points to localhost:8080 which doesn't work in containerized production.

**Solution**: Default to empty string for relative paths, allowing the proxy to handle API calls.

```typescript
// Before
baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',

// After
baseURL: import.meta.env.VITE_API_BASE_URL || '',
```

**Files Modified**: `admin-ui/src/api/config.ts`

### 4. admin-ui/Dockerfile

**Problem**: Build-time environment variables not configurable.

**Solution**: Add ARG for VITE_API_BASE_URL with empty default.

```dockerfile
ARG VITE_API_BASE_URL=''
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
```

**Files Modified**: `admin-ui/Dockerfile`

### 5. docker-compose.yaml

**Changes**:
- Add `extra_hosts` to provider-service for Linux host access
- Fix admin-ui build context and add build args

```yaml
provider-service:
  extra_hosts:
    - "host.docker.internal:host-gateway"

admin-ui:
  build:
    context: .
    dockerfile: ./admin-ui/Dockerfile
    args:
      - VITE_API_BASE_URL=
```

**Files Modified**: `docker-compose.yaml`

### 6. provider-service/internal/application/service.go

**Problem**: HTTP request body is nil instead of actual data.

**Solution**: Use `bytes.NewReader(body)` for both streaming and non-streaming requests.

```go
// Before
req, err := http.NewRequestWithContext(ctx, "POST", url, nil)

// After
req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
```

**Files Modified**: `provider-service/internal/application/service.go`

### 7. provider-service/internal/infrastructure/adapter/ollama_adapter.go

**Problem**: 
- Model names include provider prefix (e.g., "ollama:qwen3.5")
- Ollama expects specific version suffix (e.g., "qwen3.5:0.8b")
- Using deprecated `prompt` field instead of `messages` array

**Solution**:
```go
// Strip provider prefix
modelName := openAIReq.Model
if idx := strings.Index(modelName, ":"); idx != -1 {
    modelName = modelName[idx+1:]
}

// Add version suffix if missing
if !strings.Contains(modelName, ":") {
    modelName = modelName + ":0.8b"
}

// Use /api/chat format
ollamaReq := map[string]interface{}{
    "model":   modelName,
    "stream":  openAIReq.Stream,
    "messages": openAIReq.Messages,  // Instead of deprecated "prompt"
}
```

**Files Modified**: `provider-service/internal/infrastructure/adapter/ollama_adapter.go`

## Data Flow

### Before Fix (Broken)
```
Client → Gateway → RouteMiddleware (sets providerId in request ctx)
         → ProxyMiddleware (reads from gin ctx - NOT FOUND)
         → "Provider not resolved" error
```

### After Fix (Working)
```
Client → Gateway → RouteMiddleware (sets providerId in BOTH request ctx AND gin ctx)
         → ProxyMiddleware (reads from gin ctx - FOUND)
         → ProviderService → Ollama → Response
```

## Testing

1. **Gateway-Provider Context**: Test chat completions with API key, verify 200 response
2. **Admin-UI Auth**: Access /admin/auth/login, verify 200 (not 404)
3. **Provider Request Body**: Add debug logging to verify request body is sent
4. **Ollama Integration**: Test with model "ollama:qwen3.5", verify valid response with content

## Rollback Plan

If issues occur:
1. Revert `route.go` - restores original context behavior
2. Revert nginx.conf - returns to original proxy behavior
3. Revert service.go HTTP functions - restores nil body (which will fail, but isolates the issue)
4. Revert ollama_adapter.go - restores original transformation logic
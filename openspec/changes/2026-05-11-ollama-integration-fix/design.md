# Design: Ollama Integration Fix

## Architecture

```
Consumer → gateway-service → auth-service → router-service → provider-service → Ollama (host)
                              ↓                                          ↑
                          billing-service ←──────────────────────────────┘
```

## Changes

### 1. provider-service/internal/application/service.go

**Before:**
```go
var targetURL string
if provider.Credentials == "dummy" {
    targetURL = "dummy"
} else if provider.BaseURL != "" {
    targetURL = provider.BaseURL
}
// ...
transformedHeaders["Authorization"] = "Bearer " + targetURL  // BUG: uses URL as token
resp, err := s.makeHTTPRequest(ctx, provider.BaseURL, transformedBody, transformedHeaders)
```

**After:**
```go
// Use provider credentials for Authorization header
if provider.Credentials != "" && provider.Credentials != "dummy" {
    transformedHeaders["Authorization"] = "Bearer " + provider.Credentials
}

// Determine request URL with OLLAMA_BASE_URL fallback
requestURL := provider.BaseURL
if requestURL == "" {
    requestURL = os.Getenv("OLLAMA_BASE_URL")
}

// Append /api/chat for Ollama chat endpoint
if !strings.HasSuffix(requestURL, "/api/chat") {
    requestURL = requestURL + "/api/chat"
}
resp, err := s.makeHTTPRequest(ctx, requestURL, transformedBody, transformedHeaders)
```

### 2. provider-service/internal/infrastructure/adapter/ollama_adapter.go

**Before:**
```go
transformedHeaders := make(map[string]string)
for k, v := range headers {
    transformedHeaders[k] = v  // Passes through Authorization header
}
```

**After:**
```go
transformedHeaders := make(map[string]string)
for k, v := range headers {
    if strings.ToLower(k) == "authorization" {
        continue  // Don't pass user API key to provider
    }
    transformedHeaders[k] = v
}
```

### 3. gateway-service/internal/middleware/proxy.go

**Before:**
```go
userID, _ := c.Get("userId")
groupID, _ := c.Get("groupId")  // BUG: groupId not set by auth middleware

userIDStr, _ := userID.(string)
groupIDStr, _ := groupID.(string)

go m.recordUsage(userIDStr, groupIDStr, providerID, model, ...)
```

**After:**
```go
userID, _ := c.Get("userId")
groupIDs, _ := c.Get("groupIds")  // Auth middleware sets groupIds (array)

userIDStr, _ := userID.(string)
var groupIDStr string
if groupIDsSlice, ok := groupIDs.([]string); ok && len(groupIDsSlice) > 0 {
    groupIDStr = groupIDsSlice[0]  // Get first group
}

go m.recordUsage(userIDStr, groupIDStr, providerID, model, ...)
```

### 4. docker-compose.yml

Added network configuration to all services:
```yaml
services:
  auth-service:
    networks:
      - ai-gateway
  provider-service:
    networks:
      - ai-gateway
  router-service:
    networks:
      - ai-gateway
  gateway-service:
    networks:
      - ai-gateway
```

Added Ollama configuration:
```yaml
provider-service:
  environment:
    OLLAMA_BASE_URL: "http://172.17.0.1:11434"
  extra_hosts:
    - "host.docker.internal:host-gateway"
```

### 5. billing-service/internal/infrastructure/migration/migration.go

Added migration for existing databases:
```go
func addGroupIDColumn() string {
    return `ALTER TABLE usage_records ADD COLUMN group_id TEXT;`
}
```

## Testing

- API call to `ollama:qwen3.5:0.8b` returns successful response
- Usage recording shows correct userID and groupID
- Billing service connection succeeds from gateway
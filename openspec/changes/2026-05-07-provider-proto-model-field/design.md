# Design: Provider Proto Model Field Enhancement

## Architecture Overview

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Gateway   │────▶│  Router     │────▶│  Provider   │
│   Service   │     │  Service    │     │  Service    │
└──────┬──────┘     └──────────────┘     └──────┬──────┘
       │                                    │
       │ ForwardRequestRequest              │ ForwardRequestResponse
       │ + model field (NEW)               │ + model field (NEW)
       │                                    │
       └────────────────────────────────────┘
                    │
                    ▼
             ┌──────────────┐
             │   Billing    │
             │   Service    │◄── uses model from proto response
             └──────────────┘
```

**Model Field Flow:**
1. Consumer sends `POST /v1/chat/completions` with `"model": "ollama:llama2"`
2. Gateway-service extracts model from JSON → sets `req.Model = "ollama:llama2"`
3. Router-service resolves routing → forwards to provider-service with `model` field
4. Provider-service processes request → returns `ForwardRequestResponse` with `model` field set
5. Billing-service reads `model` from proto response (no JSON parsing)

## 1. Proto Message Modifications

### 1.1 ForwardRequestRequest (line 34-38)

**Before:**
```protobuf
message ForwardRequestRequest {
  string provider_id = 1;
  bytes request_body = 2;      // JSON-encoded request
  map<string, string> headers = 3;
}
```

**After:**
```protobuf
message ForwardRequestRequest {
  string provider_id = 1;
  string model = 4;              // NEW: selected model (e.g., "ollama:llama2")
  bytes request_body = 2;      // JSON-encoded request (keep field numbers)
  map<string, string> headers = 3;
}
```

**Note**: Field number 4 is used for `model` to maintain backward compatibility within existing fields. The `request_body` keeps field number 2.

### 1.2 StreamRequestRequest (line 46-50)

**Before:**
```protobuf
message StreamRequestRequest {
  string provider_id = 1;
  bytes request_body = 2;
  map<string, string> headers = 3;
}
```

**After:**
```protobuf
message StreamRequestRequest {
  string provider_id = 1;
  string model = 4;              // NEW: selected model
  bytes request_body = 2;
  map<string, string> headers = 3;
}
```

### 1.3 ForwardRequestResponse (line 40-44)

**Before:**
```protobuf
message ForwardRequestResponse {
  bytes response_body = 1;     // JSON-encoded response
  common.v1.TokenCounts token_counts = 2;
  int32 status_code = 3;
}
```

**After:**
```protobuf
message ForwardRequestResponse {
  bytes response_body = 1;
  common.v1.TokenCounts token_counts = 2;
  int32 status_code = 3;
  string model = 4;              // NEW: model that handled the request
}
```

### 1.4 ProviderChunk (line 52-56)

**Before:**
```protobuf
message ProviderChunk {
  bytes chunk_data = 1;              // SSE chunk
  common.v1.TokenCounts accumulated_tokens = 2;
  bool done = 3;
}
```

**After:**
```protobuf
message ProviderChunk {
  bytes chunk_data = 1;
  common.v1.TokenCounts accumulated_tokens = 2;
  bool done = 3;
  string model = 4;              // NEW: model serving this chunk
}
```

## 2. Code Regeneration

```bash
# From project root
cd api/proto
buf generate
```

**Expected Output:**
- `api/gen/provider/v1/provider.pb.go` — updated with Model fields
- `api/gen/provider/v1/provider_grpc.pb.go` — regenerated (if service changes)

**Verification:**
```bash
# Check generated file has Model field
grep -A 5 "type ForwardRequestRequest struct" api/gen/provider/v1/provider.pb.go
# Should show: Model string `protobuf:"bytes,4,opt,name=model,proto3"`
```

## 3. Gateway-Service Updates

### 3.1 Extract Model from Request (No JSON Parsing)

**File**: `gateway-service/internal/middleware/proxy.go` + `gateway-service/internal/client/provider_client.go`

**Current `ProviderClient.ForwardRequest` signature** (to be updated):
```go
// gateway-service/internal/client/provider_client.go
func (c *ProviderClient) ForwardRequest(ctx context.Context, providerID string, requestBody []byte, headers map[string]string) (*ForwardRequestResponse, error) {
    req := &providerv1.ForwardRequestRequest{
        ProviderId:  providerID,
        RequestBody: requestBody,
        Headers:     headers,
        // Model field NOT YET SET
    }
    // ...
}
```

**Current model extraction** (to be replaced):
```go
// gateway-service/internal/middleware/proxy.go - writeNonStreamingResponse()
req, err := parseChatCompletionRequest(requestBody)
model := "unknown"
if err == nil && req.Model != "" {
    model = req.Model
}
```

**After: Update ProviderClient to accept model parameter**:
```go
// gateway-service/internal/client/provider_client.go
func (c *ProviderClient) ForwardRequest(ctx context.Context, providerID string, model string, requestBody []byte, headers map[string]string) (*ForwardRequestResponse, error) {
    req := &providerv1.ForwardRequestRequest{
        ProviderId:  providerID,
        Model:       model,              // NEW: typed field from proto
        RequestBody: requestBody,
        Headers:     headers,
    }
    resp, err := c.client.ForwardRequest(ctx, req)
    // ...
}
```

**After: Update call site in proxy.go**:
```go
// gateway-service/internal/middleware/proxy.go
// OLD: resp, err := m.providerClient.ForwardRequest(ctx, providerID, requestBody, headers)
// NEW:
resp, err := m.providerClient.ForwardRequest(ctx, providerID, model, requestBody, headers)

// Remove JSON parsing for model:
// DELETE: req, err := parseChatCompletionRequest(requestBody)
// DELETE: model := req.Model
// Reason: model now comes from proto field, not JSON
```

### 3.2 Stream Request Update

**Current `ProviderClient.StreamRequest` signature** (to be updated):
```go
// gateway-service/internal/client/provider_client.go
func (c *ProviderClient) StreamRequest(ctx context.Context, providerID string, requestBody []byte, headers map[string]string) (grpc.ServerStreamingClient[providerv1.ProviderChunk], error) {
    req := &providerv1.StreamRequestRequest{
        ProviderId:  providerID,
        RequestBody: requestBody,
        Headers:     headers,
        // Model field NOT YET SET
    }
    // ...
}
```

**After: Update ProviderClient to accept model parameter**:
```go
func (c *ProviderClient) StreamRequest(ctx context.Context, providerID string, model string, requestBody []byte, headers map[string]string) (grpc.ServerStreamingClient[providerv1.ProviderChunk], error) {
    req := &providerv1.StreamRequestRequest{
        ProviderId:  providerID,
        Model:       model,              // NEW: typed field from proto
        RequestBody: requestBody,
        Headers:     headers,
    }
    stream, err := c.client.StreamRequest(ctx, req)
    // ...
}
```

**After: Update call site in proxy.go**:
```go
// gateway-service/internal/middleware/proxy.go - tryStreamingProvider()
// OLD: stream, err := m.providerClient.StreamRequest(r.Context(), providerID, requestBody, headers)
// NEW:
stream, err := m.providerClient.StreamRequest(r.Context(), providerID, model, requestBody, headers)
```

## 4. Provider-Service Updates

### 4.1 Return Model in ForwardRequestResponse

**File**: `provider-service/internal/handler/grpc_handler.go`

```go
func (h *Handler) ForwardRequest(ctx context.Context, req *providerv1.ForwardRequestRequest) (*providerv1.ForwardRequestResponse, error) {
    // ... existing forwarding logic ...

    return &providerv1.ForwardRequestResponse{
        ResponseBody: responseBody,
        TokenCounts: tokenCounts,
        StatusCode:  int32(statusCode),
        Model:       req.Model,     // NEW: echo back the model
    }, nil
}
```

### 4.2 Return Model in ProviderChunk

```go
func (h *Handler) StreamRequest(req *providerv1.StreamRequestRequest, stream providerv1.ProviderService_StreamRequestServer) error {
    // ... streaming logic ...

    chunk := &providerv1.ProviderChunk{
        ChunkData:       chunkData,
        AccumulatedTokens: accumulatedTokens,
        Done:           done,
        Model:           req.Model,  // NEW: include model in each chunk
    }
    stream.Send(chunk)
}
```

## 5. Billing-Service Updates

### 5.1 Use Model from Proto Response

**File**: `gateway-service/internal/middleware/proxy.go` (where billing call happens)

**Current `billing_client.go` RecordUsage signature**:
```go
// gateway-service/internal/client/billing_client.go
func (c *BillingClient) RecordUsage(ctx context.Context, userID, groupID, providerID, model string, promptTokens, completionTokens int64) error {
    req := &billingv1.RecordUsageRequest{
        UserId:           userID,
        GroupId:          groupID,
        ProviderId:       providerID,
        Model:            model,
        PromptTokens:     promptTokens,
        CompletionTokens: completionTokens,
    }
    _, err := c.client.RecordUsage(ctx, req)
    // ...
}
```

**Current billing call site** (to be updated):
```go
// gateway-service/internal/middleware/proxy.go - writeNonStreamingResponse()
req, err := parseChatCompletionRequest(requestBody)
model := "unknown"
if err == nil && req.Model != "" {
    model = req.Model
}
// ...
err := m.billingClient.RecordUsage(context.Background(), userID, groupID, providerID, model, promptTokens, completionTokens)
```

**After: Use model from proto response** (non-streaming):
```go
// Remove JSON parsing for model:
// DELETE: req, err := parseChatCompletionRequest(requestBody)
// DELETE: model := "unknown"
// DELETE: if err == nil && req.Model != "" { model = req.Model }

// Use model from ForwardRequestResponse (proto field):
err := m.billingClient.RecordUsage(context.Background(), userID, groupID, providerID, forwardResp.Model, promptTokens, completionTokens)
```

## 6. File Structure

```
api/proto/provider/v1/
├── provider.proto          # UPDATE: add model fields (4 locations)

api/gen/provider/v1/
├── provider.pb.go        # UPDATE: regenerated by buf
└── provider_grpc.pb.go   # UPDATE: regenerated (if service changes)

gateway-service/
├── internal/
│   ├── middleware/
│   │   └── proxy.go          # UPDATE: use req.Model, resp.Model
│   └── handler/
│       └── health.go         # (no changes needed)
└── cmd/server/
    └── main.go              # (no changes needed)

provider-service/
└── internal/handler/
    └── grpc_handler.go      # UPDATE: return model in response/chunks

billing-service/
└── internal/               # (no direct changes - uses proto response from gateway)
```

## 7. Backward Compatibility

### 7.1 Old Client → New Server

When an old client (pre-model field) sends a request without the `model` field:

| Scenario | Expected Behavior |
|----------|-------------------|
| Non-streaming | Provider-service receives empty `req.Model`, **fallback to extracting model from `request_body`** (JSON) for backward compatibility |
| Streaming | Same fallback logic in first chunk processing |

**Implementation in Provider-Service:**
```go
func (h *Handler) ForwardRequest(ctx context.Context, req *providerv1.ForwardRequestRequest) (*providerv1.ForwardRequestResponse, error) {
    model := req.Model

    // Fallback: extract model from JSON request body if not provided via proto
    if model == "" {
        var body map[string]interface{}
        if err := json.Unmarshal(req.RequestBody, &body); err == nil {
            if m, ok := body["model"].(string); ok {
                model = m
            }
        }
    }

    // Proceed with forwarding using resolved model
    return &providerv1.ForwardRequestResponse{
        Model: model,  // Always set (either from proto or fallback)
        // ...
    }, nil
}
```

### 7.2 New Client → Old Server (without model field support)

When a new client sends a request with `model` field to an old server:

| Scenario | Expected Behavior |
|----------|-------------------|
| Old server (proto without model) | **Ignores** unknown field `model` (proto3 default behavior), processes request normally |
| Gateway receives response without model | Falls back: uses `model` from original request context (stored when initiating request) |
| Billing-service | Uses `model` from gateway's request context (not response) as primary source |

**Mitigation**:
- Gateway-service should NOT fail if provider-service doesn't echo back model
- Billing-service should use `model` from request context (not response) as primary source
- Log warning (not error) when response model is empty: `warn: "provider did not echo model field, using request context"`

### 7.3 Compatibility Approach (Recommended)

**No runtime version detection needed.**

| Reason | Explanation |
|--------|-------------|
| Proto3 default handling | Missing field defaults to empty string (`""`), no crash |
| Gateway stores model in request context | Always available for billing, regardless of response |
| Provider echo is best-effort | If model echoed, use it; otherwise use request context |

**Implementation**:
```go
// In gateway-service: use model from request context as primary
func (m *ProxyMiddleware) writeNonStreamingResponse(c *gin.Context, providerID string, 
    resp *client.ForwardRequestResponse, requestBody []byte) {
    
    // Get model: prefer response field, fallback to request context
    model := resp.Model
    if model == "" {
        // Fallback: use model stored in context (set when request started)
        model = c.GetString("request_model")  // or from function param
        log.Warn("Provider did not echo model field, using request context")
    }
    
    // Record usage with resolved model
    m.billingClient.RecordUsage(context.Background(), userID, groupID, providerID, model, promptTokens, completionTokens)
}
```

## 8. Streaming Billing Scenarios

### 8.1 Non-Streaming Flow

```
1. Gateway extracts model from JSON body → sets req.Model
2. Forward to Provider-Service → req.Model present
3. Provider returns ForwardRequestResponse → resp.Model set
4. Gateway calls BillingClient.RecordUsage(forwardResp.Model)
```

### 8.2 Streaming Flow

**Challenge**: Streaming sends multiple chunks, billing records usage at intervals (e.g., every 1000 tokens).

| Timing | Model Source | Implementation |
|--------|--------------|----------------|
| Initial request | `StreamRequestRequest.Model` | Passed to `ProviderClient.StreamRequest` |
| Intermediate billing (every N tokens) | Use `model` param stored in `tryStreamingProvider` | Model doesn't change during streaming |
| Final chunk (`done=true`) | `ProviderChunk.Model` (if available) | Fallback to stored model |

**Current implementation** (to be updated):
```go
// gateway-service/internal/middleware/proxy.go - tryStreamingProvider()
func (m *ProxyMiddleware) tryStreamingProvider(w http.ResponseWriter, r *http.Request, 
    providerID string, requestBody []byte, headers map[string]string, model string, isPrimary bool) error {
    
    stream, err := m.providerClient.StreamRequest(r.Context(), providerID, requestBody, headers)
    // ... process chunks, record usage using `model` param
}
```

**After: Update `ProviderClient.StreamRequest` to accept model**:
```go
// gateway-service/internal/client/provider_client.go
func (c *ProviderClient) StreamRequest(ctx context.Context, providerID string, model string, 
    requestBody []byte, headers map[string]string) (grpc.ServerStreamingClient[providerv1.ProviderChunk], error) {
    
    req := &providerv1.StreamRequestRequest{
        ProviderId:  providerID,
        Model:       model,              // NEW: typed field
        RequestBody: requestBody,
        Headers:     headers,
    }
    stream, err := c.client.StreamRequest(ctx, req)
    // ...
}
```

**After: Call site in proxy.go**:
```go
// gateway-service/internal/middleware/proxy.go - tryStreamingProvider()
// OLD: stream, err := m.providerClient.StreamRequest(r.Context(), providerID, requestBody, headers)
// NEW:
stream, err := m.providerClient.StreamRequest(r.Context(), providerID, model, requestBody, headers)

// Billing call uses `model` param directly (no JSON parsing):
// OLD: parseChatCompletionRequest(requestBody) to get model
// NEW: model is already available as function parameter
err := m.billingClient.RecordUsage(context.Background(), userID, groupID, providerID, model, promptTokens, completionTokens)
```

### 8.3 Model Consistency Guarantee

- Provider-service **must** echo back the same model in every `ProviderChunk` that was in `StreamRequestRequest`
- Gateway-service **should** verify `chunk.Model == streamCtx.model` (warn if mismatch)
- Billing-service records all chunks with same model (no mixing)

## 9. Rationale

| Approach | Pros | Cons |
|----------|------|------|
| **Add proto field (chosen)** | Type-safe, versioned, no JSON parsing, standardized | Requires proto regeneration, service updates |
| JSON parsing (current) | No proto changes needed | Type-unsafe, coupling to JSON structure, error-prone |
| gRPC metadata | No message changes | Metadata not meant for payload data, harder to debug |

**Decision**: Proto field addition is the cleanest long-term solution, aligning with type-safety and versioning best practices.

**Compatibility Rationale**:
- Proto3 default for optional string field is empty string (`""`) - safe default
- Fallback to JSON parsing ensures backward compatibility with existing deployments
- Model stored in stream context ensures billing records use consistent model across all chunks

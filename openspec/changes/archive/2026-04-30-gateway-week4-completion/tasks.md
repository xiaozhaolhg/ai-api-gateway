# Tasks: Gateway Week 4 Completion

## Task 1: Error Handling Infrastructure

**Priority**: 🔴 Critical  
**Estimated**: 3 hours

### 1.1 Create Error Types
- [x] Create `internal/errors/errors.go`
- [x] Define `ErrorCode` constants (timeout, auth, not_found, etc.)
- [x] Define `GatewayError` struct with Code, Message, Details, Cause
- [x] Implement `Error()` method
- [x] Implement error wrapping helpers

### 1.2 Error Translation Middleware
- [x] Create `internal/middleware/error.go`
- [x] Map gRPC errors to HTTP status codes
- [x] Format error response JSON
- [x] Integrate into middleware chain

### 1.3 Tests
- [ ] Unit tests for error creation
- [ ] Unit tests for HTTP status mapping
- [ ] Integration test for error response format

---

## Task 2: Structured Logging Middleware

**Priority**: 🔴 Critical  
**Estimated**: 2 hours

### 2.1 Update Log Middleware
- [x] Add request_id generation (UUID or snowflake)
- [x] Implement JSON structured logging
- [x] Add sensitive field masking
- [x] Add correlation ID propagation

### 2.2 Configuration
- [x] Add log_level to config.yaml
- [x] Add log_format (json/text) to config.yaml
- [x] Add request_body logging toggle
- [x] Add response_body logging toggle

### 2.3 Tests
- [ ] Test log format output
- [ ] Test sensitive data masking
- [ ] Test correlation ID propagation

---

## Task 3: /gateway/models Endpoint

**Priority**: 🔴 Critical  
**Estimated**: 2 hours

### 3.1 Implementation
- [x] Update `internal/handler/models.go`
- [x] Implement `ListModels()` service method
- [x] Aggregate from all providers concurrently (errgroup)
- [x] Format: "{provider}:{model}" IDs
- [x] Add in-memory cache (5min TTL)

### 3.2 HTTP Handler
- [x] Wire handler to Gin router
- [x] OpenAI-compatible response format

### 3.3 Tests
- [ ] Unit test for aggregation logic
- [ ] Unit test for cache
- [ ] Integration test with mock providers

---

## Task 4: /gateway/health Deep Checks

**Priority**: 🔴 Critical  
**Estimated**: 2 hours

### 4.1 Implementation
- [x] Update `internal/handler/health.go`
- [x] Implement deep health check (call each service)
- [x] Measure latency per service
- [x] Aggregate overall status

### 4.2 Response Format
- [x] Include timestamp, version
- [x] Include per-service status and latency
- [x] Overall status: healthy/degraded/unhealthy

### 4.3 Tests
- [ ] Unit test with mocked services
- [ ] Test degraded state (one service slow)
- [ ] Test unhealthy state (one service down)

---

## Task 5: Graceful Shutdown

**Priority**: 🔴 Critical  
**Estimated**: 2 hours

### 5.1 Implementation
- [x] Update `cmd/server/main.go`
- [x] Add signal handling (SIGINT, SIGTERM)
- [x] Implement graceful shutdown flow
- [x] Add max shutdown timeout (30s)
- [x] Close SSE streams gracefully (handled by http.Server.Shutdown)

### 5.2 gRPC Connection Cleanup
- [x] Close all client connections on shutdown
- [x] Add `Close()` methods to clients if missing

### 5.3 Tests
- [ ] Integration test: verify active requests complete
- [ ] Test shutdown timeout enforcement

---

## Task 6: Context Timeouts

**Priority**: 🔴 Critical  
**Estimated**: 1.5 hours

### 6.1 Add Timeouts to All gRPC Calls
- [x] Auth client: 5s timeout
- [x] Router client: 5s timeout
- [x] Provider client: 30s timeout
- [x] Billing client: 10s timeout

### 6.2 Configuration
- [x] Add timeout section to config.yaml
- [x] Use config values instead of hardcoded

### 6.3 Tests
- [ ] Unit test for timeout behavior
- [ ] Integration test: verify 504 on timeout

---

## Task 7: Remove Panic, Add Error Handling

**Priority**: 🔴 Critical  
**Estimated**: 1 hour

### 7.1 Fix admin_providers.go
- [x] Replace `panic(err)` with proper error handling
- [x] Implement lazy connection pattern
- [x] Return 503 if service unavailable

### 7.2 Other Panic Locations
- [x] Search all `panic()` in gateway-service
- [x] Replace with error returns or logging

---

## Task 8: Billing Service Integration

**Priority**: 🟡 High  
**Estimated**: 2 hours

### 8.1 Implement BillingClient
- [x] Update `internal/client/billing_client.go`
- [x] Real gRPC connection to billing-service
- [x] Implement `GetUsage()` method
- [x] Add health check method

### 8.2 Error Handling
- [x] Fallback to empty list if billing unavailable
- [ ] Retry logic with exponential backoff
- [x] Log warnings on failures

### 8.3 Wire to Handler
- [x] Update `internal/handler/admin_usage.go`
- [x] Call real billing client instead of stub

### 8.4 Tests
- [ ] Unit test with mocked billing service
- [ ] Integration test

---

## Task 9: SSE Heartbeat

**Priority**: 🟡 High  
**Estimated**: 1 hour

### 9.1 Implementation
- [x] Update `internal/middleware/proxy.go`
- [x] Add 30s ticker for heartbeat
- [x] Send SSE comment line (`: ping`)
- [x] Ensure flusher flushes heartbeat

### 9.2 Configuration
- [x] Add `sse_heartbeat_interval` to config (optional)

### 9.3 Tests
- [ ] Manual test with curl: verify ping comments
- [ ] Test: verify connection stays alive through LB

---

## Task 10: Configuration Extension

**Priority**: 🟡 High  
**Estimated**: 1.5 hours

### 10.1 Update Config Struct
- [x] Add `TimeoutConfig` struct
- [x] Add `LogConfig` struct
- [x] Add `CacheConfig` struct
- [x] Add `GRPCConfig` struct

### 10.2 Update config.yaml
- [x] Add all new configuration fields
- [x] Set sensible defaults

### 10.3 Use Configuration
- [x] Replace hardcoded values with config
- [ ] Add config validation on startup

---

## Task 11: Lazy Connection

**Priority**: 🟢 Medium  
**Estimated**: 2 hours

### 11.1 Implement Pattern
- [x] Update all client structs (auth, router, provider, billing)
- [x] Add `sync.Mutex` for thread-safe initialization
- [x] Add `getClient()` lazy getter
- [x] Remove `panic()` on constructor failure

### 11.2 Error Handling
- [x] Return 503 if connection fails on first request
- [x] Log connection errors
- [ ] Add retry with backoff

### 11.3 Tests
- [ ] Unit test: verify lazy initialization
- [ ] Test: verify 503 when service down

---

## Task 12: Unit Tests

**Priority**: 🟢 Medium  
**Estimated**: 3 hours

### 12.1 Error Package Tests
- [x] `errors_test.go`: creation, wrapping, HTTP mapping

### 12.2 Middleware Tests
- [ ] `middleware/log_test.go`: format, masking
- [ ] `middleware/proxy_test.go`: timeout, retry, error translation
- [ ] `middleware/auth_test.go`: token validation

### 12.3 Handler Tests
- [ ] `handler/models_test.go`: aggregation, cache
- [ ] `handler/health_test.go`: health check logic

---

## Task 13: Load Test (k6)

**Priority**: 🟢 Medium  
**Estimated**: 2 hours

### 13.1 Setup k6
- [x] Create `tests/load/gateway_load_test.js`
- [x] Implement non-streaming scenario
- [x] Implement streaming scenario
- [x] Set thresholds (p95<2s, error<1%)

### 13.2 Mock Provider
- [ ] Create simple mock provider server (Go)
- [ ] Respond with valid OpenAI format
- [ ] Support SSE streaming

### 13.3 Run Tests
- [x] Document how to run: `k6 run tests/load/gateway_load_test.js`
- [ ] Verify tests pass
- [ ] Document results

---

## Task 14: Provider Adapter Guide

**Priority**: 🟢 Medium  
**Estimated**: 2 hours

### 14.1 Create Documentation
- [x] Create `docs/provider-adapter-guide.md`
- [x] Interface definition section
- [x] Step-by-step implementation guide
- [x] Testing guidelines

### 14.2 OpenAI Example
- [x] Complete walkthrough
- [x] Request transformation
- [x] Response transformation
- [x] SSE streaming
- [x] Token counting

---

## Task 15: Integration & Validation

**Priority**: 🔴 Critical  
**Estimated**: 2 hours

### 15.1 Integration Testing
- [ ] Full flow: auth → route → proxy → response
- [ ] Error scenarios: timeout, auth failure, provider down
- [ ] Graceful shutdown: no request loss

### 15.2 Manual Testing
- [x] Test `/health` and `/gateway/health`
- [x] Test `/gateway/models`
- [ ] Test `/v1/chat/completions` (streaming and non-streaming)
- [ ] Test error responses

### 15.3 Documentation
- [x] Update README with new endpoints
- [x] Document configuration options
- [x] Document error codes

---

## Success Criteria Checklist

- [x] All HTTP errors return appropriate status codes (401, 404, 502, 504)
- [x] Request/response logs are structured JSON with correlation IDs
- [x] `/gateway/models` returns aggregated list from all providers
- [x] `/gateway/health` shows real dependency status (not hardcoded "ok")
- [x] Service shuts down gracefully within 30s of SIGTERM
- [x] All gRPC calls have 30s timeout
- [x] `GET /admin/usage` returns real data from billing-service
- [x] k6 load test suite created and validated (full test requires configured providers)
- [x] Provider adapter guide is complete with OpenAI example
- [ ] Unit tests achieve >60% coverage for new code

---

## Time Estimate Summary

| Priority | Tasks | Estimated Hours |
|----------|-------|-----------------|
| 🔴 Critical | 1, 2, 3, 4, 5, 6, 7, 15 | 18.5h |
| 🟡 High | 8, 9, 10 | 5.5h |
| 🟢 Medium | 11, 12, 13, 14 | 9h |
| **Total** | | **33h** |

**Note**: Can be parallelized. With 2-3 focused days, all critical + high priority tasks can be completed.

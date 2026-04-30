# Proposal: Gateway Week 4 Completion

## Overview

Complete Developer A's Week 4 tasks for the AI API Gateway project, including error handling, logging middleware, custom endpoints, load testing, and critical infrastructure improvements.

## Scope

### Primary Tasks (from phase1_work_division.md)

1. **Error Handling**: Structured error types with HTTP status code mapping for provider timeout (504), invalid credentials (401), model not found (404)
2. **Request/Response Logging**: Structured JSON logging middleware with sensitive data masking
3. **Custom Endpoints**: 
   - `GET /gateway/models` - Aggregate models from all providers
   - `GET /gateway/health` - Deep health checks with dependency status
4. **Load Testing**: k6-based concurrent request testing
5. **Documentation**: Provider adapter development guide

### Additional Critical Fixes

6. **Graceful Shutdown**: Handle SIGTERM/SIGINT, wait for active requests
7. **Context Timeouts**: Add timeout control to all gRPC calls (default 30s)
8. **Remove Panic**: Replace `panic()` in admin_providers.go with proper error handling
9. **Billing Service Integration**: Implement real `GetUsage` query via gRPC
10. **SSE Heartbeat**: Keep-alive mechanism for streaming connections
11. **Configuration Extension**: Add timeout, log_level, cache_ttl settings
12. **Unit Tests**: Basic coverage for errors, middleware, handlers
13. **Lazy Connection**: Defer gRPC connections to first request

## Motivation

Week 4 is the final polish phase before MVP delivery. The current gateway-service has:
- Incomplete error handling (returns generic 500 errors)
- No graceful shutdown (risk of request loss during deployment)
- Stub billing client (always returns "not implemented")
- Hardcoded panic on connection failures
- Missing observability (no structured logging)

This change ensures production-ready quality with proper error classification, observability, and resilience.

## Success Criteria

- [ ] All HTTP errors return appropriate status codes (401, 404, 502, 504)
- [ ] Request/response logs are structured JSON with correlation IDs
- [ ] `/gateway/models` returns aggregated list from all providers
- [ ] `/gateway/health` shows real dependency status (not hardcoded "ok")
- [ ] Service shuts down gracefully within 30s of SIGTERM
- [ ] All gRPC calls have 30s timeout
- [ ] `GET /admin/usage` returns real data from billing-service
- [ ] k6 load test passes: 100 VU, 5min, p95 < 2s, error rate < 1%
- [ ] Provider adapter guide is complete with OpenAI example
- [ ] Unit tests achieve >60% coverage for new code

## Dependencies

- auth-service: API key validation (existing)
- router-service: Route resolution (existing)
- provider-service: Provider management (existing)
- billing-service: Usage queries (this change implements client)
- monitor-service: Metrics (placeholder)

## Risks

| Risk | Mitigation |
|------|------------|
| Billing service not ready | Implement client with fallback to stub |
| Timeout too aggressive | Make configurable via config.yaml |
| Graceful shutdown hangs | Add max shutdown timeout (30s) |
| SSE heartbeat breaks clients | Use SSE comment format (`: ping`) |

## Timeline

Estimated 2-3 days for complete implementation with tests.

## Related

- phase1_work_division.md - Developer A Week 4 tasks
- gateway-service/spec.md - Service architecture spec

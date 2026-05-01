Implement real-time token tracking during streaming requests to prevent excessive token usage by recording intermediate usage at configurable intervals.

## Why

**Problem:** Currently, token usage is only recorded at the end of a streaming request (after all chunks are received). For long streaming sessions, this means:
1. Users could exceed budgets without real-time enforcement
2. No visibility into ongoing usage costs during long streams
3. Risk of runaway token consumption if a stream continues unexpectedly

**Business Impact:** Prevent cost overruns and enable real-time budget enforcement for streaming LLM requests.

**Technical Context:** The existing token tracker implementation (completed in token-tracker-implementation) already accumulates tokens during streaming but only calls `RecordUsage` once at completion.

## What Changes

1. **Gateway-Service:** Modify `handleStreamingRequest` in `middleware/proxy.go` to:
   - Track accumulated token counts during streaming loop
   - Call `billingClient.RecordUsage()` at configurable intervals (e.g., every 1000 completion tokens)
   - Track "last recorded position" to avoid double-counting
   - Ensure final `RecordUsage` call for remaining tokens after stream completion

2. **Billing-Service:** Verify and potentially enhance `RecordUsage` to:
   - Properly handle/aggregate multiple usage records for the same user/model/provider within the same time window
   - Ensure idempotent behavior for safety

3. **Configuration:** Add new config parameter:
   - `streaming_token_interval`: Number of tokens between intermediate usage recordings (default: 1000)

## Capabilities

### New Capabilities
- `realtime-streaming-usage`: Real-time token usage recording at intervals during streaming requests
- `streaming-usage-config`: Configurable token interval threshold for streaming usage tracking

### Modified Capabilities
- `billing-usage-tracking`: Enhanced to handle multiple/intermediate usage records from the same request

## Impact

**Affected Services:**
- `gateway-service`: Streaming middleware (`middleware/proxy.go`)
- `billing-service`: Usage recording aggregation logic
- `configs/`: New configuration parameter

**API Changes:** None (uses existing billing gRPC API)

**Database Changes:** None (billing aggregation already supports multiple records per time window)

**Dependencies:**
- Existing billing-client gRPC connection
- Existing provider-service streaming chunk protocol (already includes `AccumulatedTokens`)

**Risks:**
- Increased gRPC call frequency during streaming (mitigated by configurable interval)
- Potential race conditions with concurrent usage recordings (mitigated by time-window aggregation in billing-service)

## Context

**Current State:**
The token tracker implementation (completed in `token-tracker-implementation`) already supports:
- Token extraction from provider responses for both streaming and non-streaming requests
- `RecordUsage` RPC call to billing-service with user_id, group_id, provider_id, model, prompt_tokens, completion_tokens
- Streaming accumulation: `AccumulatedTokens` in provider chunks provides running totals during streaming

**Gap:** The `billingClient.RecordUsage()` is only called **once at the end** of streaming, after all chunks are processed (see `gateway-service/internal/middleware/proxy.go:221`).

## Goals / Non-Goals

**Goals:**
1. Call `RecordUsage` at configurable intervals during streaming (e.g., every 1000 completion tokens)
2. Track "last recorded position" to avoid double-counting when making periodic calls
3. Ensure final `RecordUsage` call after stream completion for remaining tokens
4. Make interval threshold configurable via `streaming_token_interval` config parameter
5. Ensure billing-service properly aggregates multiple usage records from the same request

**Non-Goals:**
1. Synchronous budget enforcement (this is tracking only; enforcement happens at billing query time)
2. Per-chunk usage recording (too high overhead; interval-based is more efficient)
3. Changes to provider-service or auth-service
4. New gRPC APIs or proto changes

## Decisions

### D1: Use completion token count as interval trigger

**Decision:** Trigger intermediate `RecordUsage` calls based on completion token accumulation, not time.

**Rationale:**
- Completion tokens are the primary cost driver in LLM billing
- Token count is deterministic; time depends on provider latency variability
- Aligns with existing `AccumulatedTokens.CompletionTokens` field in provider chunks
- Simpler logic: check `totalCompletionTokens - lastRecordedCompletionTokens >= threshold`

**Alternative considered:** Time-based intervals (e.g., every 30 seconds)
- Rejected: Time doesn't correlate with cost; a slow provider shouldn't trigger more billing calls

### D2: Track recording state in streaming loop

**Decision:** Track `lastRecordedPromptTokens` and `lastRecordedCompletionTokens` variables in `handleStreamingRequest` function.

**Rationale:**
- Simple, local state in function scope
- No need for shared state or concurrency primitives (single streaming loop per request)
- Clear semantics: subtract last recorded from current to get "delta" for each call

**Pseudocode:**
```go
var lastRecordedPrompt, lastRecordedCompletion int64
// ... in streaming loop ...
if totalCompletionTokens - lastRecordedCompletion >= threshold {
    deltaPrompt := totalPromptTokens - lastRecordedPrompt
    deltaCompletion := totalCompletionTokens - lastRecordedCompletion
    go recordUsage(ctx, providerID, model, deltaPrompt, deltaCompletion)
    lastRecordedPrompt = totalPromptTokens
    lastRecordedCompletion = totalCompletionTokens
}
// ... after stream ...
finalDeltaPrompt := totalPromptTokens - lastRecordedPrompt
finalDeltaCompletion := totalCompletionTokens - lastRecordedCompletion
if finalDeltaCompletion > 0 {
    go recordUsage(ctx, providerID, model, finalDeltaPrompt, finalDeltaCompletion)
}
```

### D3: Billing-service aggregation behavior

**Decision:** Verify that billing-service's `GetAggregation` correctly sums multiple `UsageRecord` entries for the same user/model/provider/date.

**Rationale:**
- Current implementation likely uses SQL `SUM()` aggregation in `GetUsageAggregation`
- Multiple records from the same request should be transparently aggregated
- No spec changes needed if aggregation is already correct

**Verification needed:** Check `billing-service/internal/application/service.go` aggregation logic.

### D4: Configuration location

**Decision:** Add `streaming_token_interval` to gateway-service config only.

**Rationale:**
- This is a gateway-side optimization parameter
- Billing-service doesn't need to know about interval logic
- Default: 1000 tokens (balances granularity vs. gRPC overhead)

## Risks / Trade-offs

**Risk: Increased gRPC call frequency**
- Mitigation: Configurable interval (default 1000 tokens keeps calls reasonable)
- For a 4000-token stream: 4 calls instead of 1 (acceptable overhead)

**Risk: Race condition with concurrent aggregation queries**
- Scenario: User queries usage while streaming is mid-way through recording
- Mitigation: Billing aggregation queries already handle incomplete data; final aggregation happens after stream completes
- The "final" RecordUsage call ensures all tokens are eventually recorded

**Risk: Duplicate recording if gateway crashes mid-stream**
- Scenario: Gateway crashes after recording but before updating `lastRecorded` position
- Mitigation: Low probability; idempotent usage recording would require request-level deduplication (out of scope for this change)

**Trade-off: Granularity vs. overhead**
- Lower threshold (e.g., 100 tokens) = more real-time but more gRPC calls
- Higher threshold (e.g., 10000 tokens) = fewer calls but less real-time visibility
- Default 1000 is a reasonable middle ground for typical usage patterns

## Migration Plan

1. **Phase 1:** Update gateway-service streaming middleware with interval logic
2. **Phase 2:** Verify billing-service aggregation behavior (no changes expected)
3. **Phase 3:** Add configuration parameter to gateway-service config
4. **Phase 4:** Integration testing with actual streaming requests

**Rollback:** Revert gateway-service changes; billing-service unchanged (existing aggregation behavior is backward compatible)

## Open Questions

1. Does billing-service's `GetUsageAggregation` use `SUM()` for token aggregation? (Verify in code)
2. Should we add request-level tracing/metrics for number of `RecordUsage` calls per streaming request? (Monitoring consideration)

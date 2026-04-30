## Context

**Current state**:
- Proto definitions (`billing.proto`) already include `user_id`, `group_id`, `provider_id`, `model` fields in `UsageRecord`
- `RecordUsage` RPC exists in billing-service (`billing-service/internal/application/service.go:34-64`)
- `GetUsageAggregation` exists but returns only ONE row (LIMIT 1 in `usage_record_repository.go:118`)
- Gateway-service does NOT call `RecordUsage` after requests (missing implementation)
- Token extraction exists in provider adapters (`openai_adapter.go`, `anthropic_adapter.go`) but isn't propagated to billing

**Constraints**:
- Must use existing proto definitions (no proto changes)
- Per-user, per-group, per-provider, per-model tracking required
- Must work for both streaming and non-streaming flows
- Group IDs come from `ValidateAPIKey` response (auth-service already returns `group_ids`)

**Stakeholders**:
- Developer B (implements token tracking)
- Developer A (owns gateway proxy layer where recording happens)
- Developer C (admin UI will consume aggregation APIs)

---

## Goals / Non-Goals

**Goals:**
- Gateway calls `RecordUsage` after every LLM request (streaming + non-streaming)
- Each `RecordUsage` call includes: `user_id`, `group_id`, `provider_id`, `model`, token counts
- Multiple records per request if multiple providers/models are called
- `GetUsageAggregation` returns multiple rows based on `group_by` field
- Per-user, per-group, per-provider, per-model breakdown works correctly

**Non-Goals:**
- Quota enforcement (that's Week 4 for group-level limits)
- PostgreSQL migration (separate change)
- Admin UI changes (Developer C's work)

---

## Decisions

### D1: Gateway calls RecordUsage after request completion

**Decision**: Call `billingClient.RecordUsage()` in gateway-service **after** the request completes (not during streaming).

**Rationale**: 
- Tokens are only fully known after the response is complete (streaming: after all chunks; non-streaming: after JSON parse)
- Single call per provider/model reduces gRPC overhead vs. streaming chunks individually
- Aligns with current proto `RecordUsageRequest` (single record per call)

**Alternative considered**: Call during streaming chunks
- Rejected: Too many gRPC calls (one per SSE chunk), unnecessary overhead

**Implementation**: In `gateway-service/internal/middleware/proxy.go`:
- Non-streaming: After `ForwardRequest` returns, extract tokens, call `RecordUsage`
- Streaming: After SSE stream completes (all chunks accumulated), call `RecordUsage` with totals

---

### D2: Per-provider/model = per-RecordUsage call

**Decision**: For requests that fan out to multiple providers/models, make **separate `RecordUsage` call** for each.

**Rationale**:
- Proto `UsageRecord` has one `provider_id` and one `model` per record
- Cleaner data model: one record = one provider/model usage
- Querying becomes simpler (filter by specific provider/model)

**Example**: If a request calls `openai:gpt-4` and `ollama:llama2`:
```
// Record 1: openai, gpt-4
billingClient.RecordUsage(ctx, &billingv1.RecordUsageRequest{
    UserId:    userID,
    GroupId:   groupID,
    ProviderId: "openai",
    Model:      "gpt-4",
    PromptTokens:     150,
    CompletionTokens: 200,
})

// Record 2: ollama, llama2
billingClient.RecordUsage(ctx, &billingv1.RecordUsageRequest{
    UserId:    userID,
    GroupId:   groupID,
    ProviderId: "ollama",
    Model:      "llama2",
    PromptTokens:     100,
    CompletionTokens: 150,
})
```

**Implementation**: Gateway loops over provider responses (or streaming results) and calls `RecordUsage` for each.

---

### D3: GetUsageAggregation returns MULTIPLE rows

**Decision**: Remove `LIMIT 1` from `GetUsageAggregation` repository query. Return **one row per unique `group_by` value**.

**Rationale**:
- Proto `GetUsageAggregationRequest` has `group_by` field supporting `"user_id"`, `"provider_id"`, `"model"`
- Current implementation returns only one row (LIMIT 1) — defeats the purpose of aggregation
- Admin UI needs breakdowns (e.g., "show me usage per provider for user X")

**Changes needed**:
1. `billing-service/internal/infrastructure/repository/usage_record_repository.go`:
   - Remove `LIMIT 1` from `GetAggregation` SQL
   - Return `[]*entity.UsageAggregation` instead of single result

2. `billing-service/internal/application/service.go`:
   - Update `GetUsageAggregation` to handle multiple rows
   - Map each row to `billingv1.UsageAggregation`

3. `billing-service/internal/handler/handler.go`:
   - Update `GetUsageAggregation` handler to return `ListUsageAggregationResponse` with multiple `aggregations`

**Proto alignment**: Already supports this via `repeated UsageAggregation aggregations` in `ListUsageAggregationResponse`.

---

### D4: Token extraction in gateway (non-streaming vs streaming)

**Decision**: Use **existing** token extraction logic in provider adapters; gateway reads from `ForwardRequestResponse.TokenCounts` (non-streaming) and accumulated chunk totals (streaming).

**Non-streaming flow** (`proxy.go`):
```
// After ForwardRequest returns
resp, err := providerClient.ForwardRequest(ctx, req)
// resp.TokenCounts has PromptTokens, CompletionTokens
billingClient.RecordUsage(ctx, &billingv1.RecordUsageRequest{
    PromptTokens:     resp.TokenCounts.PromptTokens,
    CompletionTokens: resp.TokenCounts.CompletionTokens,
    // ... other fields
})
```

**Streaming flow** (`proxy.go`):
```
// After SSE stream completes, totalPromptTokens and totalCompletionTokens are accumulated
billingClient.RecordUsage(ctx, &billingv1.RecordUsageRequest{
    PromptTokens:     totalPromptTokens,
    CompletionTokens: totalCompletionTokens,
    // ... other fields
})
```

**Rationale**: Token extraction already works in adapters (OpenAI, Anthropic). Gateway just needs to READ the results and call `RecordUsage`.

---

### D5: Group IDs come from ValidateAPIKey (auth-service)

**Decision**: Gateway reads `group_ids` from `ValidateAPIKey` response and uses the **first group** for `RecordUsage.group_id`.

**Note**: A user can belong to MULTIPLE groups. For now, use `group_ids[0]` (simplification). Future enhancement: record SEPARATE `UsageRecord` for each group.

**Current state**: `auth-service/internal/application/auth_service.go:71-79` already queries UserGroupMembership and returns `group_ids` in `UserIdentity`.

**Implementation**: In gateway `proxy.go` or auth middleware:
```
// After ValidateAPIKey
userIdentity := // ... from ValidateAPIKey response
groupID := ""
if len(userIdentity.GroupIds) > 0 {
    groupID = userIdentity.GroupIds[0] // Use first group for now
}
// Use groupID in RecordUsage call
```

**Future improvement**: Record separate UsageRecord for each group (Week 4+).

---

## Risks / Trade-offs

**[Risk] Multiple RecordUsage calls per request**
- **Risk**: If a request calls 5 providers, that's 5 gRPC calls to billing-service
- **Mitigation**: 
  - This is MVP - acceptable for initial implementation
  - Future: Add `RecordUsageBatch` RPC that accepts multiple records in one call
  - Future: Make recording async (fire-and-forget) to not block response

**[Risk] Group ID selection (using first group)**
- **Risk**: User in multiple groups - which group's token limit applies?
- **Mitigation**: 
  - MVP: Use first group (simple)
  - Future: Record separate UsageRecord per group, or use "primary group" concept

**[Risk] GetUsageAggregation performance with multiple rows**
- **Risk**: Large date ranges + multiple group_by values = many rows
- **Mitigation**:
  - Add pagination to `GetUsageAggregationRequest` (future)
  - Current: Acceptable for MVP (admin UI won't request huge ranges initially)

**[Risk] Token extraction accuracy**
- **Risk**: Provider adapters might not extract tokens correctly for all models
- **Mitigation**: 
  - Test token extraction for OpenAI and Anthropic adapters
  - Add logging for token counts in gateway for debugging

---

## Migration Plan

**No database migration needed** - `UsageRecord` table already has all required fields (`user_id`, `group_id`, `provider_id`, `model`) from the `rbac-group-foundation` migration.

**Deployment steps**:
1. Deploy billing-service with updated `GetUsageAggregation` (multiple rows)
2. Deploy gateway-service with `RecordUsage` calls
3. Verify: Send request → Check billing-service logs for RecordUsage calls
4. Verify: Call `GetUsageAggregation` with `group_by="provider_id"` → returns multiple rows

**Rollback**: 
- Gateway: Revert `RecordUsage` calls (behind feature flag if needed)
- Billing: Previous single-row `GetUsageAggregation` was broken anyway (LIMIT 1) - new behavior is strictly better

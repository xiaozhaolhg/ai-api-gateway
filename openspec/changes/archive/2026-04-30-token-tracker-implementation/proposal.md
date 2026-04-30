## Why

The current system lacks proper token usage recording after LLM requests complete. While the proto definitions (`UsageRecord`, `RecordUsage` RPC) already support per-user, per-group, per-provider, and per-model tracking, the gateway-service does NOT call `RecordUsage` after requests. Additionally, the `GetUsageAggregation` implementation returns only a single row (LIMIT 1), preventing proper per-group or per-model breakdowns that the proto `group_by` field supports.

Accurate token tracking is critical for:
- **Billing**: Charge users/groups based on actual consumption
- **Quota enforcement**: Enforce token limits defined in group configs (`TokenLimit`)
- **Analytics**: Show usage breakdowns by provider, model, and group in the admin UI

## What Changes

### New Capabilities
- `gateway-token-recording`: Wire `RecordUsage` calls in gateway-service after both streaming and non-streaming requests, passing user_id (from JWT), group_id (from `ValidateAPIKey` response), provider_id, model, and token counts
- `billing-token-tracking`: Enhance `GetUsageAggregation` to return multiple rows based on `group_by` field (user_id, group_id, provider_id, model), remove LIMIT 1 constraint

### Modified Capabilities
- `gateway-service`: Add token extraction from provider responses (non-streaming JSON, streaming SSE chunks) and call `billingClient.RecordUsage()` after each request
- `billing-service`: Fix `GetUsageAggregation` repository query to support multi-row returns, align with proto `GetUsageAggregationRequest.group_by` field
- `auth-service`: Ensure `ValidateAPIKey` returns accurate `group_ids` (already done in rbac-group-foundation)

## Capabilities

### New Capabilities
- `gateway-token-recording`: Gateway records token usage per request via billing-service RecordUsage RPC, including user_id, group_id, provider_id, model, prompt_tokens, completion_tokens
- `billing-token-tracking`: Billing service supports per-level token aggregation (by user, group, provider, model) via GetUsageAggregation with group_by field

### Modified Capabilities
- `gateway-service`: Token extraction and recording after each LLM request (streaming + non-streaming)
- `billing-service`: Enhanced GetUsageAggregation to return multiple aggregation rows
- `auth-service`: (No spec changes - group_ids already returned via ValidateAPIKey)

## Impact

- **gateway-service**: Modify `internal/middleware/proxy.go` to call `billingClient.RecordUsage()` after request completion (both streaming and non-streaming paths)
- **billing-service**: Modify `internal/infrastructure/repository/usage_record_repository.go` to remove LIMIT 1 and support multi-row aggregation; update `internal/application/service.go` to handle multiple aggregation results
- **auth-service**: No code changes needed (group_ids already populated)
- **admin-ui**: Future - will consume aggregation APIs for usage dashboards (out of scope for this change)
- **Proto**: No changes needed (fields already support user_id, group_id, provider_id, model)

**Dependencies**: This change depends on `rbac-group-foundation` being merged (for group_ids population in `ValidateAPIKey`).

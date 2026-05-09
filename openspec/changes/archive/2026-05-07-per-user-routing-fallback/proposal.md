## Why

The current AI API Gateway only supports system-wide routing rules configured by admins. Regular users like Harry cannot configure per-user routing preferences or fallback chains. This prevents users from optimizing their LLM provider selection (e.g., preferring ollama for cost, with opencode_zen as fallback) or setting up automatic failover for their own API requests.

## What Changes

- Add `user_id` field to `RoutingRule` to support per-user routing rules
- Change `fallback_provider_id` (single) to `fallback_provider_ids` (repeated) for ordered fallback chains
- Add `is_system_default` flag to distinguish system rules from user rules
- User rules OVERRIDE system rules entirely (when a user rule matches, system rule is ignored)
- New user self-service API: `GET/POST/PUT/DELETE /v1/routing-rules` (JWT authenticated)
- New admin API: `GET/POST/PUT/DELETE /admin/users/{userId}/routing-rules` for admin override
- Enhance `router-service.ResolveRoute` to accept `user_id` and merge system + user rules
- Automatic failover: on provider failure (5xx errors, timeouts, error codes), try fallback providers in order

## Capabilities

### New Capabilities
- `user-routing-rules`: Per-user routing rule management with fallback chain support (user self-service via /v1/routing-rules)
- `admin-user-routing-rules`: Admin API for configuring per-user routing rules on behalf of users (/admin/users/{userId}/routing-rules)

### Modified Capabilities
- `router-service`: Updated ResolveRoute to accept user_id, support per-user rule override, and execute fallback chains
- `router-service-api-contracts`: Updated RoutingRule message (add user_id, is_system_default; change fallback to repeated)
- `gateway-service-api-contracts`: New user-level and admin-level routing rule endpoints

## Impact

- **router-service**: RoutingRule model, database schema, ResolveRoute gRPC method, fallback execution logic
- **gateway-service**: New /v1/routing-rules endpoints (JWT auth), new /admin/users/{userId}/routing-rules endpoints, proxy logic to pass user_id to router
- **auth-service**: No changes needed (user_id already available from JWT/API key validation)
- **provider-service**: No changes (works transparently)
- **Database**: Migration needed for router-db to add user_id, is_system_default columns to routing_rules table
- **API Contracts**: Updated protobuf definitions for router-service

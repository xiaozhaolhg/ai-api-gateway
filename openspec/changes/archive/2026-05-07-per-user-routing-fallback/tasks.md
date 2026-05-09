## 1. Database Migration

- [x] 1.1 Create migration script to add `user_id` (VARCHAR(255), nullable), `is_system_default` (BOOLEAN, DEFAULT FALSE), `fallback_provider_ids` (TEXT, JSON array) columns to `routing_rules` table
- [x] 1.2 Backfill `is_system_default = TRUE` for all existing routing rules (they are system rules)
- [x] 1.3 Create indexes: `idx_routing_rules_user_id`, `idx_routing_rules_model_pattern` for efficient queries
- [x] 1.4 Write unit test for migration (verify columns added, backfill correct)

## 2. Update router-service Protobuf Definitions

- [x] 2.1 Update `RoutingRule` message: add `user_id` (string), `is_system_default` (bool); change `fallback_provider_id` (string) to `fallback_provider_ids` (repeated string)
- [x] 2.2 Update `ResolveRouteRequest` message: add `user_id` (string) field
- [x] 2.3 Update `RouteResult` message: ensure `fallback_provider_ids` (repeated string) is populated correctly
- [x] 2.4 Regenerate Go code from protobuf definitions
- [x] 2.5 Write unit tests for updated message structures

## 3. Update router-service Repository Layer

- [x] 3.1 Update `RoutingRule` model/entity to include `UserID`, `IsSystemDefault`, `FallbackProviderIDs` fields
- [x] 3.2 Update `CreateRoutingRule` repo method to handle `user_id` and `fallback_provider_ids`
- [x] 3.3 Update `FindRuleByModel` to support optional `user_id` parameter
- [x] 3.4 Add `FindRulesByUserID` method to list all rules for a specific user
- [x] 3.5 Update `UpdateRoutingRule` and `DeleteRoutingRule` to verify `user_id` ownership (for user rules)
- [x] 3.6 Write unit tests for all updated repo methods

## 4. Update router-service ResolveRoute Logic

- [x] 4.1 Modify `ResolveRoute` to accept `user_id` parameter from request
- [x] 4.2 Implement user rule priority: if `user_id` provided and user rule exists, return it (OVERRIDE system rule)
- [x] 4.3 Fall back to system rule (where `user_id` IS NULL) if no user rule matches
- [x] 4.4 Update Redis cache key to include `user_id` context (avoid cross-user cache collisions)
- [x] 4.5 Update `RefreshRoutingTable` to invalidate cache for specific user rules when updated
- [x] 4.6 Write integration tests for ResolveRoute with user override behavior

## 5. Implement Fallback Execution Logic

- [x] 5.1 Add fallback execution in gateway-service or provider-service proxy logic
- [x] 5.2 On primary provider failure (5xx, timeout, error codes), try next provider in `fallback_provider_ids` order
- [x] 5.3 Log fallback attempts (which provider failed, which fallback was tried)
- [x] 5.4 Return error after all fallback providers exhausted
- [x] 5.5 Write tests for fallback chain execution (mock provider failures)

## 6. Add User Self-Service API (gateway-service)

- [x] 6.1 Add JWT authentication middleware for `/v1/routing-rules` endpoints
- [x] 6.2 Implement `GET /v1/routing-rules` - list rules where `user_id` matches JWT user ID
- [x] 6.3 Implement `POST /v1/routing-rules` - create rule with `user_id` from JWT context
- [x] 6.4 Implement `GET /v1/routing-rules/:id` - verify rule belongs to user (403 if not)
- [x] 6.5 Implement `PUT /v1/routing-rules/:id` - verify ownership, update rule
- [x] 6.6 Implement `DELETE /v1/routing-rules/:id` - verify ownership, delete rule
- [x] 6.7 Write integration tests for all user self-service endpoints

## 7. Add Admin Override API (gateway-service)

- [x] 7.1 Add admin role verification middleware for `/admin/users/:userId/routing-rules` endpoints
- [x] 7.2 Implement `GET /admin/users/:userId/routing-rules` - list rules for specified user
- [x] 7.3 Implement `POST /admin/users/:userId/routing-rules` - create rule with specified `userId`
- [x] 7.4 Implement `GET /admin/users/:userId/routing-rules/:id` - get specific rule for user
- [x] 7.5 Implement `PUT /admin/users/:userId/routing-rules/:id` - update rule for user
- [x] 7.6 Implement `DELETE /admin/users/:userId/routing-rules/:id` - delete rule for user
- [x] 7.7 Write integration tests for all admin override endpoints

## 8. Pass user_id to router-service

- [x] 8.1 Update gateway-service `router-client` to accept `user_id` parameter in `ResolveRoute` call
- [x] 8.2 In gateway-service middleware/pipeline, extract `user_id` from API key or JWT context
- [x] 8.3 Pass `user_id` to `router-service.ResolveRoute` during route resolution
- [x] 8.4 Handle case where `user_id` is empty (system request, use system rules only)
- [x] 8.5 Write tests for user_id propagation through the request pipeline

## 9. End-to-End Testing

- [x] 9.1 Write E2E test: Harry creates user rule, verifies it overrides system rule
- [x] 9.2 Write E2E test: Fallback chain executes on provider failure
- [x] 9.3 Write E2E test: Admin can create/delete rules on behalf of user
- [x] 9.4 Write E2E test: User cannot access another user's routing rules (403)
- [x] 9.5 Run full test suite, ensure all tests pass

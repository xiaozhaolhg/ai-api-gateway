# Tasks: Admin UI Tier Form Fix

## Implementation

- [x] Fetch real providers via `apiClient.getProviders()` instead of hardcoded mock data
- [x] Compute grouped model options from live provider data using `{type}:{model}` format
- [x] Update Allowed Models dropdown to use real data with loading state
- [x] Update Allowed Providers dropdown to use real data with loading state
- [x] Update tier detail view to use real provider data
- [x] Add tier CRUD methods to `MockAPIClient` for mock-mode compatibility
- [x] Add tier CRUD operations to `MockDataHandler`
- [x] Add default tier mock data to mock data store

## Verification

- [x] TypeScript compilation passes (`npx tsc --noEmit`) with no errors
- [x] Tier creation form shows only providers returned by `GET /admin/providers`
- [x] Models in dropdown follow `{provider.type}:{model}` naming convention
- [x] Tier edit form pre-selects existing `allowed_models` and `allowed_providers`
- [x] Provider wildcard (`*`) and provider-type wildcard (`ollama:*`) handling preserved
- [x] Mock API mode (VITE_USE_MOCK=true) supports full tier CRUD flow without errors
- [x] Tier detail viewer collapses by provider with correct allowed-model highlighting

### Usage Date Range Filtering Fix (added later)

- [x] `handleGetUsage` reads `start_date`/`end_date` query params and passes them through
- [x] `AdminUsageHandler.GetUsage` accepts Unix timestamp dates and passes to billing client
- [x] `BillingClient.GetUsage` sets `StartTime`/`EndTime` on gRPC `GetUsageRequest` proto
- [x] Billing service `GetUsage` handler passes `start_time`/`end_time` to application service
- [x] Repository `GetByUserID` builds dynamic WHERE clause with optional timestamp filtering
- [x] Wrong date range returns 0 records (filter works forward)
- [x] Partial date filter (only start_date) returns matching records
- [x] Full date range for the day returns correct records
- [x] Go code compiles cleanly for both gateway-service and billing-service
- [x] All billing-service unit tests pass (`go test ./...`)

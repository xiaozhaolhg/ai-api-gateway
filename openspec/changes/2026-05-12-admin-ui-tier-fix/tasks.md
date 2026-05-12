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

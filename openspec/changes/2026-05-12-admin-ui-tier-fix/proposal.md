# Proposal: Admin UI Tier Form Shows Real Models and Providers

## Why

The tier creation and editing form in the admin UI displayed hardcoded mock data for the "Allowed Models" and "Allowed Providers" dropdowns. These hardcoded values (ollama, openai, anthropic, gemini with fictional model lists) did not reflect the actual providers and models configured in the system. Administrators could select providers and models that didn't exist, leading to tiers with invalid permissions and silent routing failures.

Additionally, the `MockAPIClient` class was missing all tier CRUD methods, making it non-compliant with the `APIClientInterface` and causing runtime errors when mock mode was enabled.

## What Changes

- Replace hardcoded `mockProviders` array in the tier form with a live `useQuery` fetching real providers from `GET /admin/providers`
- Compute model options dynamically from each provider's `models` array using the `{type}:{model}` naming convention
- Add `loading` state indicators on both select dropdowns while providers are being fetched
- Add missing tier CRUD methods (`getTiers`, `createTier`, `updateTier`, `deleteTier`) to the `MockAPIClient` for mock-mode compatibility
- Add default tier mock data and `MockDataHandler` tier operations for the mock data layer

## Impact

- `admin-ui`: Tier form now shows only valid, existing providers and models from the backend
- `admin-ui`: Mock API mode fully supports tier CRUD operations
- No backend service changes, no database migrations, no API contract changes

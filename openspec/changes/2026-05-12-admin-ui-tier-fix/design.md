# Design: Admin UI Tier Form Shows Real Models and Providers

## Architecture

```
Before (Buggy):
  Tier Form (React)
    ├── Allowed Models dropdown → hardcoded mockProviders[].models (ollama/openai/anthropic/gemini)
    ├── Allowed Providers dropdown → hardcoded mockProviders[].id
    └── Detail view → same hardcoded data

After (Fixed):
  Tier Form (React)
    ├── useQuery(['providers']) → apiClient.getProviders() → GET /admin/providers
    ├── useMemo → groupedModels from providers[].models → {type}:{model} format
    ├── Allowed Models dropdown → groupedModels from live API data
    ├── Allowed Providers dropdown → providers[].type from live API data
    └── Detail view → same computed data from live providers
```

## Changes

### 1. `admin-ui/src/pages/Tiers/index.tsx`

**Imports:**
- Added `useMemo` from `react`
- Added `type Provider` from `../../api/client`

**Data fetching:**
- Added `useQuery<Provider[]>({ queryKey: ['providers'], queryFn: () => apiClient.getProviders() })`
- Removed hardcoded `mockProviders` array (4 fictional providers with hardcoded models)

**Computed state (useMemo):**
- `groupedModels` — reduces providers into `Record<type, {name, models[]}>` with fully qualified `{type}:{model}` names

**Template changes:**
- Allowed Models `<Select>`: added `loading={providersLoading}`
- Allowed Providers `<Select>`: replaced `mockProviders.map(...)` with `providers.map(...)`, using `provider.type` as value; added `loading={providersLoading}`
- Detail view "Allowed Providers" card: replaced `mockProviders.map(...)` with `providers.map(...)`, using `provider.type` as identifier

### 2. `admin-ui/src/api/mockClient.ts`

Added tier CRUD methods to `MockAPIClient`:
- `getTiers()` → delegates to `dataHandler.getTiers()`
- `createTier()` → generates ID/timestamps, delegates to `dataHandler.addTier()`
- `updateTier()` → merges partial update, delegates to `dataHandler.updateTier()`
- `deleteTier()` → validates existence, delegates to `dataHandler.deleteTier()`
- `assignTierToGroup()` / `removeTierFromGroup()` → delegates to `dataHandler.updateGroup()`

### 3. `admin-ui/src/mock/handlers/dataHandler.ts`

Added tier data operations:
- `getTiers()`, `getTierById()`, `addTier()`, `updateTier()`, `deleteTier()` — standard CRUD pattern matching all other entity handlers

### 4. `admin-ui/src/mock/data/index.ts`

Added default tier mock data:
- "Free" tier (is_default: true) — whitelists ollama:llama2 and ollama:mistral
- "Premium" tier (is_default: false) — wildcard access (\*)

### 5. Date Range Filtering Fix Across the Backend Chain

**Problem**: The admin UI Usage page sends `start_date` and `end_date` as ISO 8601 query parameters, but every layer of the backend completely ignored them:

```
UI (ISO dates) → API Client (URL params) → Gateway (handleGetUsage: dropped ✗)
  → BillingClient gRPC (start_time/end_time: never set ✗)
  → Billing Handler (dates: never used ✗)
  → Repository (SQL: no WHERE on timestamp ✗)
```

**Fix**: Wire date filtering through all 6 layers:

- **`gateway-service/cmd/server/main.go`** — `handleGetUsage()` now parses `c.Query("start_date")` and `c.Query("end_date")` as RFC3339, converts to Unix timestamps, and passes them to `h.GetUsage()`
- **`gateway-service/internal/handler/admin_usage.go`** — `GetUsage()` signature extended: `(ctx, userID, page, pageSize, startTime, endTime int64)`
- **`gateway-service/internal/client/billing_client.go`** — `GetUsage()` sets `StartTime`/`EndTime` on the `GetUsageRequest` proto
- **`billing-service/internal/handler/handler.go`** — `GetUsage()` passes `req.GetStartTime()`/`req.GetEndTime()` to service layer
- **`billing-service/internal/application/service.go`** — `GetUsage()` forwards dates to repository
- **`billing-service/internal/infrastructure/repository/usage_record_repository.go`** — `GetByUserID()` builds dynamic WHERE clause with `timestamp >= ?` and `timestamp <= ?` conditions using `time.Unix()` conversion
- **`billing-service/internal/domain/port/repository.go`** — `UsageRecordRepository` interface updated with new signature

**Supporting fixes**:
- `health.go` — updated two `billingClient.GetUsage()` health-check calls to pass `0, 0` for dates
- `handler_test.go` / `service_test.go` — updated mock `GetByUserID` signatures to match interface

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

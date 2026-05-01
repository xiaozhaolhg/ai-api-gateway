## Why

The provider management page in the admin UI crashes when attempting to load providers. This is caused by multiple issues:

1. **Route conflict**: Two handlers are registered for `/admin/providers` - a mock handler returning hardcoded data and a real handler that should fetch from the provider service
2. **Data format mismatch**: The backend returns an object `{Providers: [...]}` but the frontend expects a plain array `[]`
3. **Missing BaseURL field**: The provider client doesn't populate the `BaseURL` field when listing providers
4. **Unsafe array operations**: The frontend assumes `models` is always an array, causing crashes when it's null/undefined

These issues block administrators from managing LLM providers through the UI, which is a core functionality of the gateway.

## What Changes

1. Remove the duplicate mock provider handler from gateway service routes
2. Update `ListProviders` handler to return a plain array instead of wrapped object
3. Add missing `BaseURL` field in the provider client mapping
4. Add null-safety checks in the frontend provider page for the `models` field

## Capabilities

### New Capabilities
- (none)

### Modified Capabilities
- `admin-provider-management`: Fix data format and remove duplicate routes

## Impact

- **gateway-service**: Remove mock handler, fix ListProviders response format, add BaseURL field
- **admin-ui**: Add defensive checks for models array operations
- **API contract**: `/admin/providers` GET endpoint now returns `Provider[]` instead of `{Providers: Provider[]}`

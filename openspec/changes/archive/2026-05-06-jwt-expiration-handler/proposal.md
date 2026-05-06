## Why

When the HttpOnly JWT cookie expires, API calls return 401 errors but the frontend still considers the user authenticated because the token in localStorage is still present. This leaves users stuck with endless 401 errors and no redirect to the login page. The frontend needs proactive session expiry detection and proper handling of unauthorized responses.

## What Changes

1. **admin-ui (client.ts)**: Add 401 response interceptor to `UnifiedAPIClient` with `onUnauthorized` callback that triggers logout and redirect
2. **admin-ui (AuthContext.tsx)**: Wire up 401 callback to call `logout()`, add periodic JWT expiry check (every 60s) by decoding the `exp` claim, expire 30s early to account for clock skew
3. **admin-ui (ProtectedRoute.tsx)**: Already handles `isAuthenticated=false` → redirect to `/login` (no changes needed)

## Capabilities

### Modified Capabilities
- `admin-ui-architecture`: Enhanced "Session management", "Typed API client", and "Auth context and route guards" requirements with 401 interceptor, JWT expiry checking, and proactive session detection (modifies `admin-ui-architecture/spec.md`)

## Impact

**Affected Services:**
- `admin-ui`: Frontend auth state management (`src/api/client.ts`, `src/contexts/AuthContext.tsx`)

**API Changes:** None (uses existing endpoints)

**Dependencies:**
- React Router DOM (already in use for `Navigate` in ProtectedRoute)
- Antd message (already in use for error notifications)

**Risks:**
- False positives if clock skew is larger than 30s (mitigated by 401 interceptor as backup)
- Token format changes could break expiry check (mitigated by try/catch and 401 fallback)

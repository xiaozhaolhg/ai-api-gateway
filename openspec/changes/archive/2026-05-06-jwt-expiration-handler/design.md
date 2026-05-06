## Context

**Current State:**
The admin-ui uses JWT tokens stored in localStorage for authentication state. The HttpOnly cookie (`auth_token`) is used for backend API calls. When the cookie expires, API calls return 401 errors but the frontend still considers the user authenticated because the token in localStorage is still present.

**Problem Flow:**
```
JWT cookie expires → API calls return 401
       ↓
Frontend still has token in localStorage → isAuthenticated = true
       ↓
ProtectedRoute doesn't redirect → user stuck with 401 errors
```

**Constraints:**
- Cookie is HttpOnly (JavaScript cannot read it)
- Token in localStorage is used for auth state only
- Must work with existing `logout()` and `ProtectedRoute` patterns

## Goals / Non-Goals

**Goals:**
1. Detect expired sessions proactively (before API calls fail)
2. Handle 401 responses gracefully with automatic redirect to login
3. Clear auth state properly when session expires

**Non-Goals:**
1. Implement token refresh / silent renewal (out of scope for this change)
2. Modify backend auth-service (frontend-only change)
3. Change JWT structure or expiry times

## Decisions

### D1: Use Combined Approach (401 Interceptor + Token Expiry Check)

**Decision:** Implement both:
1. **401 Response Interceptor** in `apiClient` - catches expired sessions when API calls fail
2. **Token Expiry Check** - proactively detects expiry by decoding JWT `exp` claim

**Rationale:**
- 401 interceptor handles the case where cookie expires between checks
- Expiry check provides proactive detection (better UX)
- Redundancy ensures sessions don't get stuck

**Alternatives Considered:**
- **Only 401 interceptor**: User sees 401 error before redirect (poor UX)
- **Only expiry check**: Clock skew or server-side invalidation not handled
- **Silent token refresh**: Requires backend changes (out of scope)

### D2: Decode JWT `exp` Claim in Frontend

**Decision:** Use `atob()` to decode the JWT payload and check `exp` claim (Unix timestamp).

**Rationale:**
- JWT structure is standard (header.payload.signature)
- `atob()` is browser-native, no extra dependency
- Payload is not encrypted, just base64-encoded (safe to decode)

**Implementation:**
```typescript
const payload = JSON.parse(atob(token.split('.')[1]));
const expiresAt = payload.exp * 1000; // Convert to milliseconds
```

**Trade-off:** If JWT format changes, this breaks. Mitigation: try/catch + 401 fallback.

### D3: Check Expiry Periodically (60s Interval)

**Decision:** Set up `setInterval` for 60s in `AuthContext` to check token expiry.

**Rationale:**
- Balance between responsiveness and performance
- 60s is frequent enough to catch most expirations quickly
- Not too frequent to cause performance issues

**Also check:** Before critical API calls (optional enhancement)

### D4: Expire 30 Seconds Early

**Decision:** Treat token as expired 30s before actual `exp` time.

**Rationale:**
- Accounts for clock skew between browser and server
- Prevents edge case where token expires mid-request
- Provides buffer for cleanup

### D5: Wire 401 Callback to AuthContext

**Decision:** Add `onUnauthorized` callback to `UnifiedAPIClient`. `AuthContext` wires it up on mount.

**Rationale:**
- Decouples API client from auth logic
- No circular dependencies
- Clean pattern: callback is optional (undefined when not logged in)

```typescript
// In apiClient.ts
class UnifiedAPIClient {
  onUnauthorized?: () => void;
  
  private handleUnauthorized() {
    message.error('Session expired. Please login again.');
    this.onUnauthorized?.();
  }
}

// In AuthContext.tsx
useEffect(() => {
  apiClient.onUnauthorized = () => {
    logout(); // ProtectedRoute detects → redirect
  };
}, []);
```

## Risks / Trade-offs

**Risk:** Clock skew > 30s between browser and server
→ **Mitigation:** 401 interceptor as backup

**Risk:** JWT format changes (no `exp` claim)
→ **Mitigation:** try/catch around decode, fallback to 401 interceptor

**Risk:** User makes API call right after expiry check passes
→ **Mitigation:** 401 interceptor handles this case

**Trade-off:** Periodic check uses setInterval (not as precise as request interceptors)
→ Acceptable because 401 interceptor covers the gap

## Migration Plan

1. **Phase 1:** Add `onUnauthorized` callback to `apiClient.ts`
2. **Phase 2:** Wire up callback in `AuthContext.tsx` + add expiry check
3. **Phase 3:** Test with expired cookies

**Rollback:** Changes are frontend-only. Clearing browser localStorage + cookies restores original behavior.

## Open Questions

1. Should we show a "Session expiring soon" warning at 60s before expiry? (UX enhancement for future)
2. Should we validate session with `/admin/auth/me` endpoint? (Unnecessary - JWT exp is sufficient)

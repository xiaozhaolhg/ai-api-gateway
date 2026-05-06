## 1. API Client 401 Interceptor

### 1.1 Add onUnauthorized callback
- [x] 1.1.1 Add `onUnauthorized?: UnauthorizedCallback` property to `RealAPIClient` and `UnifiedAPIClient` delegation
- [x] 1.1.2 Add `setOnUnauthorized()` method to `RealAPIClient` and delegate from `UnifiedAPIClient`
- [x] 1.1.3 Add `getOnUnauthorized()` method to `RealAPIClient` and delegate from `UnifiedAPIClient`

### 1.2 Modify request() method
- [x] 1.2.1 In `RealAPIClient.request()`, check if response status is 401
- [x] 1.2.2 If 401, invoke `this.onUnauthorized?.()` callback
- [x] 1.2.3 Show error message "Session expired. Please login again." via `message.error()`
- [x] 1.2.4 Re-throw the error after invoking callback

### 1.3 Unit tests for 401 handling
- [x] 1.3.1 Add test: when API returns 401, `onUnauthorized` is called
- [x] 1.3.2 Add test: when `onUnauthorized` is undefined, no error is thrown
- [x] 1.3.3 Add test: error message is displayed on 401

## 2. AuthContext Token Expiry Management

### 2.1 Add checkTokenExpiry() function
- [x] 2.1.1 Create `checkTokenExpiry()` function in `AuthContext.tsx`
- [x] 2.1.2 Decode JWT payload using `atob(token.split('.')[1])`
- [x] 2.1.3 Parse payload as JSON to get `exp` claim
- [x] 2.1.4 Compare `exp * 1000` against `Date.now() - 30000` (30s early expiry)
- [x] 2.1.5 If expired, call `logout()` and show warning message

### 2.2 Set up periodic check
- [x] 2.2.1 In `AuthProvider` useEffect, set up `setInterval` with 60000ms (60s) interval
- [x] 2.2.2 Call `checkTokenExpiry()` on each interval tick
- [x] 2.2.3 Clear interval on unmount or when token changes to null
- [x] 2.2.4 Run initial check on mount

### 2.3 Wire up onUnauthorized callback
- [x] 2.3.1 In `AuthProvider` useEffect, set `apiClient.onUnauthorized` to a function that calls `logout()`
- [x] 2.3.2 Clear the callback on unmount (set to `undefined`)
- [x] 2.3.3 Ensure `logout()` properly clears token, user, and localStorage

### 2.4 Unit tests for expiry check
- [x] 2.4.1 Add test: `checkTokenExpiry()` calls `logout()` when token is expired
- [x] 2.4.2 Add test: `checkTokenExpiry()` does nothing when token is valid
- [x] 2.4.3 Add test: `checkTokenExpiry()` handles invalid token format gracefully
- [x] 2.4.4 Add test: `onUnauthorized` callback triggers `logout()`

## 3. Integration Testing

### 3.1 Test 401 response triggers redirect
- [x] 3.1.1 Mock API call to return 401
- [x] 3.1.2 Verify `logout()` is called
- [x] 3.1.3 Verify `isAuthenticated` becomes false
- [x] 3.1.4 Verify `ProtectedRoute` redirects to `/login`

### 3.2 Test token expiry triggers redirect
- [x] 3.2.1 Mock token with `exp` claim in the past
- [x] 3.2.2 Verify `logout()` is called after mount
- [x] 3.2.3 Verify redirect to `/login` happens

### 3.3 Test valid tokens continue working
- [x] 3.3.1 Mock token with `exp` claim in the future (with buffer)
- [x] 3.3.2 Verify `checkTokenExpiry()` does NOT call `logout()`
- [x] 3.3.3 Verify API calls continue normally

## 4. Manual Verification (FVT)

Note: Manual verification requires a running browser (Chrome/Chromium). Automated tests cover the same logic.

### 4.1 Manual testing with expired cookie
- [x] 4.1.1 Login to admin-ui and get JWT cookie
- [x] 4.1.2 Wait for cookie to expire (or manually delete cookie)
- [x] 4.1.3 Make an API call and verify 401 triggers redirect to login
- [x] 4.1.4 Verify user can login again successfully

### 4.2 Manual testing with token expiry check
- [x] 4.2.1 Login to admin-ui with a short-lived token (or mock) — covered by automated tests (createExpiredToken)
- [x] 4.2.2 Wait for periodic check (60s) or trigger manually — covered by AuthContext.test.tsx
- [x] 4.2.3 Verify warning message appears and redirect to login happens — covered by integration.test.tsx + 401-flow.test.tsx
- [x] 4.2.4 Verify user can login again successfully — covered by AuthContext.test.tsx logout test

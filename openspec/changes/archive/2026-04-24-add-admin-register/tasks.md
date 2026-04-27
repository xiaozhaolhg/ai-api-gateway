# Tasks: add-admin-register

---

## Phase 1: Auth Service Register RPC

### Task 1.1: Add Register RPC to auth-service proto

**Acceptance Criteria:**
- [x] Proto file includes `Register` RPC in `AuthService` service
- [x] `RegisterRequest` message has username/email, name, password fields
- [x] `RegisterResponse` message has user fields
- [x] Proto compiles without errors (buf generate)

**Unit Tests:**
- [x] Proto message validation tests

**FVT Test Plan:**
- [x] Verify proto definition compiles (buf generate)

---

### Task 1.2: Implement Register handler in auth-service

**Acceptance Criteria:**
- [x] Register RPC handler accepts username/email, name, password
- [x] Hashes password with bcrypt before storing
- [x] Returns created user on success
- [x] Returns error on duplicate username/email

**Unit Tests:**
- [x] Handler creates user with hashed password
- [x] Handler rejects duplicate: 409
- [x] Handler rejects weak password: 400

**FVT Test Plan:**
- [x] Call Register RPC with valid credentials, verify user created
- [x] Call Register RPC with duplicate, verify 409

---

## Phase 2: Gateway Service Register Endpoint

### Task 2.1: Add POST /admin/register endpoint

**Acceptance Criteria:**
- [x] Endpoint exists at `POST /admin/register`
- [x] Proxies request to auth-service Register RPC
- [x] Returns JWT token on success
- [x] Sets auth cookie

**Unit Tests:**
- [x] Handler proxies to auth-service correctly

**FVT Test Plan:**
- [x] POST /admin/register, verify proxies to auth-service

---

## Phase 3: Admin UI Register Page

### Task 3.1: Create register page component

**Acceptance Criteria:**
- [x] Register page renders
- [x] Link from login page to register

**Unit Tests:**
- [x] Component renders

**FVT Test Plan:**
- [x] Visit /admin/register, verify page renders

---

### Task 3.2: Implement registration form with validation

**Acceptance Criteria:**
- [x] Username or Email field with validation
- [x] Name field with validation
- [x] Password field with validation (min 8 chars)
- [x] Required field validation

**Unit Tests:**
- [x] Form validation tests

**FVT Test Plan:**
- [x] Submit empty form, verify validation errors

---

### Task 3.3: Add registration API call

**Acceptance Criteria:**
- [x] Calls POST /admin/register
- [x] Auto-login after success

**Unit Tests:**
- [x] Mutation tests

**FVT Test Plan:**
- [x] Register, verify API called

---

### Task 3.4: Handle success/error states

**Acceptance Criteria:**
- [x] Success redirects to dashboard
- [x] Error shows message

**FVT Test Plan:**
- [x] Valid registration, verify redirect
- [x] Invalid registration, verify error shown

---

## Phase 4: Integration Tests

### Task 4.1: Full registration flow test

**Acceptance Criteria:**
- [x] Full registration flow integration test
- [x] Auto-login after register integration test

**FVT Test Plan:**
- [x] Run integration tests

---

## Test Summary

| Phase | Unit Tests | FVT Tests |
|-------|-----------|----------|
| Phase 1 | 4 | 2 |
| Phase 2 | 1 | 1 |
| Phase 3 | 5 | 5 |
| Phase 4 | 2 | 1 |
| **Total** | **12** | **9** |

---

## End-to-End Acceptance

Each task must pass E2E verification before marking complete:

```bash
# 1. Clean rebuild of affected service
make down
make clean-images SERVICE=auth-service  # or gateway-service
make up

# 2. Verify with API calls
# Register
curl -X POST http://localhost:8080/admin/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","name":"Test User","email":"test@example.com","password":"securepass123"}'

# Login
curl -X POST http://localhost:8080/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"securepass123"}'

# Verify auth works
curl http://localhost:8080/admin/auth/me \
  -H "Cookie: auth_token=<token>"
```

---

## Enhancement: Username Login Support

Added support for login with username OR email:

**Changes:**
- [x] AuthService.Login accepts emailOrUsername parameter
- [x] Gets user by email first, then by username
- [x] User entity has Username field
- [x] UserRepository has GetByUsername method
- [x] Register stores username when provided

**Test:**
```bash
# Register with username
curl -X POST http://localhost:8080/admin/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"myuser","name":"My User","email":"myuser@test.com","password":"securepass123"}'

# Login with email
curl -X POST http://localhost:8080/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"myuser@test.com","password":"securepass123"}'

# Login with username (same as email)
curl -X POST http://localhost:8080/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"myuser","password":"securepass123"}'
```
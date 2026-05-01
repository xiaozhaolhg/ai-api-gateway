---

## API Contracts

### 5.1 Admin UI Login Endpoint

#### POST /admin/login

Authenticates user via email/password and sets JWT cookie.

**Request:**
```json
{
  "email": "admin@example.com",
  "password": "secure_password"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "user": {
    "id": "usr_abc123",
    "email": "admin@example.com",
    "name": "Admin User",
    "role": "admin"
  }
}
```

**Response (401 Unauthorized):**
```json
{
  "success": false,
  "error": "invalid_credentials",
  "message": "Invalid email or password"
}
```

**Cookies:**
- `auth_token`: JWT token
  - Name: `auth_token`
  - Path: `/admin`
  - HttpOnly: `true`
  - Secure: `true` (production)
  - SameSite: `strict`
  - Max-Age: `86400` (24 hours)

---

#### POST /admin/logout

Clears auth cookie and invalidates session.

**Request:** Empty body

**Response (200 OK):**
```json
{
  "success": true
}
```

**Cookies:**
- `auth_token`: Expired/max-age=0

---

#### GET /admin/me

Returns current authenticated user (for session restoration).

**Headers:**
```
Cookie: auth_token=<jwt>
```

**Response (200 OK):**
```json
{
  "id": "usr_abc123",
  "email": "admin@example.com",
  "name": "Admin User",
  "role": "admin"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "unauthenticated"
}
```

---

### 5.2 Auth Service RPC

#### Login RPC

```protobuf
message LoginRequest {
  string email = 1;
  string password = 2;
}

message LoginResponse {
  string token = 1;
  User user = 2;
}

service AuthService {
  rpc Login(LoginRequest) returns (LoginResponse);
}
```

---

### 5.3 JWT Token Structure

**Payload Claims:**
```json
{
  "sub": "usr_abc123",
  "email": "admin@example.com",
  "role": "admin",
  "exp": 1715616000,
  "iat": 1715530000
}
```

| Claim | Type | Description |
|-------|------|-------------|
| `sub` | string | User ID (usr_*) |
| `email` | string | User email |
| `role` | string | admin, user, or viewer |
| `exp` | int | Expiration timestamp (Unix epoch) |
| `iat` | int | Issued at timestamp |

---

### 5.4 Role-Based Access Matrix

| Endpoint | admin | user | viewer |
|----------|-------|------|--------|
| GET /admin/me | ✓ | ✓ | ✓ |
| POST /admin/logout | ✓ | ✓ | ✓ |
| GET /admin/providers | ✓ | ✗ | ✗ |
| GET /admin/users | ✓ | ✗ | ✗ |
| GET /admin/api-keys | ✓ | own only | ✗ |
| GET /admin/usage | ✓ | own only | own only |
| GET /admin/health | ✓ | ✓ | ✓ |
| GET /admin/settings | ✓ | ✓ | ✗ |

---

### 5.5 Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `invalid_credentials` | 401 | Email/password mismatch |
| `unauthenticated` | 401 | Missing or invalid JWT |
| `forbidden` | 403 | User lacks required role |
| `user_not_found` | 404 | User ID doesn't exist |
| `server_error` | 500 | Internal service error |

---

## Internal API Client

All admin-ui API calls use the APIClientInterface (src/api/types.ts), implemented by both RealAPIClient and MockAPIClient for consistent behavior.

### Requirement: Unified Client Interface
The admin-ui SHALL use a unified API client that implements the same interface for both mock and real modes.

#### Scenario: Interface compliance
- **WHEN** the APIClientInterface is inspected
- **THEN** it SHALL define typed methods for all /admin/* endpoints
- **AND** both RealAPIClient and MockAPIClient SHALL implement this interface

#### Scenario: Mode switching
- **WHEN** the mode is switched via DevTools or environment variable
- **THEN** the UnifiedAPIClient SHALL delegate to the appropriate implementation
- **AND** the API behavior SHALL remain consistent

### RealAPIClient

Implements HTTP calls to gateway-service at `VITE_API_BASE_URL`:

```typescript
class RealAPIClient implements APIClientInterface {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  async login(email: string, password: string): Promise<LoginResponse> {
    const resp = await fetch(`${this.baseURL}/admin/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });
    return resp.json();
  }

  // ... other methods matching APIClientInterface
}
```

### MockAPIClient

Implements mock data operations using MockDataHandler:

```typescript
class MockAPIClient implements APIClientInterface {
  private dataHandler: MockDataHandler;
  private delay: number;

  constructor(delay: number = 500) {
    this.dataHandler = MockDataHandler.getInstance();
    this.delay = delay;
  }

  private async simulateNetworkDelay<T>(data: T): Promise<T> {
    if (this.delay === 0) return data;
    return new Promise(resolve => setTimeout(() => resolve(data), this.delay));
  }

  async login(email: string, password: string): Promise<LoginResponse> {
    // Mock implementation with mock data
    const user = this.dataHandler.getUsers().find(u => u.email === email);
    if (!user) throw new Error('User not found');
    return this.simulateNetworkDelay({
      token: `mock-jwt-token-${Date.now()}`,
      user
    });
  }

  // ... other methods matching APIClientInterface
}
```

### APIClientInterface Definition

```typescript
export interface APIClientInterface {
  // Authentication
  login(email: string, password: string): Promise<LoginResponse>;
  register(name: string, username: string, email: string, password: string): Promise<RegisterResponse>;
  logout(): Promise<void>;
  getCurrentUser(): Promise<User>;

  // Providers
  getProviders(): Promise<Provider[]>;
  createProvider(provider: Omit<Provider, 'id' | 'created_at' | 'updated_at'>): Promise<Provider>;
  updateProvider(id: string, provider: Partial<Provider>): Promise<Provider>;
  deleteProvider(id: string): Promise<void>;

  // Users, API Keys, Usage, Routing Rules, Groups, Permissions, Budgets, Pricing Rules, Alert Rules, Alerts, Health
  // ... full CRUD interface for all entities
}
```

---

### 5.6 User Entity (Updated)

```protobuf
message User {
  string id = 1;
  string email = 2;
  string name = 3;
  string role = 4;          // "admin" | "user" | "viewer"
  string password_hash = 5;  // bcrypt hash
  int64 created_at = 6;
  int64 updated_at = 7;
}
```

---

### 5.7 Initial Admin User

For first-time setup, a bootstrap endpoint creates the initial admin:

#### POST /admin/bootstrap

**Request:**
```json
{
  "email": "admin@example.com",
  "password": "secure_password",
  "name": "Initial Admin"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "user": {
    "id": "usr_system",
    "email": "admin@example.com",
    "name": "Initial Admin",
    "role": "admin"
  }
}
```

**Security:** This endpoint should be disabled after first admin is created (via config flag or environment variable).

# Design: Mock API System

## Architecture Overview

```
┌─────────────────────────────────────────┐
│           API Client Layer              │
├─────────────────────────────────────────┤
│  ┌─────────────┐    ┌──────────────┐   │
│  │ Real API    │    │ Mock API     │   │
│  │ Client      │    │ Client       │   │
│  └─────────────┘    └──────────────┘   │
│         │                   │           │
│         ▼                   ▼           │
│  ┌─────────────┐    ┌──────────────┐   │
│  │ HTTP Fetch  │    │ Local Data   │   │
│  └─────────────┘    └──────────────┘   │
└─────────────────────────────────────────┘
```

## Components

### 1. UnifiedAPIClient (src/api/client.ts)

A unified client that switches between RealAPIClient and MockAPIClient based on configuration.

```typescript
class UnifiedAPIClient implements APIClientInterface {
  private realClient: RealAPIClient;
  private mockClient: MockAPIClient;
  private useMock: boolean;

  constructor() {
    this.realClient = new RealAPIClient(API_CONFIG.baseURL);
    this.mockClient = new MockAPIClient(API_CONFIG.mockDelay);
    this.useMock = API_CONFIG.useMock;
  }

  private getActiveClient(): APIClientInterface {
    return this.useMock ? this.mockClient : this.realClient;
  }

  // All API methods delegate to active client
  async login(email: string, password: string) {
    return this.getActiveClient().login(email, password);
  }
  // ... other methods
}
```

### 2. APIClientInterface (src/api/types.ts)

Defines the contract that both RealAPIClient and MockAPIClient must implement.

```typescript
export interface APIClientInterface {
  // Authentication
  login(email: string, password: string): Promise<LoginResponse>;
  register(name: string, username: string, email: string, password: string): Promise<RegisterResponse>;
  logout(): Promise<void>;
  getCurrentUser(): Promise<User>;

  // Providers, Users, API Keys, Usage, Routing Rules, Groups, Permissions, Budgets, Pricing Rules, Alert Rules, Alerts, Health
  // ... full CRUD interface for all entities
}
```

### 3. MockAPIClient (src/api/mockClient.ts)

Implements the full APIClientInterface with mock data operations.

- Uses MockDataHandler for all data operations
- Simulates network delay via configurable delay
- Returns realistic mock responses

### 4. MockDataHandler (src/mock/handlers/dataHandler.ts)

Singleton that manages mock data stored in localStorage.

- Load data from localStorage on initialization
- Save data to localStorage on every change
- Provide CRUD operations for all entity types
- Support import/export of mock data as JSON

### 5. MockManager (src/mock/MockManager.ts)

High-level manager for mock data operations.

- Reset to default data
- Export data as JSON file (download)
- Import data from JSON file
- Get data statistics
- Switch mock/real mode at runtime

### 6. DevTools (src/components/DevTools.tsx)

React component that provides runtime configuration.

- Only renders in development mode (hidden in production)
- Switch between Mock and Real API modes
- Adjust network delay simulation
- Reset, export, import mock data
- Display data statistics
- Show current API configuration

### 7. Default Mock Data (src/mock/data/index.ts)

Pre-configured mock data for development.

- 3 sample users (admin, user, viewer)
- 2 providers (Ollama Local, OpenCode Zen)
- Sample API keys, usage records, routing rules
- Groups, permissions, budgets, pricing rules
- Alert rules and active alerts

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_BASE_URL` | Real API base URL | `http://localhost:8080` |
| `VITE_USE_MOCK` | Enable Mock API mode | `false` |
| `VITE_MOCK_DELAY` | Mock network delay (ms) | `500` |

### Runtime Configuration

Use the DevTools panel to switch modes and configure settings without restarting the dev server.

## Data Persistence

Mock data is stored in `localStorage` under the key `mockDataStore`.

- Data survives page refreshes
- Data persists across browser sessions
- Can be exported/imported for backup or sharing

## Testing

Unit tests for MockAPIClient cover:

- Authentication flows (login, register, logout, get current user)
- CRUD operations for all entities
- Error handling (not found, duplicate, invalid input)
- Network delay simulation

# Mock API Development Guide

This guide explains how to use the Mock API system for frontend development without requiring a running backend.

## Overview

The Mock API system provides a complete implementation of all admin API endpoints, allowing you to:
- Develop frontend features independently
- Test UI components with realistic data
- Simulate various API scenarios
- Work offline without backend dependencies

## Architecture

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

## Quick Start

### 1. Enable Mock Mode

Create or edit `.env.development`:

```bash
VITE_USE_MOCK=true
VITE_MOCK_DELAY=500
```

### 2. Start Development Server

```bash
npm run dev
```

### 3. Access DevTools

Click the gear icon (⚙️) in the bottom-right corner to access the DevTools panel.

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_BASE_URL` | Real API base URL | `http://localhost:8080` |
| `VITE_USE_MOCK` | Enable Mock API mode | `false` |
| `VITE_MOCK_DELAY` | Mock network delay in ms | `500` |

### Runtime Configuration

Use the DevTools panel to:
- Switch between Mock and Real API modes
- Adjust network delay simulation
- Reset mock data to defaults
- Export/import mock data

## Mock Data Structure

### Default Data

The system includes comprehensive default mock data:

#### Users
- 3 sample users with different roles (admin, user, viewer)
- Pre-configured with email addresses and status

#### Providers
- 2 providers: Ollama Local and OpenCode Zen
- Each with multiple models and configuration

#### Routing Rules
- 3 routing rules for different model patterns
- Configured with priorities and fallback chains

#### Budgets
- 3 budgets at different scopes (global, group, user)
- With current spend tracking

#### Pricing Rules
- 3 pricing rules for different models
- With prompt and completion pricing

#### Alert Rules & Alerts
- 3 alert rules for different metrics
- 2 active alerts with different severities

### Data Persistence

Mock data is stored in `localStorage` under the key `mockDataStore`. This means:
- Data survives page refreshes
- Data persists across browser sessions
- You can export/import data for backup or sharing

## API Coverage

### Authentication
- ✅ `login(email, password)` - User authentication
- ✅ `register(name, username, email, password)` - User registration
- ✅ `logout()` - User logout
- ✅ `getCurrentUser()` - Get current user info

### Providers
- ✅ `getProviders()` - List all providers
- ✅ `createProvider(provider)` - Create new provider
- ✅ `updateProvider(id, provider)` - Update provider
- ✅ `deleteProvider(id)` - Delete provider

### Users
- ✅ `getUsers()` - List all users
- ✅ `createUser(user)` - Create new user
- ✅ `updateUser(id, user)` - Update user
- ✅ `deleteUser(id)` - Delete user

### API Keys
- ✅ `getAPIKeys(userId)` - List user's API keys
- ✅ `createAPIKey(userId, name)` - Create new API key
- ✅ `deleteAPIKey(id)` - Delete API key

### Usage
- ✅ `getUsage(userId, startDate, endDate)` - Get usage records

### Routing Rules
- ✅ `getRoutingRules()` - List all routing rules
- ✅ `createRoutingRule(rule)` - Create new routing rule
- ✅ `updateRoutingRule(id, rule)` - Update routing rule
- ✅ `deleteRoutingRule(id)` - Delete routing rule

### Groups
- ✅ `getGroups()` - List all groups
- ✅ `createGroup(group)` - Create new group
- ✅ `updateGroup(id, group)` - Update group
- ✅ `deleteGroup(id)` - Delete group
- ✅ `addGroupMember(groupId, userId)` - Add user to group
- ✅ `removeGroupMember(groupId, userId)` - Remove user from group

### Permissions
- ✅ `getPermissions()` - List all permissions
- ✅ `createPermission(permission)` - Create new permission
- ✅ `updatePermission(id, permission)` - Update permission
- ✅ `deletePermission(id)` - Delete permission

### Budgets
- ✅ `getBudgets()` - List all budgets
- ✅ `createBudget(budget)` - Create new budget
- ✅ `updateBudget(id, budget)` - Update budget
- ✅ `deleteBudget(id)` - Delete budget

### Pricing Rules
- ✅ `getPricingRules()` - List all pricing rules
- ✅ `createPricingRule(rule)` - Create new pricing rule
- ✅ `updatePricingRule(id, rule)` - Update pricing rule
- ✅ `deletePricingRule(id)` - Delete pricing rule

### Alert Rules
- ✅ `getAlertRules()` - List all alert rules
- ✅ `createAlertRule(rule)` - Create new alert rule
- ✅ `updateAlertRule(id, rule)` - Update alert rule
- ✅ `deleteAlertRule(id)` - Delete alert rule

### Alerts
- ✅ `getAlerts()` - List all alerts
- ✅ `acknowledgeAlert(id)` - Acknowledge alert

### Health
- ✅ `getProviderHealth()` - Get provider health status

## Development Workflow

### Typical Development Session

1. **Start with Mock Mode**
   ```bash
   # Enable mock mode in .env.development
   VITE_USE_MOCK=true
   npm run dev
   ```

2. **Develop UI Components**
   - Use mock data to build and test components
   - No need to run backend services

3. **Test with Real API**
   - Use DevTools to switch to Real API mode
   - Test integration with actual backend

4. **Debug Issues**
   - Switch back to Mock mode
   - Export current state for debugging
   - Create specific test scenarios

### Creating Test Scenarios

1. **Export Current State**
   - Use DevTools → Export Mock Data
   - Save as JSON file

2. **Modify Data**
   - Edit the JSON file to create specific scenarios
   - Add edge cases, error conditions, etc.

3. **Import Test Data**
   - Use DevTools → Import Mock Data
   - Load your custom test scenario

4. **Test Your Changes**
   - Verify UI behavior with the test data
   - Reset to defaults when done

## Testing

### Unit Tests

The Mock API includes comprehensive unit tests:

```bash
# Run all tests
npm run test

# Run tests in watch mode
npm run test: --watch

# Run tests with UI
npm run test:ui
```

### Test Coverage

Tests cover:
- Authentication flows
- CRUD operations for all entities
- Error handling
- Network delay simulation
- Data persistence

## Advanced Usage

### Custom Mock Data

Create custom mock data for specific test cases:

```typescript
import MockDataHandler from './mock/handlers/dataHandler';

const handler = MockDataHandler.getInstance();

// Add custom test data
handler.addUser({
  id: 'custom-1',
  name: 'Test User',
  email: 'test@example.com',
  role: 'admin',
  status: 'active',
  created_at: new Date().toISOString()
});
```

### Programmatic Control

```typescript
import MockManager from './mock/MockManager';

const manager = MockManager.getInstance();

// Switch modes programmatically
manager.setMockMode(true);

// Export data
const jsonData = manager.exportData();

// Import data
manager.importData(jsonData);

// Get statistics
const stats = manager.getDataStats();
console.log(stats);
```

### Network Simulation

Adjust network delay to test loading states:

```typescript
// Create client with custom delay
const mockClient = new MockAPIClient(2000); // 2 second delay
```

## Troubleshooting

### Mock Data Not Persisting

- Check if localStorage is enabled in your browser
- Verify you're not in private/incognito mode
- Check browser console for localStorage errors

### API Calls Failing in Mock Mode

- Verify `VITE_USE_MOCK=true` is set
- Check that DevTools shows "Mock" as current mode
- Try resetting mock data to defaults

### Data Not Updating After Changes

- Mock data changes require page reload to take effect
- Use DevTools → Reset to Defaults if data is corrupted
- Clear browser cache if issues persist

### Real API Not Working After Switching

- Verify backend services are running
- Check `VITE_API_BASE_URL` is correct
- Ensure CORS is configured on backend

## Best Practices

1. **Start with Mock Mode**: Develop UI features first with mock data
2. **Test Early**: Switch to real API early to catch integration issues
3. **Version Control**: Commit meaningful mock data configurations
4. **Clean Up**: Reset mock data before committing changes
5. **Document Scenarios**: Keep notes on important test scenarios

## Future Enhancements

Planned improvements to the Mock API system:

- [ ] Error simulation (network errors, server errors)
- [ ] Advanced data generation (randomized data)
- [ ] Request/response logging
- [ ] Performance profiling
- [ ] Automated test scenario generation
- [ ] GraphQL-style query simulation

## Contributing

When adding new API endpoints:

1. Add interface method to `APIClientInterface` in `src/api/types.ts`
2. Implement in `RealAPIClient` class
3. Implement in `MockAPIClient` class
4. Add default mock data to `src/mock/data/index.ts`
5. Add data handler methods to `MockDataHandler`
6. Write unit tests in `src/api/__tests__/mockClient.test.ts`

## Support

For issues or questions:
- Check this guide first
- Review unit tests for usage examples
- Inspect DevTools panel for runtime information
- Check browser console for error messages

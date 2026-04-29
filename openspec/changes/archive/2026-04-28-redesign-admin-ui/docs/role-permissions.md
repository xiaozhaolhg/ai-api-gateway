# Role Permissions Matrix

This document describes the access control matrix for the admin UI, including which roles can access each feature and what operations they can perform.

## Roles Overview

| Role | Description | Use Case |
|------|-------------|----------|
| **admin** | Full system access | System administrators managing all resources |
| **user** | Limited write access | Regular users managing their own resources |
| **viewer** | Read-only access | Auditors, stakeholders monitoring usage |

## Navigation Access Matrix

| Page | admin | user | viewer | Notes |
|------|-------|------|--------|-------|
| Dashboard | ✓ Full | ✓ Full | ✓ Full | All roles see system overview |
| Providers | ✓ Full | ✗ | ✗ | Admin-only provider configuration |
| Routing Rules | ✓ Full | ✗ | ✗ | Admin-only routing configuration |
| Users | ✓ Full | ✗ | ✗ | Admin-only user management |
| Groups | ✓ Full | ✗ | ✗ | Admin-only group management |
| API Keys | ✓ All | ✓ Own | ✗ | Users can only see/manage own keys |
| Permissions | ✓ Full | ✗ | ✗ | Admin-only permission management |
| Usage | ✓ All | ✓ Own | ✓ Own | Filtered by ownership for non-admin |
| Budgets | ✓ Full | ✗ | ✗ | Admin-only budget configuration |
| Pricing Rules | ✓ Full | ✗ | ✗ | Admin-only pricing configuration |
| Health | ✓ Full | ✓ Full | ✓ Full | All roles can monitor system health |
| Alerts | ✓ Full | ✗ | ✗ | Admin-only alert management |
| Settings | ✓ Full | ✓ Own | ✗ | Users can update own profile |

## Operation Permissions Matrix

### Dashboard

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| View summary cards | ✓ | ✓ | ✓ |
| View quick actions | ✓ | ✓ | ✓ |
| Execute quick actions | ✓ | ✓ (limited) | ✗ |

### Providers

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| List providers | ✓ | ✗ | ✗ |
| Create provider | ✓ | ✗ | ✗ |
| Edit provider | ✓ | ✗ | ✗ |
| Delete provider | ✓ | ✗ | ✗ |
| Test provider health | ✓ | ✗ | ✗ |

### Routing Rules

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| List routing rules | ✓ | ✗ | ✗ |
| Create routing rule | ✓ | ✗ | ✗ |
| Edit routing rule | ✓ | ✗ | ✗ |
| Delete routing rule | ✓ | ✗ | ✗ |
| Change rule priority | ✓ | ✗ | ✗ |

### Users

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| List users | ✓ | ✗ | ✗ |
| Create user | ✓ | ✗ | ✗ |
| Edit user | ✓ | ✗ | ✗ |
| Delete user | ✓ | ✗ | ✗ |
| Change user role | ✓ | ✗ | ✗ |
| Reset user password | ✓ | ✗ | ✗ |

### Groups

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| List groups | ✓ | ✗ | ✗ |
| Create group | ✓ | ✗ | ✗ |
| Edit group | ✓ | ✗ | ✗ |
| Delete group | ✓ | ✗ | ✗ |
| Add/remove members | ✓ | ✗ | ✗ |

### API Keys

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| List all keys | ✓ | ✗ | ✗ |
| List own keys | ✓ | ✓ | ✗ |
| Create key | ✓ | ✓ (own) | ✗ |
| Delete key | ✓ | ✓ (own) | ✗ |
| View key once | ✓ | ✓ (own) | ✗ |

### Permissions

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| List permissions | ✓ | ✗ | ✗ |
| Create permission | ✓ | ✗ | ✗ |
| Edit permission | ✓ | ✗ | ✗ |
| Delete permission | ✓ | ✗ | ✗ |

### Usage

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| View all usage | ✓ | ✗ | ✗ |
| View own usage | ✓ | ✓ | ✓ |
| Filter by date | ✓ | ✓ | ✓ |
| Export data | ✓ | ✓ (own) | ✗ |
| View cost analysis | ✓ | ✓ (own) | ✗ |

### Budgets

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| List budgets | ✓ | ✗ | ✗ |
| Create budget | ✓ | ✗ | ✗ |
| Edit budget | ✓ | ✗ | ✗ |
| Delete budget | ✓ | ✗ | ✗ |
| Set budget alerts | ✓ | ✗ | ✗ |

### Pricing Rules

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| List pricing rules | ✓ | ✗ | ✗ |
| Create pricing rule | ✓ | ✗ | ✗ |
| Edit pricing rule | ✓ | ✗ | ✗ |
| Delete pricing rule | ✓ | ✗ | ✗ |

### Health

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| View service status | ✓ | ✓ | ✓ |
| View provider health | ✓ | ✓ | ✓ |
| View latency metrics | ✓ | ✓ | ✓ |
| View error rates | ✓ | ✓ | ✓ |

### Alerts

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| List alert rules | ✓ | ✗ | ✗ |
| Create alert rule | ✓ | ✗ | ✗ |
| Edit alert rule | ✓ | ✗ | ✗ |
| Delete alert rule | ✓ | ✗ | ✗ |
| List active alerts | ✓ | ✗ | ✗ |
| Acknowledge alerts | ✓ | ✗ | ✗ |

### Settings

| Operation | admin | user | viewer |
|-----------|-------|------|--------|
| View own profile | ✓ | ✓ | ✗ |
| Edit own profile | ✓ | ✓ | ✗ |
| Change password | ✓ | ✓ | ✗ |
| View all users | ✓ | ✗ | ✗ |
| Edit other users | ✓ | ✗ | ✗ |

## Examples

### Example 1: Admin User Workflow

An admin user can:
1. Navigate to **Users** → Create a new user with role "user"
2. Navigate to **API Keys** → Generate an API key for the new user
3. Navigate to **Groups** → Create a group and add the user
4. Navigate to **Permissions** → Grant the group access to specific models
5. Navigate to **Budgets** → Set a monthly budget for the user

### Example 2: Regular User Workflow

A user can:
1. Navigate to **Dashboard** → View system summary
2. Navigate to **API Keys** → View and manage only their own API keys
3. Navigate to **Usage** → View their own usage history and costs
4. Navigate to **Health** → Monitor system health status
5. Navigate to **Settings** → Update their profile or change password

### Example 3: Viewer User Workflow

A viewer can:
1. Navigate to **Dashboard** → View system summary
2. Navigate to **Usage** → View their own usage (read-only)
3. Navigate to **Health** → Monitor system health status
4. Cannot access any CRUD operations or configuration pages

## Implementation Notes

### Client-Side Enforcement

- Navigation menu items are filtered based on `user.role` from `AuthContext`
- Protected routes use `ProtectedRoute` wrapper with optional `requiredRole` prop
- Action buttons (Edit, Delete) are conditionally rendered based on role

### Server-Side Enforcement

- Gateway middleware validates JWT and extracts user role
- Each admin endpoint checks role before processing write operations
- Data queries are scoped by user_id for non-admin users
- Role checks are **security boundaries**, not just UX convenience

### Error Handling

- Unauthorized access returns HTTP 403 Forbidden
- Client displays "Access Denied" message and redirects to dashboard
- Attempted API calls without proper role are rejected server-side

## Security Considerations

1. **Client-side filtering is not security**: Role-based UI filtering is for UX only. All security enforcement happens server-side.

2. **Principle of least privilege**: Users are granted minimum permissions needed for their role.

3. **Audit trail**: All write operations are logged with user_id for accountability.

4. **Session management**: JWT tokens expire after 24 hours. Logout invalidates the session.

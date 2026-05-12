# Admin UI - Hide Permissions and Budgets Pages

## Overview

Remove the standalone Permissions and Budgets pages from the admin UI navigation and routing. The `GroupPermissionsTab` component (used within the Group edit dialog) is preserved.

## Current State

### App.tsx Routes

```tsx
<Route path="permissions" element={<ProtectedRoute requiredRole="admin"><Permissions /></ProtectedRoute>} />
<Route path="budgets" element={<ProtectedRoute requiredRole="admin"><Budgets /></ProtectedRoute>} />
```

### AppShell.tsx Sidebar

```typescript
'/permissions': ['admin'],
'/budgets': ['admin'],
```

Sidebar menu items for "Permissions" and "Budgets".

## Expected State

### App.tsx Routes

Both route entries removed.

### AppShell.tsx Sidebar

Both sidebar entries removed.

## Components NOT Modified

- `GroupPermissionsTab.tsx` — used inside Group edit modal, preserved
- `pages/Permissions.tsx` — source kept for future re-enablement
- `pages/Budgets.tsx` — source kept for future re-enablement
- `App.tsx` import statements for Permissions/Budgets may remain (tree-shaken)

## Verification

- TypeScript compilation: `npx tsc --noEmit` must pass
- Group edit modal's Permissions tab must still function
- Sidebar must no longer show Permissions or Budgets entries
- Direct navigation to `/admin/permissions` or `/admin/budgets` should 404 or redirect

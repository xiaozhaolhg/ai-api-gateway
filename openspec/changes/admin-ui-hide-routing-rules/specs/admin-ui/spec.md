# Admin UI - Hide Routing Rules Page

## Overview

Remove the standalone Routing Rules page from the admin UI navigation and routing.

## Current State

### App.tsx

```tsx
import { RoutingRules } from './pages/RoutingRules';

<Route path="routing" element={<ProtectedRoute requiredRole="admin"><RoutingRules /></ProtectedRoute>} />
```

### AppShell.tsx

```typescript
import { BranchesOutlined } from '@ant-design/icons';

roleAccess: {
  '/routing': ['admin'],
}

// Sidebar "Infrastructure" group:
{
  key: '/routing',
  icon: <BranchesOutlined />,
  label: t('navigation.routingRules'),
},
```

## Expected State

Route and sidebar entry removed. Import for `RoutingRules` removed from App.tsx if it's the only usage. `BranchesOutlined` import in AppShell.tsx should be checked — if it's only used for the routing rules icon, remove it.

## Components NOT Modified

- `pages/RoutingRules.tsx` — source kept
- `components/DevTools.tsx` — references to routing rule data stats kept

## Verification

- TypeScript compilation: `npx tsc --noEmit` must pass
- Sidebar must no longer show Routing Rules entry
- Direct navigation to `/admin/routing` should fall through to SPA root

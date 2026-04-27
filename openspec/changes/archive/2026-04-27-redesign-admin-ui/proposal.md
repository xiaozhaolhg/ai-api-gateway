## Why

The admin UI started as a bare-bones prototype with no authentication, no role-based access control, and a plain interface. This change tracks the comprehensive redesign that adds proper auth flows, a professional antd-based UI, i18n support, and production-quality admin experience.

## What Changes

- Add login and registration pages with email/username + password authentication
- Add session management via gateway-service auth endpoints (`/admin/auth/login`, `/admin/auth/register`)
- Implement role-based access control with three roles: admin, user, viewer
- Redesign UI using antd 6 component library with Tailwind CSS
- Add AppShell layout with sidebar navigation and role-filtered tabs
- Add dashboard page as the landing page after login
- Add i18n support (English/Chinese)
- Upgrade data fetching to TanStack Query for caching and loading states
- Add proper form handling with react-hook-form and validation
- Implement ProtectedRoute wrapper for route guards
- Add language switcher in header

## Capabilities

### New Capabilities
- `admin-ui-auth`: Login page, registration page, session management, auth context, route guards
- `admin-ui-dashboard`: Dashboard with summary cards and key metrics
- `admin-ui-i18n`: Internationalization with English/Chinese support
- `admin-ui-groups`: User group management page
- `admin-ui-permissions`: Permission management page
- `admin-ui-budgets`: Budget management page
- `admin-ui-pricing`: Pricing rules management page
- `admin-ui-routing`: Routing rules management page
- `admin-ui-alerts`: Alert management page

### Modified Capabilities
- `admin-ui-layout`: AppShell with antd Layout, sidebar navigation, role-filtered tabs
- `admin-ui-pages`: Redesigned pages with antd components (Table, Form, Modal, Card, etc.)

## Impact

- **Code**: admin-ui (full rewrite of components, new pages, new dependencies), gateway-service (new `/admin/auth/login`, `/admin/auth/register` endpoints), auth-service (new `Login`, `Register` RPCs)
- **APIs**: New auth endpoints, JWT cookie handling
- **Dependencies**: antd 6, TanStack Query, react-hook-form, i18next, react-router-dom v7
- **Systems**: Admin UI accessibility and security posture significantly improved

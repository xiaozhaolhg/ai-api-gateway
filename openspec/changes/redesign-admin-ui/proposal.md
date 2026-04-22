## Why

The current admin UI is a bare-bones prototype with no authentication, no role-based access control, and a plain, unstyled interface. It lacks a login page, has no visual hierarchy or iconography, and every page is accessible to anyone. The UI needs a professional redesign with proper auth flows, a collapsible sidebar with role-filtered navigation, and a modern component library to deliver a production-quality admin experience.

## What Changes

- Add login page with email/password authentication (aligned with auth-service)
- Add session management via gateway-service login endpoint
- Implement role-based navigation with three roles: admin, user, viewer
- Redesign UI using shadcn/ui component library with Lucide icons
- Add collapsible sidebar with active state indicators and role-filtered tabs
- Add dashboard page as the landing page after login
- Upgrade data fetching to TanStack Query for caching and loading states
- Add proper form handling with react-hook-form and validation
- Redirect `/admin` to login (unauthenticated) or dashboard (authenticated)

## Capabilities

### New Capabilities
- `admin-ui-auth`: Login page, session management, auth context, route guards
- `admin-ui-dashboard`: Dashboard with summary cards and key metrics
- `admin-ui-rbac`: Role-based navigation filtering and access control (admin/user/viewer)

### Modified Capabilities
- `admin-ui-layout`: Collapsible sidebar, shadcn/ui components, Lucide icons, active states
- `admin-ui-pages`: Redesigned pages with shadcn/ui tables, modals, forms

## Impact

- **Code**: admin-ui (full rewrite of components, new pages, new dependencies), gateway-service (new `/admin/login` endpoint), auth-service (new `Login` RPC)
- **APIs**: New `POST /admin/login` endpoint, session cookie or JWT handling
- **Dependencies**: shadcn/ui, Lucide React, TanStack Query, react-hook-form
- **Systems**: Admin UI accessibility and security posture significantly improved

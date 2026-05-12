## Context

The admin UI includes a standalone Routing Rules page under "Infrastructure" section in the sidebar. This page allows CRUD operations on routing rules that map model patterns to providers. However, in practice:

- Routing is primarily configured through provider setup and tier permissions
- The router service handles bare model resolution automatically
- The page is not part of any current operational workflow

## Goals / Non-Goals

**Goals:**
- Remove Routing Rules route from App.tsx
- Remove Routing Rules sidebar item from AppShell.tsx
- UI-only change, zero backend impact

**Non-Goals:**
- Deleting the RoutingRules page component
- Modifying backend routing rule APIs
- Changing router service behavior
- Removing DevTools references

## Decisions

**1. Remove Route and Sidebar, Keep Component**
- **Choice**: Only unregister the route and sidebar entry
- **Rationale**: Source file remains for easy re-enablement. Follows same pattern as Permissions/Budgets removal.

**2. DevTools Reference**
- **Choice**: Leave DevTools references untouched
- **Rationale**: DevTools is a development-only component; references to routing rules data stats cause no harm.

## Scope of Changes

| File | Change |
|------|--------|
| `admin-ui/src/App.tsx` | Remove Route for `/routing` |
| `admin-ui/src/components/AppShell.tsx` | Remove sidebar nav item for Routing Rules and its roleAccess entry |

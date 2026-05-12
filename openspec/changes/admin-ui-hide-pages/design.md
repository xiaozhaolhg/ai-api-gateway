## Context

The admin UI has two standalone pages (Permissions, Budgets) that are registered as routes and displayed in the sidebar. However:

- **Permissions page**: Shows all permissions across all groups, but the actual workflow for managing permissions happens inside the Group edit dialog via `GroupPermissionsTab`. The standalone page provides no additional functionality.
- **Budgets page**: Displays budget data (list, create, edit, delete) but the underlying billing data is incomplete, making the page show empty states.

Both pages are marked as admin-only (`requiredRole="admin"`), but simply hiding unused pages is a cleaner UX than showing non-functional ones.

## Goals / Non-Goals

**Goals:**
- Remove standalone Permissions page route and sidebar entry
- Remove standalone Budgets page route and sidebar entry
- Preserve GroupPermissionsTab (used within Group edit dialog)
- UI-only change, zero backend impact

**Non-Goals:**
- Deleting page components (keep source for future re-enablement)
- Touching backend code
- Changing API contracts
- Removing API client methods

## Decisions

**1. Remove Routes, Keep Components**
- **Choice**: Only remove route registrations in App.tsx and sidebar entries in AppShell.tsx
- **Rationale**: Page components remain in the codebase for easy re-enablement. Avoids deleting code.
- **Alternative**: Delete the page components entirely — rejected, they may be useful when data integration is complete

**2. Preserve GroupPermissionsTab**
- **Choice**: Keep the inline permissions tab within the Group management dialog
- **Rationale**: This is the actual permission editing workflow. It's not a standalone page.
- **Alternative**: Remove entirely — rejected, permissions editing is still needed in group context

## Scope of Changes

| File | Change |
|------|--------|
| `admin-ui/src/App.tsx` | Remove Route for `/permissions` and `/budgets` |
| `admin-ui/src/components/AppShell.tsx` | Remove sidebar nav items for Permissions and Budgets |

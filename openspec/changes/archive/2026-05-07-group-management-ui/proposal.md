## Why

The admin UI has basic Groups and Permissions pages, but they lack critical functionality: the Groups page has no member management or permission assignment, and the Permissions page form fields don't align with the backend's resource_type/resource_id/action model. Users cannot be assigned to groups from the UI, making group-based access control impossible to configure through the admin interface. This blocks the core RBAC workflow: create group → assign permissions → add users → enforce access.

## What Changes

- Add group member management UI: view members, add/remove users from groups
- Fix Permissions page to use backend-aligned fields (resource_type, resource_id, action, effect)
- Add group assignment to Users page (assign users to groups during create/edit)
- Add group detail view showing members and permissions together
- Add search/filter to Users and Groups tables

## Capabilities

### Modified Capabilities
- `admin-ui-groups`: Add member management via expandable row (add/remove users), group permissions view, search/filter, model pattern and parent group fields
- `admin-ui-permissions`: Fix form to use resource_type/resource_id/action/effect fields matching auth-service gRPC API, update table columns, add group filter

### New Capabilities
- `admin-ui-users`: User management page with group assignment during create/edit, password field on creation, search/filter by name/email/role

## Impact

- **admin-ui**: Groups.tsx, Permissions.tsx, Users.tsx pages modified; new GroupDetail component
- **api/client.ts**: May need new API client methods for group members if not already present
- **gateway-service**: No changes — REST endpoints for groups/members/permissions already exist
- **auth-service**: No changes — gRPC handlers for groups/members/permissions already exist

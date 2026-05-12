## Why

The standalone Permissions and Budgets admin pages serve no functional purpose in the current UI. Permissions are managed inline within the Group edit dialog (via GroupPermissionsTab), making the standalone page redundant and confusing. Budgets are not yet integrated with real billing data — the page shows empty/unusable state. Keeping these pages visible creates a broken UX and confuses users.

## What Changes

- Remove the standalone Permissions page route and sidebar navigation entry
- Remove the standalone Budgets page route and sidebar navigation entry
- Keep the GroupPermissionsTab component intact (used within Group edit dialog)
- No backend code changes

## Capabilities

### Removed Capabilities
- `admin-ui-permissions-page`: Standalone Permissions management page
- `admin-ui-budgets-page`: Standalone Budgets management page

### Preserved Capabilities
- `admin-ui-group-permissions-tab`: Inline permission editing within Group management (unchanged)
- `admin-ui-budget-api`: Backend budget CRUD endpoints (unchanged, accessible via API)

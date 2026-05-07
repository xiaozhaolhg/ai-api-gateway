## 1. Setup and Foundation

- [x] 1.1 Review existing Groups.tsx, Permissions.tsx, Users.tsx components in admin-ui
- [x] 1.2 Understand current API integration patterns (React Query setup, API client)
- [x] 1.3 Verify auth-service gRPC API contracts for groups, members, and permissions

## 2. Group Expandable Row with Members and Permissions Tabs

- [x] 2.1 Add Ant Design Table `expandable` prop to Groups page with expanded row rendering
- [x] 2.2 Create expandable row layout with Tabs component (Members tab, Permissions tab)
- [x] 2.3 Implement React Query integration for `GET /admin/groups/:id/members` in expanded row
- [x] 2.4 Display members table with columns: name, email, role, actions (remove button)
- [x] 2.5 Implement "Add Member" button in Members tab that opens user selector
- [x] 2.6 Populate user selector from `GET /admin/users` API
- [x] 2.7 Implement `POST /admin/groups/:id/members` for adding members
- [x] 2.8 Implement `DELETE /admin/groups/:id/members/:userId` for removing members
- [x] 2.9 Implement React Query integration for `GET /admin/permissions?group_id=:id` in expanded row
- [x] 2.10 Display permissions table with columns: resource type, resource ID, action, effect
- [x] 2.11 Implement "Add Permission" button in Permissions tab pre-filled with group ID

## 3. Permission Form Field Alignment

- [x] 3.1 Update permission form to replace `model_pattern` with `resource_type` select (options: model/provider/system)
- [x] 3.2 Add `resource_id` input field to permission form (e.g., `gpt-4`, `ollama:*`)
- [x] 3.3 Add `action` select field to permission form (options: access/manage)
- [x] 3.4 Add `effect` select field to permission form (options: allow/deny)
- [x] 3.5 Update permission form submission to send `group_id`, `resource_type`, `resource_id`, `action`, `effect`
- [x] 3.6 Update permission table columns to display: group name, resource type, resource ID, action, effect, created at, actions
- [x] 3.7 Fix permission edit form to pre-populate with correct backend field names
- [x] 3.8 Add group filter dropdown to Permissions page using `GET /admin/groups`

## 4. User Group Assignment

- [x] 4.1 Add "Groups" multi-select field to User create form populated from `GET /admin/groups`
- [x] 4.2 Implement group assignment after user creation via `POST /admin/groups/:id/members`
- [x] 4.3 Add "Groups" multi-select field to User edit form with current groups pre-selected
- [x] 4.4 Implement group add/remove logic on user edit (compare before/after, call add/remove APIs)
- [x] 4.5 Display user groups as Ant Design Tags in Users table
- [x] 4.6 Add password field to user creation form (required, min 8 characters)
- [x] 4.7 Ensure password is included in `POST /admin/users` request body

## 5. Search and Filter

- [x] 5.1 Add search input to Groups page with client-side filtering by group name
- [x] 5.2 Add search input to Users page with client-side filtering by name or email
- [x] 5.3 Add role filter dropdown to Users page (filter by admin/user role)
- [x] 5.4 Add group filter dropdown to Permissions page (filter by group)
- [x] 5.5 Implement Ant Design Table `filter` and `sorter` props for client-side filtering

## 6. Group Form Enhanced Fields

- [x] 6.1 Add model patterns field (tags input) to Group create/edit form
- [x] 6.2 Implement tags input component for entering patterns like `gpt-*`, `ollama:*`
- [x] 6.3 Add parent group selector to Group create/edit form populated from existing groups
- [x] 6.4 Wire up new fields to `POST /admin/groups` and `PUT /admin/groups/:id` APIs

## 7. Testing

- [x] 7.1 Add unit tests for GroupMembersTab component
- [x] 7.2 Add unit tests for GroupPermissionsTab component
- [x] 7.3 Add unit tests for updated PermissionForm component
- [x] 7.4 Add unit tests for updated UserForm component (groups, password)
- [x] 7.5 Add integration tests for group member add/remove flow
- [x] 7.6 Add integration tests for permission form submission with correct fields
- [x] 7.7 Add integration tests for user group assignment flow
- [x] 7.8 Verify search/filter functionality works correctly on all pages

## Purpose

Permission management interface for the admin UI, including model access permission configuration.

## Requirements

### Requirement: Permission management page
The admin-ui SHALL provide a page at `/permissions` for managing model access permissions.

#### Scenario: List permissions
- **WHEN** the permissions page is loaded
- **THEN** it SHALL call `GET /admin/permissions` via gateway-service
- **AND** display all permissions in an antd Table with columns: group name, resource type, resource ID, action, effect, created at, actions

#### Scenario: Create permission
- **WHEN** the user fills in group, resource_type, resource_id, action, and effect and submits
- **THEN** the admin-ui SHALL call `POST /admin/permissions` with `group_id`, `resource_type`, `resource_id`, `action`, and `effect`
- **AND** the new permission SHALL appear in the table

#### Scenario: Edit permission
- **WHEN** the user clicks edit on a permission
- **THEN** the admin-ui SHALL open a modal with the permission's current values pre-populated (`resource_type`, `resource_id`, `action`, `effect`)
- **AND** on submit, call `PUT /admin/permissions/:id`

#### Scenario: Delete permission
- **WHEN** the user confirms deletion of a permission
- **THEN** the admin-ui SHALL call `DELETE /admin/permissions/:id`
- **AND** remove the permission from the table

### Requirement: Permission form field alignment
The admin-ui permission form SHALL use fields that match the auth-service gRPC API.

#### Scenario: Create permission with correct fields
- **WHEN** the user opens the permission creation form
- **THEN** the form SHALL contain fields: group (select), resource_type (select: model/provider/system), resource_id (input, e.g., `gpt-4`, `ollama:*`), action (select: access/manage), effect (select: allow/deny)

#### Scenario: Submit permission with aligned fields
- **WHEN** the user submits the permission form
- **THEN** the admin-ui SHALL call `POST /admin/permissions` with `group_id`, `resource_type`, `resource_id`, `action`, and `effect`
- **AND** the permission SHALL be created in the backend

#### Scenario: Edit permission with correct fields
- **WHEN** the user edits an existing permission
- **THEN** the form SHALL pre-populate with the permission's `resource_type`, `resource_id`, `action`, and `effect` values

### Requirement: Permission table column alignment
The admin-ui permission table SHALL display fields matching the backend data model.

#### Scenario: Permission table columns
- **WHEN** the permissions page is loaded
- **THEN** the table SHALL display columns: group name, resource type, resource ID, action, effect, created at, actions

### Requirement: Permission filter by group
The admin-ui SHALL support filtering permissions by group.

#### Scenario: Filter permissions by group
- **WHEN** the user selects a group from the filter dropdown
- **THEN** the permissions table SHALL show only permissions for that group

### Requirement: Permission group selector
The admin-ui SHALL provide a group selector when creating or editing permissions.

#### Scenario: Group dropdown populated from groups API
- **WHEN** the user opens the permission form
- **THEN** the group selector SHALL be populated from `GET /admin/groups`
- **AND** display group names as options

---

## Role-Based Access Control

### Requirement: Role-based navigation filtering
The admin UI SHALL filter navigation options based on user role.

#### Scenario: Admin role navigation
- **WHEN** user with 'admin' role loads
- **THEN** show all navigation tabs: Dashboard, Providers, Users, API Keys, Usage, Health, Settings.

#### Scenario: User role navigation
- **WHEN** user with 'user' role loads
- **THEN** show limited tabs: Dashboard, API Keys (own), Usage (own), Health, Settings.

#### Scenario: Viewer role navigation
- **WHEN** user with 'viewer' role loads
- **THEN** show read-only tabs: Dashboard, Usage (own), Health.

### Requirement: Role-based access control
The admin UI SHALL enforce access restrictions at the page level.

#### Scenario: Admin full access
- **WHEN** admin role user accesses any admin page
- **THEN** allow full CRUD operations.

#### Scenario: User limited access
- **WHEN** user role accesses Providers or Users pages
- **THEN** redirect to dashboard with access denied message.

#### Scenario: Viewer read-only access
- **WHEN** viewer role attempts create/edit/delete operations
- **THEN** disable buttons and show read-only state.

### Requirement: Role-based data filtering
The admin UI SHALL filter data based on user role and ownership.

#### Scenario: User API key access
- **WHEN** user role views API Keys page
- **THEN** only show API keys belonging to that user.

#### Scenario: User usage data
- **WHEN** user role views Usage page
- **THEN** only show usage data for that user.

#### Scenario: Viewer usage data
- **WHEN** viewer role views Usage page
- **THEN** only show usage data for that user with read-only controls.


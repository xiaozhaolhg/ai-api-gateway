## Purpose

Permission management interface for the admin UI, including model access permission configuration.

## Requirements

### Requirement: Permission management page
The admin-ui SHALL provide a page at `/permissions` for managing model access permissions.

#### Scenario: List permissions
- **WHEN** the permissions page is loaded
- **THEN** it SHALL call `GET /admin/permissions` via gateway-service
- **AND** display all permissions in an antd Table with columns: group, model pattern, effect (allow/deny), created at

#### Scenario: Create permission
- **WHEN** the user fills in group, model pattern, and effect and submits
- **THEN** the admin-ui SHALL call `POST /admin/permissions`
- **AND** the new permission SHALL appear in the table

#### Scenario: Edit permission
- **WHEN** the user clicks edit on a permission
- **THEN** the admin-ui SHALL open a modal with the permission's current values
- **AND** on submit, call `PUT /admin/permissions/:id`

#### Scenario: Delete permission
- **WHEN** the user confirms deletion of a permission
- **THEN** the admin-ui SHALL call `DELETE /admin/permissions/:id`
- **AND** remove the permission from the table

### Requirement: Permission group selector
The admin-ui SHALL provide a group selector when creating or editing permissions.

#### Scenario: Group dropdown populated from groups API
- **WHEN** the user opens the permission form
- **THEN** the group selector SHALL be populated from `GET /admin/groups`
- **AND** display group names as options

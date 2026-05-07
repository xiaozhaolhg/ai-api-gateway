## Purpose

Delta spec for changes to `openspec/specs/admin-ui-permissions/spec.md`

## ADDED Requirements

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

#### Scenario: Group dropdown populated from groups API
- **WHEN** the user opens the permission form
- **THEN** the group selector SHALL be populated from `GET /admin/groups`
- **AND** display group names as options
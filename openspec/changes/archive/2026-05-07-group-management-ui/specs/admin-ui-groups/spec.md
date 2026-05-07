## Purpose

Delta spec for changes to `openspec/specs/admin-ui-groups/spec.md`

## ADDED Requirements

### Requirement: Group detail expandable row with Members and Permissions tabs
The admin-ui SHALL show group details in an expandable row within the Groups table.

#### Scenario: Expand group row
- **WHEN** the user clicks the expand icon on a group row
- **THEN** the admin-ui SHALL render an expanded area with two tabs: Members and Permissions
- **AND** each tab SHALL load its data independently via React Query

#### Scenario: View members list in expanded row
- **WHEN** the Members tab is active
- **THEN** the admin-ui SHALL display a table of group members with columns: name, email, role, remove action
- **AND** show an "Add Member" button that opens a user selector

#### Scenario: View group members via API
- **WHEN** the user expands a group row in the Groups table
- **THEN** the admin-ui SHALL call `GET /admin/groups/:id/members`
- **AND** display a nested table of members with columns: name, email, role, actions (remove)

#### Scenario: Add member to group via expandable row
- **WHEN** the user clicks "Add Member" in the expanded group row
- **THEN** the admin-ui SHALL show a user selector populated from `GET /admin/users`
- **AND** on selection, call `POST /admin/groups/:id/members` with the user ID
- **AND** the user SHALL appear in the group's member list

#### Scenario: Remove member from group via expandable row
- **WHEN** the user clicks "Remove" on a member row
- **THEN** the admin-ui SHALL call `DELETE /admin/groups/:id/members/:userId`
- **AND** the user SHALL be removed from the group's member list

#### Scenario: View group permissions in expanded row
- **WHEN** the user expands a group row
- **THEN** the admin-ui SHALL call `GET /admin/permissions?group_id=:id`
- **AND** display a nested table of permissions with columns: resource type, resource ID, action, effect

#### Scenario: Add permission from group row
- **WHEN** the user clicks "Add Permission" in the Permissions tab
- **THEN** the admin-ui SHALL open the permission form pre-filled with the group ID

### Requirement: Group search and filter
The admin-ui SHALL support searching and filtering groups.

#### Scenario: Search groups by name
- **WHEN** the user types in the search input
- **THEN** the groups table SHALL filter to show only groups whose name matches the search term

### Requirement: Group form enhanced fields
The admin-ui SHALL support additional group configuration fields.

#### Scenario: Create group with model patterns
- **WHEN** the user creates or edits a group
- **THEN** the form SHALL include a model patterns field (tags input for patterns like `gpt-*`, `ollama:*`)

#### Scenario: Create group with parent group
- **WHEN** the user creates or edits a group
- **THEN** the form SHALL include a parent group selector populated from existing groups
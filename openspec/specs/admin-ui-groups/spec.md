## Purpose

Group management interface for the admin UI, including group configuration and membership management.

## Requirements

### Requirement: Group management page
The admin-ui SHALL provide a page at `/groups` for managing user groups.

#### Scenario: List groups
- **WHEN** the groups page is loaded
- **THEN** it SHALL call `GET /admin/groups` via gateway-service
- **AND** display all groups in an antd Table with columns: name, description, member count, created at

#### Scenario: Create group
- **WHEN** the user fills in name and description and submits
- **THEN** the admin-ui SHALL call `POST /admin/groups`
- **AND** the new group SHALL appear in the table

#### Scenario: Edit group
- **WHEN** the user clicks edit on a group
- **THEN** the admin-ui SHALL open a modal with the group's current values
- **AND** on submit, call `PUT /admin/groups/:id`

#### Scenario: Delete group
- **WHEN** the user confirms deletion of a group
- **THEN** the admin-ui SHALL call `DELETE /admin/groups/:id`
- **AND** remove the group from the table

### Requirement: Group membership management
The admin-ui SHALL allow managing members within a group.

#### Scenario: View group members
- **WHEN** the user clicks on a group row or a "Members" action
- **THEN** the admin-ui SHALL display a list of users in the group
- **AND** show available users not yet in the group for adding

#### Scenario: Add member to group
- **WHEN** the user selects a user and adds them to the group
- **THEN** the admin-ui SHALL call `POST /admin/groups/:id/members` with the user ID
- **AND** the user SHALL appear in the group's member list

#### Scenario: Remove member from group
- **WHEN** the user removes a member from the group
- **THEN** the admin-ui SHALL call `DELETE /admin/groups/:id/members/:userId`
- **AND** the user SHALL be removed from the group's member list

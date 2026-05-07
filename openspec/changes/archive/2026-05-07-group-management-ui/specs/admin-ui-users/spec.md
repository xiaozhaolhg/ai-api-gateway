## Purpose

Delta spec for new `openspec/specs/admin-ui-users/spec.md`

## ADDED Requirements

### Requirement: User group assignment
The admin-ui SHALL allow assigning users to groups during create and edit.

#### Scenario: Assign groups on user creation
- **WHEN** the user creates a new user
- **THEN** the form SHALL include a "Groups" multi-select field populated from `GET /admin/groups`
- **AND** after user creation, call `POST /admin/groups/:id/members` for each selected group

#### Scenario: Assign groups on user edit
- **WHEN** the user edits an existing user
- **THEN** the form SHALL show a "Groups" multi-select with the user's current groups pre-selected
- **AND** on submit, add the user to newly selected groups and remove from deselected groups

#### Scenario: Display user groups in table
- **WHEN** the users table is rendered
- **THEN** each user row SHALL display their group names as tags

### Requirement: User search and filter
The admin-ui SHALL support searching and filtering users.

#### Scenario: Search users by name or email
- **WHEN** the user types in the search input
- **THEN** the users table SHALL filter to show only users whose name or email matches the search term

#### Scenario: Filter users by role
- **WHEN** the user selects a role from the filter dropdown
- **THEN** the users table SHALL show only users with that role

### Requirement: User password on creation
The admin-ui SHALL include a password field when creating a user.

#### Scenario: Create user with password
- **WHEN** the user creates a new user
- **THEN** the form SHALL include a password field (required, min 8 characters)
- **AND** the password SHALL be sent to `POST /admin/users` in the request body
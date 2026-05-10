## MODIFIED Requirements

### Requirement: User group assignment
The admin-ui SHALL allow assigning users to groups during create and edit.

#### Scenario: Assign groups on user creation
- **WHEN** user creates a new user
- **THEN** form SHALL include a "Groups" multi-select field populated from `GET /admin/groups`
- **AND** after user creation, call `POST /admin/groups/:id/members` for each selected group

#### Scenario: Assign groups on user edit
- **WHEN** user edits an existing user
- **THEN** form SHALL show a "Groups" multi-select with the user's current groups pre-selected
- **AND** on submit, add the user to newly selected groups and remove from deselected groups

#### Scenario: Display user groups in table
- **WHEN** users table is rendered
- **THEN** each user row SHALL display their group names as tags
- **AND** groups are displayed alphabetically

### Requirement: User creation with username support
The admin-ui SHALL include username field in user creation and edit forms.

#### Scenario: Create user with username
- **WHEN** user creates a new user
- **THEN** form SHALL include username field (optional for existing users, required for new users)
- **AND** form validates username uniqueness before submission
- **AND** username is sent to `POST /admin/users` in the request body

#### Scenario: Edit user with username
- **WHEN** user edits an existing user
- **THEN** form SHALL show current username in username field
- **AND** allow username modification with uniqueness validation
- **AND** updated username is sent to `PUT /admin/users/:id` in the request body

#### Scenario: Username validation
- **WHEN** user types in username field
- **THEN** system provides real-time feedback on username availability
- **AND** shows error if username is already taken or invalid

### Requirement: User search and filter
The admin-ui SHALL support searching and filtering users.

#### Scenario: Search users by name, email, or username
- **WHEN** user types in the search input
- **THEN** users table SHALL filter to show only users whose name, email, or username matches the search term

#### Scenario: Filter users by role
- **WHEN** user selects a role from the filter dropdown
- **THEN** users table SHALL show only users with that role

### Requirement: User password on creation
The admin-ui SHALL include a password field when creating a user.

#### Scenario: Create user with password
- **WHEN** user creates a new user
- **THEN** form SHALL include a password field (required, min 8 characters)
- **AND** password SHALL be sent to `POST /admin/users` in the request body

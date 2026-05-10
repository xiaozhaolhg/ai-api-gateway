## MODIFIED Requirements

### Requirement: Group creation with description
The admin-ui SHALL allow creating groups with names and descriptions.

#### Scenario: Create group with description
- **WHEN** user creates a new group with name and description
- **THEN** form SHALL include description field
- **AND** description SHALL be sent to `POST /admin/groups` in the request body
- **AND** response includes description field

#### Scenario: Create group without description
- **WHEN** user creates a new group with name but no description
- **THEN** form SHALL allow empty description field
- **AND** empty description SHALL be sent to `POST /admin/groups`
- **AND** response includes empty description field

### Requirement: Group editing with description
The admin-ui SHALL allow editing group descriptions.

#### Scenario: Edit group description
- **WHEN** user edits an existing group
- **THEN** form SHALL pre-fill description field with current value
- **AND** allow modification of description
- **AND** updated description SHALL be sent to `PUT /admin/groups/:id`
- **AND** response includes updated description field

#### Scenario: Partial group update
- **WHEN** user updates group name but not description
- **THEN** form SHALL preserve existing description
- **AND** description SHALL be sent unchanged to `PUT /admin/groups/:id`
- **AND** response includes unchanged description field

### Requirement: Group member count display
The admin-ui SHALL display member counts for each group.

#### Scenario: Groups list with member counts
- **WHEN** groups table is rendered
- **THEN** each group row SHALL display member count in dedicated column
- **AND** counts are calculated from group membership data
- **AND** counts update in real-time when members are added/removed

#### Scenario: Group detail with member count
- **WHEN** user views group details
- **THEN** system SHALL show total member count
- **AND** provide list of all members with their roles

#### Scenario: Real-time member count updates
- **WHEN** user is added to or removed from a group
- **THEN** group member count SHALL update immediately
- **AND** all connected clients see updated count without refresh

### Requirement: Group management interface
The admin-ui SHALL provide comprehensive group management interface.

#### Scenario: Group list with search
- **WHEN** user views groups page
- **THEN** system SHALL display searchable list of groups
- **AND** each group shows name, description, member count, and creation date

#### Scenario: Group creation form
- **WHEN** user clicks "Add Group"
- **THEN** system SHALL show form with name, description, parent group, and tier selection
- **AND** validate required fields before submission

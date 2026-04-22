## ADDED Requirements

### Requirement: Routing rules management page
The admin-ui SHALL provide a page at `/routing` for managing routing rules.

#### Scenario: List routing rules
- **WHEN** the routing rules page is loaded
- **THEN** it SHALL call `GET /admin/routing-rules` via gateway-service
- **AND** display all rules in an antd Table with columns: model pattern, provider, adapter type, priority, fallback chain, status

#### Scenario: Create routing rule
- **WHEN** the user fills in the routing rule form (model pattern, provider ID, adapter type, priority, fallback provider IDs) and submits
- **THEN** the admin-ui SHALL call `POST /admin/routing-rules`
- **AND** the new rule SHALL appear in the table

#### Scenario: Edit routing rule
- **WHEN** the user clicks edit on a routing rule
- **THEN** the admin-ui SHALL open a modal with the rule's current values
- **AND** on submit, call `PUT /admin/routing-rules/:id`
- **AND** update the rule in the table

#### Scenario: Delete routing rule
- **WHEN** the user confirms deletion of a routing rule
- **THEN** the admin-ui SHALL call `DELETE /admin/routing-rules/:id`
- **AND** remove the rule from the table

### Requirement: Fallback chain configuration
The admin-ui SHALL allow configuring fallback provider chains for routing rules.

#### Scenario: Fallback chain input
- **WHEN** the user is creating or editing a routing rule
- **THEN** the form SHALL provide an ordered list input for fallback provider IDs
- **AND** the user SHALL be able to add, remove, and reorder fallback entries

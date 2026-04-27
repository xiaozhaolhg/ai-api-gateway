## ADDED Requirements

### Requirement: Budget management page
The admin-ui SHALL provide a page at `/budgets` for managing budgets.

#### Scenario: List budgets
- **WHEN** the budgets page is loaded
- **THEN** it SHALL call `GET /admin/budgets` via gateway-service
- **AND** display all budgets in an antd Table with columns: name, scope (user/group/global), limit, current spend, remaining, status (active/exceeded), period

#### Scenario: Create budget
- **WHEN** the user fills in name, scope, limit amount, period, and optional user/group ID and submits
- **THEN** the admin-ui SHALL call `POST /admin/budgets`
- **AND** the new budget SHALL appear in the table

#### Scenario: Edit budget
- **WHEN** the user clicks edit on a budget
- **THEN** the admin-ui SHALL open a modal with the budget's current values
- **AND** on submit, call `PUT /admin/budgets/:id`

#### Scenario: Delete budget
- **WHEN** the user confirms deletion of a budget
- **THEN** the admin-ui SHALL call `DELETE /admin/budgets/:id`
- **AND** remove the budget from the table

### Requirement: Budget status indicators
The admin-ui SHALL display visual indicators for budget status.

#### Scenario: Budget within limit
- **WHEN** a budget's current spend is below its limit
- **THEN** the status SHALL display a green "Active" badge

#### Scenario: Budget soft cap exceeded
- **WHEN** a budget's soft cap is exceeded
- **THEN** the status SHALL display a yellow "Warning" badge

#### Scenario: Budget hard cap exceeded
- **WHEN** a budget's hard cap is exceeded
- **THEN** the status SHALL display a red "Exceeded" badge

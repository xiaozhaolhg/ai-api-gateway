## ADDED Requirements

### Requirement: Alert rules management page
The admin-ui SHALL provide a page at `/alerts` for managing alert rules and viewing active alerts.

#### Scenario: List alert rules
- **WHEN** the alerts page is loaded on the "Rules" tab
- **THEN** it SHALL call `GET /admin/alert-rules` via gateway-service
- **AND** display all rules in an antd Table with columns: name, metric, condition (threshold), channel, status (enabled/disabled), created at

#### Scenario: Create alert rule
- **WHEN** the user fills in name, metric, condition, threshold value, and notification channel and submits
- **THEN** the admin-ui SHALL call `POST /admin/alert-rules`
- **AND** the new rule SHALL appear in the table

#### Scenario: Edit alert rule
- **WHEN** the user clicks edit on an alert rule
- **THEN** the admin-ui SHALL open a modal with the rule's current values
- **AND** on submit, call `PUT /admin/alert-rules/:id`

#### Scenario: Delete alert rule
- **WHEN** the user confirms deletion of an alert rule
- **THEN** the admin-ui SHALL call `DELETE /admin/alert-rules/:id`
- **AND** remove the rule from the table

### Requirement: Active alerts view
The admin-ui SHALL provide a tab for viewing active (firing) alerts.

#### Scenario: List active alerts
- **WHEN** the user switches to the "Active Alerts" tab
- **THEN** the admin-ui SHALL call `GET /admin/alerts`
- **AND** display alerts in an antd Table with columns: alert rule, severity, status (firing/acknowledged/resolved), triggered at, description

#### Scenario: Acknowledge alert
- **WHEN** the user clicks "Acknowledge" on a firing alert
- **THEN** the admin-ui SHALL call `PUT /admin/alerts/:id/acknowledge`
- **AND** the alert status SHALL change to "acknowledged"

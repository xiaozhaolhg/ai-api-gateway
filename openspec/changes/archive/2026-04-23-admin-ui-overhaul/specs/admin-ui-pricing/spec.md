## ADDED Requirements

### Requirement: Pricing rules management page
The admin-ui SHALL provide a page at `/pricing` for managing pricing rules.

#### Scenario: List pricing rules
- **WHEN** the pricing rules page is loaded
- **THEN** it SHALL call `GET /admin/pricing-rules` via gateway-service
- **AND** display all rules in an antd Table with columns: model, provider, price per prompt token, price per completion token, currency, effective date

#### Scenario: Create pricing rule
- **WHEN** the user fills in model, provider, prices, and currency and submits
- **THEN** the admin-ui SHALL call `POST /admin/pricing-rules`
- **AND** the new rule SHALL appear in the table

#### Scenario: Edit pricing rule
- **WHEN** the user clicks edit on a pricing rule
- **THEN** the admin-ui SHALL open a modal with the rule's current values
- **AND** on submit, call `PUT /admin/pricing-rules/:id`

#### Scenario: Delete pricing rule
- **WHEN** the user confirms deletion of a pricing rule
- **THEN** the admin-ui SHALL call `DELETE /admin/pricing-rules/:id`
- **AND** remove the rule from the table

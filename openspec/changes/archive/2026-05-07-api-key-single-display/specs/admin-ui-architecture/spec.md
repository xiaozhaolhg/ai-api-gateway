## MODIFIED Requirements

### Requirement: Issue API key
The admin-ui SHALL call POST /admin/api-keys and display the generated key with single-display security behavior.

#### Scenario: Issue API key with ephemeral display
- **WHEN** the user clicks "Issue API Key" for a user
- **THEN** the admin-ui SHALL call POST /admin/api-keys
- **AND** the generated key SHALL be displayed in an antd Alert with a copy button
- **AND** a prominent warning SHALL be displayed: "This API key will only be shown once. Please copy it now and store it securely."
- **AND** the key SHALL be stored only in component state (not localStorage or sessionStorage)

#### Scenario: API key disappears on navigation
- **WHEN** the user navigates away from the API keys page after key creation
- **THEN** the API key SHALL be immediately cleared from component state
- **AND** navigating back to the page SHALL NOT display the key again

#### Scenario: API key disappears on modal close
- **WHEN** the user closes the API key display modal or alert
- **THEN** the API key SHALL be immediately cleared from component state
- **AND** a "key-dismissed" flag SHALL be set in sessionStorage to prevent re-display

#### Scenario: API key not re-displayed after dismissal
- **WHEN** the user returns to the API keys page after previously dismissing a key
- **THEN** the page SHALL check the "key-dismissed" flag in sessionStorage
- **AND** SHALL NOT display the previously generated key
- **AND** SHALL show a message "Previous API key was shown once. Generate a new key if needed."

#### Scenario: Copy to clipboard functionality
- **WHEN** the user clicks the copy button next to the API key
- **THEN** the key SHALL be copied to the clipboard
- **AND** a success message "API key copied to clipboard" SHALL be displayed

#### Scenario: Revoke API key
- **WHEN** the user confirms revocation via antd Popconfirm
- **THEN** the admin-ui SHALL call DELETE /admin/api-keys/:id
- **AND** the key SHALL be removed from the list
- **AND** any displayed key SHALL be cleared from component state

#### Scenario: API key user selector populated from users API
- **WHEN** the API keys page is loaded
- **THEN** the user selector SHALL be populated from GET /admin/users
- **AND** display user names as options (NOT hardcoded values)

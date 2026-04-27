## ADDED Requirements

### Requirement: Internationalization framework
The admin-ui SHALL support Chinese and English languages using react-i18next and antd locale.

#### Scenario: Language initialization
- **WHEN** the admin-ui starts
- **THEN** it SHALL detect the browser language or use the stored preference
- **AND** configure react-i18next with `en` and `zh` namespaces
- **AND** set antd `ConfigProvider` locale to match the selected language

#### Scenario: Language switching
- **WHEN** the user selects a different language from the language switcher
- **THEN** the admin-ui SHALL update all UI text to the selected language
- **AND** update antd component locale (DatePicker, Pagination, Empty, etc.)
- **AND** persist the language preference to localStorage

### Requirement: Translation key structure
The admin-ui SHALL organize translation keys by page namespace.

#### Scenario: Namespace per page
- **WHEN** the translation files are inspected
- **THEN** each page SHALL have its own namespace (e.g., `dashboard`, `providers`, `routing`)
- **AND** common keys (navigation, actions, status labels) SHALL be in a `common` namespace

### Requirement: Language switcher component
The admin-ui SHALL provide a language switcher in the header area.

#### Scenario: Switcher display
- **WHEN** the header is rendered
- **THEN** a language toggle SHALL be visible showing the current language
- **AND** clicking it SHALL switch between English and Chinese

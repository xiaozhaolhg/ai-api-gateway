## MODIFIED Requirements

### Requirement: Internationalization framework
The admin-ui SHALL support Chinese and English languages using react-i18next and antd locale, with initial focus on English translation resolution.

#### Scenario: Language initialization
- **WHEN** the admin-ui starts
- **THEN** it SHALL detect the browser language or use the stored preference
- **AND** configure react-i18next with `en` and `zh` namespaces
- **AND** set antd `ConfigProvider` locale to match the selected language
- **AND** ensure all translation keys are properly resolved to text (not showing as stub keys)

#### Scenario: Language switching
- **WHEN** the user selects a different language from the language switcher
- **THEN** the admin-ui SHALL update all UI text to the selected language
- **AND** update antd component locale (DatePicker, Pagination, Empty, etc.)
- **AND** persist the language preference to localStorage

#### Scenario: Translation key resolution
- **WHEN** a component renders with a translation key like "dashboard.title"
- **THEN** the system SHALL resolve it to meaningful English text
- **AND** SHALL NOT display the raw translation key to users
- **AND** SHALL provide fallback text for missing keys

### Requirement: Translation key structure
The admin-ui SHALL organize translation keys by page namespace.

#### Scenario: Namespace per page
- **WHEN** the translation files are inspected
- **THEN** each page SHALL have its own namespace (e.g., `dashboard`, `providers`, `routing`)
- **AND** common keys (navigation, actions, status labels) SHALL be in a `common` namespace

#### Scenario: Stub key identification
- **WHEN** translation keys are not resolved
- **THEN** the system SHALL log the missing keys for debugging
- **AND** SHALL display a fallback text instead of the raw key

### Requirement: Language switcher component
The admin-ui SHALL provide a language switcher in the header area.

#### Scenario: Switcher display
- **WHEN** the header is rendered
- **THEN** a language toggle SHALL be visible showing the current language
- **AND** clicking it SHALL switch between English and Chinese

## ADDED Requirements

### Requirement: English translation completeness
The admin-ui SHALL have complete English translations for all UI elements currently showing stub keys.

#### Scenario: Dashboard translations
- **WHEN** the dashboard page is rendered
- **THEN** all translation keys like "dashboard.title" SHALL resolve to English text
- **AND** navigation elements SHALL have proper English labels
- **AND** action buttons SHALL display English text

#### Scenario: Provider management translations
- **WHEN** the providers page is rendered
- **THEN** all provider-related translation keys SHALL resolve to English text
- **AND** status indicators SHALL have English labels
- **AND** form fields SHALL show English placeholders and labels

#### Scenario: Common UI elements translations
- **WHEN** any page is rendered
- **THEN** common elements like navigation, buttons, and messages SHALL be in English
- **AND** error messages SHALL be properly translated
- **AND** success messages SHALL be in English

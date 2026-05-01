# admin-ui-mock-api

## Purpose

Development-only Mock API system for admin-ui, enabling frontend development without running backend services.

## Scope

- **In Scope**: admin-ui only, no backend service changes
- **Out of Scope**: Production use (DevTools component only renders in development mode)

## Requirements

### Requirement: Mock/Real Mode Switching

The admin-ui SHALL support switching between Mock and Real API modes via environment variables or DevTools at runtime.

#### Scenario: Enable via environment variable
- **WHEN** `VITE_USE_MOCK=true` is set in `.env.development`
- **THEN** the UnifiedAPIClient SHALL use MockAPIClient for all API calls

#### Scenario: Switch via DevTools
- **WHEN** the user toggles the mode switch in the DevTools panel
- **THEN** the UnifiedAPIClient SHALL switch between MockAPIClient and RealAPIClient
- **AND** reload the page to apply changes

#### Scenario: Configuration precedence
- **WHEN** both env var and runtime switch are used
- **THEN** the runtime switch SHALL take precedence after page load

### Requirement: Data Persistence

Mock data SHALL persist in localStorage across page refreshes.

#### Scenario: Data survives refresh
- **WHEN** the page is refreshed
- **THEN** the MockDataHandler SHALL reload data from localStorage
- **AND** all previous changes SHALL be preserved

#### Scenario: Data isolation
- **WHEN** multiple browser tabs are open
- **THEN** they SHALL share the same localStorage data
- **AND** changes in one tab SHALL be visible in others after refresh

### Requirement: Network Simulation

Mock API SHALL support configurable network delay via `VITE_MOCK_DELAY`.

#### Scenario: Configurable delay
- **WHEN** `VITE_MOCK_DELAY=1000` is set
- **THEN** all mock API responses SHALL be delayed by approximately 1000ms

#### Scenario: Zero delay
- **WHEN** `VITE_MOCK_DELAY=0` is set
- **THEN** mock API responses SHALL return immediately with no delay

#### Scenario: Runtime delay adjustment
- **WHEN** the user adjusts the delay in DevTools
- **THEN** subsequent API calls SHALL use the new delay value

### Requirement: Complete API Coverage

The MockAPIClient SHALL implement the full APIClientInterface.

#### Scenario: All methods implemented
- **WHEN** the MockAPIClient is inspected
- **THEN** it SHALL implement all methods defined in APIClientInterface
- **AND** each method SHALL return realistic mock data

#### Scenario: CRUD operations
- **WHEN** any CRUD operation is performed in mock mode
- **THEN** the MockDataHandler SHALL persist the change to localStorage
- **AND** subsequent queries SHALL reflect the change

### Requirement: DevTools Panel

The admin-ui SHALL provide a DevTools panel for mock API configuration.

#### Scenario: DevTools visibility
- **WHEN** the app is running in development mode
- **THEN** the DevTools button SHALL be visible in the bottom-right corner

#### Scenario: DevTools hidden in production
- **WHEN** the app is built for production
- **THEN** the DevTools component SHALL NOT be rendered

#### Scenario: Data management
- **WHEN** the user clicks "Reset to Defaults"
- **THEN** all mock data SHALL be reset to the default dataset

#### Scenario: Data export
- **WHEN** the user clicks "Export Mock Data"
- **THEN** a JSON file SHALL be downloaded with all current mock data

#### Scenario: Data import
- **WHEN** the user uploads a valid mock data JSON file
- **THEN** the mock data SHALL be replaced with the imported data
- **AND** the page SHALL reload to reflect changes

### Requirement: Unit Test Coverage

The MockAPIClient SHALL have comprehensive unit tests.

#### Scenario: Test coverage
- **WHEN** the test suite is run
- **THEN** all MockAPIClient methods SHALL be tested
- **AND** edge cases (not found, duplicate, invalid input) SHALL be covered

#### Scenario: Test environment
- **WHEN** tests are run
- **THEN** the MockAPIClient SHALL use 0 delay for fast execution
- **AND** each test SHALL start with fresh default data

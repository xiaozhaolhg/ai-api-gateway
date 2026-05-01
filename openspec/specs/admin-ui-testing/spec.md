# admin-ui-testing

## Purpose

Testing standards and requirements for admin-ui, including unit tests, integration tests, and E2E tests.

## Requirements

### Requirement: Test Coverage

The admin-ui SHALL maintain comprehensive test coverage for all critical paths.

#### Scenario: Auth context tests
- **WHEN** AuthContext is tested
- **THEN** login, logout, session restoration, and role checks SHALL be verified

#### Scenario: Page integration tests
- **WHEN** any page component is tested
- **THEN** CRUD operations, error states, loading states, and role-based access SHALL be covered

#### Scenario: Minimum coverage threshold
- **WHEN** test coverage is measured
- **THEN** minimum 80% line coverage SHALL be maintained

### Requirement: Test Utilities

The admin-ui SHALL provide reusable test utilities.

#### Scenario: Mock API client
- **WHEN** tests need API mocking
- **THEN** a mock API client SHALL be available in `src/test/mocks.ts`

#### Scenario: Role mocking
- **WHEN** tests need to simulate different roles
- **THEN** a role mock utility SHALL be available

### Requirement: Error Boundary Testing

The admin-ui SHALL test error boundary behavior.

#### Scenario: Query error handling
- **WHEN** a TanStack Query fails
- **THEN** the QueryErrorBoundary SHALL catch and display retry option

#### Scenario: Error recovery
- **WHEN** user clicks retry button
- **THEN** the query SHALL be refetched and component re-rendered

### Requirement: Test Environment

The admin-ui SHALL have a properly configured test environment.

#### Scenario: Vitest configuration
- **WHEN** running `npm run test`
- **THEN** Vitest SHALL execute all test files matching `*.test.ts(x)` pattern

#### Scenario: Test setup
- **WHEN** tests run
- **THEN** mock API client SHALL be configured as default
- **AND** auth context SHALL be mockable with different roles

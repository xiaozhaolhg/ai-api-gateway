## MODIFIED Requirements

### Requirement: Modern component library integration
The admin UI SHALL use shadcn/ui components with Lucide icons.

#### Scenario: Component usage
- **WHEN** building UI elements
- **THEN** use shadcn/ui Button, Input, Table, Card, Dialog, Form components

#### Scenario: Icon integration
- **WHEN** displaying icons
- **THEN** use Lucide React icons consistent with shadcn/ui

#### Scenario: Styling consistency
- **WHEN** applying styles
- **THEN** follow shadcn/ui design tokens and Tailwind CSS classes

### Requirement: Collapsible sidebar navigation
The admin UI SHALL have a collapsible sidebar with icon-only mode.

#### Scenario: Sidebar toggle
- **WHEN** user clicks collapse button
- **THEN** sidebar collapses to icon-only view, expanding main content area

#### Scenario: Active state indication
- **WHEN** navigation tab is active
- **THEN** highlight with accent color and background

#### Scenario: Responsive behavior
- **WHEN** screen width is limited
- **THEN** sidebar automatically collapses to icon-only mode

### Requirement: Form handling and validation
The admin UI SHALL use react-hook-form with shadcn/ui Form components.

#### Scenario: Form submission
- **WHEN** user submits create/edit forms
- **THEN** validate with react-hook-form and display errors inline

#### Scenario: Form reset
- **WHEN** form is cancelled or successfully submitted
- **THEN** reset form state and close modal/drawer

#### Scenario: Field validation
- **WHEN** user enters invalid data
- **THEN** show real-time validation feedback with proper error messages

### Requirement: Data fetching and caching
The admin UI SHALL use TanStack Query for API data management.

#### Scenario: Data loading
- **WHEN** page loads
- **THEN** TanStack Query fetches data with loading states

#### Scenario: Data caching
- **WHEN** navigating between pages
- **THEN** cached data is used when fresh, with background refetch

#### Scenario: Error handling
- **WHEN** API requests fail
- **THEN** TanStack Query provides error states with retry options

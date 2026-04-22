## 1. Auth Service Backend Changes

- [ ] 1.1 Add Login RPC to auth-service proto definition
- [ ] 1.2 Implement password hashing in auth-service (bcrypt)
- [ ] 1.3 Add Login handler in auth-service with JWT generation
- [ ] 1.4 Update User entity to support password field
- [ ] 1.5 Extend role validation to support admin/user/viewer
- [ ] 1.6 Add password update/reset functionality
- [ ] 1.7 Write unit tests for auth-service login flow

## 2. Gateway Service Auth Endpoints

- [ ] 2.1 Add POST /admin/login endpoint in gateway-service
- [ ] 2.2 Add POST /admin/logout endpoint in gateway-service
- [ ] 2.3 Implement auth middleware for admin routes
- [ ] 2.4 Configure JWT cookie handling (secure, HTTP-only, /admin path)
- [ ] 2.5 Add user context propagation to downstream services
- [ ] 2.6 Update admin API handlers to check user role
- [ ] 2.7 Write integration tests for auth flow

## 3. Admin UI Dependencies and Setup

- [ ] 3.1 Install shadcn/ui CLI and initialize components
- [ ] 3.2 Install Lucide React icons
- [ ] 3.3 Install TanStack Query for data fetching
- [ ] 3.4 Install react-hook-form for form handling
- [ ] 3.5 Update package.json with new dependencies
- [ ] 3.6 Configure Tailwind CSS for shadcn/ui components

## 4. Auth Context and Routing

- [ ] 4.1 Create AuthContext with user state and login/logout functions
- [ ] 4.2 Create ProtectedRoute wrapper for route guards
- [ ] 4.3 Update App.tsx with auth-aware routing
- [ ] 4.4 Implement /admin redirect logic (login vs dashboard)
- [ ] 4.5 Add auth context persistence and restoration
- [ ] 4.6 Create auth hooks for role-based access

## 5. Login Page Implementation

- [ ] 5.1 Create Login page component with shadcn/ui form
- [ ] 5.2 Implement email/password form with validation
- [ ] 5.3 Add login API call with TanStack Query
- [ ] 5.4 Handle login success/error states
- [ ] 5.5 Add loading states and error messages
- [ ] 5.6 Style login page with proper layout and branding

## 6. Dashboard Page

- [ ] 6.1 Create Dashboard page component
- [ ] 6.2 Implement summary cards for key metrics
- [ ] 6.3 Add quick action buttons for common tasks
- [ ] 6.4 Integrate with multiple backend services
- [ ] 6.5 Add charts or visual indicators for system health
- [ ] 6.6 Implement data refresh and real-time updates

## 7. Layout and Navigation Redesign

- [ ] 7.1 Redesign Layout component with shadcn/ui
- [ ] 7.2 Implement collapsible sidebar with toggle
- [ ] 7.3 Add Lucide icons to navigation items
- [ ] 7.4 Implement role-based navigation filtering
- [ ] 7.5 Add active state indicators
- [ ] 7.6 Add user profile dropdown with logout

## 8. Page Redesigns with New Components

- [ ] 8.1 Redesign Providers page with shadcn/ui Table and Dialog
- [ ] 8.2 Redesign Users page with role badges and status toggles
- [ ] 8.3 Redesign API Keys page with copy-once display
- [ ] 8.4 Redesign Usage page with filters and charts
- [ ] 8.5 Redesign Health page with status cards
- [ ] 8.6 Add Settings page with profile management

## 9. Role-Based Access Control

- [ ] 9.1 Implement role checks in page components
- [ ] 9.2 Add data filtering based on user ownership
- [ ] 9.3 Disable edit/delete actions for non-admin users
- [ ] 9.4 Implement viewer read-only states
- [ ] 9.5 Add access denied messages and redirects
- [ ] 9.6 Test role-based access across all pages

## 10. Data Fetching and State Management

- [ ] 10.1 Replace useState/useEffect with TanStack Query
- [ ] 10.2 Add proper error handling and retry logic
- [ ] 10.3 Implement optimistic updates for CRUD operations
- [ ] 10.4 Add loading skeletons and states
- [ ] 10.5 Configure query caching and invalidation
- [ ] 10.6 Add query error boundaries

## 11. Form Handling and Validation

- [ ] 11.1 Replace inline forms with react-hook-form
- [ ] 11.2 Add proper validation schemas
- [ ] 11.3 Implement form modals/drawers
- [ ] 11.4 Add form reset and cancel handling
- [ ] 11.5 Show inline validation errors
- [ ] 11.6 Add form submission loading states

## 12. Testing and Quality Assurance

- [ ] 12.1 Write unit tests for auth context and hooks
- [ ] 12.2 Write integration tests for login flow
- [ ] 12.3 Test role-based access control
- [ ] 12.4 Test responsive sidebar behavior
- [ ] 12.5 Test form validation and error states
- [ ] 12.6 Perform end-to-end testing of complete admin flow

## 13. Documentation and Cleanup

- [ ] 13.1 Update admin-ui README with new features
- [ ] 13.2 Document role permissions matrix
- [ ] 13.3 Document API changes and new endpoints
- [ ] 13.4 Remove old unused components and styles
- [ ] 13.5 Update deployment documentation
- [ ] 13.6 Verify all TODO comments are addressed

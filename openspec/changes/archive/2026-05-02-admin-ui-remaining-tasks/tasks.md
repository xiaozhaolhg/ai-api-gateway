# Tasks: Complete Admin-UI Remaining Work

## Phase 9: Role-Based Access Control

- [ ] **Task 9.6**: Test role-based access across all pages
  - [ ] Create test utilities for role mocking
  - [ ] Test Providers page (admin CRUD, user/viewer blocked)
  - [ ] Test Users page (admin CRUD, user own only, viewer blocked)
  - [ ] Test APIKeys page (admin CRUD, user own only)
  - [ ] Test Usage page (admin all, user/viewer own only)
  - [ ] Test RoutingRules page (admin CRUD, others blocked)
  - [ ] Test Groups page (admin CRUD, others blocked)
  - [ ] Test Permissions page (admin CRUD, others blocked)
  - [ ] Test Budgets page (admin CRUD, others blocked)
  - [ ] Test PricingRules page (admin CRUD, others blocked)
  - [ ] Test Alerts page (admin CRUD+ack, user ack, viewer read-only)
  - [ ] Test Health page (all roles read-only)

## Phase 10: Data Fetching & State Management

- [ ] **Task 10.6**: Add query error boundaries
  - [ ] Create QueryErrorBoundary component
  - [ ] Integrate with TanStack Query error handling
  - [ ] Add retry button functionality
  - [ ] Wrap all page routes in App.tsx
  - [ ] Test error boundary rendering
  - [ ] Test retry functionality

## Phase 12: Testing & Quality Assurance

- [ ] **Task 12.1**: Write unit tests for auth context and hooks
  - [ ] Test AuthContext login/logout/session restore
  - [ ] Test ProtectedRoute component
  - [ ] Test all three roles (admin/user/viewer)
  - [ ] Test login flow (unit)
  - [ ] Test logout flow (unit)
  - [ ] Test error states

- [ ] **Task 12.2**: Write integration tests for all pages
  - [ ] Dashboard page integration test
  - [ ] Providers page integration test
  - [ ] Users page integration test
  - [ ] APIKeys page integration test
  - [ ] Usage page integration test
  - [ ] RoutingRules page integration test
  - [ ] Groups page integration test
  - [ ] Permissions page integration test
  - [ ] Budgets page integration test
  - [ ] PricingRules page integration test
  - [ ] Alerts page integration test
  - [ ] Health page integration test

- [ ] **Task 12.3**: Test mobile collapse (E2E - blocked by browser deps)

- [ ] **Task 12.4**: Test desktop expand (E2E - blocked by browser deps)

- [ ] **Task 12.5**: Test all validation (unit tests)
  - [ ] Login form validation
  - [ ] Register form validation
  - [ ] Provider form validation
  - [ ] User form validation
  - [ ] Group form validation
  - [ ] Budget form validation
  - [ ] Pricing rule form validation
  - [ ] Alert rule form validation

- [ ] **Task 12.6**: Test complete user flow (integration)
  - [ ] Admin complete flow test
  - [ ] User limited flow test
  - [ ] Viewer read-only flow test

## Phase 13: Documentation & Cleanup

- [ ] **Task 13.4**: Remove old unused components and styles
  - [ ] Audit for orphaned components
  - [ ] Remove unused imports
  - [ ] Remove leftover template files
  - [ ] Verify build size optimized

- [ ] **Task 13.5**: Update deployment documentation
  - [ ] Update README with deploy instructions
  - [ ] Document all environment variables
  - [ ] Add Docker deployment notes

- [ ] **Task 13.6**: Verify all TODO comments are addressed
  - [ ] Search for TODO/FIXME comments
  - [ ] Resolve or document each TODO
  - [ ] Ensure no placeholders remain

## Summary

| Phase | Tasks | Completed | Remaining |
|-------|-------|-----------|----------|
| Phase 9 | 6 | 5 | **1** |
| Phase 10 | 6 | 5 | **1** |
| Phase 12 | 6 | 0 | **6** |
| Phase 13 | 6 | 2 | **4** |
| **Total** | **81** | **69** | **12** |

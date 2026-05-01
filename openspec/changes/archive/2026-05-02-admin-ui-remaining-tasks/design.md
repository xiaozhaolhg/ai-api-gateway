# Design: Remaining Tasks Implementation

## Phase 9: Role-Based Access Control Verification (Task 9.6)

### Approach
- Create test utilities for role-based testing
- Test each page component with mock auth context (admin/user/viewer)
- Verify edit/delete buttons disabled for non-admin
- Verify data filtering for user/viewer roles

### Test Matrix
| Page | Admin | User | Viewer |
|-------|-------|------|--------|
| /admin/dashboard | Full access | Limited data | Read-only |
| /admin/providers | CRUD | ✗ | ✗ |
| /admin/users | CRUD | Own only | ✗ |
| /admin/api-keys | CRUD | Own only | ✗ |
| /admin/usage | Full | Own only | Own only |
| /admin/routing-rules | CRUD | ✗ | ✗ |
| /admin/groups | CRUD | ✗ | ✗ |
| /admin/permissions | CRUD | ✗ | ✗ |
| /admin/budgets | CRUD | ✗ | ✗ |
| /admin/pricing-rules | CRUD | ✗ | ✗ |
| /admin/alerts | CRUD + acknowledge | Read + acknowledge | Read-only |
| /admin/health | Read | Read | Read |

## Phase 10: Query Error Boundaries (Task 10.6)

### Component Design
```typescript
// src/components/QueryErrorBoundary.tsx
import { QueryErrorResetBoundary } from '@tanstack/react-query';
import { Button, Result } from 'antd';

export const QueryErrorBoundary: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <QueryErrorResetBoundary>
      {({ reset }) => (
        <ErrorBoundary
          fallbackRender={({ error }) => (
            <Result
              status="error"
              title="Data Loading Failed"
              subTitle={error.message}
              extra={
                <Button onClick={() => reset()}>
                  Retry
                </Button>
              }
            />
          )}
        >
          {children}
        </ErrorBoundary>
      )}
    </QueryErrorResetBoundary>
  );
};
```

### Integration Point
Wrap all page components with `QueryErrorBoundary` in `App.tsx` routes.

## Phase 12: Testing & QA (6 tasks)

### Task 12.1: Auth Context and Hooks Tests

**Test Files:**
- `src/contexts/__tests__/AuthContext.test.tsx`
- `src/components/__tests__/ProtectedRoute.test.tsx`

**Test Cases:**
- Login flow: valid credentials → sets auth state
- Login flow: invalid credentials → shows error
- Logout flow: clears auth state + localStorage
- Session restoration: valid token → restores state
- Session restoration: invalid token → clears state
- Protected route: unauthenticated → redirects to /login
- Protected route: authenticated → renders children
- Role check: admin → all tabs visible
- Role check: user → limited tabs visible
- Role check: viewer → read-only tabs visible

### Tasks 12.2-12.5: Page Integration Tests

**Test Utilities:**
- `src/test/utils.tsx` — mock providers, wrapper components
- `src/test/mocks.ts` — mock API client responses

**Test Coverage per Page:**
| Page | CRUD Tests | Error States | Loading States | Role-Based |
|------|-----------|-------------|---------------|-------------|
| Dashboard | Summary cards load | API failure | Spin render | Data filtering |
| Providers | Create/Edit/Delete | Duplicate name | Empty state | Admin only |
| Users | Create/Edit/Delete | Invalid email | Empty state | Admin only |
| APIKeys | Create/Revoke | Invalid user | Empty state | Own only |
| Usage | Filter by user/date | No results | Loading | Own only |
| RoutingRules | Create/Edit/Delete | Invalid pattern | Empty state | Admin only |
| Groups | Create/Edit/Delete/AddRemoveMember | Invalid name | Empty state | Admin only |
| Permissions | Create/Edit/Delete | Invalid model | Empty state | Admin only |
| Budgets | Create/Edit/Delete | Exceeded limit | Empty state | Admin only |
| PricingRules | Create/Edit/Delete | Invalid price | Empty state | Admin only |
| Alerts | Create/Edit/Delete/Acknowledge | Invalid rule | Empty state | Role-based |
| Health | Auto-refresh | Service down | Loading | All roles |

### Task 12.6: Admin Flow Test

**E2E Flow Test (src/test/e2e/admin-flow.test.ts):**
1. Login as admin → dashboard visible
2. Create provider → appears in table
3. Create user → appears in table
4. Issue API key → key displayed
5. Create group → add user → verify membership
6. Set budget → verify tracking
7. Create alert rule → trigger → acknowledge
8. Logout → redirected to login

## Phase 13: Documentation & Cleanup (4 tasks)

### Task 13.4: Remove Unused Components

**Audit Checklist:**
- [ ] Search for orphaned components (not imported anywhere)
- [ ] Remove leftover Vite template artifacts
- [ ] Check bundle size before/after cleanup
- [ ] Verify no TypeScript errors after removal

### Task 13.5: Update Deployment Documentation

**Update Files:**
- `admin-ui/README.md` — add environment variables table
- `openspec/specs/admin-ui-deployment/spec.md` — update deployment steps

### Task 13.6: Address TODO Comments

**Search Command:**
```bash
grep -r "TODO" admin-ui/src/ --include="*.tsx" --include="*.ts"
grep -r "FIXME" admin-ui/src/ --include="*.tsx" --include="*.ts"
```

**Resolution Actions:**
- Implement missing functionality → remove TODO
- Document as known limitation → convert to comment
- Remove if no longer relevant

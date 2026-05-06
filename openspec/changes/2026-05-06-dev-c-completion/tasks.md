# Tasks: Developer C Completion

## Phase 1: BillingClient Extension

- [ ] **Task 1.1**: Add Budget CRUD methods to BillingClient
  - [ ] `ListBudgets(ctx, page, pageSize)` → `billingv1.ListBudgetsRequest`
  - [ ] `CreateBudget(ctx, userID, limit, period, softCapPct, hardCapPct)` → `billingv1.CreateBudgetRequest`
  - [ ] `UpdateBudget(ctx, id, userID, limit, period, softCapPct, hardCapPct, status)` → `billingv1.UpdateBudgetRequest`
  - [ ] `DeleteBudget(ctx, id)` → `billingv1.DeleteBudgetRequest`
  - **Acceptance**: Unit test for each method with mock gRPC server

- [ ] **Task 1.2**: Add PricingRule CRUD methods to BillingClient
  - [ ] `ListPricingRules(ctx, page, pageSize)` → `billingv1.ListPricingRulesRequest`
  - [ ] `CreatePricingRule(ctx, model, providerID, promptPrice, completionPrice, currency)` → `billingv1.CreatePricingRuleRequest`
  - [ ] `UpdatePricingRule(ctx, id, model, providerID, promptPrice, completionPrice, currency)` → `billingv1.UpdatePricingRuleRequest`
  - [ ] `DeletePricingRule(ctx, id)` → `billingv1.DeletePricingRuleRequest`
  - **Acceptance**: Unit test for each method with mock gRPC server

## Phase 2: MonitorClient

- [ ] **Task 2.1**: Create MonitorClient with AlertRule/Alert methods
  - [ ] `ListAlertRules(ctx, page, pageSize)` → `monitorv1.ListAlertRulesRequest`
  - [ ] `CreateAlertRule(ctx, metricType, condition, threshold, channel, channelConfig)` → `monitorv1.CreateAlertRuleRequest`
  - [ ] `UpdateAlertRule(ctx, id, metricType, condition, threshold, channel, channelConfig, status)` → `monitorv1.UpdateAlertRuleRequest`
  - [ ] `DeleteAlertRule(ctx, id)` → `monitorv1.DeleteAlertRuleRequest`
  - [ ] `GetAlerts(ctx, ruleID, status, page, pageSize)` → `monitorv1.GetAlertsRequest`
  - [ ] `AcknowledgeAlert(ctx, id)` → `monitorv1.AcknowledgeAlertRequest`
  - **Acceptance**: Unit test for each method; graceful fallback when monitor-service unavailable

## Phase 3: Admin REST Handlers

- [ ] **Task 3.1**: Create admin_budgets.go handler
  - [ ] `ListBudgets` — GET /admin/budgets → BillingClient.ListBudgets
  - [ ] `CreateBudget` — POST /admin/budgets → BillingClient.CreateBudget
  - [ ] `UpdateBudget` — PUT /admin/budgets/:id → BillingClient.UpdateBudget
  - [ ] `DeleteBudget` — DELETE /admin/budgets/:id → BillingClient.DeleteBudget
  - [ ] Map proto response to UI-compatible JSON (add name, scope, scope_id, current_spend, created_at, updated_at)
  - **Acceptance**: Integration test: create budget → list → update → delete

- [ ] **Task 3.2**: Create admin_pricing_rules.go handler
  - [ ] `ListPricingRules` — GET /admin/pricing-rules → BillingClient.ListPricingRules
  - [ ] `CreatePricingRule` — POST /admin/pricing-rules → BillingClient.CreatePricingRule
  - [ ] `UpdatePricingRule` — PUT /admin/pricing-rules/:id → BillingClient.UpdatePricingRule
  - [ ] `DeletePricingRule` — DELETE /admin/pricing-rules/:id → BillingClient.DeletePricingRule
  - [ ] Map proto response to UI-compatible JSON (add effective_date, created_at, updated_at)
  - **Acceptance**: Integration test: create rule → list → update → delete

- [ ] **Task 3.3**: Refactor admin_alerts.go to use MonitorClient
  - [ ] Replace in-memory `sync.RWMutex` + `[]AlertRule` / `[]Alert` with MonitorClient gRPC calls
  - [ ] Remove `initializeDefaultData()` mock
  - [ ] Keep graceful fallback: return empty list when monitor-service unavailable
  - **Acceptance**: Alert CRUD endpoints proxy to monitor-service; fallback works when service down

- [ ] **Task 3.4**: Wire new handlers in main.go
  - [ ] Initialize MonitorClient with lazy connection
  - [ ] Register `/admin/budgets` routes
  - [ ] Register `/admin/pricing-rules` routes
  - [ ] Update `/admin/alert-rules` and `/admin/alerts` to use refactored handler
  - **Acceptance**: All new endpoints accessible and return data from gRPC services

## Phase 4: UI Embedding & Single Binary

- [ ] **Task 4.1**: Add go:embed for admin-ui static files
  - [ ] Create `gateway-service/static/` directory placeholder
  - [ ] Add `//go:embed static` directive in main.go or separate file
  - [ ] Implement static file serving with `http.FS`
  - [ ] Implement SPA fallback: non-API routes serve index.html
  - [ ] Skip embedding if static/ directory is empty (dev mode)
  - **Acceptance**: Gateway serves admin-ui HTML/JS/CSS; API routes still work

- [ ] **Task 4.2**: Create single-binary build script
  - [ ] Add `build-single` target to root Makefile
  - [ ] Step 1: `cd admin-ui && npm ci && npm run build`
  - [ ] Step 2: Copy `admin-ui/dist/` to `gateway-service/static/`
  - [ ] Step 3: `cd gateway-service && go build -o bin/gateway ./cmd/server`
  - [ ] Add `clean-static` target to remove embedded files
  - **Acceptance**: `make build-single` produces working single binary

## Phase 5: Testing & QA

- [ ] **Task 5.1**: Write unit tests for auth context and hooks
  - [ ] Test AuthContext login/logout/session restore
  - [ ] Test ProtectedRoute component (admin/user/viewer)
  - **Acceptance**: All auth tests pass

- [ ] **Task 5.2**: Write handler unit tests for new endpoints
  - [ ] admin_budgets_test.go
  - [ ] admin_pricing_rules_test.go
  - [ ] admin_alerts_test.go (updated)
  - **Acceptance**: >60% coverage for new handler code

- [ ] **Task 5.3**: Integrate QueryErrorBoundary in App.tsx
  - [ ] Wrap page routes with QueryErrorBoundary
  - [ ] Test error boundary rendering
  - **Acceptance**: API errors display retry UI instead of blank page

## Phase 6: Documentation & Cleanup

- [ ] **Task 6.1**: Write quickstart README
  - [ ] Single-binary demo: `make build-single && ./bin/gateway`
  - [ ] Docker Compose production: `docker compose up -d`
  - [ ] Environment variables table
  - **Acceptance**: New user can run gateway in <5 minutes

- [ ] **Task 6.2**: Address TODO comments and cleanup
  - [ ] Search for TODO/FIXME in admin-ui/src/
  - [ ] Resolve or document each TODO
  - [ ] Remove unused components/styles
  - **Acceptance**: No unresolved TODOs remain

## Summary

| Phase | Tasks | Priority |
|-------|-------|----------|
| Phase 1: BillingClient | 2 | **High** |
| Phase 2: MonitorClient | 1 | **High** |
| Phase 3: Admin Handlers | 4 | **High** |
| Phase 4: Embedding | 2 | **High** |
| Phase 5: Testing | 3 | **Medium** |
| Phase 6: Documentation | 2 | **Medium** |
| **Total** | **14** | |

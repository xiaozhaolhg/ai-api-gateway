## 1. Project Setup & Dependencies

- [x] 1.1 Install antd@6.3.6, @ant-design/icons, @ant-design/cssinjs as dependencies
- [x] 1.2 Install react-i18next, i18next, i18next-browser-languagedetector as dependencies
- [x] 1.3 Remove leftover Vite template files (App.css hero/framework styles, unused assets in src/assets/)
- [x] 1.4 Update index.css to import antd styles and remove conflicting Tailwind base styles
- [x] 1.5 Update vite.config.ts to configure antd-compatible build settings

## 2. Core Infrastructure — Auth

- [x] 2.1 Create `src/contexts/AuthContext.tsx` — AuthProvider with login/logout/session-restore logic, JWT token storage in localStorage
- [x] 2.2 Create `src/components/ProtectedRoute.tsx` — route guard that redirects to /login if no valid session
- [x] 2.3 Create `src/pages/Login.tsx` — login page with antd Form, username/password fields, error feedback via message.error
- [x] 2.4 Update `src/App.tsx` — add AuthProvider wrapper, ProtectedRoute for admin routes, /login route

## 3. Core Infrastructure — i18n

- [x] 3.1 Create `src/i18n/index.ts` — i18next initialization with en/zh resources, browser language detection, localStorage persistence
- [x] 3.2 Create `src/i18n/locales/en/common.json` — English common translations (navigation, actions, status labels)
- [x] 3.3 Create `src/i18n/locales/zh/common.json` — Chinese common translations
- [x] 3.4 Create `src/i18n/locales/en/dashboard.json` — English dashboard page translations
- [x] 3.5 Create `src/i18n/locales/zh/dashboard.json` — Chinese dashboard page translations
- [x] 3.6 Create remaining en/zh namespace files for: providers, routing, users, groups, apiKeys, permissions, usage, budgets, pricing, health, alerts
- [x] 3.7 Create `src/components/LanguageSwitcher.tsx` — language toggle component using antd Dropdown
- [x] 3.8 Wrap App with I18nextProvider and antd ConfigProvider with dynamic locale

## 4. Core Infrastructure — API Client Overhaul

- [x] 4.1 Refactor `src/api/client.ts` — add auth header interceptor (Authorization: Bearer token), centralized error handler calling message.error
- [x] 4.2 Add authentication API methods: login(), logout(), getCurrentUser()
- [x] 4.3 Add routing rule API methods: getRoutingRules(), createRoutingRule(), updateRoutingRule(), deleteRoutingRule()
- [x] 4.4 Add group API methods: getGroups(), createGroup(), updateGroup(), deleteGroup(), addGroupMember(), removeGroupMember()
- [x] 4.5 Add permission API methods: getPermissions(), createPermission(), updatePermission(), deletePermission()
- [x] 4.6 Add budget API methods: getBudgets(), createBudget(), updateBudget(), deleteBudget()
- [x] 4.7 Add pricing rule API methods: getPricingRules(), createPricingRule(), updatePricingRule(), deletePricingRule()
- [x] 4.8 Add alert API methods: getAlertRules(), createAlertRule(), updateAlertRule(), deleteAlertRule(), getAlerts(), acknowledgeAlert()
- [x] 4.9 Add health API method: getProviderHealth()
- [x] 4.10 Add TypeScript interfaces for all new entity types (RoutingRule, Group, Permission, Budget, PricingRule, AlertRule, Alert, ProviderHealth)

## 5. Shared Layout — AppShell

- [x] 5.1 Create `src/components/AppShell.tsx` — antd Layout with Sider (collapsible), Header (user menu + language switcher), Content (Breadcrumb + Outlet)
- [x] 5.2 Create sidebar menu configuration with 4 groups (Infrastructure, Access Control, Billing, Observability) and icons
- [x] 5.3 Delete old `src/components/Layout.tsx` — replaced by AppShell
- [x] 5.4 Update `src/App.tsx` routing structure — AppShell as layout route, all page routes as children

## 6. Page — Dashboard

- [x] 6.1 Create `src/pages/Dashboard.tsx` — summary cards (users count, providers count, monthly spend, active alerts) using antd Statistic + Card
- [x] 6.2 Add quick-action links section (add provider, create user, issue API key, view alerts)
- [x] 6.3 Add loading Spin and error Alert states for dashboard data

## 7. Page — Providers (Overhaul)

- [x] 7.1 Rewrite `src/pages/Providers.tsx` using antd Table, Button, Modal, Form, Popconfirm, Tag, Spin, Empty
- [x] 7.2 Add edit provider flow — Modal with Form, PUT /admin/providers/:id
- [x] 7.3 Add proper loading, error, and empty states
- [x] 7.4 Add success/error feedback via antd message

## 8. Page — Routing Rules (New)

- [x] 8.1 Create `src/pages/RoutingRules.tsx` — antd Table listing routing rules with model pattern, provider, adapter type, priority, fallback chain, status
- [x] 8.2 Add create routing rule flow — Modal with Form including fallback chain ordered list input
- [x] 8.3 Add edit routing rule flow — Modal with pre-filled Form
- [x] 8.4 Add delete routing rule flow — Popconfirm + DELETE call
- [x] 8.5 Add loading, error, and empty states

## 9. Page — Users (Overhaul)

- [x] 9.1 Rewrite `src/pages/Users.tsx` using antd Table, Modal, Form, Popconfirm, Tag, Spin, Empty
- [x] 9.2 Add edit user flow — Modal with Form, PUT /admin/users/:id
- [x] 9.3 Add proper loading, error, and empty states
- [x] 9.4 Add success/error feedback via antd message

## 10. Page — Groups (New)

- [x] 10.1 Create `src/pages/Groups.tsx` — antd Table listing groups with name, description, member count, created at
- [x] 10.2 Add create group flow — Modal with Form
- [x] 10.3 Add edit group flow — Modal with pre-filled Form
- [x] 10.4 Add delete group flow — Popconfirm + DELETE call
- [x] 10.5 Add group membership management — expandable row or nested view showing members, with add/remove member actions
- [x] 10.6 Add loading, error, and empty states

## 11. Page — API Keys (Overhaul)

- [x] 11.1 Rewrite `src/pages/APIKeys.tsx` using antd components
- [x] 11.2 Fix user selector — populate from GET /admin/users instead of hardcoded options
- [x] 11.3 Improve API key creation display — antd Alert with copy-to-clipboard button
- [x] 11.4 Add proper loading, error, and empty states

## 12. Page — Permissions (New)

- [x] 12.1 Create `src/pages/Permissions.tsx` — antd Table listing permissions with group, model pattern, effect, created at
- [x] 12.2 Add create permission flow — Modal with Form including group selector populated from GET /admin/groups
- [x] 12.3 Add edit permission flow — Modal with pre-filled Form
- [x] 12.4 Add delete permission flow — Popconfirm + DELETE call
- [x] 12.5 Add loading, error, and empty states

## 13. Page — Usage (Overhaul)

- [x] 13.1 Rewrite `src/pages/Usage.tsx` using antd Statistic cards, Table, Form, DatePicker
- [x] 13.2 Add model and provider filter fields alongside existing user and date filters (added Input fields for model/provider, updated UsageRecord interface)
- [x] 13.3 Add proper loading, error, and empty states
- [x] 13.4 Add success/error feedback via antd message

## 14. Page — Budgets (New)

- [x] 14.1 Create `src/pages/Budgets.tsx` — antd Table listing budgets with name, scope, limit, current spend, remaining, status, period
- [x] 14.2 Add create budget flow — Modal with Form (name, scope selector, limit, period, optional user/group)
- [x] 14.3 Add edit budget flow — Modal with pre-filled Form
- [x] 14.4 Add delete budget flow — Popconfirm + DELETE call
- [x] 14.5 Add budget status badges — green Tag for active, yellow for warning, red for exceeded
- [x] 14.6 Add loading, error, and empty states

## 15. Page — Pricing Rules (New)

- [x] 15.1 Create `src/pages/PricingRules.tsx` — antd Table listing pricing rules with model, provider, prompt/completion token prices, currency, effective date
- [x] 15.2 Add create pricing rule flow — Modal with Form
- [x] 15.3 Add edit pricing rule flow — Modal with pre-filled Form
- [x] 15.4 Add delete pricing rule flow — Popconfirm + DELETE call
- [x] 15.5 Add loading, error, and empty states

## 16. Page — Health (Overhaul)

- [x] 16.1 Rewrite `src/pages/Health.tsx` — replace mock data with GET /admin/health API call
- [x] 16.2 Use antd Table with Badge/Tag for status (healthy=green, degraded=yellow, down=red)
- [x] 16.3 Add auto-refresh capability with configurable interval (Switch + InputNumber controls, default 30s)
- [x] 16.4 Add loading, error, and empty states

## 17. Page — Alerts (New)

- [x] 17.1 Create `src/pages/Alerts.tsx` — antd Tabs with "Rules" and "Active Alerts" tabs
- [x] 17.2 Add alert rules tab — Table listing rules with name, metric, condition, channel, status, actions (edit/delete)
- [x] 17.3 Add create alert rule flow — Modal with Form (name, metric, condition, threshold, channel)
- [x] 17.4 Add edit alert rule flow — Modal with pre-filled Form
- [x] 17.5 Add active alerts tab — Table listing firing/acknowledged alerts with severity, status, triggered at, description
- [x] 17.6 Add acknowledge alert action — button calling PUT /admin/alerts/:id/acknowledge
- [x] 17.7 Add loading, error, and empty states

## 18. Gateway-Service API Contracts Spec Update

- [x] 18.1 Update `openspec/specs/gateway-service/api-contracts.md` — add all new /admin/* endpoint definitions (auth, routing-rules, groups, permissions, budgets, pricing-rules, alert-rules, alerts, health)

## 19. Cleanup & Verification

- [x] 19.1 Remove all unused imports and dead code from refactored files
- [x] 19.2 Verify all routes render correctly with antd components
- [x] 19.3 Verify i18n switching works on all pages
- [x] 19.4 Verify auth flow: login → protected routes → logout → redirect to login
- [x] 19.5 Verify API client auth header injection on all endpoints
- [x] 19.6 Run `npm run build` and verify no TypeScript or build errors

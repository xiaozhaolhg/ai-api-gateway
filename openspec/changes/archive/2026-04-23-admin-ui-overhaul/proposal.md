## Why

The current admin-ui is a scaffolded Vite template with 5 flat page components that cover only ~30% of the backend domain surface. It has hardcoded mock data, broken API key management, leftover template CSS, no authentication, no error/empty/loading states, no shared component library, and no i18n. The backend microservices define far more entities (RoutingRules, Budgets, PricingRules, AlertRules, Groups, Permissions) than the UI exposes, making the admin panel unusable for real operations.

## What Changes

- **Replace raw Tailwind pages with Ant Design 6.3.6 component system** — tables, forms, modals, messages, notifications, spinners, badges
- **Add login page + session management** — authenticate against gateway-service, protect all routes, inject auth headers on API calls
- **Add Chinese/English i18n** — react-i18next + antd locale switching
- **Add Dashboard landing page** — summary cards (users, providers, spend, alerts) at `/`
- **Overhaul existing pages** — Providers (full CRUD + edit), Users (full CRUD + edit), API Keys (fix hardcoded user dropdown, fetch real users), Usage (proper filters + summary cards), Health (real API instead of mock data)
- **Add 6 new pages** — Routing Rules, Groups, Permissions, Budgets, Pricing Rules, Alerts
- **Overhaul API client** — typed methods for all `/admin/*` endpoints, auth header injection, error feedback via antd message/notification
- **BREAKING: Update gateway-service API contracts** — add new `/admin/*` endpoints for routing rules, budgets, pricing rules, alerts, groups, permissions, and authentication
- **Remove leftover Vite template artifacts** — App.css hero/framework CSS, unused assets

## Capabilities

### New Capabilities
- `admin-ui-auth`: Login page, session management, protected routes, auth header injection
- `admin-ui-i18n`: Chinese/English internationalization with antd locale support
- `admin-ui-dashboard`: Overview dashboard with summary cards and quick navigation
- `admin-ui-routing`: Routing rules management page with fallback chain configuration
- `admin-ui-groups`: Group CRUD and membership management
- `admin-ui-permissions`: Model access permission management
- `admin-ui-budgets`: Budget CRUD with spend tracking and status indicators
- `admin-ui-pricing`: Pricing rules management per model/provider
- `admin-ui-alerts`: Alert rules CRUD and active alerts lifecycle management

### Modified Capabilities
- `admin-ui-architecture`: Complete overhaul from scaffolded pages to Ant Design component system; add shared component library, auth context, i18n provider; expand from 5 to 11 pages
- `gateway-service/api-contracts`: Add new admin API endpoints for routing rules, budgets, pricing rules, alerts, groups, permissions, and authentication

## Impact

- **admin-ui/**: Full rewrite of all source files; new dependencies: antd@6.3.6, @ant-design/icons, react-i18next, i18next
- **gateway-service**: New admin HTTP handlers needed for new endpoints (backend implementation separate from this change)
- **openspec/specs/gateway-service/api-contracts.md**: New endpoint definitions
- **openspec/specs/admin-ui-architecture/spec.md**: Updated requirements reflecting new pages and component system

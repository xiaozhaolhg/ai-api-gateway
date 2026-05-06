# Proposal: Developer C Completion — Backend Wiring, Embedding & Testing

## Problem Statement

Developer C has completed ~65% of assigned work. The admin-ui frontend is fully built (13 pages, i18n, mock/real API client), but three critical gaps remain:

1. **Missing REST endpoints**: Budgets and Pricing Rules pages call `/admin/budgets` and `/admin/pricing-rules`, but gateway-service has no REST proxies for these billing-service gRPC methods. Alert endpoints use in-memory mock instead of monitor-service gRPC.
2. **No embedded UI**: The gateway-service cannot serve the admin-ui as static files. Deployment requires a separate nginx container.
3. **Insufficient testing & documentation**: No unit/integration tests for auth context, no quickstart guide, no smoke test.

## Proposed Solution

Wire gateway-service REST endpoints to billing-service and monitor-service gRPC, embed admin-ui into the Go binary, and complete testing/documentation per OpenSpec practices.

## Scope

- **In Scope**: Budget/Pricing/Alert REST proxy endpoints, BillingClient/MonitorClient gRPC methods, go:embed UI, single-binary build, tests, README
- **Out of Scope**: PostgreSQL migration (Dev B future), new UI features, proto schema changes (Dev B Week 4)

## Success Criteria

- `/admin/budgets` and `/admin/pricing-rules` CRUD endpoints return real data from billing-service
- `/admin/alert-rules` and `/admin/alerts` endpoints proxy to monitor-service (not in-memory mock)
- `make build-single` produces a single binary with embedded UI
- Gateway serves admin-ui at `/` when no API route matches
- Unit tests for auth context and protected routes pass
- Quickstart README covers single-binary demo and Docker Compose

## Dependencies

- billing-service gRPC: Budget/PricingRule handlers (✅ implemented)
- monitor-service gRPC: AlertRule/Alert handlers (⚠️ partially implemented)
- admin-ui build: `npm run build` produces dist/ (✅ working)

## 1. Proto Changes

- [x] 1.1 Add `user_id` field to `CreateBillingAccountRequest` in billing.proto
- [x] 1.2 Add `balance`, `user_id`, `balance_updated` to `UpdateBillingAccountRequest`
- [x] 1.3 Add `user_id`, `balance`, `currency` to `BillingAccount` message
- [x] 1.4 Add `GetBillingAccountByUser` RPC + `GetBillingAccountByUserRequest` message
- [x] 1.5 Regenerate Go protos with protoc

## 2. Billing Service Handler Updates

- [x] 2.1 Update `CreateBillingAccount` handler to set `UserID` from request
- [x] 2.2 Update `UpdateBillingAccount` handler to support balance changes via `balance_updated` flag
- [x] 2.3 Implement `GetBillingAccountByUser` gRPC handler
- [x] 2.4 Update `billingAccountToProto` to include UserID, Balance, Currency

## 3. Gateway Service Billing Account Handler

- [x] 3.1 Add `AdminBillingAccountsHandler` with Get/Create/AdjustBalance methods
- [x] 3.2 Add billing account routes to gateway main.go with JWT + admin middleware
- [x] 3.3 Add `GetBillingAccountByUser`, `CreateBillingAccount`, `UpdateBillingAccountBalance` to billing client
- [x] 3.4 Fix admin middleware bug (`c.Get("role")` → `c.Get("userRole")`)

## 4. Admin UI API Client

- [x] 4.1 Add `BillingAccount` type to types.ts
- [x] 4.2 Add `getBillingAccount`, `adjustBalance` to APIClientInterface
- [x] 4.3 Add real API client implementation in client.ts
- [x] 4.4 Add mock client implementation in mockClient.ts
- [x] 4.5 Add wrapper methods in UnifiedAPIClient

## 5. Admin UI Users Page Enhancement

- [x] 5.1 Add balance column to Users table
- [x] 5.2 Add "Recharge" button in actions column (admin-only)
- [x] 5.3 Implement recharge modal with amount input
- [x] 5.4 Show current balance and success/error feedback
- [x] 5.5 Refresh balance after successful recharge

## 6. Bug Fixes Discovered During Implementation

- [x] 6.1 Fix admin middleware: was checking `c.Get("role")` but JWT sets `c.Get("userRole")`
- [x] 6.2 Fix admin-ui Dockerfile: nginx.conf path missing `admin-ui/` prefix

## 7. Integration Tests (Manual)

- [x] 7.1 Test: Admin creates billing account → 200, balance=$500
- [x] 7.2 Test: Admin recharges user balance +$100 → balance=$600
- [x] 7.3 Test: Non-admin gets 403 on billing endpoints
- [x] 7.4 Test: View billing account returns correct balance

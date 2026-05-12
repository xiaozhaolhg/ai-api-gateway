## Why

The billing system currently supports usage tracking and pricing rules, but has no way for admins to create billing accounts or recharge user balances through the UI. The only way to set up a user's billing account is direct database manipulation, which is not sustainable for production use. Admins need a proper UI and API to manage user credits.

## What Changes

- Add HTTP endpoints in gateway-service for admin billing account management:
  - `POST /admin/billing/accounts` - Create billing account with initial credit
  - `PUT /admin/billing/accounts/:userId/balance` - Adjust user balance (recharge/deduct)
  - `GET /admin/billing/accounts/:userId` - Get billing account info
- Restrict all billing account endpoints to admin role only
- Add "Recharge" action to the Users page in admin UI
- Show user balance in the Users table

## Capabilities

### New Capabilities
- `admin-billing-accounts`: HTTP API for creating and managing billing accounts
- `admin-ui-recharge`: Recharge UI integrated into the user management page

### Modified Capabilities
- `billing`: Extend existing CreateBillingAccount gRPC to accept user_id
- `gateway`: Add billing account HTTP routes with admin auth
- `admin-ui`: Enhance Users page with balance display and recharge action

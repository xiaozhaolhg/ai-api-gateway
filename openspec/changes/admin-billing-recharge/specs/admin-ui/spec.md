# Admin UI - User Recharge Feature

## Overview

Add recharge functionality to the existing Users management page. Admins can view user balances and recharge user accounts directly from the user list.

## UI Changes

### Users Table - New Columns

| Column | Description |
|--------|-------------|
| Balance | Shows user's current billing account balance, formatted as `$X.XX` |
| Actions | Existing Edit/Delete + new "Recharge" button |

### Balance Column

- Shows `$0.00` when user has no billing account
- Shows formatted balance when account exists
- Styled as monospace for readability

### Recharge Button

- Only visible to admin users (check `user.role === 'admin'`)
- Opens a modal dialog
- Shows current balance for reference

### Recharge Modal

| Field | Type | Description |
|-------|------|-------------|
| User Name | Display only | Shows who is being recharged |
| Current Balance | Display only | Shows current balance before adjustment |
| Amount | Number input | Positive number to add to balance |
| Note | Text (optional) | Reason for recharge |

Contains Confirm and Cancel buttons. On confirm, calls `PUT /admin/billing/accounts/:userId/balance` with the amount.

## API Client Changes

### New Methods

```typescript
interface APIClientInterface {
  getBillingAccount(userId: string): Promise<BillingAccount>;
  adjustBalance(userId: string, amount: number): Promise<BillingAccount>;
}
```

### New Types

```typescript
interface BillingAccount {
  id: string;
  user_id: string;
  balance: number;
  currency: string;
  status: string;
}
```

## User Experience Flow

1. Admin navigates to Users page
2. Sees all users with Balance column
3. Clicks "Recharge" on a user row
4. Modal opens showing user name and current balance
5. Admin enters amount and optional note
6. Clicks Confirm
7. UI shows loading state, then success message
8. Balance column updates automatically

## Authorization

- Recharge button hidden for non-admin users
- API calls fail with 403 for non-admin tokens
- Error state displayed appropriately

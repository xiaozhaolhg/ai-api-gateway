# Gateway Service - Admin Billing Account Endpoints

## Overview

Add HTTP endpoints to the gateway service for admin management of user billing accounts. These endpoints proxy to the billing-service gRPC API with admin authorization checks.

## Endpoints

### POST /admin/billing/accounts

Create a billing account for a user with optional initial credit.

**Request:**
```json
{
  "user_id": "string (required)",
  "name": "string (optional)",
  "initial_credit": "number (optional, default 0)"
}
```

**Response (201):**
```json
{
  "id": "string",
  "user_id": "string",
  "balance": 0.0,
  "currency": "USD",
  "status": "active"
}
```

**Authorization:** Admin only (JWT + role check)

### GET /admin/billing/accounts/:userId

Get billing account info for a user.

**Response (200):**
```json
{
  "id": "string",
  "user_id": "string",
  "balance": 100.0,
  "currency": "USD",
  "status": "active"
}
```

**Response (404):** `{"error": "billing account not found"}`

**Authorization:** Admin only

### PUT /admin/billing/accounts/:userId/balance

Adjust a user's balance (positive = recharge, negative = deduct).

**Request:**
```json
{
  "amount": 50.0
}
```

**Response (200):**
```json
{
  "id": "string",
  "user_id": "string",
  "balance": 150.0,
  "currency": "USD",
  "status": "active"
}
```

**Authorization:** Admin only

## Data Model

### BillingAccountResponse (HTTP)

```go
type BillingAccountResponse struct {
    ID       string  `json:"id"`
    UserID   string  `json:"user_id"`
    Balance  float64 `json:"balance"`
    Currency string  `json:"currency"`
    Status   string  `json:"status"`
}
```

## Route Registration

Routes registered in `cmd/server/main.go`:
```go
adminBillingHandler := handler.NewAdminBillingAccountsHandler(billingClient)

// Uses existing JWT middleware + admin-only check
r.POST("/admin/billing/accounts", adminBillingHandler.CreateBillingAccount)
r.GET("/admin/billing/accounts/:userId", adminBillingHandler.GetBillingAccount)
r.PUT("/admin/billing/accounts/:userId/balance", adminBillingHandler.AdjustBalance)
```

## Authorization

All endpoints are protected by:
1. JWT authentication middleware (validates session)
2. Admin role check (c.GetString("userRole") == "admin")

Non-admin users receive 403 Forbidden with `{"error": "insufficient permissions"}`.

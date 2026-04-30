# billing-service

## Purpose

Usage and billing domain — token counting, cost estimation, pricing, budgets, invoices.

## Service Responsibility

- **Role**: Usage recording, cost calculation, budget enforcement
- **Owned Entities**: UsageRecord, PricingRule, BillingAccount, Budget, Invoice
- **Data Layer**: billing-db (SQLite/PostgreSQL)

## Dependencies

### Calls To

| Service | Methods | Purpose |
|---|---|---|
| (none) | — | Does not call other internal services |

### Called By

| Service | Methods | Purpose |
|---|---|---|
| gateway-service | `CheckBudget`, `RecordUsage` | Rate limiting, usage tracking |
| provider-service | `OnProviderResponse` | Token counts from provider callback |
| gateway-service | `GetUsage`, `GetUsageAggregation` | Usage queries |
| gateway-service | `CreateBudget`, `UpdateBudget`, `DeleteBudget` | Budget management |

### Data Dependencies

- **Database**: billing-db (UsageRecord, PricingRule, Budget, Invoice)
- **Cache**: Redis (pricing rules, usage aggregation)

## Key Design

### Usage Recording

Two paths:
1. **Provider callback** (primary): provider-service dispatches OnProviderResponse after each response
2. **Direct call** (MVP fallback): gateway-service calls RecordUsage

### Cost Calculation

- Retrieve PricingRule for model/provider
- Apply price_per_token × token_count = cost

### Budget Enforcement

- **CheckBudget**: Returns current_spend, limit, remaining
- **Soft cap exceeded**: Alert, allow request
- **Hard cap exceeded**: Block with 429

### Key Operations

- **OnProviderResponse**: Receive token counts from provider callback
- **RecordUsage**: Direct usage recording
- **CheckBudget**: Budget check for rate limiting
- **GetUsage/GetUsageAggregation**: Usage queries
- **CreateBudget/UpdateBudget/DeleteBudget**: Budget CRUD
- **CreatePricingRule/UpdatePricingRule/DeletePricingRule**: Pricing CRUD
- **GenerateInvoice**: Invoice generation (Phase 3+)

## Requirements

### Requirement: Usage recording
The billing-service SHALL record usage data from provider callbacks and direct calls.

#### Scenario: Provider callback recording
- **WHEN** provider-service dispatches OnProviderResponse callback with token counts
- **THEN** the billing-service SHALL store the usage record with provider, model, and token counts

#### Scenario: Direct usage recording
- **WHEN** gateway-service calls RecordUsage RPC
- **THEN** the billing-service SHALL store the usage record for the user and model

### Requirement: Budget enforcement
The billing-service SHALL enforce budget limits for users and groups.

#### Scenario: Budget check
- **WHEN** gateway-service calls CheckBudget RPC
- **THEN** the billing-service SHALL return current spend, limit, and remaining budget
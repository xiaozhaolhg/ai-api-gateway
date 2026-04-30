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

### Requirement: GetUsageAggregation returns multiple rows
The billing-service `GetUsageAggregation` SHALL return **one row per unique `group_by` value** (user_id, group_id, provider_id, model), NOT just a single row.

#### Scenario: Aggregation by provider_id
- **WHEN** `GetUsageAggregation` is called with `group_by="provider_id"`
- **THEN** the system returns one `UsageAggregation` row for each provider_id in the result set

#### Scenario: Aggregation by model
- **WHEN** `GetUsageAggregation` is called with `group_by="model"`
- **THEN** the system returns one `UsageAggregation` row for each model in the result set

#### Scenario: Aggregation by group_id
- **WHEN** `GetUsageAggregation` is called with `group_by="group_id"`
- **THEN** the system returns one `UsageAggregation` row for each group_id in the result set

#### Scenario: Multiple aggregation rows returned
- **WHEN** `GetUsageAggregation` repository query executes for a user with multiple providers/models
- **THEN** it returns N rows (one per unique group_by value) with `GROUP BY <group_by>`
- **AND** NOT limited to 1 row (LIMIT 1 removed)

### Requirement: GetUsageAggregation removes LIMIT 1 constraint
The `GetUsageAggregation` repository query SHALL NOT use `LIMIT 1` and SHALL return all aggregation rows matching the query.

#### Scenario: Repository returns all aggregation rows
- **WHEN** the repository executes the aggregation SQL with `GROUP BY provider_id`
- **THEN** it returns all grouped rows (not just the first one)

#### Scenario: Service handles multiple aggregation results
- **WHEN** the repository returns multiple `UsageAggregation` entities
- **THEN** the service maps ALL of them into `ListUsageAggregationResponse.aggregations`
- **AND** the handler returns the complete list

### Requirement: Per-user, per-group, per-provider, per-model token tracking
The system SHALL store `UsageRecord` with `user_id`, `group_id`, `provider_id`, and `model` fields to enable per-level breakdowns.

#### Scenario: Record with all fields
- **WHEN** `RecordUsage` is called with user_id, group_id, provider_id, model, and token counts
- **THEN** a `UsageRecord` is persisted with ALL fields populated

#### Scenario: Query by group_id
- **WHEN** `GetUsage` is called with `group_id="group-dev-1"`
- **THEN** all usage records for that group are returned

#### Scenario: Query by provider_id
- **WHEN** `GetUsage` is called with `provider_id="openai"`
- **THEN** all usage records for that provider are returned

#### Scenario: Query by model
- **WHEN** `GetUsage` is called with `model="gpt-4"`
- **THEN** all usage records for that model are returned

### Requirement: Budget enforcement
The billing-service SHALL enforce budget limits for users and groups.

#### Scenario: Budget check
- **WHEN** gateway-service calls CheckBudget RPC
- **THEN** the billing-service SHALL return current spend, limit, and remaining budget
## ADDED Requirements

### Requirement: DDD four-layer architecture
The billing-service SHALL implement four-layer Clean Architecture: Domain, Application, Infrastructure, and Handler with dependency direction from outer to inner layers.

#### Scenario: Domain layer has no external dependencies
- **WHEN** the domain layer is imported
- **THEN** it SHALL NOT import any code from application, infrastructure, or handler layers

### Requirement: UsageRecord entity and repository
The billing-service SHALL own the UsageRecord entity with fields: id, user_id, group_id, provider_id, model, prompt_tokens, completion_tokens, cost, timestamp. It SHALL provide a UsageRecordRepository interface.

#### Scenario: Record usage from provider callback
- **WHEN** an OnProviderResponse callback is received
- **THEN** the service SHALL create a UsageRecord with token counts and calculated cost

#### Scenario: Query usage with filters
- **WHEN** a GetUsage request is received with user_id, model, or date range filters
- **THEN** the service SHALL return matching UsageRecords with pagination

### Requirement: PricingRule entity and repository
The billing-service SHALL own the PricingRule entity with fields: id, model, provider_id, price_per_prompt_token, price_per_completion_token, currency. It SHALL provide a PricingRuleRepository interface.

#### Scenario: Cost calculation
- **WHEN** usage is recorded for a model with a matching PricingRule
- **THEN** the service SHALL calculate cost as (prompt_tokens × price_per_prompt_token) + (completion_tokens × price_per_completion_token)

#### Scenario: No pricing rule found
- **WHEN** usage is recorded for a model with no matching PricingRule
- **THEN** the service SHALL record cost as 0.0

### Requirement: Budget entity and repository
The billing-service SHALL own the Budget entity with fields: id, account_id, limit, period, soft_cap_pct, hard_cap_pct, status. It SHALL provide a BudgetRepository interface.

#### Scenario: Check budget
- **WHEN** a CheckBudget request is received for a user
- **THEN** the service SHALL return BudgetStatus with current_spend, budget_limit, remaining, and cap exceeded flags

#### Scenario: Hard cap exceeded
- **WHEN** current spend exceeds the hard cap percentage of the budget limit
- **THEN** BudgetStatus.hard_cap_exceeded SHALL be true

### Requirement: Direct usage recording
The billing-service SHALL support direct usage recording via RecordUsage gRPC call (MVP fallback when provider callback is not available).

#### Scenario: Record usage directly
- **WHEN** a RecordUsage request is received from gateway-service
- **THEN** the service SHALL create a UsageRecord with the provided token counts

### Requirement: Usage aggregation
The billing-service SHALL provide GetUsageAggregation that returns aggregated usage statistics grouped by user, model, or provider.

#### Scenario: Aggregate by model
- **WHEN** a GetUsageAggregation request is received with group_key "model"
- **THEN** the service SHALL return total prompt/completion tokens, cost, and request count per model

### Requirement: gRPC server implementation
The billing-service SHALL implement the BillingService gRPC server as defined in the api/ module proto definitions.

#### Scenario: All proto RPCs are implemented
- **WHEN** the billing-service starts
- **THEN** it SHALL register all RPCs defined in billing.proto: OnProviderResponse, RecordUsage, GetUsage, GetUsageAggregation, EstimateCost, CheckBudget, and CRUD operations for budgets, pricing rules, and billing accounts

### Requirement: SQLite persistence
The billing-service SHALL use SQLite as its database with tables for usage_records, pricing_rules, billing_accounts, and budgets, managed via migrations.

#### Scenario: Database initialized on startup
- **WHEN** the billing-service starts
- **THEN** it SHALL create the SQLite database file and run migrations if the database is new

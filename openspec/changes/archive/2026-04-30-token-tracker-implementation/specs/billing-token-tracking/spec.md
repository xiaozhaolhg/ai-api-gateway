## ADDED Requirements|

### Requirement: GetUsageAggregation returns multiple rows|
The billing-service `GetUsageAggregation` SHALL return **one row per unique `group_by` value** (user_id, group_id, provider_id, model), NOT just a single row.

#### Scenario: Aggregation by provider_id|
- **WHEN** `GetUsageAggregation` is called with `group_by="provider_id"`|
- **THEN** the system returns one `UsageAggregation` row for each provider_id in the result set|

#### Scenario: Aggregation by model|
- **WHEN** `GetUsageAggregation` is called with `group_by="model"`|
- **THEN** the system returns one `UsageAggregation` row for each model in the result set|

#### Scenario: Aggregation by group_id|
- **WHEN** `GetUsageAggregation` is called with `group_by="group_id"`|
- **THEN** the system returns one `UsageAggregation` row for each group_id in the result set|

#### Scenario: Aggregation by user_id|
- **WHEN** `GetUsageAggregation` is called with `group_by="user_id"`|
- **THEN** the system returns one `UsageAggregation` row for each user_id in the result set|

### Requirement: GetUsageAggregation removes LIMIT 1 constraint|
The `GetUsageAggregation` repository query SHALL NOT use `LIMIT 1` and SHALL return all aggregation rows matching the query.

#### Scenario: Multiple aggregation rows returned|
- **WHEN** `GetUsageAggregation` repository query executes for a user with multiple providers|
- **THEN** it returns N rows (one per unique provider_id) with `GROUP BY provider_id`|

#### Scenario: Date range filtering|
- **WHEN** `GetUsageAggregation` is called with `start_time` and `end_time`|
- **THEN** only records within that date range are included in the aggregation|

### Requirement: ListUsageAggregationResponse supports multiple aggregations|
The `ListUsageAggregationResponse` proto message SHALL carry a list of `UsageAggregation` messages (already does per proto definition).

#### Scenario: Response with multiple aggregations|
- **WHEN** the service returns aggregation results|
- **THEN** `ListUsageAggregationResponse.aggregations` contains one entry per group_by value|

### Requirement: Per-user, per-group, per-provider, per-model token tracking|
The system SHALL store `UsageRecord` with `user_id`, `group_id`, `provider_id`, and `model` fields to enable per-level breakdowns.

#### Scenario: Record with all fields|
- **WHEN** `RecordUsage` is called with user_id, group_id, provider_id, model, and token counts|
- **THEN** a `UsageRecord` is persisted with ALL fields populated|

#### Scenario: Query by group_id|
- **WHEN** `GetUsage` is called with `group_id="group-dev-1"`|
- **THEN** all usage records for that group are returned|

#### Scenario: Query by provider_id|
- **WHEN** `GetUsage` is called with `provider_id="openai"`|
- **THEN** all usage records for that provider are returned|

#### Scenario: Query by model|
- **WHEN** `GetUsage` is called with `model="gpt-4"`|
- **THEN** all usage records for that model are returned|

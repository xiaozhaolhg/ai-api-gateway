## MODIFIED Requirements|

### Requirement: GetUsageAggregation returns multiple rows|
The billing-service `GetUsageAggregation` SHALL return **one row per unique `group_by` value** (user_id, group_id, provider_id, model), NOT just a single row.|

#### Scenario: Aggregation by provider_id|
- **WHEN** `GetUsageAggregation` is called with `group_by="provider_id"`|
- **THEN** the system returns one `UsageAggregation` row for each provider_id in the result set|

#### Scenario: Aggregation by model|
- **WHEN** `GetUsageAggregation` is called with `group_by="model"`|
- **THEN** the system returns one `UsageAggregation` row for each model in the result set|

#### Scenario: Aggregation by group_id|
- **WHEN** `GetUsageAggregation` is called with `group_by="group_id"`|
- **THEN** the system returns one `UsageAggregation` row for each group_id in the result set|

#### Scenario: Multiple aggregation rows returned|
- **WHEN** `GetUsageAggregation` repository query executes for a user with multiple providers/models|
- **THEN** it returns N rows (one per unique group_by value) with `GROUP BY <group_by>`|
- **AND** NOT limited to 1 row (LIMIT 1 removed)|

### Requirement: GetUsageAggregation removes LIMIT 1 constraint|
The `GetUsageAggregation` repository query SHALL NOT use `LIMIT 1` and SHALL return all aggregation rows matching the query.|

#### Scenario: Repository returns all aggregation rows|
- **WHEN** the repository executes the aggregation SQL with `GROUP BY provider_id`|
- **THEN** it returns all grouped rows (not just the first one)|

#### Scenario: Service handles multiple aggregation results|
- **WHEN** the repository returns multiple `UsageAggregation` entities|
- **THEN** the service maps ALL of them into `ListUsageAggregationResponse.aggregations`|
- **AND** the handler returns the complete list|

### Requirement: Per-user, per-group, per-provider, per-model token tracking|
The system SHALL store `UsageRecord` with `user_id`, `group_id`, `provider_id`, and `model` fields to enable per-level breakdowns.|

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

## MODIFIED Requirements

### Requirement: Usage recording
The billing-service SHALL record usage data from provider callbacks and direct calls, correctly aggregating multiple records from the same request.

#### Scenario: Multiple records for same request
- **WHEN** gateway-service sends multiple `RecordUsage` calls for the same streaming request (intermediate + final)
- **THEN** the billing-service SHALL store each record independently
- **AND** `GetUsageAggregation` SHALL correctly sum all records for the same user/model/provider/date

#### Scenario: Aggregation of intermediate records
- **WHEN** `GetUsageAggregation` is queried for a time period that includes intermediate streaming records
- **THEN** the result SHALL include the sum of all intermediate and final records
- **AND** the total SHALL equal the actual token usage for that period

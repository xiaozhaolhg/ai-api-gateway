## ADDED Requirements

### Requirement: Enhanced usage query gRPC endpoints
The billing-service SHALL provide new gRPC endpoints for user-specific and group-specific usage queries with filtering capabilities.

#### Scenario: Get usage by user ID
- **WHEN** auth-service calls GetUsageByUser gRPC with user_id, start_date, end_date, page, and page_size
- **THEN** billing-service returns paginated UsageRecord entities filtered by user_id and date range

#### Scenario: Get usage by group ID
- **WHEN** auth-service calls GetUsageByGroup gRPC with group_id, start_date, end_date, page, and page_size
- **THEN** billing-service returns paginated UsageRecord entities filtered by group_id and date range

#### Scenario: Get filtered usage for export
- **WHEN** auth-service calls GetUsageForExport gRPC with optional user_id, group_id, start_date, end_date
- **THEN** billing-service returns all matching UsageRecord entities without pagination limits for export processing

#### Scenario: Usage query with date range validation
- **WHEN** any usage query gRPC is called with start_date and end_date
- **THEN** billing-service validates date format (YYYY-MM-DD) and that end_date is not before start_date

#### Scenario: Usage query with pagination validation
- **WHEN** any usage query gRPC is called with page and page_size
- **THEN** billing-service validates page >= 1 and page_size between 1 and 1000

### Requirement: Enhanced usage aggregation queries
The billing-service SHALL extend usage aggregation capabilities to support user-specific and group-specific aggregations with multiple grouping options.

#### Scenario: Get user usage aggregation
- **WHEN** auth-service calls GetUsageAggregationByUser with user_id, start_date, end_date, and group_by parameter
- **THEN** billing-service returns UsageAggregation entities for that user grouped by the specified dimension (provider_id, model, date)

#### Scenario: Get group usage aggregation
- **WHEN** auth-service calls GetUsageAggregationByGroup with group_id, start_date, end_date, and group_by parameter
- **THEN** billing-service returns UsageAggregation entities for that group grouped by the specified dimension

#### Scenario: Multi-dimensional aggregation
- **WHEN** aggregation query includes group_by="provider_id,model"
- **THEN** billing-service returns UsageAggregation entities grouped by both provider and model combinations

#### Scenario: Daily aggregation for export
- **WHEN** auth-service requests aggregation for export with group_by="date"
- **THEN** billing-service returns daily aggregated usage data suitable for chart generation

### Requirement: Usage repository query extensions
The billing-service SHALL extend the UsageRecordRepository interface to support user-specific and group-specific queries with filtering.

#### Scenario: Repository method for user queries
- **WHEN** service calls GetByUserID with user_id, start_date, end_date, page, page_size
- **THEN** repository returns paginated UsageRecord entities for that user within date range

#### Scenario: Repository method for group queries
- **WHEN** service calls GetByGroupID with group_id, start_date, end_date, page, page_size
- **THEN** repository returns paginated UsageRecord entities for that group within date range

#### Scenario: Repository method for export queries
- **WHEN** service calls GetForExport with optional filters and no pagination
- **THEN** repository returns all matching UsageRecord entities for export processing

#### Scenario: Repository method for user aggregation
- **WHEN** service calls GetAggregationByUserID with user_id, start_date, end_date, group_by
- **THEN** repository returns UsageAggregation entities for that user grouped as specified

#### Scenario: Repository method for group aggregation
- **WHEN** service calls GetAggregationByGroupID with group_id, start_date, end_date, group_by
- **THEN** repository returns UsageAggregation entities for that group grouped as specified

### Requirement: Usage data export formatting
The billing-service SHALL provide formatted usage data suitable for export in multiple formats through gRPC responses.

#### Scenario: CSV-formatted usage data
- **WHEN** auth-service requests usage data in CSV format
- **THEN** billing-service returns UsageRecord entities with metadata for CSV formatting (headers, row ordering)

#### Scenario: JSON-formatted usage data
- **WHEN** auth-service requests usage data in JSON format
- **THEN** billing-service returns UsageRecord entities in structure suitable for JSON serialization

#### Scenario: Export data with metadata
- **WHEN** usage data is requested for export
- **THEN** billing-service includes export metadata (generation timestamp, record count, filter parameters)

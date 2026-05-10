## 1. Billing Service gRPC Extensions

- [ ] 1.1 Add GetUsageByUser gRPC method to billing.proto
- [ ] 1.2 Add GetUsageByGroup gRPC method to billing.proto  
- [ ] 1.3 Add GetUsageForExport gRPC method to billing.proto
- [ ] 1.4 Add GetUsageAggregationByUser gRPC method to billing.proto
- [ ] 1.5 Add GetUsageAggregationByGroup gRPC method to billing.proto
- [ ] 1.6 Regenerate proto files with buf generate

## 2. Billing Service Repository Extensions

- [ ] 2.1 Add GetByUserID method to UsageRecordRepository interface
- [ ] 2.2 Add GetByGroupID method to UsageRecordRepository interface
- [ ] 2.3 Add GetForExport method to UsageRecordRepository interface
- [ ] 2.4 Add GetAggregationByUserID method to UsageRecordRepository interface
- [ ] 2.5 Add GetAggregationByGroupID method to UsageRecordRepository interface
- [ ] 2.6 Implement SQLite repository methods for user/group queries

## 3. Billing Service Handler Implementation

- [ ] 3.1 Implement GetUsageByUser gRPC handler with validation
- [ ] 3.2 Implement GetUsageByGroup gRPC handler with validation
- [ ] 3.3 Implement GetUsageForExport gRPC handler with validation
- [ ] 3.4 Implement GetUsageAggregationByUser gRPC handler
- [ ] 3.5 Implement GetUsageAggregationByGroup gRPC handler
- [ ] 3.6 Add date range and pagination validation utilities

## 4. Auth Service HTTP Endpoints

- [ ] 4.1 Add billing-service gRPC client to auth-service
- [ ] 4.2 Implement GET /admin/usage/users/:id HTTP handler
- [ ] 4.3 Implement GET /admin/usage/groups/:id HTTP handler
- [ ] 4.4 Implement GET /admin/usage/export HTTP handler
- [ ] 4.5 Add request validation middleware for usage endpoints
- [ ] 4.6 Add error handling for billing-service failures

## 5. Auth Service Export Functionality

- [ ] 5.1 Implement CSV formatting utility for usage data
- [ ] 5.2 Implement JSON formatting utility for usage data
- [ ] 5.3 Add export metadata generation (timestamp, record count)
- [ ] 5.4 Implement content-type headers for export responses
- [ ] 5.5 Add file download headers for export responses

## 6. Integration and Gateway Routing

- [ ] 6.1 Add usage API routes to gateway-service routing
- [ ] 6.2 Configure gateway-service to proxy usage endpoints to auth-service
- [ ] 6.3 Add usage API authentication middleware
- [ ] 6.4 Test cross-service communication (gateway → auth → billing)

## 7. Unit Tests

- [ ] 7.1 Unit tests for billing-service repository methods
- [ ] 7.2 Unit tests for billing-service gRPC handlers
- [ ] 7.3 Unit tests for auth-service HTTP handlers
- [ ] 7.4 Unit tests for export formatting utilities
- [ ] 7.5 Unit tests for request validation and error handling

## 8. Integration Tests

- [ ] 8.1 Integration test: GET /admin/usage/users/:id end-to-end
- [ ] 8.2 Integration test: GET /admin/usage/groups/:id end-to-end
- [ ] 8.3 Integration test: GET /admin/usage/export CSV format
- [ ] 8.4 Integration test: GET /admin/usage/export JSON format
- [ ] 8.5 Integration test: Date range filtering functionality
- [ ] 8.6 Integration test: Pagination functionality
- [ ] 8.7 Integration test: Error handling scenarios

## 9. Admin UI Enhanced Usage Features

- [ ] 9.1 Add new API client methods for user/group specific usage endpoints
- [ ] 9.2 Implement user-specific usage view with charts and analytics
- [ ] 9.3 Implement group-specific usage view with member breakdowns
- [ ] 9.4 Add export functionality (CSV/JSON) to usage page
- [ ] 9.5 Implement usage charts (token consumption, cost breakdown, trends)
- [ ] 9.6 Add virtual scrolling for large usage datasets
- [ ] 9.7 Implement usage data caching and background refresh
- [ ] 9.8 Add navigation links from user/group management to usage views

## 10. Admin UI Performance and UX

- [ ] 10.1 Implement lazy loading for usage charts and analytics
- [ ] 10.2 Add loading states and skeleton screens for usage data
- [ ] 10.3 Implement error handling and retry logic for usage API calls
- [ ] 10.4 Add breadcrumb navigation for usage views
- [ ] 10.5 Implement responsive design for usage analytics on mobile
- [ ] 10.6 Add export progress indicators for large datasets

## 11. Admin UI Testing

- [ ] 11.1 Unit tests for new usage API client methods
- [ ] 11.2 Component tests for user-specific usage view
- [ ] 11.3 Component tests for group-specific usage view
- [ ] 11.4 Integration tests for export functionality
- [ ] 11.5 Performance tests for virtual scrolling with large datasets
- [ ] 11.6 E2E tests for complete usage workflow

## 12. Documentation and Cleanup

- [ ] 12.1 Update API documentation with new usage endpoints
- [ ] 12.2 Add usage API examples to admin UI documentation
- [ ] 12.3 Update OpenAPI spec with usage endpoint schemas
- [ ] 12.4 Add user guide for enhanced usage analytics features
- [ ] 12.5 Code review and refactoring based on test results

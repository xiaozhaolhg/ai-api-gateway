## Why

The current usage tracking system only provides raw usage data through gRPC endpoints, limiting admin visibility into user and group-specific consumption patterns. Admin users need HTTP-based endpoints with filtering, export capabilities, and per-user/group views to effectively monitor and manage AI API usage across the organization.

## What Changes

- Add HTTP endpoints in auth-service for usage data access:
  - `GET /admin/usage/users/:id` - Get usage for specific user
  - `GET /admin/usage/groups/:id` - Get usage for specific group  
  - `GET /admin/usage/export` - Export usage data (CSV/JSON)
- Extend billing-service gRPC APIs to support user/group filtering
- Add export functionality with multiple format support (CSV, JSON)
- Implement date range filtering and pagination for all usage endpoints

## Capabilities

### New Capabilities
- `usage-api`: HTTP endpoints for accessing usage data with filtering and pagination
- `usage-export`: Export functionality for usage data in CSV/JSON formats
- `usage-analytics`: User and group-specific usage views with aggregation
- `admin-ui-usage`: Enhanced admin UI for usage analytics, charts, and export

### Modified Capabilities
- `billing`: Extend existing billing service to support enhanced usage queries and filtering
- `admin-ui`: Enhance existing usage page with advanced features and new views

## Impact

- **auth-service**: Add new HTTP handlers for usage endpoints
- **billing-service**: Extend repository interfaces and gRPC handlers for user/group filtering
- **gateway-service**: Route usage endpoints to auth-service
- **admin-ui**: Enhanced usage page with charts, export, and user/group specific views
- **Cross-service communication**: Enhanced gRPC calls between auth-service and billing-service
- **Database queries**: New query patterns for user/group-specific usage aggregation
- **Frontend performance**: Virtualization and caching for large usage datasets

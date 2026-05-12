# Design: Usage API Bugfixes

## Architecture

```
Admin UI (port 3000)
  │ GET /admin/auth/usage
  ▼
Gateway Service (port 8080)
  │ jwtAuthMiddleware → extracts userId + userRole from JWT
  │ handleGetUsage:
  │   if userRole == "admin" → userID = ""
  │   else → userID = c.GetString("userId")
  │   timestamp = time.Unix(r.Timestamp, 0).Format(time.RFC3339)
  ▼
Billing Service (gRPC :50054)
  │ GetUsage(userID, page, pageSize)
  │   if userID == "" → SELECT ... FROM usage_records ORDER BY ...
  │   else → SELECT ... FROM usage_records WHERE user_id = ? ...
  ▼
SQLite DB (/data/billing.db)
```

## Changes

### 1. billing-service/internal/infrastructure/migration/migration.go

**Before:**
```go
migrations := []string{
    createUsageRecordsTable(),
    addGroupIDColumn(),       // BUG: column already exists
    createPricingRulesTable(),
    ...
}
```

**After:**
```go
migrations := []string{
    createUsageRecordsTable(),
    createPricingRulesTable(),
    ...
}
```

`group_id` is already part of `CREATE TABLE IF NOT EXISTS usage_records (... group_id TEXT ...)`. The `ALTER TABLE ADD COLUMN group_id` is redundant and causes a crash.

### 2. billing-service/internal/infrastructure/repository/usage_record_repository.go

**Before:** Always queries with `WHERE user_id = ?`.
**After:** When `userID == ""` (admin view), queries without `WHERE` clause — returns all records.

### 3. gateway-service/internal/handler/admin_usage.go

- Added `Timestamp string \`json:"timestamp"\`` to `UsageRecord` struct
- Added `"time"` import
- `Timestamp: time.Unix(r.Timestamp, 0).Format(time.RFC3339)` — converts Unix seconds to ISO 8601

### 4. gateway-service/cmd/server/main.go

- `handleGetUsage` checks `c.GetString("userRole")`:
  - If `"admin"` → sets `userID = ""` (see all records)
  - Otherwise → uses JWT userID as before
- Added `"timestamp": r.Timestamp` to JSON response map

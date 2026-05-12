# Tasks: Usage API Bugfixes

## Implementation

- [x] Remove redundant `addGroupIDColumn()` migration from billing-service
- [x] Handle empty userID in repository GetByUserID (return all records for admin)
- [x] Add Timestamp field to gateway UsageRecord + RFC3339 conversion
- [x] Add timestamp to JSON response in handleGetUsage
- [x] Add admin role check in handleGetUsage (pass empty userID for admins)
- [x] Rebuild and restart billing-service and gateway-service containers

## Verification

- [x] billing-service starts without migration crash (logs: "listening on :50054")
- [x] Regular user usage: `{"usage": []}` (filtered by own userID)
- [x] Admin user usage: returns ALL users' records (`user_id` shows multiple distinct IDs)
- [x] Timestamp in response: RFC3339 format (`"2026-05-12T02:59:54Z"`)
- [x] Admin timestamp renders as valid date in `new Date(timestamp).toLocaleString()`

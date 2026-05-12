# Proposal: Usage API Bugfixes — Admin View, Timestamp Format, Migration Redundancy

## Why

Three defects were discovered in the usage analytics workflow:

1. **Migration crash**: `billing-service` container crashes on startup with `duplicate column name: group_id` because the `ALTER TABLE` migration runs after the column is already created in the `CREATE TABLE` statement, making the service unreachable via Docker DNS.

2. **Invalid Date on frontend**: The `timestamp` field is dropped entirely between the gRPC response and HTTP JSON serialization in `gateway-service`, causing `new Date(undefined)` → `"Invalid Date"` in the admin UI.

3. **Admin users see no data**: `handleGetUsage` always passes the JWT `userId` to `billing-service`, which filters `WHERE user_id = ?`. An admin user only sees their own (likely empty) records instead of all users' usage.

## What Changes

- Remove redundant `addGroupIDColumn()` migration from `billing-service` (column already in `CREATE TABLE`)
- Add `Timestamp` field to gateway's `UsageRecord` struct and convert Unix seconds → RFC3339 string
- Add `"timestamp"` field to JSON response in `handleGetUsage`
- Check `userRole == "admin"` in gateway handler → pass empty `userID` to billing service
- Handle empty `userID` in billing repository → return all records (no `WHERE` clause)

## Impact

- `billing-service`: Migration list reduced by one; repository adds admin query path
- `gateway-service`: Timestamp correctly serialized; admin role triggers unfiltered query
- No protobuf, API contract, or database schema changes

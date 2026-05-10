## Context

The current admin system has several usability gaps identified in Dev B week 4 tasks. The auth service already has username support in the entity layer but it's not exposed through the API or UI. Group descriptions are being lost during creation because the proto definition doesn't include the description field in CreateGroupRequest. User group memberships exist in the database but aren't displayed in the UI. Group member counts require additional queries that aren't currently implemented.

## Goals / Non-Goals

**Goals:**
- Enable mandatory username-based authentication alongside email (username immutable after creation)
- Display user group memberships in admin UI
- Fix group description persistence through full create/update flow
- Show real-time member counts for groups
- Maintain backward compatibility with existing email-based authentication

**Non-Goals:**
- Complete migration from email to username authentication
- Advanced group hierarchy features
- Bulk user/group operations
- Permission system redesign

## Decisions

### Authentication Strategy
**Decision**: Support dual authentication (email OR username) with mandatory, immutable username
**Rationale**: Existing users and API integrations rely on email authentication. Dual support maintains backward compatibility while adding new functionality. Username is mandatory for new users and immutable after creation to ensure consistent identity. The auth service already has GetByUsername method, requiring only API layer changes.

### Group Description Fix
**Decision**: Add description field to CreateGroupRequest/UpdateGroupRequest proto messages
**Rationale**: The root cause is missing description in the proto definition. The backend GroupService and entity already support description properly. Adding the field to proto ensures end-to-end data flow.

### User Group Display Approach
**Decision**: Enrich user list responses with group membership data
**Rationale**: The frontend already has group data and display logic. The backend needs to include group IDs in user responses. This approach minimizes frontend changes and leverages existing group management components.

### Member Count Implementation
**Decision**: Calculate member counts at query time with database optimization
**Rationale**: Real-time accuracy is more important than performance for admin operations. We'll add a COUNT query with proper indexing rather than maintaining denormalized counts, avoiding synchronization complexity.

## Risks / Trade-offs

**Performance Risk**: Additional group queries for user lists may slow down user management
→ Mitigation: Implement eager loading with JOIN queries and add database indexes on user_group mappings

**Data Consistency Risk**: Username uniqueness conflicts during user creation
→ Mitigation: Add database unique constraint on username field and proper validation in both auth service and frontend

**API Compatibility Risk**: Proto changes may break existing clients
→ Mitigation: Make description field optional in proto messages and ensure backward compatibility

**Migration Complexity**: Existing users won't have usernames
→ Mitigation: Make username optional for existing users, required only for new registrations. Username is immutable after creation to prevent identity changes.

## Migration Plan

1. **Phase 1**: Update proto definitions and regenerate code
2. **Phase 2**: Update auth service handlers to include username in responses
3. **Phase 3**: Update gateway admin handlers to pass group data
4. **Phase 4**: Update frontend forms and displays
5. **Phase 5**: Add database migration for username uniqueness constraint
6. **Phase 6**: Testing and validation

**Rollback Strategy**: Proto changes are backward compatible. Database migration can be rolled back. Frontend changes are additive.

## Open Questions

- Should usernames be case-sensitive or case-insensitive?
- What's the maximum length for usernames?
- Should member counts be cached or calculated real-time?
- Do we need audit logging for username changes?

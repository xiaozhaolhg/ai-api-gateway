## Why

The current admin UI and backend APIs have gaps in user management and group administration that limit usability. Users cannot authenticate with usernames, group memberships are not visible in user lists, group descriptions are not being saved correctly, and admins cannot see member counts at a glance. These issues reduce the efficiency of daily administrative operations and create friction in user onboarding.

## What Changes

- **Username Support**: Add mandatory username field to user creation/registration with immutable property, allowing login with username instead of just email
- **User Group Display**: Show assigned groups for each user in the user management page with visual tags
- **Group Description Fix**: Ensure group descriptions are properly passed from frontend to backend during group creation and updates
- **Group Member Count**: Display member count for each group in the groups management page with real-time updates

## Capabilities

### New Capabilities
- `username-auth`: Username-based authentication and user management with unique username validation
- `user-group-display`: Visual display of user group memberships in admin interface
- `group-description-persistence`: Proper handling of group descriptions in CRUD operations
- `group-member-count`: Real-time member count calculation and display for groups

### Modified Capabilities
- `auth-service`: Enhanced user entity and authentication methods to support usernames
- `admin-ui-auth`: Updated user management forms and displays to include username and groups

## Impact

**Backend Services:**
- Auth service: Update User entity, CreateUser/UpdateUser methods, and authentication logic
- Gateway service: Update admin user handlers to include username and group data

**Frontend:**
- User management page: Add username field to forms and group display to lists
- Group management page: Fix description form handling and add member count column
- Authentication: Update login to support username/email dual authentication

**API Changes:**
- User proto: Add username field to User message and CreateUserRequest
- Group proto: Ensure description field is properly handled in CreateGroupRequest/UpdateGroupRequest
- Admin endpoints: Return user group memberships and group member counts in responses

**Database:**
- User table: Add unique constraint on username field
- Group queries: Add member count calculation with proper indexing

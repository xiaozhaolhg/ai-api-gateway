## Why

The current admin UI only has a login function implemented through Phase 2. Developers need a way to register new users during development and testing. Without a registration flow, adding test users or admin accounts requires direct database manipulation or CLI seeding, which is inconvenient for rapid iteration.

## What Changes

- Add Register RPC in auth-service (handles password hashing)
- Add proxy endpoint in gateway-service (`POST /admin/register`)
- Add registration page in admin-ui with form validation
- Wire registration to auth context for auto-login after successful registration

## Capabilities

### New Capabilities
- `admin-register`: User registration endpoint and UI page

### Modified Capabilities
- `admin-ui-auth`: Add registration to existing auth flow

## Impact

- **Code**: gateway-service (new register endpoint), auth-service (register handler), admin-ui (register page)
- **APIs**: `POST /admin/register`
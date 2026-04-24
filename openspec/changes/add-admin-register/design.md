## Context

Phase 2 implemented `/admin/login` but there's no registration. Developers need to create users during development without database/CLI access.

## Goals / Non-Goals

**Goals:**
- Add Register RPC in auth-service
- Add proxy endpoint in gateway-service: `POST /admin/register`
- Add registration page in admin-ui

**Non-Goals:**
- Email verification (dev-only, no SMTP)
- Password reset (handled separately)
- Production security (development stage)

## Architecture

```
admin-ui (register page)
    ↓ POST /admin/register
gateway-service (proxy)
    ↓ Register RPC
auth-service (Register handler + bcrypt)
```

## Decisions

| Decision | Rationale |
|----------|----------|
| Register in auth-service | Password hashing must be in auth-service |
| Proxy in gateway | Keep existing pattern (login already does this) |
| Username or email | Flexible for dev use |
| Auto-login after register | Better UX |

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| User enumeration | Dev-only exposure |
| No CAPTCHA | Dev stage |
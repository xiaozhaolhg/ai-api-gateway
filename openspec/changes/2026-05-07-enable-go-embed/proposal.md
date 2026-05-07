# Proposal: Enable go:embed for Single-Binary UI Serving

## Problem Statement

The `gateway-service` currently has `go:embed` directive **commented out** in `main.go:29`:
```go
// var staticFiles embed.FS // Temporarily disabled for testing
```

This prevents the single-binary build from serving the admin-ui static files, requiring a separate nginx container for UI deployment. The embedded static files already exist in `gateway-service/static/` (HTML, CSS, JS), and the `setupStaticFiles()` function is already implemented and called at line 212.

**Impact**:
- Single-binary demo (`make build-single`) produces a binary that cannot serve admin-ui
- Deployment requires separate nginx container (increased complexity)
- Contradicts Phase 1 Work Division (Developer C, Week 4: "Embed React build into Go binary")

## Proposed Solution

Uncomment the `go:embed` directive in `main.go`, verify static file serving works correctly, and ensure `make build-single` produces a working single-binary with embedded UI.

## Scope

- **In Scope**:
  - Uncomment `//go:embed static` directive in `main.go`
  - Verify `setupStaticFiles()` function works with embedded files
  - Test `make build-single` and verify UI is accessible
  - Update `Makefile` if needed (already has `build-single`, `build-ui`, `embed-ui` targets)

- **Out of Scope**:
  - Changes to admin-ui build process
  - New features for static file serving
  - Changes to other services' deployment

## Success Criteria

- `gateway-service/cmd/server/main.go:29` has active `//go:embed static` directive
- `make build-single` completes without errors
- Single binary serves admin-ui at `/` when no API route matches
- `GET /` returns HTML (admin-ui index.html)
- `GET /assets/*` returns JS/CSS files
- API routes (`/health`, `/v1/*`, `/admin/*`) still work correctly

## Dependencies

- `gateway-service/static/` directory must contain built admin-ui files
- `Makefile` targets: `build-ui`, `embed-ui`, `build-single` (already exist)
- `setupStaticFiles()` function (already implemented at line 306)

## Owner

- **Primary**: Developer C (cynkiller)
- **Collaborators**: Developer A (testing)

## References

- Original task: `docs/phase1_work_division.md` → Developer C → Week 4 → "Embed React build into Go binary"
- Current implementation: `gateway-service/cmd/server/main.go`
- Static files: `gateway-service/static/` (already contains index.html, assets/)

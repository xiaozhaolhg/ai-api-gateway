# Tasks: Enable go:embed for Single-Binary UI Serving

**Owner**: Developer C (cynkiller)  
**Collaborators**: Developer A (testing)  
**Status**: Planning

## Phase 1: Enable go:embed Directive (High Priority)

- [ ] **Task 1.1**: Uncomment go:embed in main.go
  - [ ] Edit `gateway-service/cmd/server/main.go` line 29
  - [ ] Change `// var staticFiles embed.FS // Temporarily disabled for testing`
  - [ ] To `//go:embed static` + `var staticFiles embed.FS`
  - **Acceptance**: `grep -A1 "go:embed" gateway-service/cmd/server/main.go` shows active directive

- [ ] **Task 1.2**: Verify static files exist
  - [ ] Check `gateway-service/static/` contains `index.html`, `assets/`, `favicon.svg`
  - [ ] If empty, run `make build-ui` to populate
  - **Acceptance**: `ls gateway-service/static/` shows at least `index.html`

## Phase 2: Build & Verification (High Priority)

- [ ] **Task 2.1**: Test `make build-single`
  - [ ] Run `make build-single` from project root
  - [ ] Verify binary is created at `gateway-service/bin/gateway`
  - [ ] Check binary size (should be larger due to embedded files, ~10-20MB)
  - **Acceptance**: Build completes without errors

- [ ] **Task 2.2**: Verify static files are embedded
  - [ ] Use `strings` command or hex dump to verify static files in binary
  - [ ] Or write a simple Go test to check `staticFiles` variable
  - **Acceptance**: Embedded files detectable in binary

- [ ] **Task 2.3**: Test UI is accessible
  - [ ] Start binary: `./gateway-service/bin/gateway`
  - [ ] `curl http://localhost:8080/` → should return HTML
  - [ ] `curl http://localhost:8080/assets/index-*.js` → should return JS
  - [ ] Verify SPA fallback: `curl http://localhost:8080/some-route` → returns index.html
  - **Acceptance**: UI pages load correctly from single binary

- [ ] **Task 2.4**: Test API routes still work
  - [ ] `curl http://localhost:8080/health` → JSON response
  - [ ] `curl http://localhost:8080/gateway/health` → JSON response
  - [ ] Verify admin API routes still work (if auth is set up)
  - **Acceptance**: All API routes functional alongside static files

## Phase 3: Testing & Documentation (Medium Priority)

- [ ] **Task 3.1**: Write unit test for static file embedding
  - [ ] Create `gateway-service/cmd/server/main_test.go` (if not exists)
  - [ ] Write test: verify `staticFiles` contains expected files
  - [ ] Test `setupStaticFiles()` function
  - **Acceptance**: Test passes; >80% coverage for static file serving

- [ ] **Task 3.2**: Update documentation
  - [ ] Update `QUICKSTART.md` to mention single-binary includes UI
  - [ ] Verify `README.md` reflects embedded UI capability
  - **Acceptance**: Documentation matches implementation

- [ ] **Task 3.3**: Cleanup
  - [ ] Remove any TODO comments related to go:embed
  - [ ] Search for "temporarily disabled" comments
  - **Acceptance**: No references to disabled go:embed

## Summary

| Phase | Tasks | Priority | Dependencies |
|-------|-------|----------|--------------|
| Phase 1: Enable go:embed | 2 | **High** | None |
| Phase 2: Build & Verification | 4 | **High** | Phase 1 |
| Phase 3: Testing & Documentation | 3 | **Medium** | Phase 2 |
| **Total** | **9** | | |

## Timeline Estimate

| Phase | Estimated Time | Owner |
|-------|----------------|--------|
| Phase 1 | 0.5 hour | Developer C |
| Phase 2 | 1.5 hours | Developer C + Dev A |
| Phase 3 | 1 hour | Developer C |
| **Total** | **3 hours** | |

## Critical Path

```
Phase 1 (Enable go:embed) → Phase 2 (Build & Verify) → Phase 3 (Test & Document)
```

**Blocker**: None. Phase 1 is independent and quick.

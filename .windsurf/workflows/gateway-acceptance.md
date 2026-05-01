---
description: Verify that a proposal is fully implemented before archiving by checking implementation alignment, design adherence, testing completeness, documentation, and coding standards
---

# Gateway Acceptance Workflow

Verifies that a proposal is fully implemented before archiving.

## Steps

### 1. Sync with Upstream
Use `gateway-sync-upstream` skill to ensure branch is up-to-date with upstream/main before acceptance.

### 2. Run Automated Verification
```bash
opsx-verify <change-name>
```

Review the verification report for:
- CRITICAL issues: Must fix before archiving
- WARNING issues: Should fix, document if deferring
- SUGGESTION issues: Nice to fix, can defer

### 3. Manual Verification Checklist

**Implementation vs Proposal:**
- All "What Changes" from proposal.md are implemented
- No files modified outside proposal scope
- Proto/database changes match proposal description

**Design Adherence:**
- All "Decision:" sections from design.md are followed
- Architecture decisions match (sync vs async, direct vs callback)
- No anti-patterns from AGENTS.md violated

**Spec Requirements:**
- All ADDED/MODIFIED/REMOVED/RENAMED requirements handled
- Scenarios have corresponding implementation
- Delta specs synced to main specs (merge into service folders, no new directories)

**Task Completion:**
- All tasks marked complete in tasks.md
- Acceptance criteria met
- Implementation notes accurate

**Testing:**
- New code has unit tests
- Integration tests for cross-service interactions
- Tests use proper mocking (testutils/ packages)
- Skipped tests have TODO explanations

**Documentation:**
- Public functions have godoc comments
- Design decisions documented in design.md
- README files updated for affected services
- No commented-out code or .bak files

**Coding Guidelines:**
- Code follows Go formatting (`go fmt`)
- No ignored errors
- No hardcoded credentials
- Context properly propagated

**Build and Deployment:**
- `make build` succeeds
- `docker compose build` succeeds (Dockerfiles copy new directories like pkg/)
- `make up` starts all services successfully
- Services are healthy

**Git Hygiene:**
- Commits are atomic
- Commit messages follow conventional format
- .gitignore updated for new generated files
- No backup files (.bak, .swp) committed

### 4. Fix Issues
Address all CRITICAL issues. For WARNING issues, either fix or document deferral reason.

### 5. Sync Specs
```bash
opsx-sync-specs <change-name>
```

### 6. Final Verification
```bash
make up
```
Run tests to ensure nothing broken.

Commit changes with clear message.

### 7. Archive
```bash
opsx-archive <change-name>
```

### 8. Update Work Division

After successful archival, mark the completed task in work division:

**Identify the completed task:**
- Match the archived change to the task in `docs/phase1_work_division.md` or `docs/phase2_work_division.md`
- Find the corresponding unchecked task (`- [ ]`)

**Update the task status:**
```markdown
- [x] [Task description] (completed: <change-name>)
```

**Example:**
```markdown
- [x] Define repository interfaces: ProviderRepo, UserRepo, APIKeyRepo, UsageRepo, RoutingRuleRepo (completed: repository-interfaces)
```

### 9. Commit Work Division Update

**Summarize proposal updates from `proposal.md`:**
- What Changes: [list key changes from proposal]
- Affected Services: [services from proposal]
- Capabilities: [capabilities from proposal]

**Create commit with summary:**
```bash
git add docs/phase1_work_division.md
git commit -m "docs(work-division): mark <task-description> as completed

Summary of changes:
- [Key change 1]
- [Key change 2]
- [Key change 3]

Affected: [services]
Capabilities: [capability names]
Archived: openspec/changes/archive/<change-name>/"
```

**Example:**
```bash
git commit -m "docs(work-division): mark repository interfaces as completed

Summary of changes:
- Define ProviderRepo, UserRepo, APIKeyRepo interfaces
- Define UsageRepo, RoutingRuleRepo interfaces
- Implement SQLite repository implementations
- Add database migration scripts for SQLite

Affected: auth-service, billing-service
Capabilities: auth-repository, billing-repository
Archived: openspec/changes/archive/repository-interfaces/"
```

## Common Issues

**Spec Sync:**
- Delta specs not synced to main specs
- New spec directories created instead of merging into service folders
- Technical updates in wrong folder (architecture vs service spec)

**Docker Build:**
- Dockerfiles don't copy new directories (e.g., pkg/)
- Services fail to start in docker compose

**Testing:**
- Tests mocked too heavily (test the mock, not the code)
- Integration tests skipped or TODO

## Gateway-Specific Considerations

**Service Boundaries:**
- Gateway: HTTP entry, middleware orchestration
- Auth: Identity, API keys, model authorization
- Billing: Usage tracking, budgets
- Router: Model → provider routing
- Provider: Provider CRUD, request forwarding
- Monitor: Metrics, alerting

**Anti-Patterns (from AGENTS.md):**
- NO direct database access across service boundaries
- NO hardcoded credentials
- NO URL-based provider detection
- Model naming: `{provider}:{model}`
- Each service owns its database exclusively
- Cross-service data flows through gRPC APIs

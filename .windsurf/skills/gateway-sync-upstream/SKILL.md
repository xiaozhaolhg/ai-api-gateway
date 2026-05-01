---
name: gateway-sync-upstream
description: Fetches and rebases the current branch with the latest upstream/main changes for the AI API Gateway project. Use before starting new work to ensure branch is up-to-date. Do NOT use when branch is already synced, when working on main branch, or when rebasing would disrupt ongoing work.
---

# Gateway Sync Upstream

Fetches and rebases the current branch with the latest upstream/main changes for the AI API Gateway project.

## Purpose

Ensure the current branch is up-to-date with upstream/main before starting new work.

## Process

### 1. Fetch Upstream
```bash
git fetch upstream main
```

If user environment is in WSL, execute:
```
export GIT_DISCOVERY_ACROSS_FILESYSTEM=1
```

### 2. Check Current Branch and Working Directory
```bash
git branch --show-current
git status --short
```

Note the branch name for conflict resolution context.

**If working directory is NOT clean (has uncommitted changes):**

Ask user: "Working directory has uncommitted changes. What would you like to do?"

Options:
- **stash**: `git stash push -m "WIP before upstream sync"` - Save changes temporarily
- **commit**: `git add <files> && git commit -m "WIP: describe changes"` - Commit changes first
- **drop**: `git checkout -- . && git clean -fd` - Discard all changes (WARNING: data loss)
- **abort**: Stop sync and let user handle manually

Wait for user choice before proceeding.

### 3. Rebase onto Upstream
```bash
git rebase upstream/main
```

### 4. Handle Conflicts
If conflicts occur:
- Call `git-rebase-workflow` to resolve
- Apply learned patterns for gateway project:
  - Preserve upstream structural changes (imports, client initialization like billingClient, authClient)
  - Merge local logic into new structure
  - Verify dependencies (clients initialized before use)
  - Keep upstream's file organization

### 5. Verify Result
```bash
git log --oneline -5
git diff upstream/main
```

Ensure rebase completed successfully.

**If changes were stashed in Step 2:**
```bash
git stash pop
```
Restore stashed changes after successful rebase.

## Gateway-Specific Conflict Patterns

### Client Initialization Conflicts
When upstream adds new gRPC clients (e.g., billingClient, monitorClient):
- Accept upstream's client initialization pattern
- Merge local handler usage into new structure
- Verify client is initialized before handler routes are registered

### Import Conflicts
When upstream adds new proto imports:
- Accept all upstream imports
- Add any local-only imports that don't duplicate
- Remove unused imports after resolution

### Handler Structure Conflicts
When upstream reorganizes handler functions:
- Follow upstream's direct gRPC call pattern
- Adapt local logic to upstream structure
- Preserve local business rules

## Success Indicators

- `git status` shows clean working directory
- `git log` shows linear history from upstream/main
- No conflict markers in any file
- `go build ./...` passes for all services
- Tests pass (if available)

## Failure Recovery

If rebase fails repeatedly:
- Abort: `git rebase --abort`
- Create fresh branch from upstream/main
- Cherry-pick specific commits if needed
- Ask user for manual intervention

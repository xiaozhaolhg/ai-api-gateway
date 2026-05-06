---
description: Wrapper workflow around OpenSpec commands (opsx-explore/opsx-new/opsx-propose) that collects project-specific context and integrates it with standard OpenSpec process for the AI API Gateway
---

# Gateway New Proposal Workflow

Wrapper around OpenSpec commands that collects local context from the user and project knowledge (AGENTS.md, work division) to create a proposal aligned with gateway architecture patterns.

## Purpose

This workflow integrates:
- **User input**: Task selection, scope clarification, affected services
- **Project knowledge**: Service boundaries, anti-patterns, spec naming conventions from AGENTS.md
- **Work division**: Task context from phase1_work_division.md

The collected information is fed into OpenSpec commands to generate a proposal that reflects gateway-specific guardrails and requirements.

## Steps

### 1. Sync with Upstream
Use `gateway-sync-upstream` skill to sync branch with upstream/main.

### 2. Collect Context

**2.1 Identify Task from Work Division**

Read `docs/phase1_work_division.md` (or phase2).

Ask: "Which developer role? (A/B/C)"

Locate developer section, find first unchecked task (`- [ ]`).

Present to user:
```
Task: [description]
Acceptance Criteria: [criteria]
Affected Services: [services from context]
Dependencies: [previous tasks]

Proceed with this task? (yes/no)
```

**2.2 Gather Change Details**

Ask:
- "Change name? (kebab-case, e.g., 'auth-middleware')"
- "Scope clarification needed?"

**2.3 Explore if Uncertain**

If scope unclear:
```bash
opsx-explore "[topic - current implementation, architecture, API contracts]"
```

Summarize findings:
- Current state
- Gaps identified
- Proposed approach

Confirm with user: "Focus on X, Y, Z?"

**2.4 Identify Affected Services & Dependencies**

From task context and user input, identify:
- **Primary service**: Which service owns the change
- **Dependent services**: Which services need updates
- **External dependencies**: Libraries, proto changes, database migrations

Apply **Gateway Service Boundaries** from AGENTS.md:
- Gateway: HTTP entry, middleware orchestration
- Auth: Identity, API keys, model authorization
- Billing: Usage tracking, budgets
- Router: Model → provider routing
- Provider: Provider CRUD, request forwarding
- Monitor: Metrics, alerting

**Guardrails to apply:**
- NO direct database access across service boundaries
- NO hardcoded credentials
- NO URL-based provider detection
- Model naming: `{provider}:{model}`
- Cross-service data flows through gRPC APIs

### 3. Create Proposal

Feed collected context into OpenSpec:

**Option A - Full generation:**
```bash
opsx-propose <change-name>
```

**Option B - Step by step:**
```bash
opsx-new <change-name>
```
Then create artifacts guided by collected context.

### 4. Enrich Proposal with Gateway Context

Ensure the generated proposal reflects collected information:

**In proposal.md:**
- **Why**: Connect to gateway architecture needs (from work division context)
- **What Changes**: List all affected services identified in step 2.4
- **Capabilities**: Follow naming `{service-name}` or `{service-name}-{type}` (architecture/deployment/testing)
- **Impact**: Note gRPC API changes, database migrations, middleware updates

**In design.md:**
- **Context**: Reference microservices architecture from AGENTS.md
- **Decisions**: Apply service boundaries, document gRPC vs HTTP choices
- **Risks**: Note cross-service dependencies identified in step 2.4

**In specs/:**
Follow `openspec/config.yaml` conventions:
- Service specs: `{service-name}`
- Architecture specs: `{service-name}-architecture`
- Deployment specs: `{service-name}-deployment`
- Testing specs: `{service-name}-testing`
- System specs: `system` or `system-deployment`

**In tasks.md:**
- Group by service or layer (domain, application, handler)
- Include unit test tasks (per config.yaml rules)
- Include integration test tasks (for live service testing)
- Include FVT tasks
- **Note:** Integration tests must be run against live services in Docker Compose environment

**Live Service Integration Testing Requirements:**
When planning tasks, ensure integration tests include:
1. Environment setup: `make down && make clean-images && make up`
2. Health check verification: `docker compose ps`
3. Cross-service gRPC call verification
4. End-to-end request flow testing
5. Log verification: `docker compose logs -f [service-name]`

### 5. Validate Alignment

Verify the proposal reflects the original task context:
- Does proposal scope match work division task?
- Are all affected services identified in step 2.4 included?
- Are acceptance criteria from work division reflected in tasks.md?
- Are gateway guardrails (anti-patterns) considered?

## Output

A complete OpenSpec change where:
- `proposal.md` incorporates work division task context
- `design.md` references gateway architecture patterns
- `specs/` follow gateway naming conventions
- `tasks.md` include gateway-specific testing requirements

## Developer Roles Reference

- **Developer A**: Router, Provider, Monitor services
- **Developer B**: Auth, Data Access Layer, Token Tracker
- **Developer C**: Admin UI, Frontend, Integration

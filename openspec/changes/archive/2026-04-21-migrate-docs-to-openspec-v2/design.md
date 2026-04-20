## Context

Design docs under docs/ are flat files:
- docs/architecture_design.md - system architecture
- docs/service_interaction.md - gRPC API contracts
- docs/business_requirement.md - product phases

These lack service-level organization and dependency tracking in OpenSpec.

## Goals / Non-Goals

**Goals:**
- Create system/ folder with service relationship diagram and shared API contracts
- Create per-service spec folders with architecture.md (design + dependencies)
- Show calling relationships in system spec, service dependencies in each service spec

**Non-Goals:**
- No implementation code (spec management only)
- No deployment/testing specs — just architecture.md per service
- No duplicate router-service specs — reference existing

## Decisions

### Service Structure

Each service gets one folder under openspec/specs/:

```
openspec/specs/
├── system/
│   ├── architecture.md     ← service relationships, layered view, calling diagram
│   └── api-contracts.md  ← shared gRPC definitions
├── gateway-service/
│   └── architecture.md  ← design + calls to auth, router, provider, billing
├── auth-service/
│   └── architecture.md  ← design + returns data to gateway
├── router-service/
│   └── architecture.md  ← reference to existing specs
├── provider-service/
│   └── architecture.md  ← design + callbacks to billing, monitor
├── billing-service/
│   └── architecture.md  ← design + receives callbacks
└── monitor-service/
    └── architecture.md  ← design + receives callbacks
```

### Service Relationships (System)

System architecture.md includes:
- Layer view: Edge → Service → Provider → Data
- Calling relationships: which service calls which
- Provider callback pattern visualization

### Service Dependencies (Each Service)

Each service's architecture.md includes:
- **Calls To**: gRPC endpoints this service invokes
- **Called By**: services that call this service
- **Data Dependencies**: databases, caches used
- **Key Operations**: primary API methods

## Content Structure

### system/architecture.md

- Layered architecture diagram
- Service calling relationship diagram
- Service list with responsibilities
- Provider callback mechanism

### system/api-contracts.md

- Shared message definitions (UserIdentity, RouteResult, TokenCounts, etc.)
- Common error codes

### {service}/architecture.md

- Service responsibility summary
- **Calls To**: [service] → methods used
- **Called By**: [service] → methods received
- Key design principles

## Migration Plan

1. Create directories under openspec/specs/
2. Create system/architecture.md with service relationship diagrams
3. Create system/api-contracts.md with shared types
4. Create each service's architecture.md with dependencies section
5. Reference router-service specs (do not recreate)

## Risks / Trade-offs

| Risk | Mitigation |
|---|---|
| Duplicate content with existing router-service specs | Reference existing, do not copy |
| Inconsistent depth across services | Use same template structure for all |
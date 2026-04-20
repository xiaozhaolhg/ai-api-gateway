## Context

The router-service needs a standardized deployment workflow. Currently, developers must manually run docker build, then either docker run or kubectl/helm commands to test changes. This introduces inconsistency and friction.

The project structure:
```
router-service/
├── Dockerfile (multi-stage Go build)
├── charts/ (Helm chart)
├── configs/config.yaml
└── ...
```

## Goals / Non-Goals

**Goals:**
- Provide make targets for building docker image
- Provide make target for running in local Docker
- Provide make target for deploying to KinD via Helm
- Remove deprecated k8s.yaml manifest

**Non-Goals:**
- CI/CD pipeline (out of scope)
- Production K8s deployment (only local Docker and KinD)

## Decisions

### Decision 1: Image Name

**Choice:** `router-service:latest`

**Rationale:** Simple, consistent naming. Matches existing usage in k8s.yaml.

### Decision 2: Kind Configuration

**Choice:** Make KIND_CLUSTER configurable via environment variable

**Rationale:** Teams may have multiple KinD clusters. `?=desktop` allows:
```bash
make deploy-kind KIND_CLUSTER=my-cluster
```

### Decision 3: Helm Chart Path

**Choice:** `./charts`

**Rationale:** Chart contents (Chart.yaml, templates/, values.yaml) live directly in `charts/` — the chart name is derived from Chart.yaml, not the directory name. The `./` prefix ensures Helm treats it as a local path rather than a repository name.

### Decision 4: Namespace

**Choice:** `ai-gateway`

**Rationale:** Matches existing namespace in deleted k8s.yaml.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Kind cluster not running | Show clear error message |
| Helm not installed | Document prerequisite in Makefile header |
| Port 8080 conflicts | Document in Makefile header |

## Migration Plan

1. Create Makefile in router-service/
2. Delete k8s.yaml and deploy/ directory (Done)
3. Developers use `make deploy-docker` or `make deploy-kind`

## Open Questions

None.
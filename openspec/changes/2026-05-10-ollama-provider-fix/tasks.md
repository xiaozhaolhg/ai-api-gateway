# Tasks: Fix Ollama Provider Integration and Admin-UI Deployment

## Implementation Tasks

- [x] Fix gateway route.go - add c.Set("providerId")
- [x] Fix admin-ui nginx.conf - change /admin to /admin/auth
- [x] Fix admin-ui config.ts - default baseURL to empty string
- [x] Fix admin-ui Dockerfile - add VITE_API_BASE_URL ARG
- [x] Fix docker-compose.yaml - add extra_hosts and build context
- [x] Fix provider-service service.go - use bytes.NewReader for body
- [x] Fix ollama_adapter.go - strip prefix, add version suffix, use messages

## Verification Tasks

- [x] Test chat completions API with ollama:qwen3.5 model
- [x] Verify admin-ui auth endpoints return 200
- [x] Verify provider requests include body data
- [x] Verify ollama returns valid response with content

## Documentation Tasks

- [x] Create proposal.md
- [x] Create design.md
- [ ] Create this tasks.md

## Status: COMPLETED

All implementation and verification tasks have been completed. The system now successfully:
1. Routes requests to Ollama provider via gateway
2. Returns valid chat completion responses with actual content
3. Admin-UI auth endpoints accessible via nginx proxy
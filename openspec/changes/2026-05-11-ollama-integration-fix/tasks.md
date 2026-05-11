# Tasks: Ollama Integration Fix

## Implementation

- [x] Fix provider-service Authorization header to use credentials not URL
- [x] Add /api/chat path to Ollama URL
- [x] Support OLLAMA_BASE_URL environment variable
- [x] Fix Ollama adapter to not pass through user Authorization header
- [x] Add debug logging for troubleshooting
- [x] Fix Docker networking (all services on ai-gateway network)
- [x] Fix proxy middleware to get groupID from groupIds array
- [x] Add group_id column migration to billing-service
- [x] Test end-to-end API call with Ollama provider

## Verification

- [x] API call succeeds: `curl -X POST http://localhost:8080/v1/chat/completions -d '{"model":"ollama:qwen3.5:0.8b","messages":[{"role":"user","content":"Hi"}]}'`
- [x] Response returned with valid content
- [x] Usage recording shows userID and groupID in logs
- [x] Billing service connection succeeds (no DNS errors)

## Cleanup (Not Done - Optional)

- [ ] Remove debug logging statements
- [ ] Add unit tests for new code paths
- [ ] Document Ollama provider configuration
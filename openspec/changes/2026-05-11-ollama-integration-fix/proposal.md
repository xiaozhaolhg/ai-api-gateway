# Proposal: Ollama Integration Fix

## Summary

Fix critical issues in the Ollama provider integration to enable successful end-to-end API calls from gateway to local Ollama service, and ensure proper usage recording to billing service.

## Problem Statement

1. **Authorization Header Bug**: Provider service incorrectly uses target URL as Bearer token instead of actual credentials, causing "405 Method Not Allowed" errors
2. **Missing API Path**: Requests to Ollama missing `/api/chat` path suffix
3. **Header Passthrough**: User's API key Authorization header passes through to provider, conflicting with provider credentials
4. **Network Isolation**: Gateway and billing-service on different Docker networks, causing connection failures
5. **GroupID Not Passed**: Usage recording missing group information from authorization context

## Proposed Solution

1. Fix Authorization header to use provider credentials, not target URL
2. Append `/api/chat` to Ollama base URL
3. Filter out user Authorization header in Ollama adapter
4. Add `networks: ai-gateway` to all services in docker-compose.yml
5. Fix usage recording to get groupID from `groupIds` array

## Expected Outcomes

- Successful API calls to Ollama provider
- Usage records include user_id and group_id
- All services communicate properly over Docker network

## Timeline

- 1 hour: Fix provider-service authorization and URL path
- 1 hour: Fix Docker networking
- 1 hour: Fix usage recording with groupID

## Risks

- Network changes may affect existing deployments
- Debug logging adds verbosity (should be removed in production)
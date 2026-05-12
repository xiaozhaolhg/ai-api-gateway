# Tasks: Provider-Service Non-Streaming Timeout Fix

## Implementation

- [ ] Change `makeHTTPRequest()` HTTP client timeout from 60s to 0 in provider-service

## Verification

- [ ] Previously failing request ("为什么天是蓝色的？") completes successfully
- [ ] Previously passing requests ("你是谁？") still work
- [ ] Streaming requests unaffected
- [ ] Provider-service unit tests pass (`go test ./...`)

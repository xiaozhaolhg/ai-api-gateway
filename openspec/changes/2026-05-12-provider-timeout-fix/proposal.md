# Proposal: Increase Non-Streaming Forward Request Timeout for Long-Thinking Models

## Why

The qwen3.5:0.8b model generates extended "thinking" traces (chain-of-thought reasoning) that can take significantly longer than 60 seconds to complete via Ollama's `/api/chat` endpoint. When queried with open-ended or scientific questions (e.g., "为什么天是蓝色的？"), the model's reasoning process produces thousands of thinking tokens before the final answer.

The provider-service's non-streaming HTTP client has a hardcoded 60-second timeout (`Timeout: 60 * time.Second`), which causes the request to fail with `Client.Timeout exceeded while awaiting headers` before Ollama finishes generating the response. Streaming requests already work correctly because their HTTP client uses `Timeout: 0` (no timeout).

## What Changes

- In `provider-service/internal/application/service.go`: Change `makeHTTPRequest()`'s HTTP client timeout from `60 * time.Second` to `0` (no timeout), matching the streaming version's behavior.
- The context from the gRPC handler already provides the necessary deadline/cancellation control from upstream callers (the gateway, the HTTP client), so removing the redundant HTTP client timeout is safe.

## Impact

- **provider-service**: Non-streaming requests to Ollama no longer have an artificial 60s ceiling. Long-thinking models can complete naturally.
- **No changes** to the gateway-service, proto definitions, configuration, or deployment.
- **No database migrations** needed.

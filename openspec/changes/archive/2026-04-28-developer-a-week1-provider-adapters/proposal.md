## Why

As part of Phase 1 Work Division (Developer A - Week 1), we need to establish the foundational interfaces and implementations for provider adapters and routing. This change implements the core abstractions that enable the gateway to:

1. Transform requests/responses between OpenAI-compatible format and various provider formats (OpenAI, Anthropic)
2. Handle both streaming (SSE) and non-streaming responses with unified token counting
3. Provide a clean Router interface for model-to-provider resolution

These interfaces are prerequisites for Week 2+ integration work where gateway-service will orchestrate the full request flow.

## What Changes

### Provider Adapter Framework
- **Update `ProviderAdapter` interface** to support streaming via unified `TransformResponse` method with `isStreaming` flag
- **Create `StreamingResult` entity** to encapsulate transformed data, accumulated token counts, and finalization status
- **Extend `TokenCounts` entity** with `AccumulatedTokens` field for streaming scenarios

### Provider Implementations
- **OpenAI Adapter**: Pass-through for non-streaming, SSE chunk transformation with token accumulation for streaming
- **Anthropic Adapter**: OpenAI↔Anthropic format conversion for both streaming and non-streaming, message role mapping, stop reason conversion

### Router Interface
- **Define `Router` interface** in router-service domain layer with full CRUD operations for routing rules plus route resolution

### Testing
- Complete unit tests for all adapters covering transformation correctness, token accumulation accuracy, and error scenarios

## Capabilities

### Modified Capabilities
- `provider-service-adapter-framework`: Extend ProviderAdapter interface with unified TransformResponse supporting streaming via isStreaming flag; add StreamingResult and TokenCounts domain entities for SSE support
- `router-service-routing`: Define Router domain interface with full CRUD operations for routing rules plus route resolution

### Affected Mandatory Specs
- `provider-service-architecture`: Add StreamingResult and TokenCounts entities; document unified TransformResponse interface design
- `provider-service-testing`: Add test scenarios for SSE streaming, token accumulation, and error handling
- `router-service-architecture`: Add Router domain interface definition following Clean Architecture principles

## Impact

- **provider-service**: New domain entities, updated adapter interface, SSE implementations
- **router-service**: New domain interface file
- **api/proto**: No changes (existing protobuf definitions sufficient)
- **gateway-service**: No direct changes (uses interfaces through gRPC)

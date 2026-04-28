## ADDED Requirements

### Router Service Redis Caching

Router service shall cache resolved routes in Redis with TTL.

#### Scenario: Cache Hit
- **WHEN** a `ResolveRoute` request is received and the route exists in Redis cache
- **THEN** return the cached route immediately without querying the database

#### Scenario: Cache Miss
- **WHEN** a `ResolveRoute` request is received and the route is not in cache
- **THEN** query the database, cache the result in Redis with 5-minute TTL, then return

#### Scenario: Cache Invalidation
- **WHEN** `RefreshRoutingTable` is called after provider configuration changes
- **THEN** clear all routing-related cache keys to force fresh lookups

### Router Service Authorized Models Filtering

Router service shall filter routes based on authorized models passed from gateway.

#### Scenario: Authorized Route Resolution
- **WHEN** `ResolveRoute` is called with a model and `authorized_models` list
- **THEN** only return routes for providers serving models in the authorized list

#### Scenario: Unauthorized Model Request
- **WHEN** the requested model is not in `authorized_models`
- **THEN** return NOT_FOUND error without querying providers

### Gateway Service Streaming Proxy

Gateway service shall proxy SSE streaming responses from providers to consumers.

#### Scenario: Streaming Request
- **WHEN** a chat completion request with `stream: true` is received
- **THEN** establish SSE connection to consumer and stream chunks from provider-service

#### Scenario: Token Accumulation During Streaming
- **WHEN** processing SSE chunks from provider
- **THEN** accumulate token counts across all chunks and report final count on completion

### Gateway Service Non-Streaming Proxy

Gateway service shall proxy non-streaming requests to providers.

#### Scenario: Non-Streaming Request
- **WHEN** a chat completion request without `stream: true` (or `stream: false`)
- **THEN** call `ForwardRequest`, return complete response with token counts

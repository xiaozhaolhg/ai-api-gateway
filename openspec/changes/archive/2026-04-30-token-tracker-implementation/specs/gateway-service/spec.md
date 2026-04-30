## MODIFIED Requirements|

### Requirement: Token recording after non-streaming requests|
The gateway-service SHALL call billing-service `RecordUsage` RPC after completing a non-streaming LLM request, passing user_id (from JWT), group_id (from ValidateAPIKey response), provider_id, model, and token counts from the provider response.|

#### Scenario: Non-streaming request completes|
- **WHEN** a non-streaming request to a provider completes|
- **THEN** the gateway extracts prompt_tokens and completion_tokens from the ForwardRequestResponse.TokenCounts|
- **AND** calls `billingClient.RecordUsage()` with user_id, group_id, provider_id, model, and token counts|

#### Scenario: Multiple providers called (non-streaming)|
- **WHEN** a request fans out to multiple providers/models|
- **THEN** the gateway SHALL call `RecordUsage` separately for each provider/model combination|
- **AND** each call includes the correct provider_id and model for that specific call|

### Requirement: Token recording after streaming requests|
The gateway-service SHALL accumulate token counts across all SSE chunks and call billing-service `RecordUsage` after the stream completes.|

#### Scenario: Streaming request completes|
- **WHEN** a streaming LLM request completes (all SSE chunks received)|
- **THEN** the gateway sums totalPromptTokens and totalCompletionTokens across all chunks|
- **AND** calls `billingClient.RecordUsage()` with the accumulated totals|

#### Scenario: Multiple streaming providers|
- **WHEN** a streaming request calls multiple providers|
- **THEN** the gateway SHALL call `RecordUsage` separately for each provider after their stream completes|
- **AND** each call uses the correct provider_id and accumulated tokens for that provider|

### Requirement: Token extraction from provider responses|
The gateway-service SHALL extract token counts from provider responses for both non-streaming (JSON response) and streaming (SSE chunks) flows.|

#### Scenario: Non-streaming token extraction|
- **WHEN** a non-streaming response is received from provider-service|
- **THEN** the gateway reads `resp.TokenCounts.PromptTokens` and `resp.TokenCounts.CompletionTokens`|
- **AND** makes these available for `RecordUsage` call|

#### Scenario: Streaming token accumulation|
- **WHEN** SSE chunks are received from provider-service|
- **THEN** the gateway updates running totals: `totalPromptTokens += chunk.AccumulatedTokens.PromptTokens`|
- **AND** `totalCompletionTokens += chunk.AccumulatedTokens.CompletionTokens`|
- **AND** after stream completion, uses these totals for `RecordUsage`|

### Requirement: Group ID propagation for token recording|
The gateway-service SHALL use the user's group_id (from `ValidateAPIKey` response) when calling `RecordUsage`.|

#### Scenario: Single group user|
- **WHEN** a user belongs to one group|
- **THEN** the gateway uses that group's ID as `group_id` in `RecordUsage`|

#### Scenario: Multi-group user (MVP)|
- **WHEN** a user belongs to multiple groups|
- **THEN** the gateway uses the FIRST group ID from `UserIdentity.group_ids` as `group_id` in `RecordUsage`|
- **Note**: Future enhancement will record separate UsageRecord for each group|

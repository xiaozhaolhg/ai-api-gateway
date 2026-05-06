## MODIFIED Requirements

### Requirement: Token recording after streaming requests
The gateway-service SHALL accumulate token counts across all SSE chunks and call billing-service `RecordUsage` **at configurable intervals during streaming** AND after the stream completes.

#### Scenario: Streaming request with interval recording
- **WHEN** a streaming LLM request is in progress
- **AND** the accumulated completion tokens exceed the configured `streaming_token_interval` threshold since the last recording
- **THEN** the gateway calls `billingClient.RecordUsage()` with the delta tokens since the last recording
- **AND** updates the last-recorded position to avoid double-counting

#### Scenario: Streaming request final recording
- **WHEN** a streaming LLM request completes (all SSE chunks received or stream errors)
- **THEN** the gateway calls `billingClient.RecordUsage()` with any remaining delta tokens not yet recorded
- **AND** the final call ensures all tokens are accounted for

#### Scenario: Streaming request below threshold
- **WHEN** a streaming LLM request completes without reaching the token interval threshold
- **THEN** the gateway calls `billingClient.RecordUsage()` once with the full accumulated totals
- **AND** behavior is identical to the previous single-call approach

#### Scenario: Multiple streaming providers
- **WHEN** a streaming request calls multiple providers
- **THEN** the gateway SHALL track recording state independently per provider stream
- **AND** call `RecordUsage` separately for each provider based on their accumulated tokens

## ADDED Requirements

### Requirement: Configurable streaming token interval
The gateway-service SHALL support a `streaming_token_interval` configuration parameter that controls how frequently intermediate usage is recorded during streaming.

#### Scenario: Default interval
- **WHEN** no `streaming_token_interval` is configured
- **THEN** the default value SHALL be 1000 completion tokens

#### Scenario: Custom interval via config
- **WHEN** `streaming_token_interval` is set in the gateway-service configuration
- **THEN** intermediate `RecordUsage` calls are triggered every N completion tokens

#### Scenario: Interval set to zero
- **WHEN** `streaming_token_interval` is set to 0
- **THEN** intermediate recording is disabled
- **AND** usage is recorded only at stream completion (legacy behavior)

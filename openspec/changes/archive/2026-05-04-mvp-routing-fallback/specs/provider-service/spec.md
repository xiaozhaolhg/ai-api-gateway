## ADDED Requirements

### Requirement: OpenCode Zen Adapter Inference
Provider service SHALL automatically infer `opencode-zen` adapter type for providers with matching endpoint or type identifier.

#### Scenario: Infer opencode-zen adapter from endpoint
- **WHEN** a provider is created with endpoint containing `opencode.ai/zen`
- **THEN** set the provider's adapter type to `opencode-zen` automatically

#### Scenario: Infer opencode-zen adapter from type
- **WHEN** a provider is created with type `opencode-zen`
- **THEN** set the provider's adapter type to `opencode-zen` automatically

#### Scenario: Explicit adapter override
- **WHEN** a provider is created with explicit adapter type set to `opencode-zen`
- **THEN** use the explicitly configured adapter type regardless of endpoint

## MODIFIED Requirements

<!-- No existing requirements modified for this change -->

## REMOVED Requirements

<!-- No requirements removed for this change -->

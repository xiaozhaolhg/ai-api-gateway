## MODIFIED Requirements

### Requirement: Configuration via YAML and environment variables
The gateway-service SHALL load configuration from a YAML file with environment variable resolution, including gRPC addresses for all 5 internal services AND the `streaming_token_interval` parameter.

#### Scenario: Streaming token interval from config
- **WHEN** the config file specifies `streaming_token_interval: 500`
- **THEN** the gateway SHALL record usage every 500 completion tokens during streaming

#### Scenario: Streaming token interval from environment
- **WHEN** the environment variable `STREAMING_TOKEN_INTERVAL` is set
- **THEN** it SHALL override the YAML config value

#### Scenario: Default streaming token interval
- **WHEN** neither config nor environment variable sets `streaming_token_interval`
- **THEN** the default value of 1000 SHALL be used

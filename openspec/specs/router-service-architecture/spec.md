## Purpose

(TBD) Architecture specifications for the router-service component.

## Requirements

### Requirement: Clean Architecture Layer Structure
The router-service SHALL implement four-layer Clean Architecture: Domain, Application, Infrastructure, and Interface layers with clear dependency direction from outer to inner layers.

#### Scenario: Domain layer has no external dependencies
- **WHEN** the domain layer is imported
- **THEN** it SHALL NOT import any code from application, infrastructure, or interface layers

#### Scenario: Infrastructure implements domain interfaces
- **WHEN** a new LLM provider needs to be added
- **THEN** it SHALL be implemented by creating a struct that implements the Provider interface defined in the domain layer

### Requirement: Provider Abstraction
The router-service SHALL define a Provider interface in the domain layer that abstracts the operations: ChatCompletion, ListModels, and StreamChatCompletion.

#### Scenario: Request routing to appropriate provider
- **WHEN** a chat completion request is received with model "ollama:llama3"
- **THEN** the router SHALL route the request to the OllamaProvider implementation

#### Scenario: Request routing to OpenCode Zen
- **WHEN** a chat completion request is received with model "opencode:gpt-5.3-codex"
- **THEN** the router SHALL route the request to the OpenCodeZenProvider implementation

### Requirement: Request/Response Transformation
The router-service SHALL transform requests from OpenAI format to provider-specific formats and responses from provider format back to OpenAI format.

#### Scenario: Non-streaming request transformation
- **WHEN** a request with "stream": false is received at /v1/chat/completions
- **THEN** the service SHALL forward a correctly formatted request to the upstream provider and return the response in OpenAI format

#### Scenario: Streaming request transformation
- **WHEN** a request with "stream": true is received at /v1/chat/completions
- **THEN** the service SHALL stream SSE chunks in OpenAI format: `data: {"choices":[{"delta":{"content":"..."}}]}\n\n`

### Requirement: OpenAI-Compatible Endpoints
The router-service SHALL implement the following OpenAI-compatible endpoints:

#### Scenario: Chat completions endpoint
- **WHEN** POST /v1/chat/completions is called with valid request body
- **THEN** the service SHALL return a response in OpenAI chat completions format

#### Scenario: List models endpoint
- **WHEN** GET /v1/models is called
- **THEN** the service SHALL return a list of available models from all configured providers

#### Scenario: Health check endpoint
- **WHEN** GET /health is called
- **THEN** the service SHALL return `{"status": "ok"}` with HTTP 200

### Requirement: Configuration Management
The router-service SHALL load configuration from a YAML file mounted at /app/config/config.yaml.

#### Scenario: Configuration loads successfully
- **WHEN** the service starts with a valid config.yaml
- **THEN** all configured providers SHALL be initialized and ready to handle requests

#### Scenario: Provider uses environment variable
- **WHEN** a provider's api_key field contains "${ENV_VAR_NAME}"
- **THEN** the service SHALL resolve the value from the corresponding environment variable

### Requirement: Streaming Response Handling
The router-service SHALL properly handle streaming responses including proper SSE formatting and connection management.

#### Scenario: Streaming connection stays open
- **WHEN** a streaming request is processed
- **THEN** the connection SHALL remain open until the provider finishes streaming or the client disconnects

#### Scenario: Streaming sends correct format
- **WHEN** chunks are received from the upstream provider
- **THEN** each chunk SHALL be formatted as `data: <json>\n\n` with final `data: [DONE]`
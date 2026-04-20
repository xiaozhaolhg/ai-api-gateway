# router-service-provider-factory

## Purpose

Enable extensible provider registration system using factory pattern for the AI API Gateway router service.

## ADDED Requirements

### Requirement: Provider registration uses factory pattern
The router-service SHALL use a factory pattern for provider registration where each provider type implements a ProviderFactory interface with methods for type identification, creation, validation, defaults, and description.

#### Scenario: Register built-in providers
- **WHEN** the application starts
- **THEN** the system SHALL register OllamaFactory and OpenCodeZenFactory with the ProviderRegistry
- **AND** each factory SHALL provide a unique type identifier

#### Scenario: Create provider from factory
- **WHEN** the registry receives a provider type and settings
- **THEN** the system SHALL call the corresponding factory's Create method
- **AND** return an instantiated Provider interface

### Requirement: Provider registry manages factory lifecycle
The router-service SHALL implement a ProviderRegistry that manages factory registration, provider instantiation, and provides discovery capabilities with thread-safe operations.

#### Scenario: Register new factory
- **WHEN** a factory is registered with the registry
- **THEN** the system SHALL store the factory by its type identifier
- **AND** reject duplicate registrations with an error

#### Scenario: Create provider with valid settings
- **WHEN** creating a provider with valid settings
- **THEN** the registry SHALL validate the settings using the factory's Validate method
- **AND** instantiate the provider using the factory's Create method
- **AND** return the provider instance

#### Scenario: Create provider with invalid settings
- **WHEN** creating a provider with invalid settings
- **THEN** the registry SHALL return a validation error
- **AND** not instantiate the provider

#### Scenario: Create provider with unknown type
- **WHEN** creating a provider with an unknown type
- **THEN** the registry SHALL return an error listing available provider types
- **AND** not instantiate any provider

### Requirement: Provider factory provides validation and defaults
Each provider factory SHALL validate required settings and provide default values for optional settings to ensure proper provider configuration.

#### Scenario: Ollama factory validates endpoint
- **WHEN** validating Ollama provider settings
- **THEN** the factory SHALL require a non-empty endpoint
- **AND** return an error if endpoint is missing

#### Scenario: Ollama factory provides defaults
- **WHEN** requesting Ollama factory defaults
- **THEN** the factory SHALL return endpoint as "http://localhost:11434"
- **AND** enabled as false
- **AND** api_key as empty string

#### Scenario: OpenCode Zen factory validates endpoint
- **WHEN** validating OpenCode Zen provider settings
- **THEN** the factory SHALL require a non-empty endpoint
- **AND** return an error if endpoint is missing

#### Scenario: OpenCode Zen factory provides defaults
- **WHEN** requesting OpenCode Zen factory defaults
- **THEN** the factory SHALL return endpoint as "https://opencode.ai/zen"
- **AND** enabled as false
- **AND** api_key as empty string

### Requirement: Config uses map-based provider structure
The router-service configuration SHALL use a map-based structure for providers where keys are provider type identifiers and values are ProviderSettings with enabled, endpoint, and api_key fields.

#### Scenario: Load config with multiple providers
- **WHEN** loading config with multiple providers in the providers map
- **THEN** the system SHALL parse all provider entries
- **AND** resolve environment variables in endpoint and api_key fields
- **AND** return enabled providers as a map

#### Scenario: Load config with no providers
- **WHEN** loading config with empty providers map
- **THEN** the system SHALL return an empty enabled providers map
- **AND** not fail startup

### Requirement: Handler uses registry for provider instantiation
The router-service handler SHALL use the ProviderRegistry to instantiate providers from configuration instead of hardcoded URL-based matching.

#### Scenario: Initialize providers from config
- **WHEN** the handler setup runs with enabled providers in config
- **THEN** the system SHALL iterate through enabled providers
- **AND** create each provider using the registry
- **AND** log successful initialization per provider
- **AND** log warnings for failed provider creation

#### Scenario: Handle provider creation failure
- **WHEN** a provider fails to create during setup
- **THEN** the system SHALL log the error
- **AND** continue with remaining providers
- **AND** not fail application startup

### Requirement: Provider discovery endpoint exposes registry information
The router-service SHALL provide a GET /v1/providers endpoint that returns information about registered provider types, their configuration status, and default settings.

#### Scenario: List all registered providers
- **WHEN** a client requests GET /v1/providers
- **THEN** the system SHALL return a list of all registered provider types
- **AND** include description for each provider type
- **AND** include configuration status (configured/enabled)
- **AND** include default settings for each provider type

#### Scenario: List providers with current configuration
- **WHEN** a provider is configured in the config
- **THEN** the /v1/providers response SHALL include current endpoint
- **AND** include whether api_key is present
- **AND** include enabled status

### Requirement: Provider naming follows config key convention
Provider implementations SHALL use the same name as their config key identifier, and model prefixes SHALL use the format {provider_name}: for consistency across all layers.

#### Scenario: Ollama provider naming
- **WHEN** the OllamaFactory Type() method is called
- **THEN** it SHALL return "ollama"
- **AND** the provider SHALL use "ollama:" as model prefix
- **AND** ListModels() SHALL return models with "ollama:" prefix

#### Scenario: OpenCode Zen provider naming
- **WHEN** the OpenCodeZenFactory Type() method is called
- **THEN** it SHALL return "opencode_zen"
- **AND** the provider SHALL use "opencode_zen:" as model prefix
- **AND** ListModels() SHALL return models with "opencode_zen:" prefix

#### Scenario: Router uses provider name for routing
- **WHEN** the ModelRouter receives a model request
- **THEN** it SHALL extract the provider prefix from the model string
- **AND** match it against provider names
- **AND** route to the corresponding provider

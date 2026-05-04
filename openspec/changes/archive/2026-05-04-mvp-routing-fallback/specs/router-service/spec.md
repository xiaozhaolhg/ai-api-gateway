## ADDED Requirements

### Requirement: Fallback Route Resolution
Router service SHALL return fallback provider IDs and corresponding model mappings when resolving routes.

#### Scenario: Route resolution with fallback
- **WHEN** `ResolveRoute` is called and a matching routing rule has a configured fallback provider
- **THEN** return `RouteResult` with `fallback_provider_ids` and `fallback_models` populated

#### Scenario: Route resolution without fallback
- **WHEN** `ResolveRoute` is called and no fallback provider is configured for the matching rule
- **THEN** return `RouteResult` with empty `fallback_provider_ids` and `fallback_models`

#### Scenario: Fallback model mapping
- **WHEN** a routing rule has `fallback_model` configured
- **THEN** include the `fallback_model` in the corresponding index of `RouteResult.fallback_models`

## MODIFIED Requirements

<!-- No existing requirements modified for this change -->

## REMOVED Requirements

<!-- No requirements removed for this change -->

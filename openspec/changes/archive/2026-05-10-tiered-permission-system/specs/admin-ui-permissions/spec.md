# admin-ui-permissions (delta)

## MODIFIED Requirements

### Requirement: Tier selection in group permission assignment
The admin UI group permission assignment SHALL present a tier selector instead of individual permission rule editor for model/provider access.

#### Scenario: Group permission assignment shows tier selector
- **WHEN** an admin opens the permission assignment for a group
- **THEN** a tier selector dropdown is displayed showing available tiers (predefined and custom)

#### Scenario: Tier preview
- **WHEN** an admin selects a tier in the tier selector
- **THEN** a preview panel shows the allowed models and providers for the selected tier

#### Scenario: Assign tier to group
- **WHEN** an admin selects a tier and saves
- **THEN** the group's tier_id is updated via AssignTierToGroup RPC

### Requirement: Tier management page
The admin UI SHALL provide a tier management page for CRUD operations on tiers.

#### Scenario: List tiers
- **WHEN** an admin navigates to the tier management page
- **THEN** all tiers (predefined and custom) are listed with name, description, and model/provider counts

#### Scenario: Create custom tier
- **WHEN** an admin creates a new custom tier with name, description, allowed_models, and allowed_providers
- **THEN** the tier is created via CreateTier RPC and appears in the tier list

#### Scenario: Edit custom tier
- **WHEN** an admin edits a custom tier's allowed_models or allowed_providers
- **THEN** the tier is updated via UpdateTier RPC

#### Scenario: Predefined tier edit disabled
- **WHEN** an admin views a predefined tier (is_default=true)
- **THEN** edit and delete actions are disabled

#### Scenario: Delete custom tier
- **WHEN** an admin deletes a custom tier not referenced by any group
- **THEN** the tier is deleted via DeleteTier RPC

#### Scenario: Delete tier in use
- **WHEN** an admin attempts to delete a tier referenced by groups
- **THEN** an error message is shown indicating the tier is in use

# auth-service (delta)

## MODIFIED Requirements

### Requirement: CheckModelAuthorization tier-based enforcement
The auth-service CheckModelAuthorization RPC SHALL resolve model authorization through tier-based access control instead of the current MVP stub that returns all models as authorized.

#### Scenario: Tier-based authorization replaces MVP stub
- **WHEN** CheckModelAuthorization is called with user_id, group_ids, and model
- **THEN** the system resolves the user's groups, loads each group's tier, collects allowed_models patterns from all tiers, matches the requested model against collected patterns, and checks for deny-rules in the Permission table
- **AND** returns allowed=true only if the model matches at least one tier pattern and no deny-rule exists

#### Scenario: No tier assigned to any group
- **WHEN** CheckModelAuthorization is called for a user whose groups have no tier assigned
- **THEN** allowed=false with reason="no tier assigned"

### Requirement: Group entity tier reference
The Group entity SHALL carry an optional tier_id field referencing a Tier entity.

#### Scenario: Group with tier
- **WHEN** a Group is created or updated with tier_id="tier-123"
- **THEN** the tier_id is persisted and retrievable via GetByID/List

#### Scenario: Group without tier
- **WHEN** a Group is created without tier_id
- **THEN** tier_id is null and the group has no tier-based model access

---

## Tier Management Requirements

### Requirement: Tier entity persistence
The auth-service SHALL persist a Tier entity with id, name, description, is_default, allowed_models (list of glob patterns), allowed_providers (list of glob patterns), created_at, and updated_at fields in auth-db database.

#### Scenario: Database migration on startup
- **WHEN** auth-service starts
- **THEN** Tier, TierModelPattern, and TierProviderPattern tables are created/updated via GORM AutoMigrate

#### Scenario: Tier with model patterns
- **WHEN** a Tier is created with allowed_models=["openai:gpt-4","anthropic:claude-*"]
- **THEN** patterns are persisted and retrievable via GetTier/ListTiers

#### Scenario: Tier with provider patterns
- **WHEN** a Tier is created with allowed_providers=["openai","anthropic"]
- **THEN** provider patterns are persisted and retrievable

### Requirement: Predefined tier seeding
The auth-service SHALL seed four predefined tiers on startup: Basic, Standard, Premium, and Enterprise. Seeding is idempotent — existing tiers are not overwritten.

#### Scenario: First startup
- **WHEN** auth-service starts with an empty tiers table
- **THEN** four Tier records are created: Basic, Standard, Premium, Enterprise with predefined model/provider patterns

#### Scenario: Subsequent startup
- **WHEN** auth-service starts with existing predefined tiers
- **THEN** no duplicate tiers are created; existing tiers remain unchanged

#### Scenario: Predefined tier immutability
- **WHEN** an admin attempts to delete or update a predefined tier (is_default=true)
- **THEN** operation SHALL return an error

### Requirement: Tier CRUD RPCs
The auth-service SHALL implement CreateTier, GetTier, UpdateTier, DeleteTier, and ListTiers gRPC handlers.

#### Scenario: Create custom tier
- **WHEN** CreateTier is called with name="power-devs", allowed_models=["openai:gpt-4"], allowed_providers=["openai"]
- **THEN** a Tier entity is persisted with is_default=false and returned with generated ID and timestamps

#### Scenario: GetTier
- **WHEN** GetTier is called with an existing tier ID
- **THEN** Tier entity with its model and provider patterns is returned

#### Scenario: UpdateTier
- **WHEN** UpdateTier is called with an existing custom tier ID and new allowed_models
- **THEN** Tier entity is updated with new patterns and updated_at timestamp

#### Scenario: DeleteTier with no group references
- **WHEN** DeleteTier is called with a custom tier ID that no group references
- **THEN** Tier entity and its pattern records are deleted

#### Scenario: DeleteTier with group references
- **WHEN** DeleteTier is called with a tier ID that is referenced by one or more groups
- **THEN** operation SHALL return an error indicating tier is in use

#### Scenario: ListTiers with pagination
- **WHEN** ListTiers is called with page=1, page_size=10
- **THEN** up to 10 tiers are returned with total count

### Requirement: Group-tier assignment
The auth-service SHALL support assigning a tier to a group via group's tier_id field and provide AssignTierToGroup and RemoveTierFromGroup RPCs.

#### Scenario: Assign tier to group
- **WHEN** AssignTierToGroup is called with group_id and tier_id
- **THEN** Group entity's tier_id is set to specified tier

#### Scenario: Remove tier from group
- **WHEN** RemoveTierFromGroup is called with group_id
- **THEN** Group entity's tier_id is set to null (no tier assigned)

#### Scenario: Group with no tier
- **WHEN** a group has no tier assigned (tier_id is null)
- **THEN** members of that group have no model access via tier resolution (only deny-rules from Permission apply)

### Requirement: Tier-based model authorization
The auth-service CheckModelAuthorization SHALL resolve user's groups, load each group's tier, collect allowed_models patterns from all tiers, match the requested model, and check for deny-rules.

#### Scenario: User with single group and tier
- **WHEN** CheckModelAuthorization is called for a user in group "devs" with tier "Standard" (allowed_models=["openai:gpt-4","anthropic:claude-*"])
- **AND** the requested model is "anthropic:claude-3-opus"
- **THEN** allowed=true, authorized_models=["openai:gpt-4","anthropic:claude-*"]

#### Scenario: User with multiple groups and different tiers
- **WHEN** a user belongs to group "devs" (tier Standard) and group "leads" (tier Premium)
- **THEN** user's authorized models are the union of Standard and Premium tier patterns

#### Scenario: Deny-rule overrides tier allow
- **WHEN** a user's tier allows "openai:gpt-4" but a Permission deny-rule exists for resource_type="model", resource_id="openai:gpt-4", action="access"
- **THEN** CheckModelAuthorization returns allowed=false for that model

#### Scenario: Model not in any tier
- **WHEN** the requested model does not match any allowed_models pattern in the user's tiers
- **THEN** allowed=false with reason="model not authorized by any tier"

### Requirement: Granular permission migration
The auth-service SHALL provide a migration mechanism to convert existing Permission records with resource_type="model" and effect="allow" into tier assignments.

#### Scenario: Migrate group permissions to tier
- **WHEN** migration is triggered for a group that has allow-permissions for models
- **THEN** a custom tier is created with allowed_models derived from the group's model allow-permissions, the tier is assigned to the group, and original allow-permissions are removed

#### Scenario: Group with no model allow-permissions
- **WHEN** migration is triggered for a group with no model allow-permissions
- **THEN** no tier is created; the group's tier_id remains null

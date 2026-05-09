## 1. Domain Layer — Tier Entity & Repository

- [x] 1.1 Define Tier entity struct in `auth-service/internal/domain/entity/tier.go` with fields: ID, Name, Description, IsDefault, AllowedModels ([]string), AllowedProviders ([]string), CreatedAt, UpdatedAt. Acceptance: struct compiles, fields match spec.
- [x] 1.2 Define TierRepository interface in `auth-service/internal/domain/port/repository.go` with methods: Create, GetByID, Update, Delete, List, GetByName, GetDefaultTiers. Acceptance: interface compiles, methods match spec CRUD requirements.
- [x] 1.3 Implement TierRepository using GORM in `auth-service/internal/infrastructure/repository/tier_repository.go`. Acceptance: all interface methods work against SQLite, AutoMigrate creates tiers table.
- [x] 1.4 Add TierID field to Group entity in `auth-service/internal/domain/entity/group.go`. Acceptance: Group struct has optional TierID field.
- [x] 1.5 Unit tests for TierRepository: Create, GetByID, Update, Delete, List, GetByName, GetDefaultTiers. Acceptance: all tests pass.

## 2. Application Layer — Tier Service

- [x] 2.1 Create TierService in `auth-service/internal/application/tier_service.go` with methods: CreateTier, GetTier, UpdateTier, DeleteTier, ListTiers, AssignTierToGroup, RemoveTierFromGroup, SeedDefaultTiers. Acceptance: service compiles with all methods.
- [x] 2.2 Implement SeedDefaultTiers: create Basic, Standard, Premium, Enterprise tiers with predefined model/provider patterns; idempotent (skip if exists). Acceptance: repeated calls do not create duplicates.
- [x] 2.3 Implement tier CRUD with validation: predefined tiers (IsDefault=true) cannot be updated or deleted; custom tiers cannot be deleted if referenced by groups. Acceptance: appropriate errors returned for invalid operations.
- [x] 2.4 Implement AssignTierToGroup and RemoveTierFromGroup: update Group.TierID field. Acceptance: group's tier_id is set/cleared correctly.
- [x] 2.5 Unit tests for TierService: CRUD, seeding, assignment, validation errors. Acceptance: all tests pass.

## 3. Authorization — Tier-Based CheckModelAuthorization

- [x] 3.1 Update AuthService.CheckModelAuthorization to replace MVP stub with tier-based resolution: resolve user groups → load tiers → collect allowed_models patterns → match requested model → check deny-rules. Acceptance: returns allowed=true only when model matches tier pattern and no deny-rule exists.
- [x] 3.2 Implement glob pattern matching for model names following `{provider}:{model}` convention. Acceptance: "anthropic:claude-*" matches "anthropic:claude-3-opus", "openai:gpt-4" matches exactly "openai:gpt-4".
- [x] 3.3 Handle edge cases: user with no groups, group with no tier, model not matching any tier pattern. Acceptance: returns allowed=false with appropriate reason in each case.
- [x] 3.4 Unit tests for tier-based CheckModelAuthorization: single group/tier, multiple groups/tiers, deny-rule override, no tier, no groups, pattern matching. Acceptance: all tests pass.

## 4. Proto & gRPC Handlers

- [x] 4.1 Add Tier message, TierModelPattern, TierProviderPattern, CreateTierRequest/Response, GetTierRequest/Response, UpdateTierRequest/Response, DeleteTierRequest/Response, ListTiersRequest/Response, AssignTierToGroupRequest/Response, RemoveTierFromGroupRequest/Response to `api/proto/auth/v1/auth.proto`. Acceptance: proto compiles with buf.
- [x] 4.2 Regenerate proto Go code with `buf generate`. Acceptance: generated code includes new Tier messages and service methods.
- [x] 4.3 Implement gRPC handlers in `auth-service/internal/handler/handler.go` for CreateTier, GetTier, UpdateTier, DeleteTier, ListTiers, AssignTierToGroup, RemoveTierFromGroup. Acceptance: all handlers delegate to TierService and return correct proto responses.
- [x] 4.4 Register new gRPC handlers in auth-service server startup. Acceptance: server starts and all new RPCs are accessible.
- [x] 4.5 Add tier_id field to Group message, CreateGroupRequest, UpdateGroupRequest in proto. Acceptance: group responses include tier_id, create/update requests accept tier_id.
- [x] 4.6 Implement CreateGroup handler to read req.TierId and call AssignTierToGroup if tier_id is provided. Acceptance: new groups are created with tier assignment when tier_id is passed.
- [x] 4.7 Implement UpdateGroup handler to handle tier assignment (when tier_id non-empty) and tier removal (when tier_id empty). Acceptance: updating group's tier_id correctly assigns or removes tier.

## 4.5. Gateway-Service HTTP Endpoints

- [x] 4.5.1 Add tier CRUD HTTP endpoints in gateway-service: POST /admin/auth/tiers (create), GET /admin/auth/tiers (list), GET /admin/auth/tiers/:id (get), PUT /admin/auth/tiers/:id (update), DELETE /admin/auth/tiers/:id (delete). Acceptance: endpoints delegate to auth-service gRPC and return correct JSON responses.
- [x] 4.5.2 Add tier assignment HTTP endpoint: POST /admin/auth/groups/:id/tier (assign tier to group). Acceptance: endpoint assigns tier via AssignTierToGroup gRPC call.
- [x] 4.5.3 Add tier removal HTTP endpoint: DELETE /admin/auth/groups/:id/tier (remove tier from group). Acceptance: endpoint removes tier via RemoveTierFromGroup gRPC call.
- [x] 4.5.4 Add tier_id field to group create/update HTTP handlers. Acceptance: tier_id is accepted in request body and passed to gRPC.

## 5. Migration — Granular Permissions to Tiers

- [x] 5.1 Create migration script in `auth-service/cmd/migrate_permissions/main.go` that: groups model allow-permissions by group_id, creates a custom tier per group with collected patterns, assigns tier to group, removes original allow-permissions. Acceptance: script runs idempotently.
- [x] 5.2 Add migration trigger to auth-service startup or as a separate CLI command. Acceptance: migration can be executed via command line.
- [x] 5.3 Test migration with sample data: create groups with model allow-permissions, run migration, verify tier assignments and permission cleanup. Acceptance: migration produces correct tier assignments.

## 6. Admin UI — Tier Management Page

- [x] 6.1 Create Tier management page in `admin-ui/src/pages/Tiers/` with list view showing tier name, description, is_default badge, model/provider counts. Acceptance: page renders tier list from API.
- [x] 6.2 Implement create/edit custom tier form with fields: name, description, allowed_models (bulk multi-select interface with provider grouping), allowed_providers (tag input). Acceptance: form creates/updates tier via API with support for selecting multiple models from multiple providers in bulk.
- [x] 6.3 Disable edit/delete for predefined tiers (is_default=true). Acceptance: edit and delete buttons are disabled for predefined tiers.
- [x] 6.4 Add tier detail view showing expanded model and provider pattern lists. Acceptance: detail view shows all patterns for a tier.
- [x] 6.5 Add edit/delete buttons for custom tiers. Acceptance: edit opens form with pre-filled values, delete shows confirmation and removes tier from API.

## 7. Admin UI — Group Permission Assignment Update

- [x] 7.1 Update group permission assignment UI to show tier selector dropdown instead of individual permission rule editor. Acceptance: dropdown lists available tiers.
- [x] 7.2 Add tier preview panel: when a tier is selected, show allowed models and providers with visual grouping by provider. Acceptance: preview updates on tier selection change and displays models organized by their respective providers for easy verification.
- [x] 7.3 Wire tier selector to AssignTierToGroup/RemoveTierFromGroup API calls. Acceptance: selecting a tier and saving updates the group's tier assignment.

## 8. Integration Tests

- [x] 8.1 Integration test: create tier → assign to group → add user to group → CheckModelAuthorization returns allowed=true for tier models. Acceptance: end-to-end tier authorization works.
- [x] 8.2 Integration test: create tier → assign to group → add deny-rule → CheckModelAuthorization returns allowed=false for denied model. Acceptance: deny-rule overrides tier allow.
- [x] 8.3 Integration test: user in multiple groups with different tiers → CheckModelAuthorization returns union of tier patterns. Acceptance: multi-group tier resolution works.
- [x] 8.4 Integration test: full flow — `make down && make clean-images && make up` → health check → create custom tier → assign to group → send request → verify authorization. Acceptance: live service integration test passes.

package port

import (
	"github.com/ai-api-gateway/router-service/internal/domain/entity"
)

// Router defines the interface for model-to-provider routing operations.
// It follows Clean Architecture principles - this is a domain interface
// that defines the contract for routing logic without dependencies on
// infrastructure or transport layers.
//
// Implementations of this interface should be in the application layer
// and use the RoutingRuleRepository for persistence operations.
type Router interface {
	// ResolveRoute resolves a model name to a provider identifier.
	//
	// Parameters:
	//   - model: the requested model name (e.g., "gpt-4o", "claude-3")
	//   - authorizedModels: list of models the user is authorized to access
	//
	// Returns:
	//   - *entity.RouteResult: contains provider_id and adapter_type
	//   - error: if no matching route found or authorization denied
	//
	// The implementation should:
	//   1. Match the model against routing rules using pattern matching
	//   2. Filter results by authorizedModels list
	//   3. Select the highest priority matching rule
	//   4. Return NOT_FOUND error if no match exists
	ResolveRoute(model string, authorizedModels []string) (*entity.RouteResult, error)

	// CreateRoutingRule creates a new routing rule.
	//
	// Parameters:
	//   - rule: the routing rule to create (ID will be generated)
	//
	// Returns:
	//   - error: if validation fails or persistence error occurs
	//
	// The implementation should:
	//   1. Validate the rule (model_pattern, provider_id required)
	//   2. Generate a unique ID
	//   3. Persist to database via repository
	//   4. Optionally refresh cache
	CreateRoutingRule(rule *entity.RoutingRule) error

	// UpdateRoutingRule updates an existing routing rule.
	//
	// Parameters:
	//   - rule: the routing rule with updated fields (ID must be set)
	//
	// Returns:
	//   - error: if rule not found or validation fails
	//
	// The implementation should:
	//   1. Verify the rule exists
	//   2. Validate the updated fields
	//   3. Persist changes to database
	//   4. Optionally refresh cache
	UpdateRoutingRule(rule *entity.RoutingRule) error

	// DeleteRoutingRule deletes a routing rule by ID.
	//
	// Parameters:
	//   - id: the unique identifier of the rule to delete
	//
	// Returns:
	//   - error: if rule not found or persistence error
	//
	// The implementation should:
	//   1. Verify the rule exists
	//   2. Delete from database
	//   3. Invalidate any cached routes matching this pattern
	DeleteRoutingRule(id string) error

	// ListRoutingRules retrieves routing rules with pagination.
	//
	// Parameters:
	//   - page: page number (1-indexed)
	//   - pageSize: number of items per page
	//
	// Returns:
	//   - []*entity.RoutingRule: list of rules ordered by priority desc
	//   - int: total count of all rules
	//   - error: if retrieval fails
	ListRoutingRules(page, pageSize int) ([]*entity.RoutingRule, int, error)

	// RefreshRoutingTable invalidates the routing table cache.
	//
	// This should be called after any provider configuration changes
	// to ensure subsequent ResolveRoute calls use updated data.
	//
	// Returns:
	//   - error: if cache invalidation fails
	RefreshRoutingTable() error
}

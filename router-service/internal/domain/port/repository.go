package port

import (
	"github.com/ai-api-gateway/router-service/internal/domain/entity"
)

// RoutingRuleRepository defines the interface for routing rule persistence operations
type RoutingRuleRepository interface {
	Create(rule *entity.RoutingRule) error
	GetByID(id string) (*entity.RoutingRule, error)
	Update(rule *entity.RoutingRule) error
	UpdateWithOwnership(rule *entity.RoutingRule, requestingUserID string) error
	Delete(id string) error
	DeleteWithOwnership(id string, requestingUserID string) error
	List(page, pageSize int) ([]*entity.RoutingRule, int, error)
	FindByModel(modelPattern string, userID *string) (*entity.RoutingRule, error)
	FindByUserID(userID string) ([]*entity.RoutingRule, error)
}

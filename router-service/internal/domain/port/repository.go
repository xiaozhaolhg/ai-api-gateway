package port

import (
	"github.com/ai-api-gateway/router-service/internal/domain/entity"
)

// RoutingRuleRepository defines the interface for routing rule persistence operations
type RoutingRuleRepository interface {
	Create(rule *entity.RoutingRule) error
	GetByID(id string) (*entity.RoutingRule, error)
	Update(rule *entity.RoutingRule) error
	Delete(id string) error
	List(page, pageSize int) ([]*entity.RoutingRule, int, error)
}

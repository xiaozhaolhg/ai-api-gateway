package port

import (
	"github.com/ai-api-gateway/auth-service/internal/domain/entity"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id string) error
	List(page, pageSize int) ([]*entity.User, int, error)
}

// APIKeyRepository defines the interface for API key persistence operations
type APIKeyRepository interface {
	Create(apiKey *entity.APIKey) error
	GetByID(id string) (*entity.APIKey, error)
	GetByKeyHash(keyHash string) (*entity.APIKey, error)
	GetByUserID(userID string, page, pageSize int) ([]*entity.APIKey, int, error)
	Delete(id string) error
}

// GroupRepository defines the interface for group persistence operations
type GroupRepository interface {
	Create(group *entity.Group) error
	GetByID(id string) (*entity.Group, error)
	Update(group *entity.Group) error
	Delete(id string) error
	List(page, pageSize int) ([]*entity.Group, int, error)
}

// UserGroupRepository defines the interface for user-group membership operations
type UserGroupRepository interface {
	GetGroupIDsByUserID(userID string) ([]string, error)
	Create(membership *entity.UserGroupMembership) error
	Delete(userID, groupID string) error
	GetByUserID(userID string) ([]*entity.UserGroupMembership, error)
	GetByGroupID(groupID string, page, pageSize int) ([]*entity.UserGroupMembership, int, error)
	Exists(userID, groupID string) (bool, error)
}

// PermissionRepository defines the interface for permission persistence operations
type PermissionRepository interface {
	Create(permission *entity.Permission) error
	GetByID(id string) (*entity.Permission, error)
	Delete(id string) error
	ListByGroupID(groupID string, page, pageSize int) ([]*entity.Permission, int, error)
	FindByUserGroups(groupIDs []string, resourceType, resourceID, action string) ([]*entity.Permission, error)
}

// TierRepository defines the interface for tier persistence operations
type TierRepository interface {
	Create(tier *entity.Tier) error
	GetByID(id string) (*entity.Tier, error)
	Update(tier *entity.Tier) error
	Delete(id string) error
	List(page, pageSize int) ([]*entity.Tier, int, error)
	GetByName(name string) (*entity.Tier, error)
	GetDefaultTiers() ([]*entity.Tier, error)
}

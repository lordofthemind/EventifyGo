package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/types"
)

type SuperUserRepositoryInterface interface {
	// General CRUD methods
	Create(ctx context.Context, superUser *types.SuperUserType) error
	FindByID(ctx context.Context, id uuid.UUID) (*types.SuperUserType, error)
	FindByEmail(ctx context.Context, email string) (*types.SuperUserType, error)
	FindByUsername(ctx context.Context, username string) (*types.SuperUserType, error)
	FindByResetToken(ctx context.Context, token string) (*types.SuperUserType, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error

	// Search methods
	SearchSuperusers(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.SuperUserType, error)

	// Field updates
	Update(ctx context.Context, superUser *types.SuperUserType) error
	UpdateField(ctx context.Context, id uuid.UUID, field string, value interface{}) error

	// Specialized queries
	GetRoleByID(ctx context.Context, id uuid.UUID) (string, error)
	FindAll2FAEnabledSuperusers(ctx context.Context) ([]*types.SuperUserType, error)

	// Specific field updates
	UpdateResetToken(ctx context.Context, id uuid.UUID, token string) error
	UpdateSuperuserRole(ctx context.Context, id uuid.UUID, role string) error
}

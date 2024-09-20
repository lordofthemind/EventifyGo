package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/types"
)

type SuperUserServiceInterface interface {
	// Create a new SuperUser with validation, hashing password, and initializing defaults
	CreateSuperUser(ctx context.Context, superUser *types.SuperUserType) (*types.SuperUserType, error)

	// Find SuperUser by different identifiers
	GetSuperUserByID(ctx context.Context, id uuid.UUID) (*types.SuperUserType, error)
	GetSuperUserByEmail(ctx context.Context, email string) (*types.SuperUserType, error)
	GetSuperUserByUsername(ctx context.Context, username string) (*types.SuperUserType, error)
	GetSuperUserByResetToken(ctx context.Context, token string) (*types.SuperUserType, error)

	// Manage 2FA
	Enable2FAForSuperUser(ctx context.Context, id uuid.UUID, secret string) error
	Disable2FAForSuperUser(ctx context.Context, id uuid.UUID) error
	GetAll2FAEnabledSuperUsers(ctx context.Context) ([]*types.SuperUserType, error)

	// Manage SuperUser roles and permissions
	UpdateSuperUserRole(ctx context.Context, id uuid.UUID, role string) error
	GetRoleBySuperUserID(ctx context.Context, id uuid.UUID) (string, error)
	UpdateSuperUserPermissions(ctx context.Context, id uuid.UUID, permissions []string) error

	// Update specific fields for a SuperUser
	UpdateSuperUserDetails(ctx context.Context, superUser *types.SuperUserType) error
	UpdateSuperUserField(ctx context.Context, id uuid.UUID, field string, value interface{}) error

	// Reset token management
	GenerateAndSetResetToken(ctx context.Context, id uuid.UUID) (string, error)
	ClearResetToken(ctx context.Context, id uuid.UUID) error

	// Delete operations
	DeleteSuperUserByID(ctx context.Context, id uuid.UUID) error

	// Search SuperUsers with pagination and sorting
	SearchSuperUsers(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.SuperUserType, error)

	GetAllSuperUsers(ctx context.Context) ([]*types.SuperUserType, error)
}

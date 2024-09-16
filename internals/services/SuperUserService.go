package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/repositories"
	"github.com/lordofthemind/EventifyGo/internals/types"
	"golang.org/x/crypto/bcrypt"
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
}

type SuperUserService struct {
	repo repositories.SuperUserRepositoryInterface
}

func NewSuperUserService(repo repositories.SuperUserRepositoryInterface) SuperUserServiceInterface {
	return &SuperUserService{repo: repo}
}

// Create a new SuperUser
func (s *SuperUserService) CreateSuperUser(ctx context.Context, superUser *types.SuperUserType) (*types.SuperUserType, error) {
	// Validate fields
	if err := validateSuperUser(superUser); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(superUser.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}
	superUser.HashedPassword = string(hashedPassword)

	// Set default values for creation
	superUser.ID = uuid.New()
	superUser.CreatedAt = time.Now()
	superUser.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, superUser); err != nil {
		return nil, fmt.Errorf("failed to create superuser: %w", err)
	}

	return superUser, nil
}

// Get SuperUser by ID
func (s *SuperUserService) GetSuperUserByID(ctx context.Context, id uuid.UUID) (*types.SuperUserType, error) {
	superUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("superuser not found: %w", err)
	}
	return superUser, nil
}

// Get SuperUser by email
func (s *SuperUserService) GetSuperUserByEmail(ctx context.Context, email string) (*types.SuperUserType, error) {
	superUser, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("superuser not found by email: %w", err)
	}
	return superUser, nil
}

// Get SuperUser by username
func (s *SuperUserService) GetSuperUserByUsername(ctx context.Context, username string) (*types.SuperUserType, error) {
	superUser, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("superuser not found by username: %w", err)
	}
	return superUser, nil
}

// Get SuperUser by reset token
func (s *SuperUserService) GetSuperUserByResetToken(ctx context.Context, token string) (*types.SuperUserType, error) {
	superUser, err := s.repo.FindByResetToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("superuser not found by reset token: %w", err)
	}
	return superUser, nil
}

// Enable 2FA for SuperUser
func (s *SuperUserService) Enable2FAForSuperUser(ctx context.Context, id uuid.UUID, secret string) error {
	superUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("superuser not found: %w", err)
	}

	// Update 2FA fields
	superUser.TwoFactorSecret = &secret
	superUser.Is2FAEnabled = true

	if err := s.repo.Update(ctx, superUser); err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	return nil
}

// Disable 2FA for SuperUser
func (s *SuperUserService) Disable2FAForSuperUser(ctx context.Context, id uuid.UUID) error {
	superUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("superuser not found: %w", err)
	}

	// Clear 2FA fields
	superUser.TwoFactorSecret = nil
	superUser.Is2FAEnabled = false

	if err := s.repo.Update(ctx, superUser); err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}

	return nil
}

// Get all 2FA enabled SuperUsers
func (s *SuperUserService) GetAll2FAEnabledSuperUsers(ctx context.Context) ([]*types.SuperUserType, error) {
	superUsers, err := s.repo.FindAll2FAEnabledSuperusers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve 2FA enabled superusers: %w", err)
	}
	return superUsers, nil
}

// Update SuperUser role
func (s *SuperUserService) UpdateSuperUserRole(ctx context.Context, id uuid.UUID, role string) error {
	if err := s.repo.UpdateSuperuserRole(ctx, id, role); err != nil {
		return fmt.Errorf("failed to update superuser role: %w", err)
	}
	return nil
}

// Get role by SuperUser ID
func (s *SuperUserService) GetRoleBySuperUserID(ctx context.Context, id uuid.UUID) (string, error) {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get role by ID: %w", err)
	}
	return role, nil
}

// Update SuperUser permissions
func (s *SuperUserService) UpdateSuperUserPermissions(ctx context.Context, id uuid.UUID, permissions []string) error {
	if err := s.repo.UpdateField(ctx, id, "permission_groups", permissions); err != nil {
		return fmt.Errorf("failed to update superuser permissions: %w", err)
	}
	return nil
}

// Update SuperUser details
func (s *SuperUserService) UpdateSuperUserDetails(ctx context.Context, superUser *types.SuperUserType) error {
	superUser.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, superUser); err != nil {
		return fmt.Errorf("failed to update superuser details: %w", err)
	}
	return nil
}

// Update specific SuperUser field
func (s *SuperUserService) UpdateSuperUserField(ctx context.Context, id uuid.UUID, field string, value interface{}) error {
	if err := s.repo.UpdateField(ctx, id, field, value); err != nil {
		return fmt.Errorf("failed to update superuser field: %w", err)
	}
	return nil
}

// Generate and set reset token
func (s *SuperUserService) GenerateAndSetResetToken(ctx context.Context, id uuid.UUID) (string, error) {
	token := generateResetToken() // You can implement your token generator
	if err := s.repo.UpdateResetToken(ctx, id, token); err != nil {
		return "", fmt.Errorf("failed to set reset token: %w", err)
	}
	return token, nil
}

// Clear reset token
func (s *SuperUserService) ClearResetToken(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.UpdateResetToken(ctx, id, ""); err != nil {
		return fmt.Errorf("failed to clear reset token: %w", err)
	}
	return nil
}

// Delete SuperUser by ID
func (s *SuperUserService) DeleteSuperUserByID(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteByID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete superuser: %w", err)
	}
	return nil
}

// Search SuperUsers with pagination and sorting
func (s *SuperUserService) SearchSuperUsers(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.SuperUserType, error) {
	superUsers, err := s.repo.SearchSuperusers(ctx, searchQuery, page, limit, sortBy)
	if err != nil {
		return nil, fmt.Errorf("failed to search superusers: %w", err)
	}
	return superUsers, nil
}

// Helper function to validate super user input
func validateSuperUser(superUser *types.SuperUserType) error {
	// Add your validation logic here (e.g., check for required fields)
	if superUser.Email == "" || superUser.Username == "" {
		return errors.New("email and username are required")
	}
	return nil
}

// Generate reset token (dummy implementation)
func generateResetToken() string {
	return uuid.New().String()
}

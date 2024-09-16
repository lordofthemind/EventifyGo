package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/types"
	"gorm.io/gorm"
)

type postgresSuperUserRepository struct {
	db *gorm.DB
}

func NewPostgresSuperUserRepository(db *gorm.DB) SuperUserRepositoryInterface {
	return &postgresSuperUserRepository{
		db: db,
	}
}

// Create inserts a new super user into the PostgreSQL database
func (r *postgresSuperUserRepository) Create(ctx context.Context, superUser *types.SuperUserType) error {
	superUser.ID = uuid.New()
	superUser.CreatedAt = time.Now()
	superUser.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Create(superUser).Error
}

// FindByID finds a super user by UUID
func (r *postgresSuperUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	if err := r.db.WithContext(ctx).First(&superUser, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("superuser not found")
		}
		return nil, err
	}
	return &superUser, nil
}

// FindByEmail finds a super user by email
func (r *postgresSuperUserRepository) FindByEmail(ctx context.Context, email string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	if err := r.db.WithContext(ctx).First(&superUser, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("superuser not found")
		}
		return nil, err
	}
	return &superUser, nil
}

// FindByUsername finds a super user by username
func (r *postgresSuperUserRepository) FindByUsername(ctx context.Context, username string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	if err := r.db.WithContext(ctx).First(&superUser, "username = ?", username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("superuser not found")
		}
		return nil, err
	}
	return &superUser, nil
}

// FindByResetToken finds a super user by reset token
func (r *postgresSuperUserRepository) FindByResetToken(ctx context.Context, token string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	if err := r.db.WithContext(ctx).First(&superUser, "reset_token = ?", token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("superuser not found")
		}
		return nil, err
	}
	return &superUser, nil
}

// DeleteByID deletes a super user by UUID
func (r *postgresSuperUserRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&types.SuperUserType{}, "id = ?", id).Error
}

// SearchSuperusers searches for super users based on a query string, with pagination and sorting
func (r *postgresSuperUserRepository) SearchSuperusers(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.SuperUserType, error) {
	var superUsers []*types.SuperUserType

	query := r.db.WithContext(ctx).
		Where("full_name ILIKE ? OR username ILIKE ? OR email ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%").
		Order(sortBy).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&superUsers)

	if query.Error != nil {
		return nil, query.Error
	}

	return superUsers, nil
}

// Update updates an entire super user document
func (r *postgresSuperUserRepository) Update(ctx context.Context, superUser *types.SuperUserType) error {
	superUser.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(superUser).Error
}

// UpdateField updates a single field of a super user document
func (r *postgresSuperUserRepository) UpdateField(ctx context.Context, id uuid.UUID, field string, value interface{}) error {
	return r.db.WithContext(ctx).Model(&types.SuperUserType{}).Where("id = ?", id).Update(field, value).Error
}

// GetRoleByID retrieves the role of a super user by their UUID
func (r *postgresSuperUserRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (string, error) {
	var role string
	err := r.db.WithContext(ctx).Model(&types.SuperUserType{}).Where("id = ?", id).Pluck("role", &role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("superuser not found")
		}
		return "", err
	}
	return role, nil
}

// UpdateResetToken updates the reset token of a super user
func (r *postgresSuperUserRepository) UpdateResetToken(ctx context.Context, id uuid.UUID, token string) error {
	return r.UpdateField(ctx, id, "reset_token", token)
}

// UpdateSuperuserRole updates the role of a super user
func (r *postgresSuperUserRepository) UpdateSuperuserRole(ctx context.Context, id uuid.UUID, role string) error {
	return r.UpdateField(ctx, id, "role", role)
}

// FindAll2FAEnabledSuperusers retrieves all super users with 2FA enabled
func (r *postgresSuperUserRepository) FindAll2FAEnabledSuperusers(ctx context.Context) ([]*types.SuperUserType, error) {
	var superUsers []*types.SuperUserType
	err := r.db.WithContext(ctx).Where("is_2fa_enabled = ?", true).Find(&superUsers).Error
	if err != nil {
		return nil, err
	}
	return superUsers, nil
}

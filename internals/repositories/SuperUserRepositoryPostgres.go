package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/types"
	"gorm.io/gorm"
)

type postgresSuperUserRepository struct {
	db *gorm.DB
}

func NewPostgresSuperUserRepository(db *gorm.DB) SuperUserRepository {
	return &postgresSuperUserRepository{
		db: db,
	}
}

// Create inserts a new super user into the PostgreSQL database
func (r *postgresSuperUserRepository) Create(ctx context.Context, superUser *types.SuperUserType) error {
	return r.db.WithContext(ctx).Create(superUser).Error
}

// FindByID finds a super user by UUID
func (r *postgresSuperUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	err := r.db.WithContext(ctx).First(&superUser, "id = ?", id).Error
	if err != nil {
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
	err := r.db.WithContext(ctx).First(&superUser, "email = ?", email).Error
	if err != nil {
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
	err := r.db.WithContext(ctx).First(&superUser, "username = ?", username).Error
	if err != nil {
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
	err := r.db.WithContext(ctx).First(&superUser, "reset_token = ?", token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("superuser not found")
		}
		return nil, err
	}
	return &superUser, nil
}

// DeleteByID deletes a super user by UUID
func (r *postgresSuperUserRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&types.SuperUserType{}, id).Error
}

// SearchSuperusers searches for super users based on a query string, with pagination and sorting
func (r *postgresSuperUserRepository) SearchSuperusers(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.SuperUserType, error) {
	var superUsers []*types.SuperUserType
	offset := (page - 1) * limit
	query := "%" + searchQuery + "%"
	err := r.db.WithContext(ctx).
		Where("full_name ILIKE ? OR username ILIKE ? OR email ILIKE ?", query, query, query).
		Order(sortBy + " ASC").
		Offset(offset).
		Limit(limit).
		Find(&superUsers).Error
	if err != nil {
		return nil, err
	}
	return superUsers, nil
}

// Update updates an entire super user document
func (r *postgresSuperUserRepository) Update(ctx context.Context, superUser *types.SuperUserType) error {
	return r.db.WithContext(ctx).Save(superUser).Error
}

// UpdateField allows updating a single field of a super user document
func (r *postgresSuperUserRepository) UpdateField(ctx context.Context, id uuid.UUID, field string, value interface{}) error {
	return r.db.WithContext(ctx).Model(&types.SuperUserType{}).Where("id = ?", id).Update(field, value).Error
}

// GetRoleByID retrieves the role of a super user by their UUID
func (r *postgresSuperUserRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (string, error) {
	var result struct {
		Role string
	}
	err := r.db.WithContext(ctx).Model(&types.SuperUserType{}).Select("role").Where("id = ?", id).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("superuser not found")
		}
		return "", err
	}
	return result.Role, nil
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

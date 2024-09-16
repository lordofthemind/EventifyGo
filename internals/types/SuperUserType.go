package types

import (
	"time"

	"github.com/google/uuid"
)

type SuperUserType struct {
	ID               uuid.UUID `bson:"_id,omitempty" json:"id,omitempty"`
	Role             string    `bson:"role" json:"role" validate:"required"`
	Email            string    `bson:"email" json:"email" validate:"required,email"`
	FullName         string    `bson:"full_name" json:"full_name" validate:"required,min=3,max=32"`
	Username         string    `bson:"username" json:"username" validate:"required,min=3,max=32,alphanum"`
	HashedPassword   string    `bson:"hashed_password" json:"-" validate:"required,min=8"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
	ResetToken       *string   `bson:"reset_token,omitempty" json:"reset_token,omitempty"`
	Is2FAEnabled     bool      `bson:"is_2fa_enabled" json:"is_2fa_enabled"`
	TwoFactorSecret  *string   `bson:"two_factor_secret,omitempty" json:"-"`
	PermissionGroups []string  `bson:"permission_groups" json:"permission_groups" validate:"dive,required"`
}

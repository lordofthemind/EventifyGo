package inmemorydb

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/repositories"
	"github.com/lordofthemind/EventifyGo/internals/types"
)

type inMemorySuperUserRepository struct {
	mu         sync.RWMutex
	superUsers map[uuid.UUID]*types.SuperUserType
}

// NewInMemorySuperUserRepository initializes an in-memory repository
func NewInMemorySuperUserRepository() repositories.SuperUserRepositoryInterface {
	return &inMemorySuperUserRepository{
		superUsers: make(map[uuid.UUID]*types.SuperUserType),
	}
}

// Create creates a new super user in the in-memory store
func (r *inMemorySuperUserRepository) Create(ctx context.Context, superUser *types.SuperUserType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	superUser.ID = uuid.New()
	superUser.CreatedAt = time.Now()
	superUser.UpdatedAt = time.Now()

	r.superUsers[superUser.ID] = superUser
	return nil
}

// FindByID finds a super user by their ID
func (r *inMemorySuperUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if superUser, exists := r.superUsers[id]; exists {
		return superUser, nil
	}
	return nil, errors.New("superuser not found")
}

// FindByEmail finds a super user by their email address
func (r *inMemorySuperUserRepository) FindByEmail(ctx context.Context, email string) (*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, superUser := range r.superUsers {
		if superUser.Email == email {
			return superUser, nil
		}
	}
	return nil, errors.New("superuser not found")
}

// FindByUsername finds a super user by their username
func (r *inMemorySuperUserRepository) FindByUsername(ctx context.Context, username string) (*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, superUser := range r.superUsers {
		if superUser.Username == username {
			return superUser, nil
		}
	}
	return nil, errors.New("superuser not found")
}

// FindByResetToken finds a super user by their reset token
func (r *inMemorySuperUserRepository) FindByResetToken(ctx context.Context, token string) (*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, superUser := range r.superUsers {
		if superUser.ResetToken != nil && *superUser.ResetToken == token {
			return superUser, nil
		}
	}
	return nil, errors.New("superuser not found")
}

// DeleteByID deletes a super user by their ID
func (r *inMemorySuperUserRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.superUsers[id]; exists {
		delete(r.superUsers, id)
		return nil
	}
	return errors.New("superuser not found")
}

// SearchSuperusers searches for super users based on a search query
func (r *inMemorySuperUserRepository) SearchSuperusers(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*types.SuperUserType
	for _, superUser := range r.superUsers {
		if matchesQuery(superUser, searchQuery) {
			results = append(results, superUser)
		}
	}

	// Apply pagination
	start := (page - 1) * limit
	if start >= len(results) {
		return nil, nil
	}

	end := start + limit
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], nil
}

// Update updates a super user
func (r *inMemorySuperUserRepository) Update(ctx context.Context, superUser *types.SuperUserType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	superUser.UpdatedAt = time.Now()

	if _, exists := r.superUsers[superUser.ID]; exists {
		r.superUsers[superUser.ID] = superUser
		return nil
	}
	return errors.New("superuser not found")
}

// UpdateField updates a specific field for a super user
func (r *inMemorySuperUserRepository) UpdateField(ctx context.Context, id uuid.UUID, field string, value interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if superUser, exists := r.superUsers[id]; exists {
		superUser.UpdatedAt = time.Now()

		switch field {
		case "role":
			superUser.Role = value.(string)
		case "reset_token":
			token := value.(string)
			superUser.ResetToken = &token
			// Add more cases as needed
		}
		r.superUsers[id] = superUser
		return nil
	}
	return errors.New("superuser not found")
}

// GetRoleByID returns the role of a super user by their ID
func (r *inMemorySuperUserRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if superUser, exists := r.superUsers[id]; exists {
		return superUser.Role, nil
	}
	return "", errors.New("role not found")
}

// FindAll2FAEnabledSuperusers finds all super users with 2FA enabled
func (r *inMemorySuperUserRepository) FindAll2FAEnabledSuperusers(ctx context.Context) ([]*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var superUsers []*types.SuperUserType
	for _, superUser := range r.superUsers {
		if superUser.Is2FAEnabled {
			superUsers = append(superUsers, superUser)
		}
	}
	return superUsers, nil
}

// UpdateResetToken updates the reset token for a super user
func (r *inMemorySuperUserRepository) UpdateResetToken(ctx context.Context, id uuid.UUID, token string) error {
	return r.UpdateField(ctx, id, "reset_token", token)
}

// UpdateSuperuserRole updates the role of a super user
func (r *inMemorySuperUserRepository) UpdateSuperuserRole(ctx context.Context, id uuid.UUID, role string) error {
	return r.UpdateField(ctx, id, "role", role)
}

// Helper function to match the search query
func matchesQuery(superUser *types.SuperUserType, query string) bool {
	return superUserContainsIgnoreCase(superUser.FullName, query) ||
		superUserContainsIgnoreCase(superUser.Email, query) ||
		superUserContainsIgnoreCase(superUser.Username, query)
}

// Helper function to do a case-insensitive contains check
func superUserContainsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr)
}

package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoSuperUserRepository struct {
	collection *mongo.Collection
}

// NewMongoSuperUserRepository initializes a new MongoDB repository
func NewMongoSuperUserRepository(db *mongo.Database) SuperUserRepository {
	return &mongoSuperUserRepository{
		collection: db.Collection("superusers"),
	}
}

// Create creates a new super user in the database
func (r *mongoSuperUserRepository) Create(ctx context.Context, superUser *types.SuperUserType) error {
	superUser.ID = uuid.New()
	superUser.CreatedAt = time.Now()
	superUser.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, superUser)
	return err
}

// FindByID finds a super user by their ID
func (r *mongoSuperUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&superUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

// FindByEmail finds a super user by their email address
func (r *mongoSuperUserRepository) FindByEmail(ctx context.Context, email string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	filter := bson.M{"email": email}

	err := r.collection.FindOne(ctx, filter).Decode(&superUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

// FindByUsername finds a super user by their username
func (r *mongoSuperUserRepository) FindByUsername(ctx context.Context, username string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	filter := bson.M{"username": username}

	err := r.collection.FindOne(ctx, filter).Decode(&superUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

// FindByResetToken finds a super user by their reset token
func (r *mongoSuperUserRepository) FindByResetToken(ctx context.Context, token string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	filter := bson.M{"reset_token": token}

	err := r.collection.FindOne(ctx, filter).Decode(&superUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

// DeleteByID deletes a super user by their ID
func (r *mongoSuperUserRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// SearchSuperusers searches for super users based on a search query
func (r *mongoSuperUserRepository) SearchSuperusers(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.SuperUserType, error) {
	var superUsers []*types.SuperUserType
	skip := (page - 1) * limit

	filter := bson.M{
		"$or": []bson.M{
			{"full_name": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"email": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"username": bson.M{"$regex": searchQuery, "$options": "i"}},
		},
	}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: sortBy, Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &superUsers); err != nil {
		return nil, err
	}

	return superUsers, nil
}

// Update updates a super user
func (r *mongoSuperUserRepository) Update(ctx context.Context, superUser *types.SuperUserType) error {
	superUser.UpdatedAt = time.Now()

	filter := bson.M{"_id": superUser.ID}
	update := bson.M{"$set": superUser}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateField updates a specific field for a super user
func (r *mongoSuperUserRepository) UpdateField(ctx context.Context, id uuid.UUID, field string, value interface{}) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{field: value}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetRoleByID returns the role of a super user by their ID
func (r *mongoSuperUserRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (string, error) {
	var result struct {
		Role string `bson:"role"`
	}
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return "", errors.New("role not found")
	}
	return result.Role, err
}

// FindAll2FAEnabledSuperusers finds all super users with 2FA enabled
func (r *mongoSuperUserRepository) FindAll2FAEnabledSuperusers(ctx context.Context) ([]*types.SuperUserType, error) {
	var superUsers []*types.SuperUserType
	filter := bson.M{"is_2fa_enabled": true}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &superUsers); err != nil {
		return nil, err
	}

	return superUsers, nil
}

// UpdateResetToken updates the reset token for a super user
func (r *mongoSuperUserRepository) UpdateResetToken(ctx context.Context, id uuid.UUID, token string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"reset_token": token, "updated_at": time.Now()}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateSuperuserRole updates the role of a super user
func (r *mongoSuperUserRepository) UpdateSuperuserRole(ctx context.Context, id uuid.UUID, role string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"role": role, "updated_at": time.Now()}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

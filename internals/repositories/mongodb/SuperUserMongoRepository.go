package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/repositories"
	"github.com/lordofthemind/EventifyGo/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoSuperUserRepository struct {
	collection *mongo.Collection
}

func NewMongoSuperUserRepository(db *mongo.Database) repositories.SuperUserRepositoryInterface {
	return &mongoSuperUserRepository{
		collection: db.Collection("superusers"),
	}
}

// Create inserts a new super user into the MongoDB collection
func (r *mongoSuperUserRepository) Create(ctx context.Context, superUser *types.SuperUserType) error {
	superUser.ID = uuid.New()
	superUser.CreatedAt = time.Now()
	superUser.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, superUser)
	return err
}

// FindByID finds a super user by UUID
func (r *mongoSuperUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&superUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

// FindByEmail finds a super user by email
func (r *mongoSuperUserRepository) FindByEmail(ctx context.Context, email string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&superUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

// FindByUsername finds a super user by username
func (r *mongoSuperUserRepository) FindByUsername(ctx context.Context, username string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	filter := bson.M{"username": username}
	err := r.collection.FindOne(ctx, filter).Decode(&superUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

// FindByResetToken finds a super user by reset token
func (r *mongoSuperUserRepository) FindByResetToken(ctx context.Context, token string) (*types.SuperUserType, error) {
	var superUser types.SuperUserType
	filter := bson.M{"reset_token": token}
	err := r.collection.FindOne(ctx, filter).Decode(&superUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("superuser not found")
	}
	return &superUser, err
}

// DeleteByID deletes a super user by UUID
func (r *mongoSuperUserRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// SearchSuperusers searches for super users based on a query string, with pagination and sorting
func (r *mongoSuperUserRepository) SearchSuperusers(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.SuperUserType, error) {
	var superUsers []*types.SuperUserType
	filter := bson.M{
		"$or": []bson.M{
			{"full_name": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"username": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"email": bson.M{"$regex": searchQuery, "$options": "i"}},
		},
	}

	findOptions := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.M{sortBy: 1}) // 1 for ascending, -1 for descending

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var superUser types.SuperUserType
		if err := cursor.Decode(&superUser); err != nil {
			return nil, err
		}
		superUsers = append(superUsers, &superUser)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return superUsers, nil
}

// Update updates an entire super user document
func (r *mongoSuperUserRepository) Update(ctx context.Context, superUser *types.SuperUserType) error {
	superUser.UpdatedAt = time.Now()

	filter := bson.M{"_id": superUser.ID}
	update := bson.M{"$set": superUser}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateField allows updating a single field of a super user document
func (r *mongoSuperUserRepository) UpdateField(ctx context.Context, id uuid.UUID, field string, value interface{}) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{field: value}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetRoleByID retrieves the role of a super user by their UUID
func (r *mongoSuperUserRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (string, error) {
	var result struct {
		Role string `bson:"role"`
	}
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return "", errors.New("superuser not found")
	}
	return result.Role, err
}

// FindAll2FAEnabledSuperusers retrieves all super users with 2FA enabled
func (r *mongoSuperUserRepository) FindAll2FAEnabledSuperusers(ctx context.Context) ([]*types.SuperUserType, error) {
	var superUsers []*types.SuperUserType
	filter := bson.M{"is_2fa_enabled": true}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var superUser types.SuperUserType
		if err := cursor.Decode(&superUser); err != nil {
			return nil, err
		}
		superUsers = append(superUsers, &superUser)
	}

	if err := cursor.Err(); err != nil {
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

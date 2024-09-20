package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/repositories"
	"github.com/lordofthemind/EventifyGo/internals/types"
)

type mongoEventRepository struct {
	collection *mongo.Collection
}

// NewMongoEventRepository creates a new instance of mongoEventRepository.
func NewMongoEventRepository(db *mongo.Database) repositories.EventRepositoryInterface {
	return &mongoEventRepository{
		collection: db.Collection("events"),
	}
}

func (r *mongoEventRepository) CreateEvent(ctx context.Context, event *types.EventType) error {
	_, err := r.collection.InsertOne(ctx, event)
	return err
}

func (r *mongoEventRepository) GetEventByID(ctx context.Context, eventID uuid.UUID) (*types.EventType, error) {
	var event types.EventType
	filter := bson.M{"_id": eventID}
	err := r.collection.FindOne(ctx, filter).Decode(&event)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &event, err
}

func (r *mongoEventRepository) UpdateEvent(ctx context.Context, event *types.EventType) error {
	filter := bson.M{"_id": event.EventID}
	update := bson.M{
		"$set": bson.M{
			"name":         event.Name,
			"description":  event.Description,
			"date":         event.Date,
			"location":     event.Location,
			"capacity":     event.Capacity,
			"updated_at":   time.Now(),
			"organizer_id": event.OrganizerID,
			"attendees":    event.Attendees,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoEventRepository) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	filter := bson.M{"_id": eventID}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *mongoEventRepository) SearchEvents(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.EventType, error) {
	var events []*types.EventType
	skip := (page - 1) * limit

	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"description": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"location": bson.M{"$regex": searchQuery, "$options": "i"}},
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

	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *mongoEventRepository) ListEvents(ctx context.Context, page, limit int, sortBy string) ([]*types.EventType, error) {
	var events []*types.EventType
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: sortBy, Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *mongoEventRepository) CountEvents(ctx context.Context, searchQuery string) (int64, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"description": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"location": bson.M{"$regex": searchQuery, "$options": "i"}},
		},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	return count, err
}

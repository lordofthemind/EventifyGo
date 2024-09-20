package postgresdb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/repositories"
	"github.com/lordofthemind/EventifyGo/internals/types"
	"gorm.io/gorm"
)

type postgresEventRepository struct {
	db *gorm.DB
}

// NewPostgresEventRepository creates a new instance of postgresEventRepository.
func NewPostgresEventRepository(db *gorm.DB) repositories.EventRepositoryInterface {
	return &postgresEventRepository{db: db}
}

func (r *postgresEventRepository) CreateEvent(ctx context.Context, event *types.EventType) error {
	event.EventID = uuid.New() // Assign a new UUID
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *postgresEventRepository) GetEventByID(ctx context.Context, eventID uuid.UUID) (*types.EventType, error) {
	var event types.EventType
	if err := r.db.WithContext(ctx).First(&event, "event_id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *postgresEventRepository) UpdateEvent(ctx context.Context, event *types.EventType) error {
	event.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Model(event).Where("event_id = ?", event.EventID).Updates(event).Error
}

func (r *postgresEventRepository) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("event_id = ?", eventID).Delete(&types.EventType{}).Error
}

func (r *postgresEventRepository) SearchEvents(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.EventType, error) {
	var events []*types.EventType
	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).Where("name ILIKE ? OR description ILIKE ? OR location ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%").
		Order(sortBy).Offset(offset).Limit(limit).Find(&events).Error

	return events, err
}

func (r *postgresEventRepository) ListEvents(ctx context.Context, page, limit int, sortBy string) ([]*types.EventType, error) {
	var events []*types.EventType
	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).Order(sortBy).Offset(offset).Limit(limit).Find(&events).Error

	return events, err
}

func (r *postgresEventRepository) CountEvents(ctx context.Context, searchQuery string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&types.EventType{}).Where("name ILIKE ? OR description ILIKE ? OR location ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%").Count(&count).Error
	return count, err
}

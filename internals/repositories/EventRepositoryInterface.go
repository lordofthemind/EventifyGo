package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/internals/types"
)

type EventRepositoryInterface interface {
	// CreateEvent creates a new event in the repository.
	CreateEvent(ctx context.Context, event *types.EventType) error

	// GetEventByID retrieves an event by its ID.
	GetEventByID(ctx context.Context, eventID uuid.UUID) (*types.EventType, error)

	// UpdateEvent updates an existing event.
	UpdateEvent(ctx context.Context, event *types.EventType) error

	// DeleteEvent deletes an event by its ID.
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error

	// SearchEvents searches for events based on the search query, pagination, and sorting.
	SearchEvents(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.EventType, error)

	// ListEvents retrieves a list of events with pagination and sorting.
	ListEvents(ctx context.Context, page, limit int, sortBy string) ([]*types.EventType, error)

	// CountEvents returns the count of events based on the search query.
	CountEvents(ctx context.Context, searchQuery string) (int64, error)
}

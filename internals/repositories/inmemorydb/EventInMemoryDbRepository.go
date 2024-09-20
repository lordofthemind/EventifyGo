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

type inMemoryEventRepository struct {
	mu     sync.RWMutex
	events map[uuid.UUID]*types.EventType
}

// NewInMemoryEventRepository creates a new instance of inMemoryEventRepository.
func NewInMemoryEventRepository() repositories.EventRepositoryInterface {
	return &inMemoryEventRepository{
		events: make(map[uuid.UUID]*types.EventType),
	}
}

func (r *inMemoryEventRepository) CreateEvent(ctx context.Context, event *types.EventType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	event.EventID = uuid.New() // Assign a new UUID
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()
	r.events[event.EventID] = event
	return nil
}

func (r *inMemoryEventRepository) GetEventByID(ctx context.Context, eventID uuid.UUID) (*types.EventType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, exists := r.events[eventID]
	if !exists {
		return nil, errors.New("event not found")
	}
	return event, nil
}

func (r *inMemoryEventRepository) UpdateEvent(ctx context.Context, event *types.EventType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[event.EventID]; !exists {
		return errors.New("event not found")
	}

	event.UpdatedAt = time.Now()
	r.events[event.EventID] = event
	return nil
}

func (r *inMemoryEventRepository) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[eventID]; !exists {
		return errors.New("event not found")
	}

	delete(r.events, eventID)
	return nil
}

func (r *inMemoryEventRepository) SearchEvents(ctx context.Context, searchQuery string, page, limit int, sortBy string) ([]*types.EventType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Basic search by name or description (case insensitive)
	var result []*types.EventType
	for _, event := range r.events {
		if matchEvent(event, searchQuery) {
			result = append(result, event)
		}
	}

	// Pagination
	start := (page - 1) * limit
	if start > len(result) {
		return []*types.EventType{}, nil
	}

	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], nil
}

func (r *inMemoryEventRepository) ListEvents(ctx context.Context, page, limit int, sortBy string) ([]*types.EventType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*types.EventType
	for _, event := range r.events {
		result = append(result, event)
	}

	// Pagination
	start := (page - 1) * limit
	if start > len(result) {
		return []*types.EventType{}, nil
	}

	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], nil
}

func (r *inMemoryEventRepository) CountEvents(ctx context.Context, searchQuery string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var count int64
	for _, event := range r.events {
		if matchEvent(event, searchQuery) {
			count++
		}
	}
	return count, nil
}

// Helper function to match an event with a search query
func matchEvent(event *types.EventType, query string) bool {
	return eventContainsIgnoreCase(event.Name, query) ||
		eventContainsIgnoreCase(event.Description, query) ||
		eventContainsIgnoreCase(event.Location, query)
}

func eventContainsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]storage.Event
}

func (s *Storage) Connect(ctx context.Context, connect string) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func New() *Storage {
	return &Storage{
		mu:     sync.RWMutex{},
		events: make(map[uuid.UUID]storage.Event),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) (uuid.UUID, error) {
	s.mu.RLock()
	for _, e := range s.events {
		if event.DateTime.Equal(e.DateTime) {
			return uuid.Nil, storage.ErrDateBusy
		}
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[event.ID] = *event
	return event.ID, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	s.mu.RLock()
	for _, e := range s.events {
		if event.DateTime.Equal(e.DateTime) && (event.ID != e.ID) {
			return storage.ErrDateBusy
		}
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[event.ID] = *event
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, id)
	return nil
}

func (s *Storage) GetEventID(ctx context.Context, id uuid.UUID) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.events[id]; ok {
		return &e, nil
	}
	return nil, storage.ErrEventNotFound
}

func (s *Storage) GetListEvents(ctx context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]storage.Event, 0)
	for _, e := range s.events {
		e := e
		events = append(events, e)
	}

	return events, nil
}

func (s *Storage) GetListEventsDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	var events []storage.Event
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.events {
		e := e
		if (e.DateTime.Equal(day) || e.DateTime.After(day)) && e.DateTime.Before(day.AddDate(0, 0, 1)) {
			events = append(events, e)
		}
	}
	return events, nil
}

func (s *Storage) GetListEventsWeek(ctx context.Context, beginDate time.Time) ([]storage.Event, error) {
	var events []storage.Event
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.events {
		e := e
		if (e.DateTime.Equal(beginDate) || e.DateTime.After(beginDate)) && e.DateTime.Before(beginDate.AddDate(0, 0, 7)) {
			events = append(events, e)
		}
	}
	return events, nil
}

func (s *Storage) GetListEventsMonth(ctx context.Context, beginDate time.Time) ([]storage.Event, error) {
	var events []storage.Event
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.events {
		e := e
		if (e.DateTime.Equal(beginDate) || e.DateTime.After(beginDate)) && e.DateTime.Before(beginDate.AddDate(0, 1, 0)) {
			events = append(events, e)
		}
	}
	return events, nil
}

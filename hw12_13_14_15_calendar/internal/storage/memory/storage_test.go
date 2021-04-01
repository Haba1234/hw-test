package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

//nolint
func TestStorage(t *testing.T) {
	t.Run("create new event", func(t *testing.T) {
		s := New()
		id := uuid.New()
		event := &storage.Event{
			ID:          id,
			Title:       "Event1",
			Description: "Test event",
		}
		val, err := s.CreateEvent(context.Background(), event)
		require.NoError(t, err)
		require.Equal(t, id, val)
	})

	t.Run("date busy", func(t *testing.T) {
		s := New()
		tn := time.Now()
		id := uuid.New()
		event := &storage.Event{
			ID:          id,
			Title:       "Event1",
			DateTime:    tn,
			Description: "Test event",
		}
		val, err := s.CreateEvent(context.Background(), event)
		require.NoError(t, err)
		require.Equal(t, id, val)

		event = &storage.Event{
			ID:          uuid.New(),
			Title:       "Event2",
			DateTime:    tn,
			Description: "Test event 2",
		}
		_, err = s.CreateEvent(context.Background(), event)
		require.Error(t, err, storage.ErrDateBusy)
	})

	t.Run("update event and get event", func(t *testing.T) {
		s := New()
		id := uuid.New()
		event := &storage.Event{
			ID:          id,
			Title:       "Event1",
			Description: "Test event",
		}
		val, err := s.CreateEvent(context.Background(), event)
		require.NoError(t, err)
		require.Equal(t, id, val)

		event = &storage.Event{
			ID:    id,
			Title: "Event2",
		}
		err = s.UpdateEvent(context.Background(), event)
		require.NoError(t, err)

		e, err := s.GetEventID(context.Background(), id)
		require.NoError(t, err)
		require.Equal(t, "Event2", e.Title)
	})

	t.Run("event not found", func(t *testing.T) {
		s := New()
		_, err := s.GetEventID(context.Background(), uuid.New())
		require.Error(t, err, storage.ErrEventNotFound)
	})

	t.Run("delete event", func(t *testing.T) {
		s := New()
		id := uuid.New()
		event := &storage.Event{
			ID:    id,
			Title: "Event1",
		}
		val, err := s.CreateEvent(context.Background(), event)
		require.NoError(t, err)
		require.Equal(t, id, val)

		err = s.DeleteEvent(context.Background(), id)
		require.NoError(t, err)

		_, err = s.GetEventID(context.Background(), id)
		require.Error(t, err, storage.ErrEventNotFound)
	})

	t.Run("list event", func(t *testing.T) {
		s := New()
		id := uuid.New()
		event := &storage.Event{
			ID:          id,
			Title:       "Event1",
			Description: "Test event",
		}
		val, err := s.CreateEvent(context.Background(), event)
		require.NoError(t, err)
		require.Equal(t, id, val)

		id = uuid.New()
		event = &storage.Event{
			ID:          id,
			Title:       "Event2",
			DateTime:    time.Now(),
			Description: "Test event",
		}
		val, err = s.CreateEvent(context.Background(), event)
		require.NoError(t, err)
		require.Equal(t, id, val)

		e, err := s.GetListEvents(context.Background())
		require.NoError(t, err)
		require.Equal(t, 2, len(e))
	})

	t.Run("event lists", func(t *testing.T) {
		s := New()
		dt1, err := time.Parse("02/01/2006 15:04:05", "01/04/2021 10:00:00")
		if err != nil {
			panic(err)
		}
		event := &storage.Event{
			ID:          uuid.New(),
			Title:       "Событие 1",
			DateTime:    dt1,
			Description: "Test event",
		}
		_, err = s.CreateEvent(context.Background(), event)
		require.NoError(t, err)

		dt2, err := time.Parse("02/01/2006 15:04:05", "01/04/2021 11:00:00")
		if err != nil {
			panic(err)
		}
		event = &storage.Event{
			ID:          uuid.New(),
			Title:       "Событие 2",
			DateTime:    dt2,
			Duration:    time.Duration(2) * time.Hour,
			Description: "Что то будет!",
			UserID:      uuid.New(),
		}
		_, err = s.CreateEvent(context.Background(), event)
		require.NoError(t, err)

		dt3, err := time.Parse("02/01/2006 15:04:05", "02/04/2021 10:00:00")
		if err != nil {
			panic(err)
		}
		event = &storage.Event{
			ID:          uuid.New(),
			Title:       "Событие 3",
			DateTime:    dt3,
			Duration:    time.Duration(24) * time.Hour,
			Description: "Что то будет!",
			UserID:      uuid.New(),
		}
		_, err = s.CreateEvent(context.Background(), event)
		require.NoError(t, err)

		dt4, err := time.Parse("02/01/2006 15:04:05", "02/04/2021 12:00:00")
		if err != nil {
			panic(err)
		}
		event = &storage.Event{
			ID:          uuid.New(),
			Title:       "Событие 4",
			DateTime:    dt4,
			Duration:    time.Duration(5) * time.Hour,
			Description: "Что то будет!",
			UserID:      uuid.New(),
		}
		_, err = s.CreateEvent(context.Background(), event)
		require.NoError(t, err)

		e, err := s.GetListEvents(context.Background())
		require.NoError(t, err)
		require.Equal(t, 4, len(e))

		day, err := time.Parse("02/01/2006 15:04:05", "01/04/2021 00:00:00")
		if err != nil {
			panic(err)
		}
		e, err = s.GetListEventsDay(context.Background(), day)
		require.NoError(t, err)
		require.Equal(t, 2, len(e))

		e, err = s.GetListEventsWeek(context.Background(), day)
		require.NoError(t, err)
		require.Equal(t, 4, len(e))

		e, err = s.GetListEventsMonth(context.Background(), day)
		require.NoError(t, err)
		require.Equal(t, 4, len(e))
	})
}

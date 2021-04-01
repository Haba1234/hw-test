package sqlstorage

import (
	"context"
	"testing"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	dt, err := time.Parse("02/01/2006 15:04:05", "01/04/2021 12:00:00")
	if err != nil {
		panic(err)
	}
	id := uuid.New()

	db := New()

	t.Run("connect db, create date and date busy", func(t *testing.T) {
		err := db.Connect(context.Background(), "user=postgres password=123qwe sslmode=disable")
		require.NoError(t, err)
		defer db.Close(context.Background())

		event := &storage.Event{
			ID:          id,
			Title:       "Событие 1",
			DateTime:    dt,
			Duration:    time.Duration(2) * time.Hour,
			Description: "Тестовое событие",
			UserID:      uuid.New(),
		}
		val, err := db.CreateEvent(context.Background(), event)
		require.NoError(t, err)
		require.Equal(t, id, val)

		_, err = db.CreateEvent(context.Background(), event)
		require.Equal(t, err, storage.ErrDateBusy)
	})

	t.Run("update event and get event", func(t *testing.T) {
		err := db.Connect(context.Background(), "user=postgres password=123qwe sslmode=disable")
		require.NoError(t, err)
		defer db.Close(context.Background())

		event := &storage.Event{
			ID:           id,
			Title:        "Событие 2",
			DateTime:     dt,
			Duration:     time.Duration(2) * time.Hour,
			Description:  "Тестовое событие",
			UserID:       uuid.New(),
			NotifyBefore: 0,
		}

		err = db.UpdateEvent(context.Background(), event)
		require.NoError(t, err)

		e, err := db.GetEventID(context.Background(), id)
		require.NoError(t, err)
		require.Equal(t, "Событие 2", e.Title)
	})

	t.Run("event not found", func(t *testing.T) {
		err := db.Connect(context.Background(), "user=postgres password=123qwe sslmode=disable")
		require.NoError(t, err)
		defer db.Close(context.Background())

		_, err = db.GetEventID(context.Background(), uuid.New())
		require.Error(t, err, storage.ErrEventNotFound)
	})

	t.Run("list events", func(t *testing.T) {
		err := db.Connect(context.Background(), "user=postgres password=123qwe sslmode=disable")
		require.NoError(t, err)
		defer db.Close(context.Background())

		e, err := db.GetListEvents(context.Background())
		require.NoError(t, err)
		require.Equal(t, 1, len(e))
	})

	t.Run("delete event", func(t *testing.T) {
		err := db.Connect(context.Background(), "user=postgres password=123qwe sslmode=disable")
		require.NoError(t, err)
		defer db.Close(context.Background())

		err = db.DeleteEvent(context.Background(), id)
		require.NoError(t, err)

		_, err = db.GetEventID(context.Background(), id)
		require.Error(t, err, storage.ErrEventNotFound)
	})
}

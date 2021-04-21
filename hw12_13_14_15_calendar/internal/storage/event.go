package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("time is busy")
	ErrEventStruct   = errors.New("errors in the event structure")
)

type Storage interface {
	Connect(ctx context.Context, connect string) error
	Close(ctx context.Context) error
	CreateEvent(ctx context.Context, event *Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEventID(ctx context.Context, id uuid.UUID) (*Event, error)
	GetListEvents(ctx context.Context) ([]Event, error)
	GetListEventsDay(ctx context.Context, day time.Time) ([]Event, error)
	GetListEventsWeek(ctx context.Context, beginDate time.Time) ([]Event, error)
	GetListEventsMonth(ctx context.Context, beginDate time.Time) ([]Event, error)
}

type Event struct {
	ID           uuid.UUID     `db:"id"`            // уникальный идентификатор события.
	Title        string        `db:"title"`         // Заголовок - короткий текст.
	DateTime     time.Time     `db:"start_date"`    // Дата и время события.
	Duration     time.Duration `db:"duration"`      // Длительность события (или дата и время окончания).
	Description  string        `db:"description"`   // Описание события - длинный текст, опционально.
	UserID       uuid.UUID     `db:"user_id"`       // ID пользователя, владельца события.
	NotifyBefore int64         `db:"notify_before"` // За сколько времени высылать уведомление, опционально.
}

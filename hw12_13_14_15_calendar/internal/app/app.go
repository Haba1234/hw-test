package app

import (
	"context"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct {
	logger  Logger
	storage storage.Storage
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg, pkg string)
}

type Event struct {
	ID           uuid.UUID     `json:"id"`            // уникальный идентификатор события.
	Title        string        `json:"title"`         // Заголовок - короткий текст.
	DateTime     time.Time     `json:"start_date"`    // Дата и время события.
	Duration     time.Duration `json:"duration"`      // Длительность события (или дата и время окончания).
	Description  string        `json:"description"`   // Описание события - длинный текст, опционально.
	UserID       uuid.UUID     `json:"user_id"`       // ID пользователя, владельца события.
	NotifyBefore int64         `json:"notify_before"` // За сколько времени высылать уведомление, опционально.
}

type Application interface {
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

func New(logger Logger, storage storage.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) Connect(ctx context.Context, dsn string) error {
	return a.storage.Connect(ctx, dsn)
}

func (a *App) Close(ctx context.Context) error {
	a.logger.Info("storage closing...")
	return a.storage.Close(ctx)
}

func (a *App) CreateEvent(ctx context.Context, event *Event) (uuid.UUID, error) {
	if event.Title == "" {
		return uuid.Nil, storage.ErrEventStruct
	}
	if time.Now().After(event.DateTime) {
		return uuid.Nil, storage.ErrEventStruct
	}
	if event.Duration <= 0 {
		return uuid.Nil, storage.ErrEventStruct
	}
	if event.Description == "" {
		return uuid.Nil, storage.ErrEventStruct
	}

	return a.storage.CreateEvent(ctx, &storage.Event{
		ID:           uuid.New(),
		Title:        event.Title,
		DateTime:     event.DateTime,
		Duration:     event.Duration,
		Description:  event.Description,
		UserID:       event.UserID,
		NotifyBefore: event.NotifyBefore,
	})
}

func (a *App) UpdateEvent(ctx context.Context, event *Event) error {
	if event.Title == "" {
		return storage.ErrEventStruct
	}
	if time.Now().After(event.DateTime) {
		return storage.ErrEventStruct
	}
	if event.Duration <= 0 {
		return storage.ErrEventStruct
	}
	if event.Description == "" {
		return storage.ErrEventStruct
	}

	return a.storage.UpdateEvent(ctx, &storage.Event{
		ID:           event.ID,
		Title:        event.Title,
		DateTime:     event.DateTime,
		Duration:     event.Duration,
		Description:  event.Description,
		UserID:       event.UserID,
		NotifyBefore: event.NotifyBefore,
	})
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEventID(ctx context.Context, id uuid.UUID) (*storage.Event, error) {
	return a.storage.GetEventID(ctx, id)
}

func (a *App) GetListEvents(ctx context.Context) ([]storage.Event, error) {
	return a.storage.GetListEvents(ctx)
}

func (a *App) GetListEventsDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	return a.storage.GetListEventsDay(ctx, day)
}

func (a *App) GetListEventsWeek(ctx context.Context, beginDate time.Time) ([]storage.Event, error) {
	return a.storage.GetListEventsWeek(ctx, beginDate)
}

func (a *App) GetListEventsMonth(ctx context.Context, beginDate time.Time) ([]storage.Event, error) {
	return a.storage.GetListEventsMonth(ctx, beginDate)
}

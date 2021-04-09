package app

import (
	"context"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg, pkg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event *storage.Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEventID(ctx context.Context, id uuid.UUID) (*storage.Event, error)
	GetListEvents(ctx context.Context) ([]*storage.Event, error)
	GetListEventsDay(ctx context.Context, day time.Time) ([]*storage.Event, error)
	GetListEventsWeek(ctx context.Context, beginDate time.Time) ([]*storage.Event, error)
	GetListEventsMonth(ctx context.Context, beginDate time.Time) ([]*storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) (uuid.UUID, error) {
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event *storage.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEventID(ctx context.Context, id uuid.UUID) (*storage.Event, error) {
	return a.storage.GetEventID(ctx, id)
}

func (a *App) GetListEvents(ctx context.Context) ([]*storage.Event, error) {
	return a.storage.GetListEvents(ctx)
}

func (a *App) GetListEventsDay(ctx context.Context, day time.Time) ([]*storage.Event, error) {
	return a.storage.GetListEventsDay(ctx, day)
}

func (a *App) GetListEventsWeek(ctx context.Context, beginDate time.Time) ([]*storage.Event, error) {
	return a.storage.GetListEventsWeek(ctx, beginDate)
}

func (a *App) GetListEventsMonth(ctx context.Context, beginDate time.Time) ([]*storage.Event, error) {
	return a.storage.GetListEventsMonth(ctx, beginDate)
}

func (a *App) Info(msg string) {
	a.logger.Info(msg)
}

func (a *App) Warn(msg string) {
	a.logger.Warn(msg)
}

func (a *App) Error(msg string) {
	a.logger.Error(msg)
}

func (a *App) Debug(msg, pkg string) {
	a.logger.Debug(msg, pkg)
}

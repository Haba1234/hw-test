package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //nolint
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) (uuid.UUID, error) {
	_, err := s.GetEventByDate(ctx, event.DateTime)
	if err != nil {
		if errors.Is(err, storage.ErrEventNotFound) {
			query := `INSERT INTO events(id, title, start_date, duration, description, user_id, notify_before)
 				VALUES(:id, :title, :start_date, :duration, :description, :user_id, :notify_before) RETURNING id;`
			rows, err := s.db.NamedQueryContext(ctx, query, event)
			if err != nil {
				return uuid.Nil, err
			}
			defer rows.Close()
			return event.ID, nil
		}
		return uuid.Nil, err
	}
	return uuid.Nil, storage.ErrDateBusy
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	_, err := s.GetEventID(ctx, event.ID)
	if err != nil {
		return err
	}

	query := `UPDATE events SET title=:title, start_date=:start_date, duration=:duration, description=:description, 
					user_id=:user_id, notify_before=:notify_before WHERE id=:id;`
	_, err = s.db.NamedExecContext(ctx, query, event)
	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)
	return err
}

func (s *Storage) GetEventID(ctx context.Context, id uuid.UUID) (*storage.Event, error) {
	event := &storage.Event{}
	err := s.db.GetContext(ctx, event, `SELECT * FROM events WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrEventNotFound
		}
		return nil, err
	}
	return event, nil
}

func (s *Storage) GetEventByDate(ctx context.Context, day time.Time) (*storage.Event, error) {
	event := &storage.Event{}
	err := s.db.GetContext(ctx, event, `SELECT * FROM events WHERE start_date = $1`, day)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrEventNotFound
		}
		return nil, err
	}
	return event, nil
}

func (s *Storage) GetListEvents(ctx context.Context) ([]*storage.Event, error) {
	var events []*storage.Event
	err := s.db.SelectContext(ctx, &events, `SELECT * FROM events`)
	if err != nil {
		return nil, err
	}
	return events, err
}

func (s *Storage) GetListEventsDay(ctx context.Context, day time.Time) ([]*storage.Event, error) {
	var events []*storage.Event
	err := s.db.SelectContext(ctx, &events, `SELECT * FROM events WHERE start_date BETWEEN $1 AND $1 + (interval '1d')`, day)
	if err != nil {
		return nil, err
	}
	return events, err
}

func (s *Storage) GetListEventsWeek(ctx context.Context, beginDate time.Time) ([]*storage.Event, error) {
	var events []*storage.Event
	err := s.db.SelectContext(ctx, &events,
		`SELECT * FROM events WHERE start_date BETWEEN $1 AND $1 + (interval '7 weeks')`, beginDate)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Storage) GetListEventsMonth(ctx context.Context, beginDate time.Time) ([]*storage.Event, error) {
	var events []*storage.Event
	err := s.db.SelectContext(ctx, &events,
		`SELECT * FROM events WHERE start_date BETWEEN $1 AND $1 + (interval '1 months')`, beginDate)
	if err != nil {
		return nil, err
	}
	return events, nil
}

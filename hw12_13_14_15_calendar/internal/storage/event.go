package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("time is busy")
)

type Event struct {
	ID           uuid.UUID     `db:"id"`            // уникальный идентификатор события.
	Title        string        `db:"title"`         // Заголовок - короткий текст.
	DateTime     time.Time     `db:"start_date"`    // Дата и время события.
	Duration     time.Duration `db:"duration"`      // Длительность события (или дата и время окончания).
	Description  string        `db:"description"`   // Описание события - длинный текст, опционально.
	UserID       uuid.UUID     `db:"user_id"`       // ID пользователя, владельца события.
	NotifyBefore int64         `db:"notify_before"` // За сколько времени высылать уведомление, опционально.
}

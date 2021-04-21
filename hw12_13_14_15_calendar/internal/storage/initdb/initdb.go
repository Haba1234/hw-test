package initdb

import (
	"context"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
)

func New(ctx context.Context, typeDB, connect string) (storage.Storage, error) {
	var db storage.Storage
	switch typeDB {
	case "sql":
		db = sqlstorage.New()
		err := db.Connect(ctx, connect)
		if err != nil {
			return nil, err
		}
	default:
		db = memorystorage.New()
	}

	return db, nil
}

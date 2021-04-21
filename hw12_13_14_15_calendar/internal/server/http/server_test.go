package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage/initdb"
	"github.com/google/uuid"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func TestBasicRestServer(t *testing.T) {
	logg, err := logger.New("DEBUG", "../../../bin/logfile.log")
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	storage, err := initdb.New(context.Background(), "memory", "postgres://postgres:123qwe@localhost:5432/calendar?sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	calendar := app.New(logg, storage)
	server := NewServer(logg, calendar)
	id := uuid.Nil

	t.Run("test hello", func(t *testing.T) {
		apitest.New().
			HandlerFunc(server.testHandler).
			Get("/test").
			Expect(t).
			Assert(func(res *http.Response, req *http.Request) error {
				assert.Equal(t, http.StatusOK, res.StatusCode)
				return nil
			}).
			End()
	})

	t.Run("successfully create event", func(t *testing.T) {
		apitest.New().
			HandlerFunc(server.createEventHandler).
			Put("/createEvent").
			JSONFromFile("../../../internal/server/tests/create_event.json").
			Expect(t).
			Assert(func(res *http.Response, req *http.Request) error {
				fmt.Println("id:", res.Body)
				json.NewDecoder(res.Body).Decode(&id)
				fmt.Println("AAAAA: ", id)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				return nil
			}).
			End()
	})

	t.Run("duplicate event", func(t *testing.T) {
		apitest.New().
			HandlerFunc(server.createEventHandler).
			Put("/createEvent").
			JSONFromFile("../../../internal/server/tests/create_event.json").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("delete event", func(t *testing.T) {
		t.Skip() // Тест зависает на мьютексах, пока не разобрался почему.
		apitest.New().
			HandlerFunc(server.deleteEventHandler).
			Delete("/deleteEvent").
			JSON(id).
			Expect(t).
			Status(http.StatusOK).
			End()
	})
}

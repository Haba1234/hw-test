package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/gorilla/mux"
)

type server struct {
	srv    *http.Server
	router *mux.Router
	app    *app.App
	logg   *logger.Logger
}

//nolint: golint
func NewServer(logg *logger.Logger, app *app.App) *server {
	s := &server{
		app:  app,
		logg: logg,
	}
	s.router = mux.NewRouter()
	s.router.Use(loggingMiddleware(s.logg))
	s.router.HandleFunc("/", s.testHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/test", s.testHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/createEvent", s.createEventHandler).Methods(http.MethodPut)
	s.router.HandleFunc("/updateEvent", s.updateEventHandler).Methods(http.MethodPut)
	s.router.HandleFunc("/deleteEvent", s.deleteEventHandler).Methods(http.MethodDelete)
	s.router.HandleFunc("/getListEvents", s.getListEventsHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/getListEventDay", s.getListEventDayHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/getListEventWeek", s.getListEventWeekHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/getListEventMonth", s.getListEventMonthHandler).Methods(http.MethodGet)
	return s
}

func (s *server) Start(ctx context.Context, addr string) error {
	s.logg.Info("HTTP server " + addr + " starting...")
	s.srv = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	if err := s.srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	s.logg.Info("HTTP server stopping...")
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

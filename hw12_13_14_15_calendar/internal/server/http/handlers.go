package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/google/uuid"
)

func (s *server) testHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" || req.URL.Path == "/test" {
		fmt.Fprintf(w, "Welcome to the home page!\n")
		return
	}
	http.NotFound(w, req)
}

func (s *server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	event := app.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := s.app.CreateEvent(r.Context(), &event)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to create event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error in sending response: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	event := app.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.app.UpdateEvent(r.Context(), &event)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to update event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id := uuid.Nil
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.app.DeleteEvent(r.Context(), id)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to delete event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *server) getListEventsHandler(w http.ResponseWriter, r *http.Request) {
	listEvents, err := s.app.GetListEvents(r.Context())
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get evens list: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(listEvents)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error in sending response: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//nolint: dupl
func (s *server) getListEventDayHandler(w http.ResponseWriter, r *http.Request) {
	var day time.Time
	err := json.NewDecoder(r.Body).Decode(&day)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	listEvents, err := s.app.GetListEventsDay(r.Context(), day)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get evens list day: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(listEvents)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error in sending response: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//nolint: dupl
func (s *server) getListEventWeekHandler(w http.ResponseWriter, r *http.Request) {
	var beginDate time.Time
	err := json.NewDecoder(r.Body).Decode(&beginDate)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	listEvents, err := s.app.GetListEventsWeek(r.Context(), beginDate)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get evens list week: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(listEvents)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error in sending response: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//nolint: dupl
func (s *server) getListEventMonthHandler(w http.ResponseWriter, r *http.Request) {
	var beginDate time.Time
	err := json.NewDecoder(r.Body).Decode(&beginDate)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	listEvents, err := s.app.GetListEventsMonth(r.Context(), beginDate)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get evens list month: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(listEvents)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error in sending response: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

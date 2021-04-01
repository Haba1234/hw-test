package internalhttp

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	srv *http.Server
	app Application
}

type Application interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg, pkg string)
}

func NewServer(addr string, app Application) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(app, http.HandlerFunc(TestHandler)))
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		app: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.app.Debug("Server starting...", "internalhttp")
	if err := s.srv.ListenAndServe(); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.app.Debug("Server stopping...", "internalhttp")
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func TestHandler(w http.ResponseWriter, req *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "Welcome to the home page!")
}

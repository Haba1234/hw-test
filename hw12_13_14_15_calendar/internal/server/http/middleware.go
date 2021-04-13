package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/gorilla/mux"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func loggingMiddleware(logg *logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			wrapperWriter := newLoggingResponseWriter(w)
			next.ServeHTTP(wrapperWriter, r)
			latency := time.Since(startTime)
			logg.Info(fmt.Sprintf("%s %s %s %s %s %d %dms %s",
				r.RemoteAddr,
				fmt.Sprint("[", startTime.Format("02/Jan/2006 15:04:05 -0700"), "]"),
				r.Method,
				r.RequestURI,
				r.Proto,
				wrapperWriter.statusCode,
				latency.Microseconds(),
				r.UserAgent()))
		})
	}
}

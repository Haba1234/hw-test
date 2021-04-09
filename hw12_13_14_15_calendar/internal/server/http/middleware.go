package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func loggingMiddleware(app Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wrapperWriter := NewLoggingResponseWriter(w)
		next.ServeHTTP(wrapperWriter, r)
		latency := time.Since(startTime)
		app.Info(fmt.Sprintf("%s %s %s %s %s %d %dms %s",
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

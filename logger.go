package main

import (
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		inner.ServeHTTP(rec, r)

		slog.Info("request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"uri", r.RequestURI,
			"route", name,
			"status", rec.status,
			"duration", time.Since(start),
		)
	})
}

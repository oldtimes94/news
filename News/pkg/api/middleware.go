package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ipAddress := r.RemoteAddr
		requestID, _ := r.Context().Value("request_id").(string)

		lrw := &loggingResponseWriter{w, http.StatusOK}
		next.ServeHTTP(lrw, r)

		logEntry := fmt.Sprintf("URI %s Request duration %s - IP Address %s - Status Code %d - Request ID %s",
			r.RequestURI,
			time.Since(startTime),
			ipAddress,
			lrw.statusCode,
			requestID,
		)
		log.Println(logEntry)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

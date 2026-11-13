package main

import (
	"log"
	"net/http"
)

type BufferedResponseWriter struct {
	statusCode int
	http.ResponseWriter
}

func NewBufferedResponseWriter(w http.ResponseWriter) *BufferedResponseWriter {
	return &BufferedResponseWriter{statusCode: http.StatusOK, ResponseWriter: w}
}

func (brw *BufferedResponseWriter) WriteHeader(statusCode int) {
	brw.statusCode = statusCode
	brw.ResponseWriter.WriteHeader(statusCode)
}

func addLoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("Request:", r.Method, r.URL.Path)
			brw := NewBufferedResponseWriter(w)
			h.ServeHTTP(brw, r)
			log.Println("Response:", brw.statusCode)
		},
	)
}

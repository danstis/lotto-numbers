// Package middleware provides HTTP middleware functions for request logging and processing.
package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// LoggingMiddleware defines the middleware for logging HTTP requests.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now() // Capture the start time
		next.ServeHTTP(w, r)
		duration := time.Since(startTime) // Calculate the duration

		// Log the request details with the time taken to serve the page, without the port number
		ipAddress := strings.Split(r.RemoteAddr, ":")[0]
		log.Printf("Request from %s: %s %s, Duration: %v",
			ipAddress, r.Method, r.URL.Path, duration)
	})
}

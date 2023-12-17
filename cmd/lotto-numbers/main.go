package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danstis/lotto-numbers/internal/handlers" // Import the handlers package
	"github.com/danstis/lotto-numbers/internal/version"
	"github.com/gorilla/mux"
)

// Main entry point for the app.
func main() {
	log.Printf("Lotto Numbers started - v%v", version.Version)

	// Set up the HTTP server using Gorilla Mux
	r := mux.NewRouter()
	// Define the logging middleware with timing
	loggingMiddleware := func(next http.Handler) http.Handler {
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

	// Wrap the router with the logging middleware
	r.Use(loggingMiddleware)

	handlers.SetupRoutes(r)

	// Retrieve the port number from an environment variable or use 8080 as default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8091"
	}

	// Start the server
	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

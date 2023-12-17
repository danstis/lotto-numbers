package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danstis/lotto-numbers/internal/handlers" // Import the handlers package
	"github.com/danstis/lotto-numbers/internal/version"
	"github.com/gorilla/mux"
)

// Main entry point for the app.
func main() {
	log.Printf("Lotto Numbers started - v%v", version.Version)

	// Set up the HTTP server using Gorilla Mux
	r := mux.NewRouter()
	// Define the logging middleware
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wrap the ResponseWriter to capture the status code
			wrapper := &statusResponseWriter{ResponseWriter: w}
			next.ServeHTTP(wrapper, r)

			// Log the request details
			log.Printf("Request from %s: %s %s, Response code: %d",
				r.RemoteAddr, r.Method, r.URL.Path, wrapper.statusCode)
		})
	}

	// Create a new type that wraps http.ResponseWriter to capture the status code
	type statusResponseWriter struct {
		http.ResponseWriter
		statusCode int
	}

	// Override the WriteHeader method to capture the status code
	func (w *statusResponseWriter) WriteHeader(code int) {
		w.statusCode = code
		w.ResponseWriter.WriteHeader(code)
	}

	// Initialize the status code with 200 in case WriteHeader is not called
	func (w *statusResponseWriter) Write(b []byte) (int, error) {
		if w.statusCode == 0 {
			w.statusCode = http.StatusOK
		}
		return w.ResponseWriter.Write(b)
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

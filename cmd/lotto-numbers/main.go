package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danstis/lotto-numbers/internal/handlers"
	"github.com/danstis/lotto-numbers/internal/middleware"
	"github.com/danstis/lotto-numbers/internal/version"
	"github.com/gorilla/mux"
)

// Main entry point for the app.
func main() {
	log.Printf("Lotto Numbers started - v%v", version.Version)

	// Set up the HTTP server using Gorilla Mux
	r := mux.NewRouter()

	// Wrap the router with the logging middleware
	r.Use(middleware.LoggingMiddleware)

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

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
	log.Printf("Version %q", version.Version)

	// Set up the HTTP server using Gorilla Mux
	r := mux.NewRouter()
	r.HandleFunc("/lottery-numbers", handlers.GetLotteryNumbers).Methods("GET")

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

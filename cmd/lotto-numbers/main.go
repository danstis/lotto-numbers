package main

import (
	"log"
	"net/http"

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

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

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

	num := generateLotteryNumbers([]int{1, 1, 1, 1, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29}, 5, 5)
	log.Printf("Numbers: %v", num)

	// Set up the HTTP server using Gorilla Mux
	r := mux.NewRouter()
	r.HandleFunc("/lottery-numbers", handlers.GetLotteryNumbers).Methods("GET")

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

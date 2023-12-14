package main

import (
	"log"
	"net/http"

	"github.com/danstis/lotto-numbers/internal/handlers" // Import the handlers package
	"github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    // Set up routes
    r.HandleFunc("/lottery-numbers", handlers.GetLotteryNumbers).Methods("GET")

    // Start the server
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal("ListenAndServe error: ", err)
    }
}

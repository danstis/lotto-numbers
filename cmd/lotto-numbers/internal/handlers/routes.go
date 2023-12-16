package handlers

import (
	"github.com/gorilla/mux"
)

// SetupRoutes configures the routes for the application.
func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/numbers", GetLotteryNumbers).Methods("GET")
}

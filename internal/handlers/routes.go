package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configures the routes for the application.
func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})
	r.HandleFunc("/numbers", GetLotteryNumbers).Methods("GET")
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
}

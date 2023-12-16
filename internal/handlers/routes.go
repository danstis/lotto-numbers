package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configures the routes for the application.
func SetupRoutes(r *mux.Router) {
	webDir := "../../web"
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, webDir+"/index.html")
	})
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(webDir+"/assets"))))
	r.HandleFunc("/numbers", GetLotteryNumbers).Methods("GET")
}
